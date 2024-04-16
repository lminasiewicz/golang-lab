package main

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

func time_fib(n int, prev time.Duration, show_duration bool, rec bool) time.Duration {
	start := time.Now()
	if rec {
		fibonacci_rec(n)
	} else {
		fibonacci(n)
	}
	dur := time.Since(start)
	if show_duration {fmt.Println("Time Elapsed: ", dur)}
	if prev != 0 {
		fmt.Println("Rate of change: ", float64(dur)/float64(prev))
	}
	return dur
}

func float_average(arr []float64) float64 {
	var result float64 = 0
	for i := 1; i < len(arr); i++ {
		result += arr[i]
	}
	result = result / float64(len(arr)-1)
	return result
}

func measure_average_fib_rate(lower int, upper int) float64 {
	rates := make([]float64, upper-lower+1)
	var prev time.Duration = 1
	for i := 0; i <= upper-lower; i++ {
		start := time.Now()
		fibonacci_rec(lower + i)
		dur := time.Since(start)
		rates[i] = float64(dur) / float64(prev)
		prev = dur
	}
	avg := float_average(rates)
	fmt.Println("Average rate of elongation of execution time between recursive fibonacci executions with consecutive parameters from", lower, "to", upper, ":\n", avg)
	return avg
} // zwraca mniej więcej złotą proporcję, więc zakładam że chodzi o to i korzystam ze stałej zamiast tej funkcji.

func make_bigfloat_duration_readable(f *big.Float) string {
	result := big.NewFloat(0)
	if f.Cmp(big.NewFloat(1000000)) >= 0 {
		result.Quo(f, big.NewFloat(1000000))
	} else {
		return f.Text('f', 2) + " nanoseconds"
	}

	if result.Cmp(big.NewFloat(0).Quo(big.NewFloat(31557600000), big.NewFloat(12))) >= 0 {
		result.Quo(result, big.NewFloat(31557600000))
		return "approximately " + result.Text('f', 2) + " years"
	}
	if result.Cmp(big.NewFloat(86400000)) >= 0 {
		result.Quo(result, big.NewFloat(86400000))
		return "approximately " + result.Text('f', 2) + " days"
	}
	if result.Cmp(big.NewFloat(3600000)) >= 0 {
		result.Quo(result, big.NewFloat(3600000))
		return "approximately " + result.Text('f', 2) + " hours"
	}
	if result.Cmp(big.NewFloat(60000)) >= 0 {
		result.Quo(result, big.NewFloat(60000))
		return "approximately " + result.Text('f', 2) + " minutes"
	}
	if result.Cmp(big.NewFloat(1000)) >= 0 {
		result.Quo(result, big.NewFloat(1000))
		return "approximately " + result.Text('f', 2) + " seconds"
	}
	return result.Text('f', 2) + " milliseconds"
}

func predicted_fib_computation_time(n int, rec bool) string {
	if rec {
		if n <= 40 {
			return fmt.Sprint("recursive fib(", n, ") computation time: "+make_bigfloat_duration_readable(big.NewFloat(float64(time_fib(n, 0, false, true)))))
		} else {
			t := time_fib(40, 0, false, true)
			// avg := measure_average_fib_rate(27, 40)
			// fmt.Println("approximate rate of change:", avg)
			avg := (1 + math.Sqrt(5)) / 2
			return fmt.Sprint("recursive fib(", n, ") computation time: "+make_bigfloat_duration_readable(big.NewFloat(float64(t)*math.Pow(avg, float64(n-40)))))
		}
	} else {
		return fmt.Sprint("iterative fib(", n, ") computation time: "+make_bigfloat_duration_readable(big.NewFloat(float64(time_fib(n, 0, false, false)))))
	}
}