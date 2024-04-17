package tools

import (
	"strconv"
)

func StrToIntOrDefault(s string, defaultValue int) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return res
}
