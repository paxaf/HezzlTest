package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paxaf/HezzlTest/config"
	"github.com/paxaf/HezzlTest/internal/controller"
	natsClient "github.com/paxaf/HezzlTest/internal/infrastructure/nats"
	"github.com/paxaf/HezzlTest/internal/logger"
	"github.com/paxaf/HezzlTest/internal/repository"
	clickHouse "github.com/paxaf/HezzlTest/internal/repository/clickhouse"
	"github.com/paxaf/HezzlTest/internal/repository/events"
	"github.com/paxaf/HezzlTest/internal/repository/postgres"
	redisClient "github.com/paxaf/HezzlTest/internal/repository/redis"
	"github.com/paxaf/HezzlTest/internal/usecase"
	"github.com/paxaf/HezzlTest/internal/worker"
)

type App struct {
	config    *config.Config
	apiServer *http.Server
	closer    *closer
	router    *gin.Engine
	logger    *logger.Logger
	work      *worker.ClickHouseWorker
}

func New(cfg *config.Config) (*App, error) {
	app := &App{}
	app.config = cfg
	app.router = gin.Default()

	app.logger = logger.New(cfg.Logger.Level)

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.GetDSN())
	if err != nil {
		app.logger.Error(err, "database connection error: %v")
	}
	pgpool := postgres.New(pool)
	rclient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	redisClient := redisClient.New(rclient)

	ns, err := natsClient.New(app.config.Nats)
	if err != nil {
		logger.Fatal("failed create conn to nats", err)
	}
	event := events.New(ns)
	repo := repository.New(redisClient, pgpool, event)
	service := usecase.New(repo)
	handler := controller.New(service)

	app.router.GET("/goods", handler.GetAll)
	app.router.GET("/goods/:id", handler.GetItem)
	app.router.GET("/goods/search/:name", handler.GetItemsByName)
	app.router.GET("/:project_id/goods", handler.GetItemsByProject)
	app.router.PATCH("/goods", handler.UpdateItem)
	app.router.POST("/goods", handler.CreateItem)
	app.router.DELETE("/goods/:id", handler.DeleteItem)
	app.router.GET("/projects", handler.GetProjects)
	app.router.GET("/projects/:id", handler.GetProject)
	app.router.POST("/projects", handler.CreateProject)
	app.router.PATCH("/projects/:id", handler.UpdateProject)
	app.router.DELETE("/projects/:id", handler.DeleteProject)

	host := app.config.APIServer.Host
	port := app.config.APIServer.Port
	addr := net.JoinHostPort(host, port)
	app.apiServer = &http.Server{
		Addr:              addr,
		Handler:           app.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ch, err := clickHouse.NewClickHouse(app.config.Clickhouse)
	if err != nil {
		logger.Fatal("failed conn ch", err)
	}
	work, err := worker.NewClickHouseWorker(ns.Conn, ch.Conn, "db.events.>")
	if err != nil {
		logger.Fatal("failed init worker", err)
	}
	app.closer = NewCloser(pgpool, redisClient, event, work)
	app.work = work
	app.logger.Info("Application initialized successfully")
	return app, nil
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		app.logger.Info("API server started successfully", "address", app.apiServer.Addr)
		if err := app.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal(err, "failed to start the server: %v")
		}
	}()
	go app.work.Start(ctx)

	<-ctx.Done()
	app.logger.Info("Received shutdown signal")

	return nil
}

func (app *App) Close() error {
	err := app.closer.Close(app)
	if err != nil {
		return err
	}
	return nil
}
