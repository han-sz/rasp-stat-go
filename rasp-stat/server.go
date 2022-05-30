package main

import (
	"fmt"
	"os"
	"rasp-stat/rasp-stat/util"
	"sync"

	"github.com/gin-gonic/gin"
)

type RaspStatServerInterface interface {
	StartServer()
}

type RaspStatServer struct {
	RaspStatServerInterface

	Port       int
	StaticPath string

	Engine   *gin.Engine
	Service  *InstantStatService
	ReadLock *sync.Mutex
}

var none = makeDataInt(-1, "")

func (server *RaspStatServer) tryAcquireLock(critical func(failed bool)) {
	if locked := server.ReadLock.TryLock(); !locked {
		critical(true)
	} else {
		critical(false)
		server.ReadLock.Unlock()
	}
}

func (server *RaspStatServer) StartServer() bool {
	if server.Engine != nil {
		log.Log("Attempted to start server while it was already running")
		return false
	}
	server.Engine = gin.New()
	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
		log.Log("Running version:", VERSION)
	} else {
		server.Engine.Use(gin.Logger())
		log.Log("Running in debug mode, version:", VERSION)
	}
	server.Engine.Use(gin.Recovery())
	server.createApiRouter()
	server.createAdminRouter()

	// Static
	if _, err := os.Stat(server.StaticPath); !os.IsNotExist(err) {
		server.createStaticRouter()
		log.Log("Serving web-app from path:", server.StaticPath)
	}

	log.Log("Starting server on port", server.Port)
	err := server.Engine.Run(fmt.Sprintf(":%d", server.Port))
	if err == nil {
	} else {
		log.Log("Could not start server", err.Error())
		return false
	}
	return true
}

func (server *RaspStatServer) createApiRouter() {
	api := server.Engine.Group("/")
	// Used as a prev val cache for when the read lock is already acquired by another request
	var cached sync.Map
	api.GET("/temp", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.temp) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestTemp := server.Service.temp[len(server.Service.temp)-1]
			response := makeDataFloat(latestTemp.temperature, latestTemp.unit, 1)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	api.GET("/gpu", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.gpu) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestGpu := server.Service.gpu[len(server.Service.gpu)-1]
			response := makeDataFloat(latestGpu.gpuSpeed, latestGpu.unit, 3)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	api.GET("/cpu", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.cpu) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestCpu := server.Service.cpu[len(server.Service.cpu)-1]
			response := makeDataFloat(latestCpu.cpuSpeed, latestCpu.unit, 3)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	api.GET("/volts", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.volts) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestVoltage := server.Service.volts[len(server.Service.volts)-1]
			response := makeDataFloat(latestVoltage.volts, latestVoltage.unit, 2)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	api.GET("/throttled", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.throttled) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestThrottleVal := server.Service.throttled[len(server.Service.throttled)-1]
			response := DataPoint{Value: latestThrottleVal.value}
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	// TODO as /mem/free
	api.GET("/memFree", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.memory) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestMemoryVal := server.Service.memory[len(server.Service.memory)-1]
			response := makeDataInt(latestMemoryVal.memFree, latestMemoryVal.unit)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	// TODO as /mem/total
	api.GET("/memTotal", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.memory) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestMemoryVal := server.Service.memory[len(server.Service.memory)-1]
			response := makeDataInt(latestMemoryVal.memTotal, latestMemoryVal.unit)
			c.JSON(200, response)
			server.ReadLock.Unlock()

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	// TODO as /mem/swap
	api.GET("/memSwap", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.memory) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestMemoryVal := server.Service.memory[len(server.Service.memory)-1]
			response := makeDataInt(latestMemoryVal.memSwap, latestMemoryVal.unit)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	// TODO as /mem/swap-total
	api.GET("/memSwapTotal", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed || len(server.Service.memory) < 1 {
				c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
				return
			}
			latestMemoryVal := server.Service.memory[len(server.Service.memory)-1]
			response := makeDataInt(latestMemoryVal.memSwapTotal, latestMemoryVal.unit)
			c.JSON(200, response)

			cached.Store(c.FullPath(), response)
		}
		server.tryAcquireLock(critical)
	})
	api.GET("/raw", func(c *gin.Context) {
		critical := func(failed bool) {
			if failed {
				c.JSON(200, none)
				return
			}
			c.JSON(200, gin.H{
				"data": gin.H{
					"cpu":       server.Service.RawCpuDataPoints(),
					"gpu":       server.Service.RawGpuDataPoints(),
					"temp":      server.Service.RawTemperatureDataPoints(),
					"volts":     server.Service.RawCpuVoltageDataPoints(),
					"throttled": server.Service.RawCpuThrottledDataPoints(),
				},
			})
		}
		server.tryAcquireLock(critical)
	})
}

func (server *RaspStatServer) createAdminRouter() {
	r := server.Engine.Group("/power")
	r.GET("/off/:min", func(c *gin.Context) {
		minutes := util.ToInt(c.Param("min"))
		if minutes < 0 {
			minutes = 0
		}
		_, err := commandOutput(fmt.Sprintf("shutdown --poweroff %d", minutes))
		log.Log("Power-off command received..")
		if err == nil {
			log.Log("Power-off has been scheduled for", minutes, "minutes")
			c.JSON(200, gin.H{"status": "triggered"})
		} else {
			log.Log("Power-off was cancelled", err.Error())
			c.JSON(500, gin.H{"status": "failed", "error": err.Error()})
		}
	})
	r.GET("/reboot", func(c *gin.Context) {
		_, err := commandOutput("reboot")
		log.Log("Reboot command received..")
		if err == nil {
			c.JSON(200, gin.H{"status": "triggered"})
		} else {
			log.Log("Reboot was cancelled", err.Error())
			c.JSON(500, gin.H{"status": "failed", "error": err.Error()})
		}
	})
}

func (server *RaspStatServer) createStaticRouter() {
	r := server.Engine.Group("/app")
	r.Static("/", server.StaticPath)
}

type DataPoint struct {
	Value string `json:"data"`
}

func makeDataFloat(value float32, unit string, rounding int) DataPoint {
	return DataPoint{Value: fmt.Sprintf("%."+fmt.Sprint(rounding)+"f %s", value, unit)}
}

func makeDataInt(value int, unit string) DataPoint {
	return DataPoint{Value: fmt.Sprintf("%d %s", value, unit)}
}
