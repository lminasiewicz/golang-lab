package main

import (
	"fmt";
	"math/rand/v2"
	"math"
)

type Field struct {
	tree bool
	age int
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
const MAX_BURN_RESISTANCE = 0.2

func attempt_burn(age int) bool {
	if age > BURN_RESISTANCE_THRESHOLD {
		return true
	} else {
		var chance float32 = MAX_BURN_RESISTANCE + (float32(age) * float32((1 - MAX_BURN_RESISTANCE) / BURN_RESISTANCE_THRESHOLD))
		if rand.Float32() < chance {
			return true
		}
		return false
	}
}

func initialize_forest(width int, length int, rate float32) [][]Field {
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
			if coords[0] + i >= 0 && coords[0] + i < len(forest) && coords[1] + j >= 0 && coords[1] + j < len(forest[0]) {
				if !(i == 0 && j == 0) {
					result = append(result, []int{coords[0]+i, coords[1]+j})
				}
			}
		}
	}
	return result
}

func burn(forest [][]Field, current [2]int, wind Direction) {
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
	fld := forest[lightning[0]][lightning[1]]
	if fld.tree && !fld.burned {
		fld.burned = true
		burn(forest, lightning, wind)
	}
}

func lightning_strike_until_successful(forest [][]Field, wind Direction) [2]int {
	lightning := [2]int{rand.IntN(len(forest)), rand.IntN(len(forest[0]))}
	for {
		lightning = [2]int{rand.IntN(len(forest)), rand.IntN(len(forest[0]))}
		fld := forest[lightning[0]][lightning[1]]
		if fld.tree && !fld.burned {
			lightning_strike(forest, lightning, wind)
			break
		}
	}
	return lightning
}

func print_forest(forest [][]Field, lightning [2]int) {
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
	// Funkcja do wypisywania statystyk dot. stanu lasu po trafieniu piorunem
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
	forestation := math.Round(10000 * (float64(trees) / float64(length * width))) / 100
	burned_rate := math.Round(10000 * (float64(burned) / float64(trees))) / 100
	fmt.Print("Forest ", length, "x", width, " (", length*width, " spaces)\n")
	fmt.Print("Forestation: ", forestation, "%\n")
	fmt.Print("Trees: ", trees, "  |  Burned: ", burned, "\n")
	fmt.Print(burned_rate, "% of the forest has been burned.\n")
	if burned == 0 {
		fmt.Println("Lightning Strike missed.")
	}
}

func simulate_once(length int, width int, forestation_rate float32, wind Direction, lightning_is_accurate bool) {
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
	fmt.Println(lightning)
}

func main() {
	simulate_once(40, 20, 0.4, None, false)
}