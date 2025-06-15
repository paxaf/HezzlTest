package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/paxaf/HezzlTest/internal/entity"
	"github.com/paxaf/HezzlTest/internal/repository"
	redisClient "github.com/paxaf/HezzlTest/internal/repository/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisSuite struct {
	suite.Suite
	redisContainer testcontainers.Container
	Client         *redis.Client
	repo           repository.Redis
}

func TestRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, new(RedisSuite))
}

func (s *RedisSuite) SetupSuite() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(s.T(), err, "Failed to start Redis container")
	s.redisContainer = redisContainer

	endpoint, err := redisContainer.Endpoint(ctx, "")
	require.NoError(s.T(), err, "Failed to get Redis endpoint")

	s.Client = redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	s.repo = redisClient.New(s.Client)
	_, err = s.Client.Ping().Result()
	require.NoError(s.T(), err, "Failed to connect to Redis")
}

func (s *RedisSuite) TearDownTest() {
	_ = s.repo.CleanCache()
}

func (s *RedisSuite) TearDownSuite() {
	if s.redisContainer != nil {
		_ = s.redisContainer.Terminate(context.Background())
	}
}

func (s *RedisSuite) TestItemRedis() {
	item := &entity.Goods{
		ProjectId: 1,
		Name:      "redis test",
	}
	err := s.repo.RedisSetItem("1", item)
	require.NoError(s.T(), err)
	getItem, err := s.repo.RedisGetItem("1")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), item, getItem)
}

func (s *RedisSuite) TestItemsRedis() {
	goods := []entity.Goods{
		{ProjectId: 1,
			Name: "redis test 1"},
		{ProjectId: 1,
			Name: "redis test 2"},
		{ProjectId: 1,
			Name: "redis test 3"},
		{ProjectId: 1,
			Name: "redis test 4"},
	}
	err := s.repo.RedisSetItem("1", goods)
	require.NoError(s.T(), err)

	result, err := s.repo.RedisGetItems("1")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), result, goods)
}
