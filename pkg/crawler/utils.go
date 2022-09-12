package crawler

import (
	"strconv"
)

func increaseNum(s string) (string, error) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return "", err
	}

	d++
	return strconv.Itoa(d), nil
}
