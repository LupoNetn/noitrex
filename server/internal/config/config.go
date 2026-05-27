package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl            string
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	DBUrl, err := ExtractEnvKey("DATABASE_URL", "")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			return nil, errors.New("dburl env value not found")
		}
	}

	Port, err := ExtractEnvKey("PORT", "6060")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			return nil, errors.New("port env value not found")
		}
	}

	JWTAccessSecret, err := ExtractEnvKey("JWT_ACCESS_SECRET", "")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			return nil, errors.New("jwt access secret env value not found")
		}
	}

	JWTRefreshSecret, err := ExtractEnvKey("JWT_REFRESH_SECRET", "")
	if err != nil {
		if errors.Is(err, ErrEnvironmentVariableNotFound) {
			return nil, errors.New("jwt refresh secret env value not found")
		}
	}

	return &Config{
		DBUrl:            DBUrl,
		Port:             Port,
		JWTAccessSecret:  JWTAccessSecret,
		JWTRefreshSecret: JWTRefreshSecret,
	}, nil

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
