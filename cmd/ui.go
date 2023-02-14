package main

import (
	"strconv"
	"strings"

	"github.com/Saied74/cli"
)

var uiItems = cli.Items{
	OrderList: []string{"iterations", "calibration", "sensors",
		"getOffsets", "updateOffsets", "fileName", "openFile",
		"closeFile", "az_start", "az_end", "az_inc", "el_start",
		"el_end", "el_inc", "record", "quit"},
	ItemList: map[string]*cli.Item{
		"iterations": &cli.Item{
			Name:      "Iterations",
			Prompt:    "How many iterations do you want to run",
			Value:     "1",
			Validator: intValidator,
		},
		"calibration": &cli.Item{
			Name:      "Calibration",
			Prompt:    "Show calibration data",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"sensors": &cli.Item{
			Name:      "Sensors",
			Prompt:    "Show sensor data",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"getOffsets": &cli.Item{
			Name:      "getOffsets",
			Prompt:    "Get Offsets",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"updateOffsets": &cli.Item{
			Name:      "updateOffsets",
			Prompt:    "Update Offsets",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"fileName": &cli.Item{
			Name:      "fileName",
			Prompt:    "Enter output file name (ending in csv)",
			Value:     "sensordata.csv",
			Validator: filenameValidator,
		},
		"openFile": &cli.Item{
			Name:      "openFile",
			Prompt:    "Open csv file for writing sensor data",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"closeFile": &cli.Item{
			Name:      "closeFile",
			Prompt:    "Close csv file when done",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"az_start": &cli.Item{
			Name:      "az_start",
			Prompt:    "Azimuth start",
			Value:     "55.0",
			Validator: floatValidator,
		},
		"az_end": &cli.Item{
			Name:      "az_end",
			Prompt:    "Azimuth end",
			Value:     "300.0",
			Validator: floatValidator,
		},
		"az_inc": &cli.Item{
			Name:      "az_inc",
			Prompt:    "Azimuth increment",
			Value:     "5.0",
			Validator: floatValidator,
		},
		"el_start": &cli.Item{
			Name:      "el_start",
			Prompt:    "Elevation start",
			Value:     "5.0",
			Validator: floatValidator,
		},
		"el_end": &cli.Item{
			Name:      "el_end",
			Prompt:    "Elevation end",
			Value:     "90.0",
			Validator: floatValidator,
		},
		"el_inc": &cli.Item{
			Name:      "el_inc",
			Prompt:    "Elevation increment",
			Value:     "5.0",
			Validator: floatValidator,
		},
		"record": &cli.Item{
			Name:      "record",
			Prompt:    "Record Sensor Data",
			Value:     "",
			Validator: cli.ItemValidator(func(x string) bool { return true }),
		},
		"quit": &cli.Item{
			Name:      "quit",
			Prompt:    "Quit",
			Value:     "",
			Validator: quitValidator,
		},
	},
}

var intValidator = cli.ItemValidator(func(x string) bool {
	_, err := strconv.Atoi(x)
	if err != nil {
		return false
	}
	return true
})

var floatValidator = cli.ItemValidator(func(x string) bool {
	_, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return false
	}
	return true
})

var filenameValidator = cli.ItemValidator(func(x string) bool {
	filenameBits := strings.Split(x, ".")
	if len(filenameBits) != 2 {
		return false
	}
	if filenameBits[1] != "csv" {
		return false
	}
	return true
})

var quitValidator = cli.ItemValidator(func(x string) bool {
	switch x {
	case "quit":
		return true
	case "Quit":
		return true
	case "q":
		return true
	case "Q":
		return true
	}
	return false
})
