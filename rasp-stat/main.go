package main

import (
	"os"
	"rasp-stat/rasp-stat/util"
	"sync"
)

const (
	VERSION string = "1.0.2"
	DEBUG   bool   = false
)

const (
	DEFAULT_PORT     = 4322
	DEFAULT_BUFFER   = 10
	DEFAULT_INTERVAL = 5
)

var log util.Loggable = &util.Log{Name: "rasp-stat", Version: VERSION}

func env() (port, interval, buffer int) {
	port = DEFAULT_PORT
	buffer = DEFAULT_BUFFER
	interval = DEFAULT_INTERVAL

	portEnv, found := os.LookupEnv("PORT")
	if found {
		parsedPort := util.ToInt(portEnv)
		if parsedPort != -1 {
			port = parsedPort
		}
		log.Log("Setting port from environment:", parsedPort)
	}
	intervalEnv, found := os.LookupEnv("FETCH_INTERVAL")
	if found {
		parsedInterval := util.ToInt(intervalEnv)
		if parsedInterval != -1 {
			interval = parsedInterval
		}
		log.Log("Setting fetch interval from environment:", interval)
	}
	bufferEnv, found := os.LookupEnv("BUFFER")
	if found {
		parsedBuffer := util.ToInt(bufferEnv)
		if parsedBuffer != -1 {
			buffer = parsedBuffer
		}
		log.Log("Setting data point buffer from environment:", interval)
	}
	return
}

func startFetchServiceAndServer() {
	port, interval, buffer := env()

	var rwMutex sync.Mutex
	var wait sync.WaitGroup
	var server RaspStatServer
	var iss InstantStatService

	wait.Add(1)
	iss = NewInstantStatService(
		uint16(interval),
		uint16(buffer),
	)
	iss.ReadWriteLock = &rwMutex
	iss.FetchAndCacheStats()

	wait.Add(1)
	server.Port = port
	server.Service = &iss
	server.ReadLock = &rwMutex
	server.StartServer()

	wait.Wait()
}

func main() {
	startFetchServiceAndServer()

}
