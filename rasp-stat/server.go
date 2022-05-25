package main

import (
	"fmt"
	"rasp-stat/rasp-stat/util"
	"sync"

	"github.com/gin-gonic/gin"
)

type RaspStatServer struct {
	Port int

	Service  *InstantStatService
	ReadLock *sync.Mutex
}

type DataPoint struct {
	Value string `json:"data"`
}

func makeData(value float32, unit string, rounding int) DataPoint {
	return DataPoint{Value: fmt.Sprintf("%."+fmt.Sprint(rounding)+"f %s", value, unit)}
}

var none = makeData(-1, "", 0)

func (server *RaspStatServer) StartServer() bool {
	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
		log.Log("Running version:", VERSION)
	} else {
		log.Log("Running in debug mode, version:", VERSION)
	}
	r := gin.New()
	if DEBUG {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	// Used as a prev val cache for when the read lock is already acquired by another request
	var cached sync.Map
	r.GET("/temp", func(c *gin.Context) {
		if locked := server.ReadLock.TryLock(); locked && len(server.Service.temp) > 0 {
			latestTemp := server.Service.temp[len(server.Service.temp)-1]
			response := makeData(latestTemp.temperature, latestTemp.unit, 1)
			c.JSON(200, response)
			server.ReadLock.Unlock()

			cached.Store(c.FullPath(), response)
		} else {
			c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
		}
	})
	r.GET("/gpu", func(c *gin.Context) {
		if locked := server.ReadLock.TryLock(); locked && len(server.Service.gpu) > 0 {
			latestGpu := server.Service.gpu[len(server.Service.gpu)-1]
			response := makeData(latestGpu.gpuSpeed, latestGpu.unit, 3)
			c.JSON(200, response)
			server.ReadLock.Unlock()

			cached.Store(c.FullPath(), response)
		} else {
			c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
		}
	})
	r.GET("/cpu", func(c *gin.Context) {
		if locked := server.ReadLock.TryLock(); locked && len(server.Service.cpu) > 0 {
			latestCpu := server.Service.cpu[len(server.Service.cpu)-1]
			response := makeData(latestCpu.cpuSpeed, latestCpu.unit, 3)
			c.JSON(200, response)
			server.ReadLock.Unlock()

			cached.Store(c.FullPath(), response)
		} else {
			c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
		}
	})
	r.GET("/volts", func(c *gin.Context) {
		if locked := server.ReadLock.TryLock(); locked && len(server.Service.volts) > 0 {
			latestVoltage := server.Service.volts[len(server.Service.volts)-1]
			response := makeData(latestVoltage.volts, latestVoltage.unit, 2)
			c.JSON(200, response)
			server.ReadLock.Unlock()

			cached.Store(c.FullPath(), response)
		} else {
			c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
		}
	})
	r.GET("/throttled", func(c *gin.Context) {
		if locked := server.ReadLock.TryLock(); locked && len(server.Service.throttled) > 0 {
			latestThrottleVal := server.Service.throttled[len(server.Service.throttled)-1]
			response := DataPoint{Value: latestThrottleVal.value}
			c.JSON(200, response)
			server.ReadLock.Unlock()

			cached.Store(c.FullPath(), response)
		} else {
			c.JSON(200, util.GetOrDefaultSafeMap(&cached, c.FullPath(), none))
		}
	})
	r.GET("/raw", func(c *gin.Context) {
		if locked := server.ReadLock.TryLock(); locked {
			c.JSON(200, gin.H{
				"data": gin.H{
					"cpu":       server.Service.RawCpuDataPoints(),
					"gpu":       server.Service.RawGpuDataPoints(),
					"temp":      server.Service.RawTemperatureDataPoints(),
					"volts":     server.Service.RawCpuVoltageDataPoints(),
					"throttled": server.Service.RawCpuThrottledDataPoints(),
				},
			})
			server.ReadLock.Unlock()
		} else {
			c.JSON(200, none)
		}
	})
	r.GET("/mem", func(c *gin.Context) {
		c.JSON(200, none)
	})
	r.GET("/memFree", func(c *gin.Context) {
		c.JSON(200, none)
	})
	r.GET("/memTotal", func(c *gin.Context) {
		c.JSON(200, none)
	})
	r.GET("/memSwap", func(c *gin.Context) {
		c.JSON(200, none)
	})
	r.GET("/memSwapTotal", func(c *gin.Context) {
		c.JSON(200, none)
	})
	r.GET("/wifi", func(c *gin.Context) {
		c.JSON(200, none)
	})

	log.Log("Starting server on port", server.Port)
	err := r.Run(fmt.Sprintf(":%d", server.Port))
	if err == nil {
	} else {
		log.Log("Could not start server", err.Error())
		return false
	}
	return true
}
