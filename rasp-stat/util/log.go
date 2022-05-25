package util

import "fmt"

type Loggable interface {
	Log(s ...interface{})
}

type Log struct {
	Loggable
	name, version string
}

func (l *Log) Log(s ...interface{}) {
	fmt.Print("[", GetOrDefault(&l.name, "")+":", GetOrDefault(&l.version, ""), "] ")
	for _, log := range s {
		fmt.Print(log, " ")
	}
	fmt.Println()
}
