package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type (
	logger struct {
		*slog.Logger
	}

	Attr slog.Attr
)

var l *logger

func SetUpLogger(level, handlerType string) {
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

	var handler slog.Handler
	switch strings.ToLower(handlerType) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	}

	sl := slog.New(handler)

	l = &logger{sl}
}

func SetUpDefaultLogger() {
	SetUpLogger("info", "text")
}

func Any(key string, value interface{}) Attr {
	return Attr(slog.Any(key, value))
}

func Debug(message string, args ...Attr) {
	log(slog.LevelDebug, message, args...)
}

func Info(message string, args ...Attr) {
	log(slog.LevelInfo, message, args...)
}

func Warn(message string, args ...Attr) {
	log(slog.LevelWarn, message, args...)
}

func Error(message string, args ...Attr) {
	log(slog.LevelError, message, args...)
}

func OPError(op string, err error, args ...Attr) {
	Error(fmt.Sprintf("%s - %s", op, err), args...)
}

func Fatal(message string, args ...Attr) {
	log(slog.LevelError, message, args...)
	os.Exit(1)
}

func log(level slog.Level, message string, args ...Attr) {
	if l == nil {
		SetUpDefaultLogger()
	}

	slogArgs := make([]slog.Attr, len(args))
	for key, val := range args {
		slogArgs[key] = slog.Attr(val)
	}

	if len(args) == 0 {
		l.Log(context.Background(), level, message)
	} else {
		l.LogAttrs(context.Background(), level, message, slogArgs...)
	}
}
