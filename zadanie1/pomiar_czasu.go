package main

import(
	"fmt"
	"time"
	"math"
	// "math/big"
)

func time_fib(n int, prev time.Duration) time.Duration {
	start := time.Now()
	fibonacci_rec(n)
	dur := time.Since(start)
	fmt.Println("Time Elapsed: ", dur)
	if prev != 0 {
		fmt.Println("Rate of change: ", float64(dur) / float64(prev))
	}
	return dur
}

func float_average(arr []float64) float64 {
	var result float64 = 0
	for i := 1; i < len(arr); i++ {
		result += arr[i]
	}
	result = result / float64(len(arr) - 1)
	return result
}

func measure_average_fib_rate(lower int, upper int) float64 {
	rates := make([]float64, upper-lower+1)
	var prev time.Duration = 1
	for i := 0; i <= upper-lower; i++ {
		start := time.Now()
		fibonacci_rec(lower + i)
		dur := time.Since(start)
		fmt.Println(lower + i, dur)
		rates[i] = float64(dur) / float64(prev)
		prev = dur
	}
	fmt.Println()
	avg := float_average(rates)
	fmt.Println("Average rate of elongation of execution time between recursive fibonacci executions with consecutive parameters from", lower, "to", upper, ":\n", avg)
	return avg
}

func predicted_fib_rec_computation_time(n int) time.Duration {
	if n <= 40 {
		return time_fib(n, 0)
	} else {
		t := time_fib(40, 0)
		avg := measure_average_fib_rate(10, 40)
		return time.Duration(float64(t) * math.Pow(avg, float64(n - 40)))
	}
}