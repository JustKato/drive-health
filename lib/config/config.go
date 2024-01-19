package config

type DHConfig struct {
	DiskFetchFrequency  int `json:"diskFetchFrequency" comment:"How often should a snapshot be taken of the current state of the disks"`
	MemoryDumpFrequency int `json:"memoryDumpFrequency" comment:"How often should we save the snapshots from memory to disk"`
	MaxHistoryAge       int
}

func GetConfiguration() DHConfig {

	// TODO: Read from os.environment or simply load the defaults

	return DHConfig{
		DiskFetchFrequency:  5,
		MemoryDumpFrequency: 16,
		MaxHistoryAge:       2592000,
	}

}
