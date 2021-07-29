package main

import (
	"fmt"
)

var (
	visited []bool
	roads   [][]int
)

func main() {
	/*
		test := make([][]int, 5)

		test[1] = append(test[1], 3)
		test[1] = append(test[1], 2)

		fmt.Println("test : ", test)
	*/
	solution(5, [][]int{{1, 2, 1}, {2, 3, 3}, {5, 2, 2}, {1, 4, 2}, {5, 3, 1}, {5, 4, 2}}, 3)
	//solution(6, [][]int{{1, 2, 1}, {1, 3, 2}, {2, 3, 2}, {3, 4, 3}, {3, 5, 2}, {3, 5, 3}, {5, 6, 1}}, 4)
}

func solution(N int, road [][]int, k int) int {
	answer := 0
	roads = make([][]int, len(road))
	visited = make([]bool, N+1)

	copy(roads, road)
	fmt.Println("roads : ", roads)
	answer = searchisland(1, k)

	fmt.Println("answer : ", answer)

	return answer
}

func searchisland(landnum int, k int) int {
	var ret = 0
	if k <= 0 {
		visited[landnum] = true
		ret += 1
		return ret
	}

	for _, road := range roads {
		if visited[landnum] {
			continue
		}

		if road[0] == landnum {
			if 0 <= k-road[2] {
				visited[road[0]] = true
				fmt.Println("road1 : ", road)
				ret = searchisland(road[1], k-road[2])
				ret += 1

			}
		} else if road[1] == landnum {
			if 0 <= k-road[2] {
				visited[road[1]] = true
				fmt.Println("road2 : ", road)
				ret = searchisland(road[0], k-road[2])
				ret += 1

			}
		}
		visited[landnum] = false
	}

	fmt.Println("ret : ", ret)
	return ret
}
