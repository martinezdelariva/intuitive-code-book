package infra

import (
	"log/slog"
	"time"
)

type MonitoringSLog struct {
	Logger *slog.Logger
}

func (m *MonitoringSLog) Monitor(name string, accountID int) func(errp *error) {
	start := time.Now()
	return func(errp *error) {
		duration := time.Since(start)
		logData := []any{"name", name, "type", "write_operation", "account_id", accountID, "duration", duration}
		if err := *errp; err != nil {
			m.Logger.Error(err.Error(), logData...)
		} else {
			m.Logger.Info("executed", logData...)
		}
	}
}
