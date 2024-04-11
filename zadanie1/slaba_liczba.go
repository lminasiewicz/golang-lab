package main

import (
	"math"
)

func fibonacci(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	prev := 0
	curr := 1
	temp := 0
	for i := 2; i <= n; i++ {
		temp = curr
		curr = curr + prev
		prev = temp
	}
	return curr
}

func fibonacci_rec(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacci_rec(n-1) + fibonacci_rec(n-2)
	}
}

func fibonacci_executions(n int) []int {
	results := make([]int, n)
	results[0] = fibonacci(n - 1)
	for i := 1; i < n; i++ {
		results[i] = fibonacci(n - i + 1)
	}
	return results
}

func calculate_weak_number(strong int) int {
	var executions []int = fibonacci_executions(strong)
	distance := math.Abs(float64(executions[0]) - float64(strong))
	curr := executions[0]
	for i := 0; i < len(executions); i++ {
		if math.Abs(float64(executions[i])-float64(strong)) < distance {
			distance = math.Abs(float64(float64(executions[i]) - float64(strong)))
			curr = i
		}
	}
	return curr
}
