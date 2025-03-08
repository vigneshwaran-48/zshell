package utils

import "strconv"

func IsNumber(value string) bool {
	if _, err := strconv.Atoi(value); err != nil {
		return false
	}
	return true
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	ret = []T{}
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func Map[T, U any](data []T, f func(T) U) []U {
	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func Contains[T comparable](ss []T, elementToFind T) bool {
	for _, s := range ss {
		if s == elementToFind {
			return true
		}
	}
	return false
}
