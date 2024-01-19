package svc

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"

	"tea.chunkbyte.com/kato/drive-health/lib/config"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

// The path to where the snapshot database is located
const SNAPSHOT_LIST_PATH = "./snapshots.dat"

// A simple in-memory buffer for the history of snapshots
var snapShotBuffer []*HardwareSnapshot

// A snapshot in time of the current state of the harddrives
type HardwareSnapshot struct {
	TimeStamp time.Time
	HDD       []*hardware.HardDrive
}

type Snapshots struct {
	List []*HardwareSnapshot
}

// The function itterates through all hard disks and takes a snapshot of their state,
// returns a struct which contains metadata as well as the harddrives themselves.
func TakeHardwareSnapshot() (*HardwareSnapshot, error) {
	drives, err := hardware.GetSystemHardDrives()
	if err != nil {
		return nil, err
	}

	snapShot := &HardwareSnapshot{
		TimeStamp: time.Now(),
		HDD:       []*hardware.HardDrive{},
	}

	for _, hdd := range drives {
		hdd.GetTemperature(true)
		snapShot.HDD = append(snapShot.HDD, hdd)
	}

	// Append to the in-memory listing
	snapShotBuffer = append(snapShotBuffer, snapShot)

	// Return the snapshot just in case there is any need to modify it,
	// any modification to it will also affect the current buffer from memory.
	return snapShot, nil
}

// The function wil check if the `.dat` file is present, if it is then it will load it into memory
func UpdateHardwareSnapshotsFromFile() {
	file, err := os.Open(SNAPSHOT_LIST_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			return // File does not exist, no snapshots to load
		}
		panic(err) // Handle error according to your error handling policy
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var snapshots Snapshots
	if err := decoder.Decode(&snapshots); err != nil {
		if err == io.EOF {
			return // End of file reached
		}
		panic(err) // Handle error according to your error handling policy
	}

	snapShotBuffer = snapshots.List

	fmt.Printf("Loaded %v snapshots from .dat", len(snapShotBuffer))
}

// Get the list of snapshots that have been buffered in memory
func GetHardwareSnapshot() []*HardwareSnapshot {
	return snapShotBuffer
}

// Dump the current snapshot history from memory to file
func SaveSnapshotsToFile() error {
	file, err := os.Create(SNAPSHOT_LIST_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	snapshots := Snapshots{List: snapShotBuffer}
	if err := encoder.Encode(snapshots); err != nil {
		return err
	}

	return nil
}

func RunService() {

	// Snapshot taking routine
	go func() {
		for {
			time.Sleep(time.Duration(config.GetConfiguration().DiskFetchFrequency) * time.Second)
			data, err := TakeHardwareSnapshot()
			if err != nil {
				fmt.Printf("Hardware Fetch Error: %s", err)
			} else {
				fmt.Println("Got Snapshot for " + data.TimeStamp.Format("02/01/2006"))
				for _, hdd := range data.HDD {
					fmt.Printf("%s[%s]: %v\n", hdd.Model, hdd.Size, hdd.Temperature)
				}
			}
		}
	}()

	// Periodic saving routine
	go func() {
		for {
			time.Sleep(time.Duration(config.GetConfiguration().MemoryDumpFrequency) * time.Second)
			err := SaveSnapshotsToFile()
			if err != nil {
				fmt.Printf("Memory Dump Error: %s", err)
			}

			fmt.Println("Saved Snapshots to file")
		}
	}()
}
