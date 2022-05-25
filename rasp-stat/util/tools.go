package util

import (
	"strconv"
	"strings"
	"sync"
)

func MapToRawData[T interface{}, U int | float32 | string](o *[]T, mapper func(T) U) []U {
	rawRes := make([]U, len(*o))
	for i, d := range *o {
		rawRes[i] = mapper(d)
	}
	return rawRes
}

func GetOrDefault[T interface{}](val *T, defaultIfNil T) T {
	if val == nil {
		return defaultIfNil
	}
	return *val
}

func GetOrDefaultMap[K comparable, V interface{}, T map[K]V](m *T, k K, defaultIfNil V) V {
	val, found := (*m)[k]
	if !found {
		return defaultIfNil
	}
	return val
}

func GetOrDefaultSafeMap[K comparable, V interface{}](m *sync.Map, k K, defaultIfNil V) V {
	val, found := (*m).Load(k)
	if !found {
		return defaultIfNil
	}
	return val.(V)
}

func Is[T interface{}](c bool, t T, f T) T {
	if c {
		return t
	}
	return f
}

func SplitEqual(s string) (key, val string) {
	pair := strings.SplitN(s, "=", 2)
	if len(pair) == 1 {
		key = pair[0]
		val = ""
	} else {
		key = strings.Trim(GetOrDefault(&pair[0], ""), "\n\t ")
		val = strings.Trim(GetOrDefault(&pair[1], ""), "\n\t ")
	}
	return key, val
}

func ToInt(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return res
}

func ToFloat(s string) float32 {
	res, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return -1
	}
	return float32(res)
}
