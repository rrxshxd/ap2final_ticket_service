package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (cfg *Config) clientOptions() *options.ClientOptions {
	opts := options.Client().ApplyURI(cfg.connectionUrl())

	return opts
}

func (cfg *Config) connectionUrl() string {
	var url string

	if cfg.Username == "" || cfg.Password == "" {
		url = fmt.Sprintf("mongodb://%s", cfg.URI)
	} else {
		url = fmt.Sprintf("mongodb://%s:%s@%s", cfg.URI, cfg.Username, cfg.Password)
	}

	return url
}
