package convert

import (
	"log"
	"strconv"
)

func ToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return 0
	}
	return val
}
