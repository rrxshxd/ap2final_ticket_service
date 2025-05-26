package app

import (
	"ap2final_ticket_service/internal/adapter/cache"
	grpcserver "ap2final_ticket_service/internal/adapter/grpc"
	mongorepo "ap2final_ticket_service/internal/adapter/mongo"
	"ap2final_ticket_service/internal/config"
	"ap2final_ticket_service/internal/usecase"
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
	mongocfg "github.com/sorawaslocked/ap2final_base/pkg/mongo"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "ticket service"

type App struct {
	grpcServer *grpcserver.Server
	cache      *cache.RedisCache
	log        *slog.Logger
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log *slog.Logger,
) (*App, error) {
	const op = "App.New"

	newLog := log.With(slog.String("op", op))
	newLog.Info("starting service", slog.String("service", serviceName))

	newLog.Info("connecting to redis cache", slog.String("host", cfg.Redis.Host))
	redisCache := cache.NewRedisCache(cfg.Redis)

	if err := redisCache.Ping(ctx); err != nil {
		newLog.Error("error connecting to redis cache", logger.Err(err))
		return nil, err
	}

	newLog.Info("connecting to mongo database", slog.String("uri", cfg.Mongo.URI))

	db, err := mongocfg.NewDB(ctx, cfg.Mongo)
	if err != nil {
		newLog.Error("error connecting to mongo database", logger.Err(err))
		return nil, err
	}

	ticketRepo := mongorepo.NewTicket(db.Connection)

	ticketUseCase := usecase.NewTicketUseCase(ticketRepo, redisCache, log)

	grpcServer := grpcserver.New(cfg.Server.GRPC, log, ticketUseCase)

	return &App{
		grpcServer: grpcServer,
		cache:      redisCache,
		log:        log,
	}, nil
}

func (a *App) stop() {
	a.grpcServer.Stop()
	if err := a.cache.Close(); err != nil {
		a.log.Error("error closing redis cache", logger.Err(err))
	}
}

func (a *App) Run() {
	a.grpcServer.MustRun()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-shutdownCh

	a.log.Info("received system shutdown signal", slog.Any("signal", s.String()))
	a.log.Info("stopping the application")
	a.stop()
	a.log.Info("graceful shutdown complete")
}
