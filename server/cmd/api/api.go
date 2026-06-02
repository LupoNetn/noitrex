package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luponetn/nexusmq/pkg/broker"
	"github.com/luponetn/noitrex/internal/auth"
	"github.com/luponetn/noitrex/internal/customers"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/internal/events"
	operator "github.com/luponetn/noitrex/internal/operators"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()

	return router
}

func CreateRoutes(router *gin.Engine, queries db.Querier, JWTAccessSecret, JWTRefreshSecret string, broker broker.Broker) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "healthy",
		})
	})

	operatorService := operator.NewService(queries)
	operatorHandler := operator.NewHandler(operatorService)
	operator.NewRouter(router, operatorHandler)

	customerService := customers.NewService(queries)
	customerHandler := customers.NewHandler(customerService)
	customers.NewRouter(router, customerHandler, JWTAccessSecret)

	eventsService := events.NewService(queries, broker)
	eventsHandler := events.NewHandler(eventsService)
	events.NewRouter(router, eventsHandler, JWTAccessSecret)

	authService := auth.NewService(queries, JWTAccessSecret, JWTRefreshSecret)
	authHandler := auth.NewHandler(authService)
	auth.NewRouter(router, authHandler)
}

func StartupServer(router *gin.Engine, port string) error {
	errChan := make(chan error, 1)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Server starting up... at port", "port", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server was unable to startup", "error", err)
			errChan <- fmt.Errorf("server was unable to startup successfully")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case errMesssage := <-errChan:
		return errMesssage
	case <-quit:
		slog.Info("shutdown signal received, shutting server gracefully...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown ungracefully", "error", err)
		return err
	}

	slog.Info("server exited properly")

	return nil
}
