package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func loadRedisClient(ctx context.Context) *redis.Client {
	clients := redis.NewClient(&redis.Options{})
	_, err := clients.Ping(ctx).Result()
	if err != nil {
		slog.Error("unable to ping the deployment",
			slog.Any("err", err.Error()))
		os.Exit(-1)
	}

	slog.Info("redis client connected successfully", slog.Any("host", Conf.Redis.Host))
	return clients
}
