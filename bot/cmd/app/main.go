package main

import (
	"bot/config"
	"bot/internal/app"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	mng, err := connectMongo(cfg.Mongo.URL, cfg.Mongo.DB)
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	defer disconnectMongo(mng)

	lg := setUpLogger(cfg.Log.Level, cfg.Log.Type)

	app, err := app.New(cfg, mng, lg)
	if err != nil {
		log.Fatalf("app: %v", err)
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		app.Run()
	}()
	log.Print("app: started")
	<-stopChan
	log.Print("app: shutting down")
	app.Shutdown()
	log.Print("app: stopped")
}

func setUpLogger(lvl, format string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(lvl) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	default:
		level = slog.LevelError
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler)
}

func connectMongo(url, name string) (*mongo.Database, error) {
	clientOPT := options.Client().ApplyURI(url).SetMaxPoolSize(1)

	var err error
	var client *mongo.Client

	attempts := 10
	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, clientOPT)
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				break
			}
		}
		log.Printf("mongo: attempts left - %d", attempts)
		time.Sleep(time.Second)
		attempts--
	}

	if err != nil {
		return nil, err
	}

	return client.Database(name), nil
}

func disconnectMongo(mng *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	err := mng.Client().Disconnect(ctx)
	if err != nil {
		log.Printf("mongo: %v", err)
	}
}
