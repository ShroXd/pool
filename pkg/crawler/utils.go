package crawler

import (
	"github.com/djimenez/iconv-go"
	"log"
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

func convertChinese(s string) string {
	// TODO: replace iconv since it causes build fail on Windows: https://github.com/djimenez/iconv-go/issues/42
	output, err := iconv.ConvertString(s, "GB2312", "utf-8")
	if err != nil {
		log.Println(err)
	}

	return output
}
