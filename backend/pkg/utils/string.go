package utils

func StrInArray(str string, arr []string, equal func(src, dst string) bool) bool {
	if equal == nil {
		return false
	}
	for _, v := range arr {
		if equal(str, v) {
			return true
		}
	}
	return false
}
