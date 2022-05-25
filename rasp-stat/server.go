package main

import (
	"fmt"
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

func (params *RaspStatServer) StartServer() {
	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
		log.Log("Running version:", VERSION)
	} else {
		log.Log("Running in debug mode, version:", VERSION)
	}
	r := gin.Default()
	r.GET("/temp", func(c *gin.Context) {
		if locked := params.ReadLock.TryLock(); locked && len(params.Service.temp) > 0 {
			latestTemp := params.Service.temp[len(params.Service.temp)-1]
			c.JSON(200, DataPoint{Value: latestTemp.temperature})
			params.ReadLock.Unlock()
		} else {
			c.JSON(200, none)
		}
	})
	r.GET("/gpu", func(c *gin.Context) {
		if locked := params.ReadLock.TryLock(); locked && len(params.Service.gpu) > 0 {
			latestGpu := params.Service.gpu[len(params.Service.gpu)-1]
			c.JSON(200, makeData(latestGpu.gpuSpeed, latestGpu.unit, 3))
			params.ReadLock.Unlock()
		} else {
			c.JSON(200, none)
		}
	})
	r.GET("/cpu", func(c *gin.Context) {
		if locked := params.ReadLock.TryLock(); locked && len(params.Service.cpu) > 0 {
			latestCpu := params.Service.cpu[len(params.Service.cpu)-1]
			c.JSON(200, makeData(latestCpu.cpuSpeed, latestCpu.unit, 3))
			params.ReadLock.Unlock()
		} else {
			c.JSON(200, none)
		}
	})
	r.GET("raw", func(c *gin.Context) {
		if locked := params.ReadLock.TryLock(); locked {
			c.JSON(200, gin.H{
				"data": gin.H{
					"cpu":  params.Service.RawCpuDataPoints(),
					"gpu":  params.Service.RawGpuDataPoints(),
					"temp": params.Service.RawTemperatureDataPoints(),
				},
			})
			params.ReadLock.Unlock()
		} else {
			c.JSON(200, none)
		}
	})
	r.GET("/volts", func(c *gin.Context) {
		c.JSON(200, none)
	})
	r.GET("/throttled", func(c *gin.Context) {
		c.JSON(200, none)
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

	log.Log("Started server on port", params.Port)
	r.Run(fmt.Sprintf(":%d", params.Port))
}
