package main

import (
	"fmt"
)

func main() {
	var n = 4

	d := make([]int, 100)

	arr := make([]int, n)

	for i := 0; i < n; i++ {
		var num int
		fmt.Scanln(&num)
		arr[i] = num

	}

	d[0] = arr[0]
	d[1] = max(arr[0], arr[1])

	for i := 2; i < n; i++ {
		d[i] = max(d[i-1], d[i-2]+arr[i])
	}
	fmt.Println(d[n-1])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
