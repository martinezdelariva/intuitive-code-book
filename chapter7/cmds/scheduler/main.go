package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"chapter7/activation"
	"chapter7/activation/infra"
)

func main() {
	path := flag.String("path", "./data.sqlite", "path to the database file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	query, err := infra.NewQuerySQLite(*path)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}

	defer func() {
		if err := query.Close(); err != nil {
			logger.Error(err.Error())
			os.Exit(-1)
		}
	}()

	store, err := infra.NewStoreSQLite(*path)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}

	defer func() {
		if err := store.Close(); err != nil {
			logger.Error(err.Error())
			os.Exit(-1)
		}
	}()

	app := activation.App{
		Store:      store,
		Clock:      &infra.ClockSystem{},
		Monitoring: &infra.MonitoringSLog{Logger: logger},
	}

	scheduler := infra.Scheduler{
		Query:  query,
		App:    app,
		Logger: logger,
	}
	scheduler.Run(ctx)
}
