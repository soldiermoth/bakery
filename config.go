package main

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Config holds all the configuration for this service
type Config struct {
	Listen   string `default:":8080"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	Client   HTTPClient
}

// HTTPClient will issue requests to the manifest
type HTTPClient struct {
	Timeout time.Duration `envconfig:"CLIENT_TIMEOUT" default:"2s"`
}

func (h HTTPClient) new() *http.Client {
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	client := &http.Client{
		Timeout: h.Timeout,
	}
	return client
}

// LoadConfig loads the configuration with environment variables injected
func LoadConfig() Config {
	return Config{}
}

// GetLogger generates a logger
func (c Config) GetLogger() *logrus.Logger {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}

	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = level
	return logger
}
