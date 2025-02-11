package config

import (
	"context"
	"log/slog"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func LoadFirebaseApp(ctx context.Context) {
	opt := option.WithCredentialsFile(Conf.ServiceAccount.Path)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		slog.Error("Failed to create firebase app", slog.Any("error", err))
		os.Exit(-1)
	}
	FirebaseApp = app
}

func LoadFirebaseAuth(ctx context.Context, app *firebase.App) *auth.Client {
	client, err := app.Auth(ctx)
	if err != nil {
		slog.Error("Failed to create firebase auth client", slog.Any("error", err))
		os.Exit(-1)
	}
	slog.Info("Firebase auth client created")
	return client
}

func LoadFirestore(ctx context.Context, app *firebase.App) *firestore.Client {
	client, err := app.Firestore(ctx)
	if err != nil {
		slog.Error("Failed to create firestore client", slog.Any("error", err))
		os.Exit(-1)
	}
	slog.Info("Firestore client created")
	return client
}
