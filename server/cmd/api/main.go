package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luponetn/noitrex/internal/config"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/internal/logger"
)

type App struct {
	Cfg    *config.Config
	Router *gin.Engine
	Db     *pgxpool.Pool
}

func main() {
	logger.NewLogger("development")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", slog.String("error", err.Error()))
	}

	dbPool, err := db.NewDBPool(cfg.DBUrl)
	if err != nil {
		slog.Error("could not startup db", slog.String("error", err.Error()))
	}
	defer dbPool.Close()

	queries := db.New(dbPool)

	router := CreateRouter()

	CreateRoutes(router, queries, cfg.JWTAccessSecret, cfg.JWTRefreshSecret)

	app := &App{
		Cfg:    cfg,
		Router: router,
		Db:     dbPool,
	}

	if err := StartupServer(router, app.Cfg.Port); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
