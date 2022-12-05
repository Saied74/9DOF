package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Saied74/cli"
)

//for injecting data into handlers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	debugOption   bool
	templateCache map[string]*template.Template
}

var iterCount = 0

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		//		templateCache: templateCache,
	}
	iterCount, _ = strconv.Atoi(uiItems.ItemList["iterations"].Value)
	c := cli.Command(&uiItems)
	for {
		item := <-c
		switch item.Name {
		case "quit":
			os.Exit(0)
		case "iterations":
			iterCount, _ = strconv.Atoi(item.Value)
		case "calibration":
			for i := 0; i < iterCount; i++ {
				err := app.showResults("c")
				if err != nil {
					log.Fatal(err)
				}
			}
		case "sensors":
			for i := 0; i < iterCount; i++ {
				err := app.showResults("s")
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
		default:
			fmt.Println("Can't do that: ", item.Name, item.Value)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (app *application) showResults(option string) error {
	result, err := app.getRemote(option)
	if err != nil {
		return err
	}
	result = strings.TrimLeft(result, "<!DOCTYPE HTML><html>")
	result = strings.TrimRight(result, "</html>\r\n")
	fmt.Println(result)
	return nil
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
