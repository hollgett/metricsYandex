package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hollgett/metricsYandex.git/internal/server/api"
	"github.com/hollgett/metricsYandex.git/internal/server/config"
	"github.com/hollgett/metricsYandex.git/internal/server/logger"
	"github.com/hollgett/metricsYandex.git/internal/server/repository"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/file"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/memory"
	"github.com/hollgett/metricsYandex.git/internal/server/repository/postgres"
	"github.com/hollgett/metricsYandex.git/internal/server/server"
	"github.com/hollgett/metricsYandex.git/internal/server/services"
)

func main() {

	//init logger
	logger, err := logger.New()
	if err != nil {
		log.Fatalf("can't init zap logger: %v", err)
	}
	defer logger.Flush()

	//initialisation main ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//shut down func
	go handlerShutDown(cancel, logger)
	//init config
	cfg, err := config.New()
	if err != nil {
		logger.LogErr("config init", err)
		cancel()
	}
	logger.LogAny("server start", "config", cfg)
	// db, _ := database.Connect(cfg.DataBaseDSN)
	repo, err := initStorage(ctx, cfg, logger)
	if err != nil {
		logger.LogErr("repository init", err)
		cancel()
	}
	handlers := services.New(repo)
	//presentation layer
	api := api.New(handlers, logger)
	server := server.New(api, logger, cfg.Addr)
	//start server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.LogErr("server", err)
			cancel()
		}
	}()
	//close server
	defer func() {
		if err := server.Close(); err != nil {
			logger.LogErr("server close", err)
		}
		logger.LogMess("close server")
	}()
	defer func() {
		if err := repo.Close(); err != nil {
			logger.LogErr("storage close", err)
		}
		logger.LogMess("storage close")
	}()
	<-ctx.Done()
}

func handlerShutDown(cancel context.CancelFunc, log logger.Logger) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	sig := <-signalChan
	log.LogAny("close app", "signal", sig)
	cancel()
}

func initStorage(ctx context.Context, cfg *config.Config, log logger.Logger) (repository.Repository, error) {
	switch {
	case cfg.DataBaseDSN != "":
		postgresDB, err := postgres.New(ctx, cfg.DataBaseDSN)
		if err != nil {
			return nil, err
		}
		log.LogMess("postgresql mode")
		return postgresDB, nil
	case cfg.PathFileStorage != "":
		fStorage, err := file.New(ctx, log, cfg.PathFileStorage, cfg.StorageInterval, cfg.Restore)
		if err != nil {
			return nil, err
		}
		log.LogMess("file mode")
		return fStorage, nil
	default:
		memStorage := memory.New()
		log.LogMess("memory mode")
		return memStorage, nil
	}
}
