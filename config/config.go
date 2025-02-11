package config

import (
	"context"
	"flag"
	"fmt"
	"go-scratch/apis"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var locJakarta *time.Location
var Profile string
var Conf Config
var Cli Clients

func loadTimezone() {
	os.Setenv("TZ", "Asia/Jakarta")
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		slog.Info("Failed to load timezone", slog.Any("error", err))
		os.Exit(-1)
	}
	time.Local = loc
	locJakarta = loc
}

func Load(ctx context.Context) {
	loadTimezone()
	loadConfigFile()
	Cli = Clients{
		MongoDB:      loadMongoDB(ctx),
		Firestore:    LoadFirestore(ctx, FirebaseApp),
		FirebaseAuth: LoadFirebaseAuth(ctx, FirebaseApp),
	}
}

func loadConfigFile() {
	profile := flag.String("active.profile", "no-file", "state the active profile")
	flag.Parse()

	if *profile == "no-file" {
		p, ok := os.LookupEnv("ACTIVE_PROFILE")
		if !ok {
			slog.Info("config file not found")
			os.Exit(-1)
		}
		profile = &p

	}
	Profile = *profile

	path := fmt.Sprintf("config-%s.yml", Profile)

	viper := viper.New()
	viper.AddConfigPath(".")
	viper.SetConfigName(path)
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		slog.Info("Failed to read config file", slog.Any("error", err.Error()))
		os.Exit(-1)
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		slog.Info("Failed to unmarshal config file", slog.Any("error", err.Error()))
		os.Exit(-1)
	}

	InitLogger()
	slog.Info("Application Running using config file", slog.String("config", "config-"+Profile))
}

func LoadEcho() *echo.Echo {
	e := echo.New()
	e.Use(SlogMiddleware())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%v", Conf.Server.Port)},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin},
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/api.yml", func(c echo.Context) error {
		data, err := apis.Api.ReadFile("api.yml")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Could not read api.yml")
		}
		return c.Blob(http.StatusOK, "application/yaml", data)
	})

	pprof.Register(e)
	return e

}
