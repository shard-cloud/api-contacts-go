package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	Port        string
	Environment string
}

func Load() *Config {
	return &Config{
		DatabaseURL: normalizeDatabaseURL(os.Getenv("DATABASE")),
		Port:        os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
	}
}

// normalizeDatabaseURL fixes common SSL parameter issues in PostgreSQL connection strings
func normalizeDatabaseURL(url string) string {
	if url == "" {
		return url
	}

	// Replace ssl=true with sslmode=require (PostgreSQL standard)
	url = strings.ReplaceAll(url, "ssl=true", "sslmode=require")
	url = strings.ReplaceAll(url, "ssl=false", "sslmode=disable")
	
	// If no sslmode is specified and it's a remote database, add sslmode=require
	if !strings.Contains(url, "sslmode=") && 
	   !strings.Contains(url, "localhost") && 
	   !strings.Contains(url, "127.0.0.1") {
		if strings.Contains(url, "?") {
			url += "&sslmode=require"
		} else {
			url += "?sslmode=require"
		}
	}

	return url
}
