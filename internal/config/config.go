package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sorawaslocked/ap2final_base/pkg/grpc"
	"github.com/sorawaslocked/ap2final_base/pkg/mongo"
	"os"
	"time"
)

type (
	Config struct {
		Env    string       `yaml:"env" env-required:"true"`
		Mongo  mongo.Config `yaml:"mongo" env-required:"true"`
		Redis  Redis        `yaml:"redis" env-required:"true"`
		Server Server       `yaml:"server" env-required:"true"`
	}

	Server struct {
		GRPC grpc.Config `yaml:"grpc" env-required:"true"`
	}

	Redis struct {
		Host     string        `yaml:"host" env-default:"localhost"`
		Port     string        `yaml:"port" env-default:"6379"`
		Password string        `yaml:"password" env-default:""`
		DB       int           `yaml:"db" env-default:"0"`
		TTL      time.Duration `yaml:"ttl" env-default:"24h"`
	}
)

func MustLoad() *Config {
	cfgPath := fetchConfigPath()

	if cfgPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("config file does not exist: " + cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("failed to load config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = "./config/local.yml"
	}

	return res
}
