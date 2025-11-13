package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/flohansen/coffee-table/internal/controller"
	"github.com/flohansen/coffee-table/internal/repository"
	"github.com/flohansen/coffee-table/pkg/app"
	pkghttp "github.com/flohansen/coffee-table/pkg/http"
	"github.com/flohansen/coffee-table/sql/migrations"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/sync/errgroup"
)

type config struct {
	ListenAddr  string
	MetricsAddr string
	Database    string
}

func main() {
	var config config
	flag.StringVar(&config.ListenAddr, "listen-addr", ":8080", "The listen address for the API server")
	flag.StringVar(&config.MetricsAddr, "metrics-addr", ":9090", "The listen address for the metrics server")
	flag.StringVar(&config.Database, "database", "postgres://localhost:5432/postgres", "The connection string for the database")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	g, ctx := errgroup.WithContext(app.SignalContext())

	pool, err := pgxpool.New(ctx, config.Database)
	if err != nil {
		log.Error("failed to create pgx pool", "error", err)
		os.Exit(1)
	}
	if err := migrations.Up(pool); err != nil {
		log.Error("failed to run migration", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserPostgres(pool)
	userController := controller.NewUserController(userRepo)

	g.Go(func() error {
		log.Info("api server started", "addr", config.ListenAddr)
		if err := pkghttp.NewServer(
			pkghttp.WithListenAddr(config.ListenAddr),
			pkghttp.WithHandler("POST /users/register", http.HandlerFunc(userController.Register)),
		).Serve(ctx); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		log.Info("metrics server started", "addr", config.MetricsAddr)
		if err := pkghttp.NewServer(
			pkghttp.WithListenAddr(config.MetricsAddr),
			pkghttp.WithHandler("/metrics", promhttp.Handler()),
		).Serve(ctx); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("application error", "error", err)
		os.Exit(1)
	}
	log.Info("shutdown complete")
}
