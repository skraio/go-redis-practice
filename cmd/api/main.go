package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv" // use for testing locally
	"github.com/redis/go-redis/v9"
	"github.com/skraio/go-redis-practice/internal/data"
)

type config struct {
	port  int
	redis struct {
		addr     string
		password string
		db       int
	}
}

type application struct {
	config  config
	logger  *slog.Logger
	records data.RecordModel
	wg      sync.WaitGroup
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// use for testing locally
	// err := godotenv.Load(".test.env")
	// if err != nil {
	// 	logger.Error("Error loading .test.env file")
	// 	os.Exit(1)
	// }

	portEnv := os.Getenv("API_SERVER_PORT")
	if portEnv == "" {
		logger.Info("API_SERVER_PORT environment variable is not set. Using default port 8080.")
		portEnv = "8080"
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		logger.Error(err.Error())
		os.Exit(1)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		logger.Error(err.Error())
		os.Exit(1)
	}

	redisDBEnv := os.Getenv("REDIS_DB")
	if redisDBEnv == "" {
		logger.Info("REDIS_DB environment variable is not set. Using default value 0.")
		redisDBEnv = "0"
	}

	redisDB, err := strconv.Atoi(redisDBEnv)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var cfg config
	cfg.port = port
	cfg.redis.addr = redisAddr
	cfg.redis.password = redisPassword
	cfg.redis.db = redisDB

	rd, err := openRedisDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer rd.Close()
	logger.Info("redis connetion pool established")

	app := &application{
		config:  cfg,
		logger:  logger,
		records: data.RecordModel{RedisDB: rd},
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openRedisDB(cfg config) (*redis.Client, error) {
	rd := redis.NewClient(&redis.Options{
		Addr:     cfg.redis.addr,
		Password: cfg.redis.password,
		DB:       cfg.redis.db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rd.Ping(ctx).Result()
	if err != nil {
		rd.Close()
		return nil, err
	}

	return rd, nil
}
