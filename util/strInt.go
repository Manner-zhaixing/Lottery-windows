package util

import "strconv"

func StrToInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}
