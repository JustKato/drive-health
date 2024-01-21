package hardware

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/JustKato/drive-health/lib/config"
	"gorm.io/gorm"
)

func GetSystemHardDrives(db *gorm.DB, olderThan *time.Time, newerThan *time.Time) ([]*HardDrive, error) {
	var systemHardDrives []*HardDrive

	// List all block devices
	devices, err := os.ReadDir("/sys/block/")
	if err != nil {
		return nil, fmt.Errorf("failed to list block devices: %w", err)
	}

	for _, device := range devices {
		deviceName := device.Name()

		// Skip non-physical devices (like loop and ram devices)
		// TODO: Read more about this, there might be some other devices we should or should not skip
		if strings.HasPrefix(deviceName, "loop") || strings.HasPrefix(deviceName, "ram") {
			continue
		}

		// Read device details
		model, _ := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/model", deviceName))
		serial, _ := os.ReadFile(fmt.Sprintf("/sys/block/%s/device/serial", deviceName))
		sizeBytes, _ := os.ReadFile(fmt.Sprintf("/sys/block/%s/size", deviceName))

		size := convertSizeToString(sizeBytes)
		transport := getTransportType(deviceName)

		// TODO: Maybe find a better way?
		if size == "0 Bytes" {
			// This looks like an invalid device, skip it.
			if config.GetConfiguration().DebugMode {
				fmt.Printf("[🟨] Igoring device:[/dev/%s], reported size of 0\n", deviceName)
			}
			continue
		}

		hwid, err := getHardwareID(deviceName)
		if err != nil {
			if config.GetConfiguration().DebugMode {
				fmt.Printf("[🟨] No unique identifier found for device:[/dev/%s] unique identifier\n", deviceName)
			}
			continue
		}

		hd := &HardDrive{
			Name:      deviceName,
			Transport: transport,
			Model:     strings.TrimSpace(string(model)),
			Serial:    strings.TrimSpace(string(serial)),
			Size:      size,
			Type:      getDriveType(deviceName),
			HWID:      hwid,
		}

		systemHardDrives = append(systemHardDrives, hd)
	}

	var updatedHardDrives []*HardDrive

	for _, sysHDD := range systemHardDrives {
		var existingHD HardDrive
		q := db.Where("hw_id = ?", sysHDD.HWID)

		if newerThan != nil && olderThan != nil {
			q = q.Preload("Temperatures", "time_stamp < ? AND time_stamp > ?", newerThan, olderThan)
		}

		result := q.First(&existingHD)

		if result.Error == gorm.ErrRecordNotFound {
			// Hard drive not found, create new
			db.Create(&sysHDD)
			updatedHardDrives = append(updatedHardDrives, sysHDD)
		} else {
			// Hard drive found, update existing
			existingHD.Name = sysHDD.Name
			existingHD.Transport = sysHDD.Transport
			existingHD.Size = sysHDD.Size
			existingHD.Model = sysHDD.Model
			existingHD.Type = sysHDD.Type
			db.Save(&existingHD)
			updatedHardDrives = append(updatedHardDrives, &existingHD)
		}
	}

	return updatedHardDrives, nil
}

func getTransportType(deviceName string) string {
	transportLink, err := filepath.EvalSymlinks(fmt.Sprintf("/sys/block/%s/device", deviceName))
	if err != nil {
		return "Unknown"
	}

	if strings.Contains(transportLink, "/usb/") {
		return "USB"
	} else if strings.Contains(transportLink, "/ata") {
		return "SATA"
	} else if strings.Contains(transportLink, "/nvme/") {
		return "NVMe"
	}

	return "Other"
}

func convertSizeToString(sizeBytes []byte) string {
	// Convert the size from a byte slice to a string, then to an integer
	sizeStr := strings.TrimSpace(string(sizeBytes))
	sizeSectors, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return "Unknown"
	}

	// Convert from 512-byte sectors to bytes
	sizeInBytes := sizeSectors * 512

	// Define size units
	const (
		_          = iota // ignore first value by assigning to blank identifier
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
	)

	var size float64 = float64(sizeInBytes)
	var unit string

	// Determine the unit to use
	switch {
	case size >= TB:
		size /= TB
		unit = "TB"
	case size >= GB:
		size /= GB
		unit = "GB"
	case size >= MB:
		size /= MB
		unit = "MB"
	case size >= KB:
		size /= KB
		unit = "KB"
	default:
		unit = "Bytes"
	}

	// Return the formatted size
	return fmt.Sprintf("%.2f %s", size, unit)
}

// Look throug /sys/block/device/ and try and find the unique identifier of the device.
func getHardwareID(deviceName string) (string, error) {
	// Define potential ID file paths
	idFilePaths := []string{
		"/sys/block/" + deviceName + "/device/wwid",
		"/sys/block/" + deviceName + "/device/wwn",
		"/sys/block/" + deviceName + "/device/serial",
	}

	// Try to read each file and return the first successful read
	for _, path := range idFilePaths {
		if idBytes, err := os.ReadFile(path); err == nil {
			return strings.TrimSpace(string(idBytes)), nil
		}
	}

	// Return an empty string if no ID is found
	return "", fmt.Errorf("could not find unique identifier for %s", deviceName)
}

// Figure out what kind of device this is by reading if it's rotational or not
func getDriveType(deviceName string) string {
	// Check if the drive is rotational (HDD)
	if isRotational, _ := os.ReadFile(fmt.Sprintf("/sys/block/%s/queue/rotational", deviceName)); string(isRotational) == "1\n" {
		return "HDD"
	}

	// Check if the drive is NVMe
	if strings.HasPrefix(deviceName, "nvme") {
		return "NVMe"
	}

	// Default to SSD for non-rotational and non-NVMe drives
	return "SSD"
}
