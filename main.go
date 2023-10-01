package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	dbUri, port := getEnvForStarup()
	logFile := loggerStartup()
	cancel, err := connect(dbUri)
	defer close(cancel)
	defer logFile.Close()
	if err != nil {
		panic(err)
	}
	ping()
	fmt.Print("Application Started complete log available in application.log file !")
	router := AllHandlers()
	log.Fatal(http.ListenAndServe(port, router))
}
