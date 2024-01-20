package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetSystemHardDrives() ([]*HardDrive, error) {

	// Execute the lsblk command to get detailed block device information
	cmd := exec.Command("lsblk", "-d", "-o", "NAME,TRAN,SIZE,MODEL,SERIAL,TYPE")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return nil, err
	}

	var hardDrives []*HardDrive

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
		if cols[1] != "nvme" && cols[5] != "Device" {
			hd := &HardDrive{
				Name:        cols[0],
				Transport:   cols[1],
				Size:        cols[2],
				Model:       cols[3],
				Serial:      cols[4],
				Type:        cols[5],
				Temperature: 0,
			}
			hardDrives = append(hardDrives, hd)
		}
	}

	// Handle error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hardDrives, nil
}
