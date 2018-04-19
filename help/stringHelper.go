package help

import "unicode"

//Переводит первую букву в строке в нижний регистр
func LowerCaseFirstLetter(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}