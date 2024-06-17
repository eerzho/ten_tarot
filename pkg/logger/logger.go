package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type logger struct {
	l *slog.Logger
}

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
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}),
	)

	return &logger{l: l}
}

func (l *logger) Debug(message interface{}, args ...interface{}) {
	l.msg(slog.LevelDebug, message, args...)
}

func (l *logger) Info(message string, args ...interface{}) {
	l.msg(slog.LevelInfo, message, args...)
}

func (l *logger) Warn(message string, args ...interface{}) {
	l.msg(slog.LevelWarn, message, args...)
}

func (l *logger) Error(message interface{}, args ...interface{}) {
	l.msg(slog.LevelError, message, args...)
}

func (l *logger) Fatal(message interface{}, args ...interface{}) {
	l.msg(slog.LevelError, message, args...)
	os.Exit(1)
}

func (l *logger) msg(level slog.Level, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(level, msg.Error(), args...)
	case string:
		l.log(level, msg, args...)
	default:
		l.log(level, fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

func (l *logger) log(level slog.Level, message string, args ...interface{}) {
	if len(args) == 0 {
		l.l.Log(context.Background(), level, message)
	} else {
		l.l.Log(context.Background(), level, message, args...)
	}
}
