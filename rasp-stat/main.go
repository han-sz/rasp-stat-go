package main

import (
	"os"
	"rasp-stat/rasp-stat/util"
	"strconv"
	"sync"
)

const (
	VERSION string = "1.0.1"
	DEBUG   bool   = false
)

var log util.Loggable = &util.Log{Name: "rasp-stat", Version: VERSION}

func env() (port int) {
	port = 4322

	portEnv, found := os.LookupEnv("PORT")
	if found {
		parsedPort, err := strconv.Atoi(portEnv)
		if err == nil {
			port = parsedPort
		}
		log.Log("Using port from environment", parsedPort)
	}
	return
}

func startFetchServiceAndServer() {
	var rwMutex sync.Mutex
	var wait sync.WaitGroup
	var iss = NewInstantStatService(5, 2)
	var server RaspStatServer = RaspStatServer{
		Service:  &iss,
		ReadLock: &rwMutex,
	}

	port := env()

	wait.Add(1)
	iss.ReadWriteLock = &rwMutex
	iss.FetchAndCacheStats()

	wait.Add(1)
	server.Port = port
	server.StartServer()

	wait.Wait()
}

func main() {
	startFetchServiceAndServer()

}
