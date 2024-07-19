package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type (
	Logger interface {
		Debug(message string, args ...Attr)
		Info(message string, args ...Attr)
		Warn(message string, args ...Attr)
		Error(message string, args ...Attr)
		Fatal(message string, args ...Attr)
	}
	logger struct {
		l *slog.Logger
	}

	Attr slog.Attr
)

func New(level string) Logger {
	var lvl slog.Level

	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	l := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}),
	)

	return &logger{l: l}
}

func Any(key string, value interface{}) Attr {
	return Attr(slog.Any(key, value))
}

func (l *logger) Debug(message string, args ...Attr) {
	l.log(slog.LevelDebug, message, args...)
}

func (l *logger) Info(message string, args ...Attr) {
	l.log(slog.LevelInfo, message, args...)
}

func (l *logger) Warn(message string, args ...Attr) {
	l.log(slog.LevelWarn, message, args...)
}

func (l *logger) Error(message string, args ...Attr) {
	l.log(slog.LevelError, message, args...)
}

func (l *logger) Fatal(message string, args ...Attr) {
	l.log(slog.LevelError, message, args...)
	os.Exit(1)
}

func (l *logger) log(level slog.Level, message string, args ...Attr) {
	slogArgs := make([]slog.Attr, len(args))
	for key, val := range args {
		slogArgs[key] = slog.Attr(val)
	}

	if len(args) == 0 {
		l.l.Log(context.Background(), level, message)
	} else {
		l.l.LogAttrs(context.Background(), level, message, slogArgs...)
	}
}
