package hardware

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type HardDrive struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string
	Transport    string
	Size         string
	Model        string
	Serial       string
	Type         string
	HWID         string
	Temperatures []HardDriveTemperature `gorm:"foreignKey:HardDriveID"`
}

type HardDriveTemperature struct {
	gorm.Model
	HardDriveID uint
	TimeStamp   time.Time
	Temperature int
}

// A snapshot in time of the current state of the harddrives
type HardwareSnapshot struct {
	TimeStamp time.Time
	HDD       []*HardDrive
}

type Snapshots struct {
	List []*HardwareSnapshot
}

func (h *HardDrive) GetTemperature() int {

	possiblePaths := []string{
		"/sys/block/" + h.Name + "/device/hwmon/",
		"/sys/block/" + h.Name + "/device/",
		"/sys/block/" + h.Name + "/device/generic/device/",
	}

	for _, path := range possiblePaths {
		// Try HDD/SSD path
		temp, found := h.getTemperatureFromPath(path)
		if found {
			return temp
		}
	}

	fmt.Printf("[🛑] Failed to get temperature for %s\n", h.Name)
	return -1
}

func (h *HardDrive) getTemperatureFromPath(basePath string) (int, bool) {
	hwmonDirs, err := os.ReadDir(basePath)
	if err != nil {
		return 0, false
	}

	for _, dir := range hwmonDirs {
		if strings.HasPrefix(dir.Name(), "hwmon") {
			tempPath := filepath.Join(basePath, dir.Name(), "temp1_input")
			if _, err := os.Stat(tempPath); err == nil {
				tempBytes, err := os.ReadFile(tempPath)
				if err != nil {
					continue
				}

				tempStr := strings.TrimSpace(string(tempBytes))
				temperature, err := strconv.Atoi(tempStr)
				if err != nil {
					continue
				}

				// Convert millidegree Celsius to degree Celsius
				return temperature / 1000, true
			}
		}
	}

	return 0, false
}
