package config

import (
	"sync"
	"test_db_server/pkg/logging"
)

type Config struct {
	DatabaseConfig
	PublisherConfig
	SubscriberConfig
}

type DatabaseConfig struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	MaxAttempts int    `json:"max_attempts"`
}

type PublisherConfig struct {
	Name    string
	Channel string
	Cluster string
	Timeout string
}

type SubscriberConfig struct {
	Name    string
	Channel string
	Cluster string
	Timeout string
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Configure project")
		instance = &Config{
			DatabaseConfig{
				Username:    "postgres",
				Password:    "postgres",
				Host:        "localhost",
				Port:        "5432",
				Database:    "postgres",
				MaxAttempts: 5,
			},
			PublisherConfig{
				Name:    "publisher1",
				Channel: "test-channel",
				Cluster: "test-cluster",
				Timeout: "10s",
			},
			SubscriberConfig{
				Name:    "subscriber1",
				Channel: "test-channel",
				Cluster: "test-cluster",
				Timeout: "10s",
			},
		}
	})
	return instance
}
