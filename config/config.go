package config

import (
	"os"
)

const (
	dbConnectionString = "DB_CONNECTION_STRING"
	redisAddr          = "REDIS_ADDR"
	redisPassword      = "REDIS_PASSWORD"
	redisDB            = "REDIS_DB"
	imagePath          = "IMGPATH"
	accountKey         = "CLOUDINARY_ACCOUNT_KEY"
	secretKey          = "CLOUDINARY_SECRET_KEY"
	cloudName          = "CLOUDINARY_CLOUD_KEY"
)

// Config contains application configuration
type Config struct {
	DBConnectionString string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	ImagePath          string
	AccountKey         string
	SecretKey          string
	CloudName          string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}
	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}
	// default configuration
	config := &Config{
		DBConnectionString: getEnvOrDefault(dbConnectionString, "postgres://postgres:qweasd123@localhost:9811/evermos?sslmode=disable"),
	}

	return config, nil
}
