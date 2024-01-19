package hardware

import (
	"fmt"

	"github.com/anatol/smart.go"
)

type HardDrive struct {
	Name        string
	Transport   string
	Size        string
	Model       string
	Serial      string
	Type        string
	Temperature int
}

// Fetch the temperature of the device, optinally update the reference object
func (h *HardDrive) GetTemperature(updateTemp bool) int {
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

	// Optionally update the reference object's temperature
	if updateTemp {
		h.Temperature = temperature
	}

	// Return the found value
	return temperature
}
