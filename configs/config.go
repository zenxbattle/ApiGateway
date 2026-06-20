package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	// Environment
	Environment            string
	JWTSecretKey           string
	LokiURL string

	// Microservices
	APIGATEWAYPORT     string
	NatsURL           string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	FrontendURL string

	UserGRPCURL      string
	CompilerGRPCURL  string
	ProblemGRPCURL   string

}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() Config {
	// Load .env file if present
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables or defaults ", err)
	}

	return Config{
		Environment:        getEnv("ENVIRONMENT", "development"),
		JWTSecretKey:       getEnv("JWTSECRETKEY", "secretLeetcode"),
		GoogleClientID:     getEnv("GOOGLECLIENTID", ""),
		GoogleClientSecret: getEnv("GOOGLECLIENTSECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLEREDIRECTURL", ""),

		LokiURL: getEnv("LOKIURL", "http://localhost:3100"),
		FrontendURL:            getEnv("FRONTENDURL", "http://localhost:8080"),

		APIGATEWAYPORT:     getEnv("APIGATEWAYPORT", "7000"),
		NatsURL:           getEnv("NATSURL", "nats://localhost:4222"),
		UserGRPCURL:        getEnv("USERGRPCURL", "localhost:50051"),
		ProblemGRPCURL:     getEnv("PROBLEMGRPCURL", "localhost:50055"),

	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}
