package main

import (
	"fmt"
)

func main() {
	matrix := [][]int{
		{1, 2, 3},
	}
	rotateRight(matrix)
	// fmt.Println(matrix)
}

func rotateRight(matrix [][]int) {
	length := len(matrix)
	for i := 0; i < length; i++ {
		length1 := len(matrix[i])
	he:
		for j := 0; j < length1; j++ {
			for y := length1 - 1; y == 2; y-- {
				fmt.Println(matrix[i][j], matrix[i][y])
				break he
			}
		}
	}
	fmt.Printf("%#v\n", matrix)
}
