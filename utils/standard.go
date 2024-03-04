package utils

// 规范字符,超过某长度,用...代替
func ReplaceLength(str string, length int) string {
	nameRune := []rune(str)
	if len(str) > length {
		return string(nameRune[:length]) + "..."
	}
	return string(nameRune)
}
