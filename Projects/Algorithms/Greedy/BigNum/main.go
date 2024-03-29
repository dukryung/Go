package main

import (
	"fmt"
	"strconv"
)

func main() {
	//var number = "1231234"
	//var k = 3
	//var number = "1924"
	//var k = 2
	var number = "4177252841"
	var k = 4

	result := solution(number, k)
	fmt.Println("result : ", result)

}

/*
func solution(number string, k int) string {
	var result string
	var extractcnt, compareextractcnt int
	var index int

	extractcnt = len(number) - k
	compareextractcnt = len(number) - k
	for i := 0; i < extractcnt; i++ {
		var extractstr byte
		for p := 0; p <= len(number)-compareextractcnt; p++ {
			if extractstr >= number[p] {

			} else {
				index = p
				extractstr = number[p]
			}
		}
		compareextractcnt -= 1
		if len(number) != index+1 {
			number = string(number[index+1:])
		}
		result += string(extractstr)

	}

	return result
}
*/

func solution(number string, k int) string {
	var length int
	var s int
	var limit = k

	num := make([]int, len(number))
	var result string
	for i, v := range number {
		n, _ := strconv.Atoi(string(v))
		num[i] = n
	}
	for {

		if length == len(number)-k {
			break
		}

		var value int
		for i := s; i <= limit; i++ {

			if num[i] > value {
				value = num[i]
				s = i + 1
			}
		}
		result += strconv.Itoa(value)
		limit++
		length++

	}

	return result
}
