package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luponetn/noitrex/internal/billing"
	"github.com/luponetn/noitrex/internal/broker"
	"github.com/luponetn/noitrex/internal/config"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/internal/logger"
	"github.com/luponetn/noitrex/internal/metering"
	"github.com/luponetn/noitrex/internal/redis"
	webhook "github.com/luponetn/noitrex/internal/webhooks"
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

	broker := broker.NewBroker()
	defer broker.Shutdown()

	redisClient, err := redis.NewClient(cfg.RedisUrl)
	if err != nil {
		slog.Error("could not startup redis", "error", err.Error())
	}
	defer redisClient.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	meteringEngine := metering.NewMeteringEngine(broker, redisClient, queries)
	go meteringEngine.Start(ctx)

	billingEngine := billing.NewBilling(queries, broker)
	go billingEngine.Start(ctx)

	webhookCtx, webhookCancel := context.WithCancel(ctx)
	defer webhookCancel()
	webhookEngine := webhook.NewWebhookEngine(queries, broker)
	go webhookEngine.Start(webhookCtx)

	CreateRoutes(router, queries, cfg.JWTAccessSecret, cfg.JWTRefreshSecret, broker)

	app := &App{
		Cfg:    cfg,
		Router: router,
		Db:     dbPool,
	}

	if err := StartupServer(router, app.Cfg.Port); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
