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

func makeData(value string) DataPoint {
	return DataPoint{Value: value}
}

func (params *RaspStatServer) StartServer() {
	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
		log.Log("Running version:", VERSION)
	} else {
		log.Log("Running in debug mode, version:", VERSION)
	}
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		// if res, err := commandOutput("/opt/vc/bin/vcgencmd measure_temp"); err == nil {
		// 	log.Log("Temp", res)
		// 	c.JSON(200, makeData(res))
		// } else {
		// 	c.JSON(200, makeData("Unknown"))
		// }
	})
	r.GET("/temp", func(c *gin.Context) {
		if locked := params.ReadLock.TryLock(); locked && len(params.Service.temp) > 0 {
			latestTemp := params.Service.temp[len(params.Service.temp)-1]
			c.JSON(200, makeData(latestTemp.temperatureCelsius))
			params.ReadLock.Unlock()
		} else {
			c.JSON(200, makeData("-1"))
		}
	})
	r.GET("/gpu", func(c *gin.Context) {
		if len(params.Service.gpu) > 0 {
			latestGpu := params.Service.gpu[len(params.Service.gpu)-1]
			c.JSON(200, makeData(latestGpu.gpuSpeed))
		} else {
			c.JSON(200, makeData("-1"))
		}
	})
	r.GET("/cpu", func(c *gin.Context) {
		if len(params.Service.cpu) > 0 {
			latestCpu := params.Service.cpu[len(params.Service.cpu)-1]
			c.JSON(200, makeData(latestCpu.cpuSpeed))
		} else {
			c.JSON(200, makeData("-1"))
		}
	})
	r.GET("/volts", func(c *gin.Context) {
		c.JSON(200, makeData("0.850v"))
	})
	r.GET("/throttled", func(c *gin.Context) {
		c.JSON(200, makeData("Yes"))
	})
	r.GET("/mem", func(c *gin.Context) {
		c.JSON(200, makeData("1000Mb"))
	})
	r.GET("/memFree", func(c *gin.Context) {
		c.JSON(200, makeData("2000Mb"))
	})
	r.GET("/memTotal", func(c *gin.Context) {
		c.JSON(200, makeData("4000Mb"))
	})
	r.GET("/memSwap", func(c *gin.Context) {
		c.JSON(200, makeData("99Mb"))
	})
	r.GET("/memSwapTotal", func(c *gin.Context) {
		c.JSON(200, makeData("99Mb"))
	})
	r.GET("/wifi", func(c *gin.Context) {
		c.JSON(200, makeData("600Mhz"))
	})

	log.Log("Started server on port", params.Port)
	r.Run(fmt.Sprintf(":%d", params.Port))
}
