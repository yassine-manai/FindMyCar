package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigFile struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	App struct {
		JSecret string
	}
}

// Load the .env file and initialize the config
func (c *ConfigFile) Load() error {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Server configuration
	c.Server.Host = c.getEnv("SERVER_HOST", "127.0.0.1")
	c.Server.Port, err = strconv.Atoi(c.getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid server port: %v", err)
	}

	// Database configuration
	c.Database.Host = c.getEnv("DB_HOST", "127.0.0.1")
	c.Database.Port, err = strconv.Atoi(c.getEnv("DB_PORT", "5432"))
	if err != nil {
		return fmt.Errorf("invalid database port: %v", err)
	}
	c.Database.User = c.getEnv("DB_USER", "root")
	c.Database.Password = c.getEnv("DB_PASSWORD", "password")
	c.Database.Name = c.getEnv("DB_NAME", "fyc")

	c.App.JSecret = c.getEnv("JWT_Secret", "0")

	return nil
}

// Helper function to get the environment variable with a default value
func (c *ConfigFile) getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
