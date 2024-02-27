package utils

import (
	"fmt"
	"regexp"
)

func Trimname(name string) string {

	if len(name) > 15 {
		var finalname string
		var name1 string
		if havespace(name) {
			name1 = firstword(name)
			if len(name1) > 15 {
				fmt.Println("1")
				finalname = name1[:15]
				return finalname

			} else {
				fmt.Println("2")
				finalname = name1
				return finalname
			}
		} else {
			fmt.Println("3")
			finalname = name[:15]
			return finalname
		}
	}
	return name
}

func IsAlphaSpace(input string) bool {
	pattern := `^[a-zA-Z0-9\s]{2,15}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

func havespace(str string) bool {
	for i := range str {
		if str[i] == ' ' {
			return true
		}
	}
	return false
}

func firstword(str string) string {
	res_str := ""
	for i := range str {
		if str[i] != ' ' {
			res_str += string(str[i])
		} else {
			break
		}
	}
	return res_str
}
