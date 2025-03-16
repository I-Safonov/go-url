package main

import (
	"url-shortener/internal/api"
	"url-shortener/internal/auth"
	"url-shortener/internal/storage"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	shards := storage.InitShards()
	rdb := storage.NewRedisClient()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(api.MetricsMiddleware())

	authGroup := e.Group("/api")
	authGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(auth.SecretKey),
		TokenLookup: "header:Authorization",
	}))

	authGroup.POST("/shorten", api.ShortenURL(shards, rdb))
	e.GET("/:code", api.RedirectURL(shards, rdb))

	e.GET("/metrics", api.PrometheusHandler())

	e.Logger.Fatal(e.Start(":1323"))
}
