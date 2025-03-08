package utils

import "strconv"

func IsNumber(value string) bool {
	if _, err := strconv.Atoi(value); err != nil {
		return false
	}
	return true
}
