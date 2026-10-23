package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/luism2302/RateLog/internal/client"
	"github.com/luism2302/RateLog/internal/data"
	"golang.org/x/crypto/bcrypt"
)

type application struct {
	logger    *slog.Logger
	cfg       config
	client    *client.Client
	models    data.Models
	dummyHash []byte
}
type config struct {
	addr int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  time.Duration
	}
}

func main() {
	var cfg config
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dummyHash, err := bcrypt.GenerateFromPassword([]byte(rand.Text()), 12)
	if err != nil {
		logger.Error("couldn't generate hash")
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		logger.Error("couldn't load .env")
		os.Exit(1)
	}

	flag.IntVar(&cfg.addr, "addr", 8000, "port for the app")
	flag.StringVar(&cfg.env, "env", "dev", "environment to use dev|prod")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("RATELOG_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	pool, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	apiKey := os.Getenv("RAWG_API_KEY")
	if apiKey == "" {
		logger.Error("couldn't find api RAWG_API_KEY in .env")
		os.Exit(1)
	}

	app := &application{
		logger:    logger,
		cfg:       cfg,
		client:    client.New(apiKey),
		models:    data.NewModels(pool),
		dummyHash: dummyHash,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.addr),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("server started", "addr", app.cfg.addr, "env", app.cfg.env)

	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(cfg.db.maxOpenConns)
	poolCfg.MaxConnIdleTime = cfg.db.maxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
