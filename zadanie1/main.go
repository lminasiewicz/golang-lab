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
	weak := calculate_weak_number(strong)
	fmt.Println(strong)
	fmt.Println(weak)
}

// silna: 509, słaba: 495
