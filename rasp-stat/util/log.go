package util

import "fmt"

type Loggable interface {
	Log(s ...interface{})
}

type Log struct {
	Loggable
	Name, Version string
}

func (l *Log) Log(s ...interface{}) {
	fmt.Print("[", GetOrDefault(&l.Name, "")+":", GetOrDefault(&l.Version, ""), "] ")
	for _, log := range s {
		fmt.Print(log, " ")
	}
	fmt.Println()
}
