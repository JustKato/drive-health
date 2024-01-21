package svc

import (
	"bytes"
	"fmt"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"tea.chunkbyte.com/kato/drive-health/lib/config"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

var db *gorm.DB

func InitDB() {
	var err error
	dbPath := config.GetConfiguration().DatabaseFilePath
	if dbPath == "" {
		dbPath = "./data.sqlite"
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&hardware.HardDrive{}, &hardware.HardDriveTemperature{})
}

func GetDatabaseRef() *gorm.DB {
	return db
}

func LogDriveTemps() error {
	drives, err := hardware.GetSystemHardDrives(db, nil, nil)
	if err != nil {
		return err
	}

	for _, hdd := range drives {
		temp := hdd.GetTemperature()
		db.Create(&hardware.HardDriveTemperature{
			HardDriveID: hdd.ID,
			TimeStamp:   time.Now(),
			Temperature: temp,
		})
	}

	return nil
}

func RunLoggerService() {
	fmt.Println("Initializing Temperature Logging Service...")

	tickTime := time.Duration(config.GetConfiguration().DiskFetchFrequency) * time.Second

	// Snapshot taking routine
	go func() {
		for {
			time.Sleep(tickTime)
			err := LogDriveTemps()
			if err != nil {
				fmt.Printf("🛑 Temperature logging failed: %s\n", err)
			}
		}
	}()
}

func GetDiskGraphImage(hddID int, newerThan *time.Time, olderThan *time.Time) (*bytes.Buffer, error) {
	var hdd hardware.HardDrive
	// Fetch by a combination of fields
	q := db.Where("id = ?", hddID)

	if newerThan == nil && olderThan == nil {
		q = q.Preload("Temperatures")
	} else {
		q = q.Preload("Temperatures", "time_stamp < ? AND time_stamp > ?", newerThan, olderThan)
	}

	// Query for the instance
	result := q.First(&hdd)
	if result.Error != nil {
		return nil, result.Error
	}

	// Prepare slices for X (time) and Y (temperature) values
	var xValues []time.Time
	var yValues []float64
	for _, temp := range hdd.Temperatures {
		xValues = append(xValues, temp.TimeStamp)
		yValues = append(yValues, float64(temp.Temperature))
	}

	// Allocate a buffer for the graph image
	graphImageBuffer := bytes.NewBuffer([]byte{})

	// TODO: Graph dark theme

	// Generate the chart
	graph := chart.Chart{
		Title: fmt.Sprintf("%s:%s[%s]", hdd.Name, hdd.Serial, hdd.Size),
		TitleStyle: chart.Style{
			FontSize: 14,
		},

		// TODO: Implement customizable sizing
		Width: 1280,

		Background: chart.Style{
			Padding: chart.Box{
				Top: 20, Right: 20, Bottom: 20, Left: 20,
			},
		},

		XAxis: chart.XAxis{
			Name: "Time",
			ValueFormatter: func(v interface{}) string {
				if ts, isValidTime := v.(float64); isValidTime {
					t := time.Unix(int64(ts/1e9), 0)

					return t.Format("Jan 2 2006, 15:04")
				}

				return ""
			},
			Style: chart.Style{},
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 0.5,
			},
			GridMinorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray.WithAlpha(64),
				StrokeWidth: 0.25,
			},
		},
		YAxis: chart.YAxis{
			Name:  "Temperature (C)",
			Style: chart.Style{},
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 0.5,
			},
			GridMinorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray.WithAlpha(64),
				StrokeWidth: 0.25,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Temperature",
				XValues: xValues,
				YValues: yValues,
				Style: chart.Style{
					StrokeColor: chart.ColorCyan,
					StrokeWidth: 2.0,
				},
			},
		},
	}

	// Add a legend to the chart
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph, chart.Style{
			Padding: chart.Box{
				Top: 5, Right: 5, Bottom: 5, Left: 5,
			},
			FontSize: 10,
		}),
	}

	// Render the chart into the byte buffer
	err := graph.Render(chart.PNG, graphImageBuffer)
	if err != nil {
		return nil, err
	}

	return graphImageBuffer, nil
}
