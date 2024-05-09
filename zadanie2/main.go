package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Field struct {
	tree   bool
	age    int
	burned bool
}

// pseudo-enum na kierunki wiatru
type Direction int

const (
	None Direction = iota
	North
	East
	South
	West
)

const MAX_AGE = 600
const BURN_RESISTANCE_THRESHOLD = 400
const MAX_BURN_RESISTANCE = 0.5

func attempt_burn(age int) bool {
	// Function, which for a given age of a tree returns a boolean value stating whether or not the attempt to incinerate the tree was successful.
	// Pseudo-random chance, influenced by a linear function that depends on constants BURN_RESISTANCE_THRESHOLD and MAX_BURN_RESISTANCE.
	if age > BURN_RESISTANCE_THRESHOLD {
		return true
	} else {
		var chance float32 = MAX_BURN_RESISTANCE + (float32(age) * float32((1-MAX_BURN_RESISTANCE)/BURN_RESISTANCE_THRESHOLD))
		if rand.Float32() < chance {
			return true
		}
		return false
	}
}

func initialize_forest(width int, length int, rate float32) [][]Field {
	// Initialize a {width} X {length} field of structs (type Field), the parameter {rate} determines what percentage of these fields will have trees on them.
	forest := make([][]Field, width)
	for i := 0; i < width; i++ {
		forest[i] = make([]Field, length)
	}
	for i := 0; i < width; i++ {
		for j := 0; j < length; j++ {
			if rand.Float32() > rate {
				var f Field
				f.tree = false
				f.age = 0
				f.burned = false
				forest[i][j] = f
			} else {
				var f Field
				f.tree = true
				f.age = rand.IntN(600)
				f.burned = false
				forest[i][j] = f
			}
		}
	}
	return forest
}

func get_coords_around(forest [][]Field, coords [2]int, wind Direction) [][]int {
	// Helper function to get the coords of all the available fields around a given field.
	// Pays attention to bounds of the map.
	result := make([][]int, 0, 11)
	i_start, j_start := -1, -1
	i_end, j_end := 2, 2

	switch wind {
	case North:
		i_start -= 1
		break
	case East:
		j_end += 1
		break
	case South:
		i_end += 1
		break
	case West:
		j_start -= 1
		break
	}
	for i := i_start; i < i_end; i++ {
		for j := j_start; j < j_end; j++ {
			if coords[0]+i >= 0 && coords[0]+i < len(forest) && coords[1]+j >= 0 && coords[1]+j < len(forest[0]) {
				if !(i == 0 && j == 0) {
					result = append(result, []int{coords[0] + i, coords[1] + j})
				}
			}
		}
	}
	return result
}

func burn(forest [][]Field, current [2]int, wind Direction) {
	// Recursive function to attempt to burn a given {current} tree, and, if successful, execute itself on surrounding fields.
	available_spaces := get_coords_around(forest, current, wind)
	next_burners := make([][2]int, 0, 11)
	for _, space := range available_spaces {
		fld := forest[space[0]][space[1]]
		if fld.tree && !fld.burned {
			if attempt_burn(fld.age) {
				forest[space[0]][space[1]].burned = true
				space_arr := [2]int{space[0], space[1]}
				next_burners = append(next_burners, space_arr)
			}
		}
	}
	for _, burner := range next_burners {
		burn(forest, burner, wind)
	}
}

func lightning_strike(forest [][]Field, lightning [2]int, wind Direction) {
	// Initial lightning strike that starts the wildfire. One-time non-recursive burn().
	fld := forest[lightning[0]][lightning[1]]
	if fld.tree && !fld.burned {
		forest[lightning[0]][lightning[1]].burned = true
		burn(forest, lightning, wind)
	}
}

