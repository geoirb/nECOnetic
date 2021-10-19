package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/http2"

	"github.com/nECOnetic/data-service/internal/mongo"
	"github.com/nECOnetic/data-service/internal/response"
	"github.com/nECOnetic/data-service/internal/service"
	transport "github.com/nECOnetic/data-service/internal/service/http"
)

type configuration struct {
	HttpPort string `envconfig:"HTTP_PORT" default:"8081"`

	StorageHosts                  []string `envconfig:"STORAGE_HOSTS" default:"127.0.0.1:27017"`
	StorageUser                   string   `envconfig:"STORAGE_USER" default:"neconetic"`
	StoragePassword               string   `envconfig:"STORAGE_PASSWORD" default:"neconetic"`
	StorageDatabase               string   `envconfig:"STORAGE_DATABASE" default:"neconetic"`
	StorageOperationInTransaction int      `envconfig:"STORAGE_TRANSACTION" default:"4000"`

	StationCollectionName      string `envconfig:"STATION_COLLECTION_NAME" default:"station"`
	EcoDataCollectionName      string `envconfig:"ECO_DATA_COLLECTION_NAME" default:"eco-data"`
	ProfilerDataCollectionName string `envconfig:"PROFILER_DATA_COLLECTION_NAME" default:"profiler-data"`
}

const (
	prefixCfg   = ""
	serviceName = "data-service"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.WithPrefix(logger, "service", serviceName)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	if time.Since(time.Date(2020, time.August, 15, 0, 0, 0, 0, time.Now().Location())) < 0 {
		level.Error(logger).Log("msg", "trial version")
		return
	}

	var cfg configuration
	if err := envconfig.Process(prefixCfg, &cfg); err != nil {
		level.Error(logger).Log("msg", "configuration", "err", err)
		os.Exit(1)
	}

	level.Error(logger).Log("msg", "initialization")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f := mongo.StorageFabric{
		StationCollectionName:      cfg.StationCollectionName,
		EcoDataCollectionName:      cfg.EcoDataCollectionName,
		ProfilerDataCollectionName: cfg.ProfilerDataCollectionName,
	}

	storage, err := f.NewStorage(
		ctx,
		cfg.StorageHosts,
		cfg.StorageUser,
		cfg.StoragePassword,
		cfg.StorageDatabase,
		cfg.StorageOperationInTransaction,
	)
	if err != nil {
		level.Error(logger).Log("msg", "init mongo", "err", err)
		os.Exit(1)
	}

	svc := service.New(
		storage,

		logger,
	)

	router := mux.NewRouter()
	transport.Routing(router, svc, response.Build)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HttpPort),
		Handler: router,
	}

	if err = http2.ConfigureServer(&httpServer, &http2.Server{}); err != nil {
		level.Error(logger).Log("msg", "configurate http2 server", "err", err)
		os.Exit(1)
	}

	go func() {
		level.Info(logger).Log("msg", "http server turn on", "port", cfg.HttpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "http server turn on", "err", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal", "signal", <-c)

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		level.Info(logger).Log("msg", "http server shoutdown", "err", err)
	}
}