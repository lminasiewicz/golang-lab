package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func monty_hall(stay_with_choice bool) bool {
	true_index := rand.IntN(3)
	choices := [3]bool{}
	switch true_index {
	case 0:
		choices[0] = true
	case 1:
		choices[1] = true
	case 2:
		choices[2] = true
	}
	choice := rand.IntN(3)

	stop := false
	throwaway := rand.IntN(3)
	for stop == false {
		throwaway := rand.IntN(3)
		if throwaway != choice && choices[throwaway] != true {
			stop = true
		}
	}

	if !stay_with_choice {
		stop := false
		for stop == false {
			new_choice := rand.IntN(3)
			if new_choice != choice && new_choice != throwaway {
				choice = new_choice
				stop = true
			}
		}
	}

	return choices[choice]
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
}
