package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl            string
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
}

func LoadConfig() *Config {
	godotenv.Load()

	DBUrl, err := ExtractEnvKey("DATABASE_URL", "")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			log.Fatal("dburl env value not found")
		}
	}

	Port, err := ExtractEnvKey("PORT", "6060")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			log.Fatal("port env value not found")
		}
	}

	JWTAccessSecret, err := ExtractEnvKey("JWT_ACCESS_SECRET", "")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			log.Fatal("jwt access secret env value not found")
		}
	}

	JWTRefreshSecret, err := ExtractEnvKey("JWT_REFRESH_SECRET", "")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			log.Fatal("jwt refresh secret env value not found")
		}
	}

	return &Config{
		DBUrl:            DBUrl,
		Port:             Port,
		JWTAccessSecret:  JWTAccessSecret,
		JWTRefreshSecret: JWTRefreshSecret,
	}

}

func ExtractEnvKey(key, fallback string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok && fallback != "" {
		return fallback, nil
	} else if !ok && fallback == "" {
		return "", ErrEnvironmentVariableNotFound
	}

	return val, nil
}
