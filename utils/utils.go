package utils

import (
	"strconv"
	"strings"
)

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

func BreakStringIntoLines(s string, lineLength int) string {
	words := strings.Fields(s)
	var lines []string
	currentLine := ""

	for _, word := range words {
		if len(word) > lineLength {
			// Handle long words by splitting them
			for len(word) > lineLength {
				lines = append(lines, word[:lineLength])
				word = word[lineLength:]
			}
			if len(word) > 0 {
				if currentLine != "" {
					lines = append(lines, currentLine)
				}
				currentLine = word
			}
		} else if len(currentLine)+len(word)+1 <= lineLength {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
