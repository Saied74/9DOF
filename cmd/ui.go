package main

import (
	"strconv"

	"github.com/Saied74/cli"
)

var uiItems = cli.Items{
	OrderList: []string{"iterations", "calibration", "sensors",
		"getOffsets", "updateOffsets", "quit"},
	ItemList: map[string]*cli.Item{
		"iterations": &cli.Item{
			Name:      "Iterations",
			Prompt:    "How many iterations do you want to run",
			Value:     "10",
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
