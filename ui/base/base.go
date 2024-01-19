package base

import "log/slog"

var Logger *slog.Logger

func init() {
	Logger = slog.Default().With("context", "ui")
}
