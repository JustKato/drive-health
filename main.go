package main

import (
	"fmt"

	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

func main() {

	hardDrives, err := hardware.GetSystemHardDrives()
	if err != nil {
		panic(err)
	}

	for _, hdd := range hardDrives {
		fmt.Printf("%s %s [%s]: %vC\n", hdd.Model, hdd.Serial, hdd.Size, hdd.GetTemperature(true))
	}

}
