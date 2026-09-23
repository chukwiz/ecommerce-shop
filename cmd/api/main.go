package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chukwiz/go-shop/internal/config"
	"github.com/chukwiz/go-shop/internal/database"
	"github.com/chukwiz/go-shop/internal/interfaces"
	"github.com/chukwiz/go-shop/internal/logger"
	"github.com/chukwiz/go-shop/internal/providers"
	"github.com/chukwiz/go-shop/internal/server"
	"github.com/chukwiz/go-shop/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database instance")
	}

	defer func() {
		if err := mainDB.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	gin.SetMode(cfg.Server.GinMode)

	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	cartService := services.NewCartService(db)
	orderService := services.NewOrderService(db)

	var uploadProvider interfaces.UploadProvider

	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	uploadService := services.NewUploadService(uploadProvider)

	srv := server.New(cfg, db, &log, authService, userService, productService, uploadService, cartService, orderService)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting http server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown hyyp server")
		return
	}

	log.Info().Msg("shutting down database")

}
