package main

import "math/big"

func factorial_string(n int) string {
	result := big.NewInt(1)
	var i int
	for i = 1; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result.Text(10)
}

func contains_substring(str string, sub string) bool {
	required := len(sub)
	current := 0
	for i := 0; i < len(str); i++ {
		if str[i] == sub[current] {
			current += 1
		} else {
			current = 0
		}
		if current == required {
			return true
		}
	}
	return false
}

func calculate_strong_number(codes [6]string) int {
	strong_number := 0
	breaker := false
	for i := 0; i >= 0; i++ {
		breaker = true
		for code := 0; code < len(codes); code++ {
			if !contains_substring(factorial_string(i), codes[code]) {
				breaker = false
			}
		}
		if breaker {
			strong_number = i
			break
		}
	}
	return strong_number
}
