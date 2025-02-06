package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoDBConnectionString() string {
	connURL := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(Conf.MongoDB.User, Conf.MongoDB.Pass),
		Host:   fmt.Sprintf("%s:%d", Conf.MongoDB.Host, Conf.MongoDB.Port),
		Path:   Conf.MongoDB.Database,
	}
	q := connURL.Query()
	q.Add("authSource", "admin")
	q.Add("readPreference", "primary")
	q.Add("ssl", "false")
	connURL.RawQuery = q.Encode()

	return connURL.String()
}

func loadMongoDB(ctx context.Context) *mongo.Database {
	connectionString := getMongoDBConnectionString()

	ctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()
	opts := options.Client()
	opts.ApplyURI(connectionString)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("unable to create connection pool",
			slog.Any("err", err.Error()))
		os.Exit(-1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("unable to ping mongo db client",
			slog.Any("err", err.Error()))
		os.Exit(-1)
	}

	slog.Info("successfully connected to mongodb",
		"mongodb", Conf.MongoDB.Database)

	return client.Database(Conf.MongoDB.Database)
}
