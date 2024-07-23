package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Mongo struct {
	maxPoolSize  uint64
	connAttempts int
	connTimeout  time.Duration

	Client *mongo.Client
	DB     *mongo.Database
}

func New(url, dbName string, opts ...Option) (*Mongo, error) {
	const op = "mongo"

	mg := &Mongo{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(mg)
	}

	clientOptions := options.Client().ApplyURI(url).SetMaxPoolSize(mg.maxPoolSize)

	var err error
	for mg.connAttempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), mg.connTimeout)
		defer cancel()

		mg.Client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			err = mg.Client.Ping(ctx, nil)
			if err == nil {
				break
			}
		}
		log.Printf("%s: attempts left - %d", op, mg.connAttempts)
		time.Sleep(mg.connTimeout)
		mg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mg.DB = mg.Client.Database(dbName)
	return mg, nil
}

func (m *Mongo) Close() {
	if m.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = m.Client.Disconnect(ctx)
	}
}
