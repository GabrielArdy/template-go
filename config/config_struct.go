package config

import (
	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	App            App            `mapstructure:"app"`
	Server         Server         `mapstructure:"server"`
	ServiceAccount ServiceAccount `mapstructure:"serviceAccount"`
	Atlas          Atlas          `mapstructure:"atlas"`
	Redis          Redis          `mapstructure:"redis"`
}

type App struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type MongoDB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"username"`
	Pass     string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Atlas struct {
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	User string `mapstructure:"user"`
}

type ServiceAccount struct {
	Path string `mapstructure:"path"`
}

type Clients struct {
	MongoDB      *mongo.Database
	Firestore    *firestore.Client
	FirebaseAuth *auth.Client
	Redis        redis.UniversalClient
}
