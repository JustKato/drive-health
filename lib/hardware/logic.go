package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"gorm.io/gorm"
)

func GetSystemHardDrives(db *gorm.DB, olderThan *time.Time, newerThan *time.Time) ([]*HardDrive, error) {

	// Execute the lsblk command to get detailed block device information
	cmd := exec.Command("lsblk", "-d", "-o", "NAME,TRAN,SIZE,MODEL,SERIAL,TYPE")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return nil, err
	}

	var systemHardDrives []*HardDrive

	// Scan the output line by line
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip the header line
		if strings.Contains(line, "NAME") {
			continue
		}

		// Split the line into columns
		cols := strings.Fields(line)
		if len(cols) < 6 {
			continue
		}

		// Filter out nvme drives (M.2)
		if cols[1] != "nvme" && cols[5] != "Device" && cols[1] != "usb" {
			hd := &HardDrive{
				Name:      cols[0],
				Transport: cols[1],
				Size:      cols[2],
				Model:     cols[3],
				Serial:    cols[4],
				Type:      cols[5],
			}
			systemHardDrives = append(systemHardDrives, hd)
		}
	}

	// Handle error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var updatedHardDrives []*HardDrive

	for _, sysHDD := range systemHardDrives {
		var existingHD HardDrive
		q := db.Where("serial = ? AND model = ? AND type = ?", sysHDD.Serial, sysHDD.Model, sysHDD.Type)

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
