package infra

import (
	"context"
	"log/slog"
	"time"

	"chapter7/activation"
)

type Scheduler struct {
	Query  *QueryDB
	App    activation.App
	Logger *slog.Logger
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		if err := s.run(); err != nil {
			s.Logger.Error(err.Error(), "type", "scheduler")
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) run() error {
	before := time.Now().Add(24 * time.Hour).UTC().Unix()

	s.Logger.Info("running", "type", "scheduler", "before", before)

	accounts, err := s.Query.AccountIDsAboutExpire(before)
	if err != nil {
		return err
	}

	for _, a := range accounts {
		_, _ = s.App.ExtendTrial(a)
	}
	return nil
}
