package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/luism2302/RateLog/internal/client"
)

type application struct {
	logger *slog.Logger
	cfg    config
	client *client.Client
}
type config struct {
	addr int
	env  string
}

func main() {
	var cfg config
	flag.IntVar(&cfg.addr, "addr", 8000, "port for the app")
	flag.StringVar(&cfg.env, "env", "dev", "environment to use dev|prod")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := godotenv.Load(); err != nil {
		logger.Error("couldn't load .env")
		os.Exit(1)
	}

	apiKey := os.Getenv("RAWG_API_KEY")
	if apiKey == "" {
		logger.Error("couldn't find api RAWG_API_KEY in .env")
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		cfg:    cfg,
		client: client.New(apiKey),
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
