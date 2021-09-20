package libs

import (
	"fmt"
	"math/rand"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ArrayInterfaceToMapString(arr []interface{}) map[string]string {
	result := make(map[string]string)
	for _, v := range arr {
		asMongoD, ok := v.(primitive.D)
		if !ok {
			return make(map[string]string)
		}
		m := asMongoD.Map()
		for k, v := range m {
			if result[k], ok = v.(string); !ok {
				return make(map[string]string)
			}
		}
	}
	return result
}

func Max(a float64, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func Min(a float64, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

func IsSliceIncludesString(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func IsStringContains(s string, list []string) int {
	cnt := 0
	for _, v := range list {
		if strings.Contains(s, v) {
			cnt += 1
		}
	}
	return cnt
}

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
