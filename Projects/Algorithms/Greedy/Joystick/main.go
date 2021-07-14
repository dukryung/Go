package main

import (
	"fmt"
)

func main() {
	//name := "ABAAAAAAAAABA"
	name := "JAN"

	solution(name)
}

func solution(name string) int {
	var result int
	var left, right int
	var leftandright int
	var leftside, rightside int
	var index int
	//위 아래 이동 값
	for i := 0; i < len(name); i++ {
		if int(name[i]) < 78 {
			result += int(name[i]) - 65
		} else {
			result += 91 - int(name[i])
		}
	}

	for i := 0; i < len(name); i++ {
		if name[i] != 'A' {
			left = i
		}
	}

	for i := len(name) - 1; i > 0; i-- {
		if name[i] != 'A' {
			right = len(name) - i
		}
	}

	for i := 1; i < len(name); i++ {
		if name[i] != 'A' {
			leftside = 2 * i
			index = i
			break
		}
	}

	for i := len(name) - 1; i > index; i-- {
		if name[i] != 'A' {
			rightside = len(name) - i
		}
	}
	leftandright = leftside + rightside

	if left >= right {
		if right >= leftandright {
			result += leftandright
		} else {
			result += right
		}
	} else {
		if left >= leftandright {
			result += leftandright
		} else {
			result += left
		}
	}

	fmt.Println("result : ", result)

	return result
}
