package main

import (
	"context"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/luism2302/RateLog/internal/client"
	"github.com/luism2302/RateLog/internal/data"
	"github.com/luism2302/RateLog/migrations"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

type application struct {
	logger    *slog.Logger
	cfg       config
	client    *client.Client
	models    data.Models
	dummyHash []byte
	wg        sync.WaitGroup
}
type config struct {
	addr string
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  time.Duration
	}
	migrate bool
}

func main() {
	var cfg config
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			break
		default:
			logger.Error("couldn't load .env")
			os.Exit(1)
		}
	}

	flag.StringVar(&cfg.addr, "addr", "8000", "port for the app")
	flag.StringVar(&cfg.env, "env", "prod", "environment to use dev|prod")
	migrate := flag.Bool("migrate", false, "determines if migrations will run")
	health := flag.Bool("healthcheck", false, "runs a healthcheck")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	if *migrate {
		migrationsDSN := os.Getenv("RATELOG_MIGRATE_DSN")
		if migrationsDSN == "" {
			logger.Error("couln' find RATELOG_MIGRATE_DSN")
			os.Exit(1)
		}
		results, err := runMigrations(migrationsDSN)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		for _, result := range results {
			logger.Info("migration result", "result", result.String())
		}
		os.Exit(0)
	}

	if *health {
		tempClient := http.Client{
			Timeout: 3 * time.Second,
		}
		url := fmt.Sprintf("http://127.0.0.1:%s/healthz", cfg.addr)
		response, err := tempClient.Get(url)
		if err != nil || response.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if cfg.env != "prod" && cfg.env != "dev" {
		logger.Error("unknown env value, only prod or dev are valid")
		os.Exit(1)
	}

	apiKey := os.Getenv("RAWG_API_KEY")
	if apiKey == "" {
		logger.Error("couldn't find RAWG_API_KEY")
		os.Exit(1)
	}

	dsn := os.Getenv("RATELOG_DB_DSN")
	if dsn == "" {
		logger.Error("couldn't find RATELOG_DB_DSN")
		os.Exit(1)
	}
	cfg.db.dsn = dsn

	pool, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dummyHash, err := bcrypt.GenerateFromPassword([]byte(rand.Text()), 12)
	if err != nil {
		logger.Error("couldn't generate hash")
		pool.Close()
		os.Exit(1)
	}

	app := &application{
		logger:    logger,
		cfg:       cfg,
		client:    client.New(apiKey),
		models:    data.NewModels(pool),
		dummyHash: dummyHash,
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		pool.Close()
		os.Exit(1)
	}
	pool.Close()
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
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func runMigrations(dsn string) ([]*goose.MigrationResult, error) {
	conn, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDB(*conn)
	defer db.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.MigrationFS)
	if err != nil {
		return nil, err
	}
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	if err := provider.Ping(ctxPing); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	results, err := provider.Up(ctx)

	return results, err
}
