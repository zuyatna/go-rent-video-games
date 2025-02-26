package utils

import "strconv"

func StringToInt(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return result
}
