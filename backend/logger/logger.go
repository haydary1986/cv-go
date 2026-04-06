package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	var handler slog.Handler
	if os.Getenv("GIN_MODE") == "release" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	Log = slog.New(handler)
	slog.SetDefault(Log)
}
