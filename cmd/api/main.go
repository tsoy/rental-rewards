package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tsoy/rental-rewards/internal/data"
	"github.com/tsoy/rental-rewards/internal/events"
	"github.com/tsoy/rental-rewards/internal/service"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		minOpenConns int
		maxIdleTime  time.Duration
	}
	gcpProjectId string
}

type application struct {
	config   config
	logger   *slog.Logger
	models   data.Models
	services service.Services
}

func main() {
	var cfg config
	//port := os.Getenv("PORT")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("RR_DB_DSN"), "Postgres DSN")
	flag.IntVar(&cfg.db.minOpenConns, "db-min-open-conns", 1, "Postgres min open connections")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "Postgres max open connections")
	//flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "Postgres max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "Postgres max connection idle time")

	flag.StringVar(&cfg.gcpProjectId, "gcp_project_id", os.Getenv("GCP_PROJECT_ID"), "Google Cloud Platform Project ID")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("DB connection pool established")
	models := data.NewModels(db)

	pubsubClient, err := pubsub.NewClient(context.Background(), cfg.gcpProjectId)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pubsubClient.Close()

	app := &application{
		config:   cfg,
		logger:   logger,
		models:   models,
		services: service.NewServices(models, events.NewPubSubPublisher(pubsubClient)),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 25 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(cfg config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	conf.MinConns = int32(cfg.db.minOpenConns)
	conf.MaxConns = int32(cfg.db.maxOpenConns)
	conf.MaxConnIdleTime = cfg.db.maxIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
