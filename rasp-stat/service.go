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

// TODO
type memory struct {
	unit         string
	memFree      string
	memTotal     string
	memSwap      string
	memSwapTotal string
}

type cpu struct {
	unit     string
	cpuSpeed float32
}

type gpu struct {
	unit     string
	gpuSpeed float32
}

type temperature struct {
	unit        string
	temperature float32
}

type throttled struct {
	value       string
	rawValue    string
	isThrottled bool
}

type voltage struct {
	unit  string
	volts float32
}

type InstantStatServiceInterface interface {
	FetchAndCacheStats()

	FetchCurrentGpuSpeed()
	FetchCurrentCpuSpeed()
	FetchCurrentCpuVoltage()
	FetchCurrentCpuThrottled()
	FetchCurrentTemperature()

	RawCpuDataPoints()
	RawGpuDataPoints()
	RawCpuVoltageDataPoints()
	RawCpuThrottledDataPoints()
	RawTemperatureDataPoints()
}

type InstantStatService struct {
	InstantStatServiceInterface

	FetchIntervalSeconds   uint16
	InMemDataPointsPerStat uint16

	ReadWriteLock *sync.Mutex

	cpu       []cpu
	gpu       []gpu
	temp      []temperature
	volts     []voltage
	throttled []throttled
}

// Helper functions

func addDataPoint[T interface{}](t *[]T, dp T, max int) {
	if len(*t)+1 > max {
		*t = (*t)[1:]
	}
	*t = append(*t, dp)
}

func (iss *InstantStatService) acquireLock(critical func()) {
	if DEBUG {
		log.Log("Acquiring lock")
	}
	iss.ReadWriteLock.Lock()
	critical()
	iss.ReadWriteLock.Unlock()
	if DEBUG {
		log.Log("Released lock")
	}
}

// InstantStatService

func NewInstantStatService(fetchIntervalSeconds uint16, inMemDataPointsPerStat uint16) InstantStatService {
	var iss InstantStatService = InstantStatService{
		// If the data points are not initialised beforehand, there will be a silent deadlock in the
		// main goroutine when there is a request from the server and the lock is being acquired
		FetchIntervalSeconds:   fetchIntervalSeconds,
		InMemDataPointsPerStat: inMemDataPointsPerStat,

		cpu:       make([]cpu, 1, inMemDataPointsPerStat),
		gpu:       make([]gpu, 1, inMemDataPointsPerStat),
		temp:      make([]temperature, 1, inMemDataPointsPerStat),
		volts:     make([]voltage, 1, inMemDataPointsPerStat),
		throttled: make([]throttled, 1, inMemDataPointsPerStat),
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

	iss.acquireLock(func() {
		addDataPoint(&iss.cpu, cpu{cpuSpeed: cpuSpeed, unit: "GHz"}, int(iss.InMemDataPointsPerStat))
	})
}

func (iss *InstantStatService) FetchCurrentCpuVoltage() {
	res, err := commandOutput(CPUVOLTS_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)

	volts := util.ToFloat(strings.ReplaceAll(v, "V", ""))

	iss.acquireLock(func() {
		addDataPoint(&iss.volts, voltage{volts: volts, unit: "volts"}, int(iss.InMemDataPointsPerStat))
	})
}

func (iss *InstantStatService) FetchCurrentCpuThrottled() {
	res, err := commandOutput(CPUTHROTTLED_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)

	iss.acquireLock(func() {
		addDataPoint(
			&iss.throttled,
			throttled{rawValue: v, value: util.Is(v == "0x0", "No", "Yes"), isThrottled: v != "0x0"},
			int(iss.InMemDataPointsPerStat),
		)
	})
}

func (iss *InstantStatService) FetchCurrentGpuSpeed() {
	res, err := commandOutput(GPUSPEED_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)
	gpuSpeed := util.ToFloat(v) / 1000.0 / 1000.0 / 1000.0

	iss.acquireLock(func() {
		addDataPoint(&iss.gpu, gpu{gpuSpeed: gpuSpeed, unit: "GHz"}, int(iss.InMemDataPointsPerStat))
	})
}

func (iss *InstantStatService) FetchCurrentTemperature() {
	res, err := commandOutput(TEMP_CMD)
	if err != nil {
		return
	}
	_, v := util.SplitEqual(res)
	formattedTemp := util.ToFloat(strings.ReplaceAll(v, "'C", ""))

	iss.acquireLock(func() {
		addDataPoint(&iss.temp, temperature{temperature: formattedTemp, unit: "C"}, int(iss.InMemDataPointsPerStat))
	})
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

func (iss *InstantStatService) RawTemperatureDataPoints() []float32 {
	return util.MapToRawData(&iss.temp, func(val temperature) float32 {
		return val.temperature
	})
}

func (iss *InstantStatService) RawCpuVoltageDataPoints() []float32 {
	return util.MapToRawData(&iss.volts, func(val voltage) float32 {
		return val.volts
	})
}

func (iss *InstantStatService) RawCpuThrottledDataPoints() []int {
	return util.MapToRawData(&iss.throttled, func(val throttled) int {
		return util.Is(val.isThrottled, 1, 0)
	})
}

func (iss *InstantStatService) FetchAndCacheStats() {
	go func() {
		for {
			iss.FetchCurrentGpuSpeed()
			iss.FetchCurrentCpuSpeed()
			iss.FetchCurrentCpuVoltage()
			iss.FetchCurrentCpuThrottled()
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
	if DEBUG {
		log.Log("Running:", split, string(output))
	}
	return string(output), nil
}
