package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"
)

//for injecting data into handlers
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	debugOption   bool
	templateCache map[string]*template.Template
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		//		templateCache: templateCache,
	}
	for {
		// options := []string{"s", "c", "r", "w"}
		// for _, q := range options {
		result, err := app.getRemote("c")
		if err != nil {
			log.Fatal(err)
		}
		result = strings.TrimLeft(result, "<!DOCTYPE HTML><html>")
		result = strings.TrimRight(result, "</html>\r\n")
		fmt.Println(result)
		// }
		time.Sleep(500 * time.Millisecond)
	}
}
