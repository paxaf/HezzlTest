package integration_test

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresSuite struct {
	suite.Suite
	pgContainer testcontainers.Container
	PgPool      *pgxpool.Pool
	repo        *postgres.PgPool
}

func TestPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(PostgresSuite))
}

func (s *PostgresSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(30 * time.Second),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(s.T(), err)
	s.pgContainer = pgContainer

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	connStr := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	err = applyGooseMigrations(connStr, "../migrations")
	assert.NoError(s.T(), err, "Failed to apply migrations")

	pool, err := pgxpool.New(ctx, connStr)
	s.PgPool = pool
	assert.NoError(s.T(), err)

	s.repo = postgres.New(pool)
}

func (s *PostgresSuite) TearDownTest() {
	_, _ = s.PgPool.Exec(context.Background(), "TRUNCATE TABLE GOODS RESTART IDENTITY")
}

func (s *PostgresSuite) TearDownSuite() {
	if s.pgContainer != nil {
		_ = s.pgContainer.Terminate(context.Background())
	}
}

func applyGooseMigrations(connStr string, migrationsDir string) error {
	cmd := exec.Command(
		"goose",
		"-dir", migrationsDir,
		"postgres", connStr,
		"up",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("goose failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (s *PostgresSuite) TestWriteGoods() {
	ctx := context.Background()
	item := &entity.Goods{
		ProjectId: 1,
		Name:      "First item",
	}
	err := s.repo.CreateItem(ctx, item)
	require.NoError(s.T(), err, "create item")
	getItem, err := s.repo.GetItem(ctx, 1)
	require.NoError(s.T(), err, "get item")
	assert.Equal(s.T(), getItem.ProjectId, item.ProjectId)
	assert.Equal(s.T(), getItem.Name, item.Name)
	assert.Equal(s.T(), getItem.Priority, 1)
	updItem := &entity.Goods{
		ProjectId:   1,
		Name:        "Updated name",
		Description: "New description",
		Priority:    10,
	}
	err = s.repo.UpdateItem(ctx, updItem)
	require.NoError(s.T(), err)
	getItem, err = s.repo.GetItem(ctx, item.ProjectId)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), getItem.ProjectId, updItem.ProjectId)
	assert.Equal(s.T(), getItem.Name, updItem.Name)
	assert.Equal(s.T(), getItem.Priority, updItem.Priority)
	err = s.repo.DeleteItem(ctx, getItem.Id)
	require.NoError(s.T(), err)
	err = s.repo.UpdateItem(ctx, updItem)
	require.ErrorAs(s.T(), err, &pgx.ErrNoRows)
	nilItem, err := s.repo.GetItem(ctx, 1)
	require.Nil(s.T(), nilItem)
	require.ErrorAs(s.T(), err, &pgx.ErrNoRows)
	err = s.repo.DeleteItem(ctx, 10)
	require.ErrorAs(s.T(), err, &pgx.ErrNoRows)
	item.ProjectId = 2
	err = s.repo.CreateItem(ctx, item)
	require.Error(s.T(), err)
}

func (s *PostgresSuite) TestReadGoods() {
	ctx := context.Background()
	goods := []*entity.Goods{
		{
			ProjectId: 1,
			Name:      "First item",
		},
		{
			ProjectId:   1,
			Name:        "Second item",
			Description: "new desc",
		},
		{
			ProjectId: 1,
			Name:      "Three item",
		},
	}
	for _, val := range goods {
		err := s.repo.CreateItem(ctx, val)
		require.NoError(s.T(), err)
	}
	allItems, err := s.repo.GetAllItems(ctx)
	require.NoError(s.T(), err)
	assert.Len(s.T(), allItems, 3)
	namedItems, err := s.repo.GetItemsByName(ctx, "item")
	require.NoError(s.T(), err)
	assert.Len(s.T(), namedItems, 3)
	namedItems, err = s.repo.GetItemsByName(ctx, "three")
	require.NoError(s.T(), err)
	assert.Len(s.T(), namedItems, 1)
	assert.Equal(s.T(), namedItems[0].Name, goods[2].Name)
}
