package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

const (
	VERSION string = "1.0.0"
	DEBUG   bool   = false
)

var port string = "4322"
var log Loggable = Log{}

type DataPoint struct {
	Value string `json:"data"`
}

func makeData(value string) DataPoint {
	return DataPoint{Value: value}
}

func env() {
	portEnv, found := os.LookupEnv("PORT")
	if found {
		port = portEnv
		log.Log("Using port from environment", portEnv)
	}
}

const BIN = "/opt/vc/bin/"
const WIFI_CMD = "iwconfig"
const TEMP_CMD = "vcgencmd measure_temp"
const GPUSPEED_CMD = "vcgencmd measure_clock core"
const CPUSPEED_CMD = "vcgencmd measure_clock arm"
const CPUTHROTTLED_CMD = "vcgencmd get_throttled"
const CPUVOLTS_CMD = "vcgencmd measure_volts"
const MEMORY_CMD = "free -m | awk 'NR==2{print $7,$2} NR==3{print $2,$3}'"

func main() {
	env()
	if !DEBUG {
		gin.SetMode(gin.ReleaseMode)
		log.Log("Running version:", VERSION)
	} else {
		log.Log("Running in debug mode, version:", VERSION)
	}
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		if res, err := commandOutput("/opt/vc/bin/vcgencmd measure_clock core -a"); err == nil {
			log.Log("CPU", res)
			c.JSON(200, makeData(res))
		} else {
		}
		if res, err := commandOutput("/opt/vc/bin/vcgencmd measure_temp"); err == nil {
			log.Log("Temp", res)
			c.JSON(200, makeData(res))
		} else {
			c.JSON(200, makeData("Unknown"))
		}
	})
	r.GET("/temp", func(c *gin.Context) {
		c.JSON(200, makeData("33c"))
	})
	r.GET("/gpu", func(c *gin.Context) {
		c.JSON(200, makeData("600Mhz"))
	})
	r.GET("/cpu", func(c *gin.Context) {
		c.JSON(200, makeData("1.00Ghz"))
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
	r.Run(":" + port)
}
