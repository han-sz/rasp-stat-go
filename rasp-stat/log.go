package main

import "fmt"

type Loggable interface {
	Log(s ...interface{})
}

type Log struct {
	Loggable
}

func (l Log) Log(s ...interface{}) {
	fmt.Print("[rasp-stat] ")
	for _, log := range s {
		fmt.Print(log, " ")
		// switch log.(type) {
		// case string:
		// 	fmt.Print(log, " ")
		// case bool:
		// 	if (bool(log)) {

		// 	}
		// default:
		// 	fmt.Print("? ")
		// }
	}
	fmt.Println()
}
