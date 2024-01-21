package svc

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tea.chunkbyte.com/kato/drive-health/lib/config"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

var db *gorm.DB

func InitDB() {
	var err error
	dbPath := config.GetConfiguration().DatabaseFilePath
	if dbPath == "" {
		dbPath = "./data.sqlite"
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&hardware.HardDrive{}, &hardware.HardDriveTemperature{})
}

func GetDatabaseRef() *gorm.DB {
	return db
}

func LogDriveTemps() error {
	drives, err := hardware.GetSystemHardDrives(db, nil, nil)
	if err != nil {
		return err
	}

	for _, hdd := range drives {
		temp := hdd.GetTemperature()
		db.Create(&hardware.HardDriveTemperature{
			HardDriveID: hdd.ID,
			TimeStamp:   time.Now(),
			Temperature: temp,
		})
	}

	return nil
}

func RunLoggerService() {
	fmt.Println("Initializing Temperature Logging Service...")

	tickTime := time.Duration(config.GetConfiguration().DiskFetchFrequency) * time.Second

	// Snapshot taking routine
	go func() {
		for {
			time.Sleep(tickTime)
			err := LogDriveTemps()
			if err != nil {
				fmt.Printf("🛑 Temperature logging failed: %s\n", err)
			}
		}
	}()
}
