// Add these to your config/env.go file

package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	
	// Pasis Payment Gateway Configuration
	PasisAppKey    string
	PasisSecretKey string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              getEnv("DB_ADDRESS", "localhost:5432"),
		DBName:                 getEnv("DB_NAME", "apartmentdb"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		
		// Pasis Configuration
		PasisAppKey:    getEnv("PASIS_APP_KEY", ""),
		PasisSecretKey: getEnv("PASIS_SECRET_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}