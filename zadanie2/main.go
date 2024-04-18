package main

import ("fmt";
		"math/rand/v2")

type Field struct {
	tree bool
	age int
	burned bool
}

func attempt_burn(age int) bool {
	if age > 400 {
		return true
	} else {
		var chance float32 = 0.2 + (float32(age) * float32(0.002))
		if rand.Float32() < chance {
			return true
		}
		return false
	}
}

func initialize_forest(length int, width int, rate float32) [][]Field {
	forest := make([][]Field, length)
	for i := 0; i < width; i++ {
		forest[i] = make([]Field, width)
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

func main() {
	forest := initialize_forest(10, 10, 0.3)
	fmt.Println(forest)
}