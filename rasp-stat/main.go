package main

import (
	"os"
	"rasp-stat/rasp-stat/util"
	"sync"
)

const (
	VERSION string = "1.0.5"
	DEBUG   bool   = false
)

const (
	DEFAULT_APP_PATH = "/var/rasp-stat/"
	DEFAULT_PORT     = 4322
	DEFAULT_BUFFER   = 10
	DEFAULT_INTERVAL = 5
)

var log util.Loggable = &util.Log{Name: "rasp-stat", Version: VERSION}

func startFetchServiceAndServer() {
	port, interval, buffer, path := env()

	var rwMutex sync.Mutex // TODO use semaphore instead for read block
	var server RaspStatServer
	var iss InstantStatService

	iss = NewInstantStatService(
		uint16(interval),
		uint16(buffer),
	)
	iss.ReadWriteLock = &rwMutex
	iss.FetchAndCacheStats()

	server.Port = port
	server.Service = &iss
	server.StaticPath = path
	server.ReadLock = &rwMutex
	server.StartServer()
}

func main() {
	startFetchServiceAndServer()
}

func env() (port, interval, buffer int, path string) {
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
		log.Log("Setting data point buffer from environment:", buffer)
	}
	appPathEnv, found := os.LookupEnv("APP_PATH")
	if found {
		if appPathEnv != "" {
			path = appPathEnv
		}
		log.Log("Setting app path from environment:", path)
	}
	return
}
