package handler

import (
	"fmt"
	"math"
	"sort"
)

func ExecuteRPCMethod(method string, params []interface{}) (interface{}, string, error) {
	switch method {
	case "floor":
		return floor(params)
	case "nroot":
		return nroot(params)
	case "reverse":
		return reverse(params)
	case "validAnagram":
		return validAnagram(params)
	case "sort":
		return sortStrings(params)
	default:
		return nil, "", fmt.Errorf("unknown method: %s", method)
	}
}


func floor(params []interface{}) (interface{}, string, error) {
	if len(params) != 1 {
		return nil, "", fmt.Errorf("invalid number of arguments")
	}

	x, ok := params[0].(float64)
	if !ok {
		return nil, "", fmt.Errorf("argument must be a float")
	}

	return math.Floor(x), "int", nil
}

func nroot(params []interface{}) (interface{}, string, error) {
	if len(params) != 2 {
		return nil, "", fmt.Errorf("invalid number of arguments")
	}

	n, ok1 := params[0].(float64)
	x, ok2 := params[1].(float64)
	if !ok1 || !ok2 || n <= 0 || x < 0 {
		return nil, "", fmt.Errorf("arguments must be positive numbers")
	}

	result := math.Pow(x, 1/n)
	return result, "float", nil
}

func reverse(params []interface{}) (interface{}, string, error) {
	if len(params) != 1 {
		return nil, "", fmt.Errorf("invalid number of arguments")
	}

	s, ok := params[0].(string)
	if !ok {
		return nil, "", fmt.Errorf("argument must be a string")
	}

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), "string", nil
}

func validAnagram(params []interface{}) (interface{}, string, error) {
	if len(params) != 2 {
		return nil, "", fmt.Errorf("invalid number of arguments")
	}

	str1, ok1 := params[0].(string)
	str2, ok2 := params[1].(string)
	if !ok1 || !ok2 {
		return nil, "", fmt.Errorf("arguments must be strings")
	}

	count := func(s string) map[rune]int {
		m := make(map[rune]int)
		for _, r := range s {
			m[r]++
		}
		return m
	}

	isAnagram := func(a, b map[rune]int) bool {
		if len(a) != len(b) {
			return false
		}
		for k, v := range a {
			if b[k] != v {
				return false
			}
		}
		return true
	}

	return isAnagram(count(str1), count(str2)), "bool", nil
}

func sortStrings(params []interface{}) (interface{}, string, error) {
    if len(params) == 0 {
        return nil, "", fmt.Errorf("params cannot be empty")
    }

    // `params` を文字列スライスに変換
    stringSlice := make([]string, len(params))
    for i, param := range params {
        str, ok := param.(string)
        if !ok {
            return nil, "", fmt.Errorf("all elements in params must be strings")
        }
        stringSlice[i] = str
    }

    // ソート処理
    sort.Strings(stringSlice)

    return stringSlice, "string[]", nil
}
