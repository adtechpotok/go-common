package help

import "unicode"

//Convert the first letter to lower case
func LowerCaseFirstLetter(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}