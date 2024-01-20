package hardware

import (
	"fmt"
	"time"

	"github.com/anatol/smart.go"
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

// Fetch the temperature of the device, optinally update the reference object
func (h *HardDrive) GetTemperature() int {
	// Fetch the device by name
	disk, err := smart.Open("/dev/" + h.Name)
	if err != nil {
		fmt.Printf("Failed to open device %s: %s\n", h.Name, err)
		return -1
	}
	defer disk.Close()

	// Fetch SMART data
	smartInfo, err := disk.ReadGenericAttributes()
	if err != nil {
		fmt.Printf("Failed to get SMART data for %s: %s\n", h.Name, err)
		return -1
	}

	// Parse the temperature
	temperature := int(smartInfo.Temperature)

	// Return the found value
	return temperature
}
