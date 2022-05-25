package util

func GetOrDefault[T interface{}](val *T, defaultIfNil T) T {
	if val == nil {
		return defaultIfNil
	}
	return *val
}
