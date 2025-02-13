package config

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoDBConnectionString() string {
	// Format: mongodb+srv://<username>:<password>@<cluster>.mongodb.net
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		url.QueryEscape(Conf.Atlas.User),
		url.QueryEscape(Conf.Atlas.Password),
		Conf.Atlas.Host,
		url.QueryEscape(Conf.Atlas.Database),
	)
}

func loadMongoDB(ctx context.Context) *mongo.Database {
	connectionString := getMongoDBConnectionString()

	ctx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()
	opts := options.Client().
		ApplyURI(connectionString).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("unable to create connection pool",
			slog.Any("err", err.Error()))
		os.Exit(-1)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		slog.Error("unable to ping the deployment",
			slog.Any("err", err.Error()))
		os.Exit(-1)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	slog.Info("mongo db client connected successfully to atlas", slog.Any("db", Conf.Atlas.Database))
	return client.Database(Conf.Atlas.Database)
}
