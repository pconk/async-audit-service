package config

import (
	"os"
)

type Config struct {
	AppPort   string
	MongoURI  string
	DbName    string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		AppPort:   getEnv("APP_PORT", "50051"), // Port default gRPC biasanya 50051
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DbName:    getEnv("DB_NAME", "audit_db"),
		JWTSecret: getEnv("JWT_SECRET", "rahasia-super-aman"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
