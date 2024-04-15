package main

import (
	"fmt"
	"strconv"
)

func main() {
	name := "lukmin"
	var codes = [6]string{}
	for i := 0; i < 6; i++ {
		codes[i] = strconv.Itoa(int(name[i]))
	}

	strong := calculate_strong_number(codes)
	weak := calculate_weak_number(30)
	fmt.Println("strong =", strong)
	fmt.Println("weak =", weak)
	fmt.Println()

	fmt.Println(predicted_fib_rec_computation_time(46))
}

// silna: 509, słaba: 22
