package util

import (
	"strconv"
	"strings"
)

func GetOrDefault[T interface{}](val *T, defaultIfNil T) T {
	if val == nil {
		return defaultIfNil
	}
	return *val
}

func SplitEqual(s string) (key, val string) {
	// kv := make(map[string]string)
	pair := strings.SplitN(s, "=", 2)
	if len(pair) == 1 {
		key = pair[0]
		val = ""
	} else {
		key = GetOrDefault(&pair[0], "")
		val = GetOrDefault(&pair[1], "")
	}
	// fmt.Println(pair)
	// for i := 0; i < len(pair)-1; i++ {
	// fmt.Println(pair[i], pair[i+1])
	// kv[pair[i]] = pair[i+1]
	// }
	return key, val
}

func ToFloat(s string) float64 {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1
	}
	return res
}
