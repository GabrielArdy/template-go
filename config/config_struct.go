package config

import "go.mongodb.org/mongo-driver/mongo"

type Config struct {
	App     App     `mapstructure:"app"`
	Server  Server  `mapstructure:"server"`
	MongoDB MongoDB `mapstructure:"mongodb"`
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

type Clients struct {
	MongoDB *mongo.Database
}
