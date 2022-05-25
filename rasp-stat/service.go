package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type memory struct {
	memFree      string `json:"memFree"`
	memTotal     string `json:"memTotal"`
	memSwap      string `json:"memSwap"`
	memSwapTotal string `json:"memSwapTotal"`
}

type cpucore struct {
	speed string `json:"speed"`
}

type cpu struct {
	cpuSpeed  string `json:"cpu"`
	throttled bool   `json:"throttled"`
}

type gpu struct {
	gpuSpeed string `json:"gpu"`
}

type temperature struct {
	temperatureCelsius string `json:"temp"`
}

func commandOutput(command string) (string, error) {
	split := strings.Fields(command)
	// fmt.Println("Running:", split[0], "Rest", split...)
	output, err := exec.Command(split[0], split[1:]...).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(output), err
}
