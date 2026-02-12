
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBHost                 string
	DBPort                 string
	DBName                 string
	DBAddress              string // Computed from DBHost:DBPort
	JWTExpirationInSeconds int64
	JWTSecret              string
	
	// Pasis Payment Gateway Configuration
	PasisAppKey    string
	PasisSecretKey string
}

var Envs = initConfig()

func initConfig() Config {
	// Load .env file (won't exist on Fly.io, which is fine)
	godotenv.Load()

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBHost:                 dbHost,
		DBPort:                 dbPort,
		DBAddress:              fmt.Sprintf("%s:%s", dbHost, dbPort), 
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