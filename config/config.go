package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseUrl   string
	ApiKey    string
	ModelName string
}

func Load() (*Config, error) {
    err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("не найден файл .env %s", err)
	}

    baseUlr, err := getEnvVar("BASE_URL")
    if err != nil {
        return nil, err
    }

    apiKey, err := getEnvVar("API_KEY")
    if err != nil {
        return nil, err
    }

    modelName, err := getEnvVar("MODEL_NAME")
    if err != nil {
        return nil, err
    }
    
	return &Config{
		BaseUrl:   baseUlr,
		ApiKey:    apiKey,
		ModelName: modelName,
	}, nil
}

func getEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
        return "", fmt.Errorf("переменная %s не установлена", key)
	}

	return value, nil
}
