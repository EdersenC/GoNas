package helper

import "strings"

func Contains(list []string, val string) bool {
	for _, v := range list {
		if strings.Contains(val, v) {
			return true
		}
	}
	return false
}
