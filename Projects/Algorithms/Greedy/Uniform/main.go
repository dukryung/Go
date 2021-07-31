package main

import "fmt"

func main() {
	//result := solution(5, []int{2, 4}, []int{1, 3, 5})
	//result := solution(5, []int{2, 4}, []int{3})
	//result := solution(3, []int{1, 2}, []int{2, 3})
	//result := solution(3, []int{2, 3}, []int{1})
	//result := solution(24, []int{12, 13, 16, 17, 19, 20, 21, 22}, []int{1, 22, 16, 18, 9, 10})
	result := solution(8, []int{2, 3, 4}, []int{1})
	fmt.Println("result :", result)
}

func solution(n int, lost []int, reserve []int) int {
	var length = make([]int, n+1)

	for i := 1; i <= n; i++ {
		length[i] = 1
	}
	length[0] = -2

	for _, value := range lost {
		length[value] -= 1
	}

	for _, value := range reserve {
		length[value] += 1
	}

	for i := range length {
		if length[i] != 2 {
			continue
		}

		if length[i-1] == 0 && i-1 > 0 {
			length[i] = 1
			length[i-1] = 1
			continue
		}

		if i+1 >= len(length) {
			continue
		} else {
			if length[i+1] == 0 {
				length[i] = 1
				length[i+1] = 1
				continue
			}
		}

	}
	fmt.Println(length)
	var ret int

	for _, value := range length {
		if value >= 1 {
			ret++
		}
	}
	return ret
}
