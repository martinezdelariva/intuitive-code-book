package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"chapter7/activation/infra"
)

func main() {
	path := flag.String("path", "./data.sqlite", "path to the database file")
	flag.Parse()

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

	rate, err := query.TrialExtendedActivationRate()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(-1)
	}

	fmt.Printf("The rate of trial extended activation is %2.f%%\n", rate)
}
