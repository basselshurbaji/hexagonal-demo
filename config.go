package main

import "os"

type DBConfig struct {
	HOST     string
	PORT     string
	DATABASE string
	USERNAME string
	PASSWORD string
}

type HTTPConfig struct {
	PORT string
}

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
}

func ConfigFromEnv() Config {
	return Config{
		DB: DBConfig{
			HOST:     getFromEnvOrDefault("DB_HOST", "mysql.hexagonal-demo.orb.local"),
			PORT:     getFromEnvOrDefault("DB_PORT", "3306"),
			DATABASE: getFromEnvOrDefault("DB_DATABASE", "hexagonal"),
			USERNAME: getFromEnvOrDefault("DB_USERNAME", "app"),
			PASSWORD: getFromEnvOrDefault("DB_PASSWORD", "app"),
		},
		HTTP: HTTPConfig{
			PORT: getFromEnvOrDefault("HTTP_PORT", "8080"),
		},
	}
}

func getFromEnvOrDefault(env string, def string) string {
	value := os.Getenv(env)
	if value == "" {
		return def
	}
	return value
}
