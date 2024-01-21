package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DHConfig struct {
	DiskFetchFrequency int `json:"diskFetchFrequency"`
	MaxHistoryAge      int `json:"maxHistoryAge"`
	DatabaseFilePath   string
	Listen             string
}

func GetConfiguration() DHConfig {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config := DHConfig{
		DiskFetchFrequency: 5,       // default value
		MaxHistoryAge:      2592000, // default value
		DatabaseFilePath:   "./data.sqlite",

		Listen: ":8080",
	}

	if val, exists := os.LookupEnv("DISK_FETCH_FREQUENCY"); exists {
		if intValue, err := strconv.Atoi(val); err == nil {
			config.DiskFetchFrequency = intValue
		}
	}

	if val, exists := os.LookupEnv("MAX_HISTORY_AGE"); exists {
		if intValue, err := strconv.Atoi(val); err == nil {
			config.MaxHistoryAge = intValue
		}
	}

	if val, exists := os.LookupEnv("LISTEN"); exists {
		config.Listen = val
	}

	if val, exists := os.LookupEnv("DATABASE_FILE_PATH"); exists {
		config.DatabaseFilePath = val
	}

	return config
}
