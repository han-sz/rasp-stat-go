package main

import (
	"os/exec"
	"rasp-stat/rasp-stat/util"
	"strings"
	"sync"
	"time"
)

const BIN = "/opt/vc/bin/"
const WIFI_CMD = "iwconfig"
const TEMP_CMD = "vcgencmd measure_temp"
const GPUSPEED_CMD = "vcgencmd measure_clock core"
const CPUSPEED_CMD = "vcgencmd measure_clock arm"
const CPUTHROTTLED_CMD = "vcgencmd get_throttled"
const CPUVOLTS_CMD = "vcgencmd measure_volts"
const MEMORY_CMD = "free -m | awk 'NR==2{print $7,$2} NR==3{print $2,$3}'"

type memory struct {
	unit         string
	memFree      string `json:"memFree"`
	memTotal     string `json:"memTotal"`
	memSwap      string `json:"memSwap"`
	memSwapTotal string `json:"memSwapTotal"`
}

type cpu struct {
	unit      string
	cpuSpeed  float32 `json:"cpu"`
	throttled bool    `json:"throttled"`
}

type gpu struct {
	unit     string
	gpuSpeed float32 `json:"gpu"`
}

type temperature struct {
	unit        string
	temperature string `json:"temp"`
}
type InstantStatServiceInterface interface {
	FetchAndCacheStats()
	FetchCurrentCpuSpeed()
	FetchCurrentGpuSpeed()
	FetchCurrentTemperature()

	RawCpuDataPoints()
	RawGpuDataPoints()
	RawTemperatureDataPoints()
}

type InstantStatService struct {
	InstantStatServiceInterface

	FetchIntervalSeconds   uint16
	InMemDataPointsPerStat uint16

	ReadWriteLock *sync.Mutex

	cpu  []cpu
	gpu  []gpu
	temp []temperature
}

func addDataPoint[T interface{}](t *[]T, dp T, max int) {
	if len(*t)+1 >= max {
		*t = (*t)[1:]
	}
	*t = append(*t, dp)
}

func NewInstantStatService(fetchIntervalSeconds uint16, inMemDataPointsPerStat uint16) InstantStatService {
	var iss InstantStatService = InstantStatService{
		cpu:                    make([]cpu, inMemDataPointsPerStat),
		gpu:                    make([]gpu, inMemDataPointsPerStat),
		temp:                   make([]temperature, inMemDataPointsPerStat),
		FetchIntervalSeconds:   fetchIntervalSeconds,
		InMemDataPointsPerStat: inMemDataPointsPerStat,
	}
	return iss
}

func (iss *InstantStatService) FetchCurrentCpuSpeed() {
	res, err := commandOutput(CPUSPEED_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)
	cpuSpeed := util.ToFloat(v) / 1000.0 / 1000.0 / 1000.0

	(*iss.ReadWriteLock).Lock()
	addDataPoint(&iss.cpu, cpu{cpuSpeed: cpuSpeed, unit: "GHz"}, int(iss.InMemDataPointsPerStat))
	if DEBUG {
		log.Log(iss.cpu)
	}
	(*iss.ReadWriteLock).Unlock()
}

func (iss *InstantStatService) FetchCurrentGpuSpeed() {
	res, err := commandOutput(GPUSPEED_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)
	gpuSpeed := util.ToFloat(v) / 1000.0 / 1000.0 / 1000.0

	(*iss.ReadWriteLock).Lock()
	addDataPoint(&iss.gpu, gpu{gpuSpeed: gpuSpeed, unit: "GHz"}, int(iss.InMemDataPointsPerStat))
	if DEBUG {
		log.Log(iss.gpu)
	}
	(*iss.ReadWriteLock).Unlock()
}

func (iss *InstantStatService) FetchCurrentTemperature() {
	res, err := commandOutput(TEMP_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)

	(*iss.ReadWriteLock).Lock()
	addDataPoint(&iss.temp, temperature{temperature: v}, int(iss.InMemDataPointsPerStat))
	if DEBUG {
		log.Log(iss.gpu)
	}
	(*iss.ReadWriteLock).Unlock()
}

func (iss *InstantStatService) RawCpuDataPoints() []float32 {
	return util.MapToRawData(&iss.cpu, func(val cpu) float32 {
		return val.cpuSpeed
	})
}

func (iss *InstantStatService) RawGpuDataPoints() []float32 {
	return util.MapToRawData(&iss.gpu, func(val gpu) float32 {
		return val.gpuSpeed
	})
}

func (iss *InstantStatService) RawTemperatureDataPoints() []string {
	return util.MapToRawData(&iss.temp, func(val temperature) string {
		return val.temperature
	})
}

func (iss *InstantStatService) FetchAndCacheStats() {
	go func() {
		for {
			iss.FetchCurrentCpuSpeed()
			iss.FetchCurrentGpuSpeed()
			iss.FetchCurrentTemperature()

			time.Sleep(time.Duration(iss.FetchIntervalSeconds) * time.Second)
		}
	}()
}

func commandOutput(command string) (string, error) {
	split := strings.Fields(command)
	output, err := exec.Command(split[0], split[1:]...).Output()
	if err != nil {
		if DEBUG {
			log.Log("Error running command:", err.Error())
		}
		return "", err
	}
	// if DEBUG {
	// 	log.Log("Running:", split, string(output))
	// }
	return string(output), nil
}
