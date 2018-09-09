package util

import (
	"fmt"
	"strings"
)

//把map转化为string  "a:b|"
func Map2String(msg map[string]string) string {
	if len(msg) == 0 {
		return ""
	}
	data := ""
	for i, v := range msg {
		data += fmt.Sprintf("%v:%v|", i, v)
	}
	return data[:len(data)-1]
}

//把string转化为map
func String2Map(msg string) map[string]string {
	data := make(map[string]string)
	tempData := strings.Split(msg, "|")

	if len(tempData) < 1 {
		return data
	}

	for _, v := range tempData {
		tempValue := strings.Split(v, ":")
		if len(tempValue) <= 1 {
			continue
		}
		data[tempValue[0]] = tempValue[1]
	}

	return data
}
