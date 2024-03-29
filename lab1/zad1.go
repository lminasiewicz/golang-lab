package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)


func in(item int, slice []int) bool {
	for i:=0; i<len(slice); i++ {
		if slice[i] == item {
			return true
		}
	}
	return false
}


func extended_monty_hall(stay_with_choice bool, n int) bool {
	true_index := rand.IntN(n)
	choices := make([]bool, n)
	choices[true_index] = true
	choice := rand.IntN(n)

	remainer := 0
	if choice == true_index {
		candidate := 0
		for true {
			candidate = rand.IntN(n)
			if !(candidate == choice) {
				break
			}
		}
		remainer = candidate
	} else {
		remainer = true_index
	}

	if !stay_with_choice {
		choice = remainer
	}

	return choice == true_index
}


func monty_hall(stay_with_choice bool) bool {
	true_index := rand.IntN(3)
	choices := [3]bool{}
	choices[true_index] = true
	choice := rand.IntN(3)


	throwaway := rand.IntN(3)
	for true {
		throwaway = rand.IntN(3)
		if throwaway != choice && throwaway != true_index {
			break
		}
	}

	if !stay_with_choice {
		for true {
			new_choice := rand.IntN(3)
			if new_choice != choice && new_choice != throwaway {
				choice = new_choice
				break
			}
		}
	}
	return choices[choice]
}


func analyse_extended_monty_hall(rounds int, stay_with_choice bool, n int) {
	wins := 0
	for i := 0; i < rounds; i++ {
		game := extended_monty_hall(stay_with_choice, n)
		if game {
			wins += 1
		}
	}
	fmt.Println("Attempts:", rounds, "| Successes:", wins, "| Success rate:", math.Round(float64(wins)/float64(rounds)*1000)/10)
}


func analyse_monty_hall(rounds int, stay_with_choice bool) {
	wins := 0
	for i := 0; i < rounds; i++ {
		game := monty_hall(stay_with_choice)
		if game {
			wins += 1
		}
	}
	fmt.Println("Attempts:", rounds, "| Successes:", wins, "| Success rate:", math.Round(float64(wins)/float64(rounds)*1000)/10)
}



func main() {
	analyse_monty_hall(500, true)
	analyse_monty_hall(500, false)
	analyse_extended_monty_hall(500, true, 10)
	analyse_extended_monty_hall(500, false, 10)
	analyse_extended_monty_hall(500, true, 100)
	analyse_extended_monty_hall(500, false, 100)
}
