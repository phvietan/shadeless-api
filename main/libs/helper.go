package libs

import (
	"fmt"
	"math/rand"
)

func ArrayInterfaceToArrayString(arr []interface{}) []string {
	result := make([]string, len(arr))
	for i, v := range arr {
		result[i] = fmt.Sprint(v)
	}
	return result
}

var (
	HexChars = []rune("0123456789abcdef")
)

func RandomString(length int) string {
	res := make([]rune, length)
	for i := 0; i < length; i++ {
		res[i] = HexChars[rand.Intn(len(HexChars))]
	}
	return string(res)
}
