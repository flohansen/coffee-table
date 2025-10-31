package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/flohansen/coffee-table/pkg/app"
	"github.com/flohansen/coffee-table/pkg/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/sync/errgroup"
)

type config struct {
	ListenAddr  string
	MetricsAddr string
}

func main() {
	var config config
	flag.StringVar(&config.ListenAddr, "listen-addr", ":8080", "The listen address for the API server")
	flag.StringVar(&config.MetricsAddr, "metrics-addr", ":9090", "The listen address for the metrics server")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	g, ctx := errgroup.WithContext(app.SignalContext())

	g.Go(func() error {
		log.Info("api server started", "addr", config.ListenAddr)
		if err := http.NewServer(
			http.WithListenAddr(config.ListenAddr),
		).Serve(ctx); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		log.Info("metrics server started", "addr", config.MetricsAddr)
		if err := http.NewServer(
			http.WithListenAddr(config.MetricsAddr),
			http.WithHandler("/metrics", promhttp.Handler()),
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
