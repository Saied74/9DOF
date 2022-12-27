package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Saied74/cli"
)

//for injecting data into handlers
type application struct {
	iterCount     int
	magneticDecl  float64
	fileName      string
	file          *os.File
	azimuth       float64
	elevation     float64
	errorLog      *log.Logger
	infoLog       *log.Logger
	debugOption   bool
	templateCache map[string]*template.Template
}

//var iterCount = 1

const magneticDecl = 12.25 //degrees - subtract from yaw to get true north

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		magneticDecl: magneticDecl,
		//		templateCache: templateCache,
	}
	app.iterCount, _ = strconv.Atoi(uiItems.ItemList["iterations"].Value)
	c := cli.Command(&uiItems)
	for {
		item := <-c
		//log.Fatal("Name: ", item.Name, "Value: ", item.Value)
		switch item.Name {
		case "quit":
			os.Exit(0)
		case "iterations":
			app.iterCount, _ = strconv.Atoi(item.Value)
		case "calibration":
			var result string
			var err error
			for i := 0; i < app.iterCount; i++ {
				result, err = app.showCalResults("c")
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Println(result)
		case "sensors":
			for i := 0; i < app.iterCount; i++ {
				_, err := app.showSensorResults("s")
				if err != nil {
					log.Fatal(err)
				}
			}
		case "getOffsets":
			err := app.storeOffsets()
			if err != nil {
				log.Fatal(err)
			}
		case "updateOffsets":
			err := app.updateOffsets()
			if err != nil {
				log.Fatal(err)
			}
		case "fileName":
			app.fileName = filepath.Join(".", item.Value)
		case "openFile":
			if app.fileName != "" {
				f, err := os.OpenFile(app.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				app.file = f
			} else {
				fmt.Println("Filename is blank, try again")
			}
		case "closeFile":
			app.file.Close()
		case "azimuth":
			az, _ := strconv.ParseFloat(item.Value, 64)
			app.azimuth = az
		case "elevation":
			el, _ := strconv.ParseFloat(item.Value, 64)
			app.elevation = el
		case "record":
			err := app.recordData()
			if err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Println("Can't do that: ", item.Name, item.Value)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (app *application) showCalResults(option string) (string, error) {
	result, err := app.getRemote(option)
	if err != nil {
		return "", err
	}
	result = strings.TrimLeft(result, "<!DOCTYPE HTML><html>")
	result = strings.TrimRight(result, "</html>\r\n")

	return result, nil
}

func (app *application) showSensorResults(option string) ([]float64, error) {
	angles := []string{"Alpha(Yaw)", "Beta(Roll)", "Gamma(Pitch)"}
	angleResults := []float64{}
	result, err := app.getRemote(option)
	if err != nil {
		return []float64{}, err
	}
	result = strings.TrimLeft(result, "<!DOCTYPE HTML><html>")
	result = strings.TrimRight(result, "</html>\r\n")
	results := strings.Split(result, ",")
	if len(results) != 3 {
		log.Fatal("did not get back 3 points of sensor data")
	}
	var raw, processed, azEl string
	for n, res := range results {
		r := strings.Split(res, ":")
		if len(r) != 2 {
			return []float64{}, fmt.Errorf("this sensor data %v did not split into two", res)
		}
		raw += fmt.Sprintf("%d = %s\t", n, r[1])
		switch n {
		case 0:
			yaw, err := strconv.ParseFloat(r[1], 64)
			if err != nil {
				return []float64{}, err
			}
			angleResults = append(angleResults, yaw)
			yaw = 360 - yaw + app.magneticDecl
			if yaw >= 360.0 {
				yaw = yaw - 360
			}
			processed += fmt.Sprintf("%s = %0.2f", angles[n], yaw)
			azEl += fmt.Sprintf("Azimuth: %0.2f\t", 360.0-yaw)
		case 1:
			roll, err := strconv.ParseFloat(r[1], 64)
			if err != nil {
				return []float64{}, err
			}
			angleResults = append(angleResults, roll)
			processed += fmt.Sprintf("\t%s = %0.2f", angles[n], -roll)
			azEl += fmt.Sprintf("Elevation: %0.2f", roll)
		case 2:
			pitch, err := strconv.ParseFloat(r[1], 64)
			if err != nil {
				return []float64{}, err
			}
			angleResults = append(angleResults, pitch)
			processed += fmt.Sprintf("\t%s = %0.2f", angles[n], -pitch)
		}
	}
	fmt.Println(raw)
	fmt.Println(processed)
	fmt.Println(azEl)
	return angleResults, nil
}

func (app *application) storeOffsets() error {
	result, err := app.getRemote("r")
	if err != nil {
		return err
	}
	result = strings.TrimLeft(result, "<!DOCTYPE HTML><html>")
	result = strings.TrimRight(result, "</html>\r\n")
	err = os.WriteFile("offsets.txt", []byte(result), 0666)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) updateOffsets() error {
	offsetBytes, err := os.ReadFile("offsets.txt")
	if err != nil {
		return err
	}
	offsets := "w" + string(offsetBytes) + "$"
	_, err = app.getRemote(offsets)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) recordData() error {
	calResults, err := app.showCalResults("c")
	if err != nil {
		return err
	}
	calParts := strings.Split(calResults, ",")
	if len(calParts) != 4 {
		return fmt.Errorf("returned bad calibration data %v", calResults)
	}
	for _, cal := range calParts {
		c := strings.Split(cal, ":")
		if len(c) != 2 {
			return fmt.Errorf("bad calibration compoinent %v", cal)
		}
		calData, err := strconv.Atoi(c[1])
		if err != nil {
			return fmt.Errorf("calibration data did not convert to int %v", c[1])
		}
		if calData != 3 {
			app.recalibrate(c[0], calData)
		}
	}
	s, err := app.showSensorResults("s")
	if err != nil {
		return err
	}
	if len(s) != 3 {
		return fmt.Errorf("Bad sensor data %v", s)
	}
	csvLine := fmt.Sprintf("%0.1f,%0.1f,%0.2f,%0.2f,%0.2f\n",
		app.azimuth, app.elevation, s[0], s[1], s[2])
	_, err = app.file.WriteString(csvLine)
	if err != nil {
		return err
	}
	fmt.Printf("\n0 = %0.2f\t1 = %0.2f\t2 = %0.2f\n", s[0], s[1], s[2])
	return nil
}

func (app *application) recalibrate(s string, n int) {
	fmt.Printf("warning: out of calibration, the value is %s:%d\n", s, n)
}
