package svc

import (
	"fmt"
	"time"

	"tea.chunkbyte.com/kato/drive-health/lib/config"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

// Delete all thermal entries that are older than X amount of seconds
func CleanupOldData() error {
	cfg := config.GetConfiguration()

	beforeDate := time.Now().Add(-1 * time.Duration(cfg.MaxHistoryAge) * time.Second)

	deleteResult := db.Where("time_stamp < ?", beforeDate).Delete(&hardware.HardDriveTemperature{})
	if deleteResult.Error != nil {
		fmt.Printf("[🛑] Error during cleanup: %s\n", deleteResult.Error)
		return db.Error
	}

	if deleteResult.RowsAffected > 0 {
		fmt.Printf("[🛑] Cleaned up %v entries before %s\n", deleteResult.RowsAffected, beforeDate)
	}

	return nil
}

func RunCleanupService() {
	fmt.Println("[🦝] Initializing Log Cleanup Service...")

	tickTime := time.Duration(config.GetConfiguration().CleanupServiceFrequency) * time.Second

	// Snapshot taking routine
	go func() {
		for {
			time.Sleep(tickTime)
			err := CleanupOldData()
			if err != nil {
				fmt.Printf("🛑 Cleanup process failed: %s\n", err)
			}
		}
	}()
}
