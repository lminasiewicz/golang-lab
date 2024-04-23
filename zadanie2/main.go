package main

import ("fmt";
		"math/rand/v2")

type Field struct {
	tree bool
	age int
	burned bool
}

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

func get_coords_around(forest [][]Field, coords [2]int) [][]int {
	result := make([][]int, 0, 8)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if coords[0] + i >= 0 && coords[0] + i < len(forest) && coords[1] + j >= 0 && coords[1] + j < len(forest[0]) {
				if !(i == 0 && j == 0) {
					result = append(result, []int{coords[0]+i, coords[1]+j})
				}
			}
		}
	}
	return result
}

func burn(forest [][]Field, current [2]int) {
	available_spaces := get_coords_around(forest, current)
	for _, space := range available_spaces {
		fld := forest[space[0]][space[1]]
		if fld.tree && !fld.burned {
			if attempt_burn(fld.age) {
				fld.burned = true
				space_arr := [2]int{space[0], space[1]}
				burn(forest, space_arr)
			}
		}
	}
}

func print_forest(forest [][]Field) {
	for row := 0; row < len(forest); row++ {
		for _, fld := range forest[row] {
			if !fld.tree {
			fmt.Print(" ")
			} else if fld.tree && !fld.burned {
				fmt.Print("◯")
			} else if fld.tree && fld.burned {
				fmt.Print("⬤")
			}
		}
		fmt.Println()
	}
}

func main() {
	length := 20
	width := 10
	forest := initialize_forest(width, length, 0.3)
	lightning := [2]int{rand.IntN(width), rand.IntN(length)}
	fmt.Println(lightning)
	burn(forest, lightning)
	print_forest(forest)
}