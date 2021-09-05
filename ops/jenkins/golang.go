package main

import "fmt"

func main() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
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
			for y := length1 - 1; y >= 2; y-- {
				var so int
				so = matrix[i][j]
				matrix[i][j] = matrix[i][y]
				matrix[i][y] = so
				break he
			}
		}
	}
	fmt.Printf("%#v\n", matrix)
}
