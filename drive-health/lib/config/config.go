package config

import (
	"os"
	"strconv"
)

type DHConfig struct {
	CleanupServiceFrequency int `json:"cleanupServiceFrequency"`
	DiskFetchFrequency      int `json:"diskFetchFrequency"`
	MaxHistoryAge           int `json:"maxHistoryAge"`

	DatabaseFilePath string `json:"databaseFilePath"`

	Listen string `json:"listen"`

	IdentityUsername string `json:"identityUsername"`
	IdentityPassword string `json:"identityPassword"`

	IsDocker bool `json:isDocker`

	DebugMode bool `json:"debugMode"`
}

var config *DHConfig = nil

func GetConfiguration() *DHConfig {

	if config != nil {
		return config
	}

	config = &DHConfig{
		DiskFetchFrequency:      5,
		CleanupServiceFrequency: 3600,
		MaxHistoryAge:           2592000,
		DatabaseFilePath:        "./data.sqlite",
		IdentityUsername:        "admin",
		IdentityPassword:        "admin",

		IsDocker: false,

		Listen: ":8080",
	}

	if val, exists := os.LookupEnv("DISK_FETCH_FREQUENCY"); exists {
		if intValue, err := strconv.Atoi(val); err == nil {
			config.DiskFetchFrequency = intValue
		}
	}

	if val, exists := os.LookupEnv("CLEANUP_SERVICE_FREQUENCY"); exists {
		if intValue, err := strconv.Atoi(val); err == nil {
			config.CleanupServiceFrequency = intValue
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

	if val, exists := os.LookupEnv("IDENTITY_USERNAME"); exists {
		config.IdentityUsername = val
	}

	if val, exists := os.LookupEnv("IDENTITY_PASSWORD"); exists {
		config.IdentityPassword = val
	}

	if val, exists := os.LookupEnv("DEBUG_MODE"); exists {
		if isDebug, err := strconv.ParseBool(val); err == nil {
			config.DebugMode = isDebug
		}
	}

	if val, exists := os.LookupEnv("IS_DOCKER"); exists {
		if isDocker, err := strconv.ParseBool(val); err == nil {
			config.IsDocker = isDocker

			config.DatabaseFilePath = "/data/data.sqlite"
		}
	}

	return config
}
