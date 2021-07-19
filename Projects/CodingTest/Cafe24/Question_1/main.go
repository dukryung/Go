package main

import (
	"fmt"
)

func main() {
	arr := []int{5, 4, 2, 1}
	solution(arr)
}

func solution(arr []int) int {

	for i := len(arr) - 1; i > 0; i-- {
		lcm := getLCM(arr[i-1], arr[i])
		arr = arr[:i-1]

		arr = append(arr, lcm)
		fmt.Println(arr)
		if len(arr) == 1 {
			break
		}
	}

	return arr[0]
}

// 최소 공배수 구하기 (유클리드 알고리즘 이용)
func getLCM(first, second int) (lcm int) {
	return first * second / getGCD(first, second) // 두 수의 곱 나누기 최대공약수
}

// 최대 공약수 구하기 (유클리드 알고리즘 이용)
func getGCD(first, second int) (gcd int) {
	if first < second { // fn에 큰 값을 오게 하기
		second, first = first, second
	}

	for second != 0 { // second가 0이 될때까지 반복
		first, second = second, first%second
	}
	return first
}
