package main

import (
	"github.com/chukwiz/go-shop/internal/config"
	"github.com/chukwiz/go-shop/internal/database"
	"github.com/chukwiz/go-shop/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	db, err := database.New(cfg.Database)
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

	log.Info().Msg("starting server")

}
