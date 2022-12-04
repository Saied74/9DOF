package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"syscall"
	"time"
)

const remoteAddr = "192.168.4.75"

//<++++++++++++++++++++++   Query head end   +++++++++++++++++++++++++++>

//Makes an HTTP call to the far end with the parameter given.  In most cases,
//it does not return an error (see below) but returns the string from the far end
func (app *application) getRemote(q string) (string, error) {
	/*
	  q = s for get sensors
	  q = c for get calibration data
	  q = r for get offset data
	  q = w for write offset data
	*/
	client := &http.Client{
		Timeout: 100 * time.Second,
	}
	// TODO: replace remote address with server name
	//better yet, make the server name or IP address a command line flag
	url := fmt.Sprintf("http://%s/?q=%s:80", remoteAddr, q)
	response, err := client.Get(url)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			app.errorLog.Printf("%v", err)
			return "Context deadline exceeded (1)", nil
		}
		if e, ok := err.(net.Error); ok && e.Timeout() { //timeout error
			app.errorLog.Printf("%v", err)
			return "Connection to the head end is lost", nil
		}
		if errors.Is(err, syscall.EACCES) { //Access denied
			app.errorLog.Printf("%v", err)
			return "Connection access denied", nil
		}
		if errors.Is(err, syscall.ECONNREFUSED) { //connection refused
			app.errorLog.Printf("%v", err)
			return "Connection to the head end refused", nil
		}
		if errors.Is(err, syscall.ECONNRESET) { //connecton reset
			app.errorLog.Printf("%v", err)
			return "Connection reset by the head end", nil
		}
		if errors.Is(err, syscall.EHOSTDOWN) { //host down
			app.errorLog.Printf("%v", err)
			return "Host is down", nil
		}
		return "", err
	}
	defer response.Body.Close()
	buf := make([]byte, 1024)
	buf, err = io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

//<++++++++++++++++   centralized error handling   +++++++++++++++++++>

//This is straight out of Alex Edward's Let's Go book
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) //to not get the helper file...
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}
