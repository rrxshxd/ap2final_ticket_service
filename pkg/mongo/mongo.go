package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Database string `yaml:"database" env:"MONGO_DB" env-required:"true"`
	URI      string `yaml:"uri" env:"MONGO_URI" env-required:"true"`
	Username string `yaml:"username" env:"MONGO_USERNAME"`
	Password string `yaml:"password" env:"MONGO_PASSWORD"`
}

type DB struct {
	Connection *mongo.Database
	Client     *mongo.Client
}

func NewDB(ctx context.Context, cfg Config) (*DB, error) {
	clientOpts := cfg.clientOptions()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, connectionError(err)
	}

	db := &DB{
		Connection: client.Database(cfg.Database),
		Client:     client,
	}

	err = db.ping(ctx)
	if err != nil {
		return nil, connectionError(err)
	}

	return db, nil
}

func (db *DB) ping(ctx context.Context) error {
	err := db.Client.Ping(ctx, nil)
	if err != nil {
		return connectionError(err)
	}

	return nil
}