func lightning_strike_until_successful(forest [][]Field, wind Direction) [2]int {
	// A separate mode for the simulate_once() function, written in cases where the random lightning strike has a low chance to hit a tree.
	// Ensures that the lightning will strike SOME tree, so you can see the results of the simulation without retrying a lot.
	lightning := [2]int{rand.IntN(len(forest)), rand.IntN(len(forest[0]))}
	for {
		lightning = [2]int{rand.IntN(len(forest)), rand.IntN(len(forest[0]))}
		fld := forest[lightning[0]][lightning[1]]
		fmt.Println(lightning, fld)
		if fld.tree && !fld.burned {
			lightning_strike(forest, lightning, wind)
			break
		}
	}
	return lightning
}

func print_forest(forest [][]Field, lightning [2]int) {
	// Prints the forest in a grid with symbols symbolizing different aspects of the simulation.
	// Empty circles are standing trees. Full circles are burned trees.
	// Star symbols mean that the lightning struck there. It has a tree and non-tree variant.
	for row := 0; row < len(forest); row++ {
		for elem := 0; elem < len(forest[0]); elem++ {
			fld := forest[row][elem]

			if !fld.tree {
				if row == lightning[0] && elem == lightning[1] {
					fmt.Print("★")
				} else {
					fmt.Print(" ")
				}
			} else if fld.tree && !fld.burned {
				fmt.Print("○︎")
			} else if fld.tree && fld.burned {
				if row == lightning[0] && elem == lightning[1] {
					fmt.Print("✪")
				} else {
					fmt.Print("●︎")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func print_forest_stats(forest [][]Field) {
	// For displaying statistics about the (singular) simulation after the lightning strike and potential subsequent wildfire.
	trees := 0
	burned := 0
	for row := 0; row < len(forest); row++ {
		for _, fld := range forest[row] {
			if fld.tree {
				trees += 1
				if fld.burned {
					burned += 1
				}
			}
		}
	}
	length := len(forest[0])
	width := len(forest)
	forestation := math.Round(10000*(float64(trees)/float64(length*width))) / 100
	burned_rate := math.Round(10000*(float64(burned)/float64(trees))) / 100
	fmt.Print("Forest ", length, "x", width, " (", length*width, " spaces)\n")
	fmt.Print("Forestation: ", forestation, "%\n")
	fmt.Print("Trees: ", trees, "  |  Burned: ", burned, "\n")
	fmt.Print(burned_rate, "% of the forest has been burned.\n")
	if burned == 0 {
		fmt.Println("Lightning Strike missed.")
	}
}

func get_quality_index(forest [][]Field) float64 {
	// A custom index calculated for a forest after its lightning struck. Since a "survival rate" index could easily be meaningless
	// due to a tiny forestation rate making it unlikely that the lightning strike will hit anything, this index measures how
	// many trees remain after the lightning strike and its consequences, weighted slightly with survival rate.
	trees := 0
	burned := 0
	for row := 0; row < len(forest); row++ {
		for _, fld := range forest[row] {
			if fld.tree {
				trees += 1
				if fld.burned {
					burned += 1
				}
			}
		}
	}
	var survival_rate float64 = float64(1) - (float64(burned) / float64(trees))
	// modifier is a linear modifier for the result, which weighs it depending on its survival rate.
	// for any given number of surviving trees, with survival rate of 0% the modifier is 0.8, and with survival rate of 100% the modifier is 1.2.
	var modifier float64 = float64(0.8) + (survival_rate * float64(0.4))
	return float64(trees) * survival_rate * modifier
}

func wind_int_to_name(wind Direction) string {
	switch wind {
	case None:
		return "None"
	case North:
		return "North"
	case East:
		return "East"
	case South: 
		return "South"
	case West:
		return "West"
	default:
		return "Error"
	}
}

func simulate_once(length int, width int, forestation_rate float32, wind Direction, lightning_is_accurate bool) {
	// Run a single simulation and display its end result and statistics. {lightning_is_accurate} parameter lets you guarantee that
	// the lightning will strike some tree, so you don't have to run it multiple times to get any result where something burned.
	forest := initialize_forest(width, length, forestation_rate)

	var lightning = [2]int{}
	if lightning_is_accurate {
		lightning = lightning_strike_until_successful(forest, wind)
	} else {
		lightning = [2]int{rand.IntN(width), rand.IntN(length)}
		lightning_strike(forest, lightning, wind)
	}
	print_forest(forest, lightning)
	print_forest_stats(forest)
}

func simulate_many(sample_size int, length int, width int, forestation_rate float32, wind Direction) float64 {
	// Run {sample_size} simulations, and instead of displaying statistics, simply return the average Forest Quality Index of these simulations.
	forest := initialize_forest(width, length, forestation_rate)

	var forest_quality_index float64 = 0
	for i := 0; i < sample_size; i++ {
		lightning := [2]int{rand.IntN(width), rand.IntN(length)}
		lightning_strike(forest, lightning, wind)
		forest_quality_index += get_quality_index(forest)
	}
	forest_quality_index = forest_quality_index / float64(sample_size)
	return forest_quality_index
}

func conduct_test(sample_size int, step float32, length int, width int, wind Direction) {
	// Run simulate_many() many times with a different {forestation_rate} parameter every time.
	// the {step} parameter dictates how many times simulate_many() will be run, and with which forestation rates.
	// {step} = 2 means simulate_many will run with {forestation_rate} of 0, then 0.02, then 0.04, ... , then 1.
	// {sample_size}, {length}, {width}, and {wind} are for the same things as in simulate_many().
	step = step / 100
	var best_score float64 = 0
	var best_rate float32 = 0
	var current_score float64 = 0
	var current_rate float32 = 0

	var rate float32 = 0
	for ; rate <= 1; rate += step {
		current_rate = rate
		current_score = simulate_many(sample_size, length, width, rate, wind)
		if current_score > best_score {
			best_rate = current_rate
			best_score = current_score
		}
	}
	if rate < 1 {
		current_rate = 1
		current_score = simulate_many(sample_size, length, width, 1, wind)
		if current_score > best_score {
			best_rate = current_rate
			best_score = current_score
		}
	}

	best_rate = float32(math.Round(float64(best_rate*10000)) / 100)
	best_score = math.Round(float64(best_score*100)) / 100
	fmt.Print("Conducted a series of simulations on the optimal forestation rate of a ", length, "x", width, " forest.\n")
	fmt.Print("For a step of ", step*100, "% in forestation rate from 0% to 100%, simulations were conducted with a sample size of ", sample_size, " per simulation.\n")
	if wind != None {
		fmt.Println("Additionally, wind is blowing to the", wind_int_to_name(wind), "in all of the simulations.")
	}
	fmt.Print("Under these conditions, the most optimal forestation rate is ", best_rate, "%, with a Forest Quality Index of ", best_score, ".\n")
}

func main() {
	// simulate_once(40, 40, 0.2, None, true)
	conduct_test(100000, 0.5, 80, 20, None)
	fmt.Println()
	conduct_test(100000, 0.5, 80, 20, North)
	fmt.Println()
	conduct_test(100000, 0.5, 80, 20, South)
	fmt.Println()
	conduct_test(100000, 0.5, 80, 20, East)
	fmt.Println()
	conduct_test(100000, 0.5, 80, 20, West)
	fmt.Println()
	conduct_test(100000, 0.5, 40, 40, None)
	fmt.Println()
	conduct_test(100000, 0.5, 40, 40, North)
	fmt.Println()
	conduct_test(100000, 0.5, 40, 40, South)
	fmt.Println()
	conduct_test(100000, 0.5, 40, 40, East)
	fmt.Println()
	conduct_test(100000, 0.5, 40, 40, West)
	fmt.Println()
	conduct_test(100000, 0.5, 20, 80, None)
	fmt.Println()
	conduct_test(100000, 0.5, 20, 80, North)
	fmt.Println()
	conduct_test(100000, 0.5, 20, 80, South)
	fmt.Println()
	conduct_test(100000, 0.5, 20, 80, East)
	fmt.Println()
	conduct_test(100000, 0.5, 20, 80, West)
}
