package main

import (
	"fmt"
	"strconv"
)

func main() {
	//result := solution(5, 31168)
	result := solution(5, 3600)
	//result := solution(8, 53)
	//result := solution(1, 1121)
	//result := solution(5, 12)
	fmt.Println("result :", result)
}

func solution(N int, number int) int {
	var DP = make([][]int, 9)
	var Nstr string
	for i := 1; i < 9; i++ {

		Nstr += strconv.Itoa(N)
		Nint, _ := strconv.Atoi(Nstr)
		var arr []int
		arr = append(arr, Nint)
		DP[i] = arr

		for j := 1; j < i; j++ {
			for _, num := range DP[i-j] {
				for _, num2 := range DP[j] {

					arr = append(arr, num+num2)
					arr = append(arr, num*num2)
					arr = append(arr, num-num2)

					if num2 != 0 {
						arr = append(arr, num/num2)
					}
				}
			}
		}

		for _, num := range arr {
			if num == number {

				//8보다 클 때 이므로:  > 8
				if i > 8 {
					return -1
				} else {
					return i
				}

			}
		}

		DP[i] = arr

	}

	return -1
}
