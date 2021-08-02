package main

import (
	"sort"
)

func main() {

	//a := [][]int{{0, 1, 5}, {1, 2, 3}, {2, 3, 3}, {3, 1, 2}, {3, 0, 4}, {2, 4, 6}, {4, 0, 7}}
	a := [][]int{{0, 1, 1}, {3, 4, 1}, {1, 2, 2}, {2, 3, 4}}

	solution(5, a)

}

/*
var isnums []int


func solution(n int, costs [][]int) int {

	var result int

	for i := 0; i < n; i++ {
		isnums = append(isnums, i)
	}

	//오름 차순 정렬 (최소 비용으로 연결 해야하기 때문에 우선 오름차순으로 정렬)
	for i := 0; i < len(costs)-1; i++ {
		for k := i + 1; k < len(costs); k++ {

			if costs[i][2] > costs[k][2] {
				temp := costs[i]
				costs[i] = costs[k]
				costs[k] = temp
			}
		}
	}

	for i := 0; i < len(costs); i++ {

		if merge(costs[i][0], costs[i][1]) == true {
			result += costs[i][2]
		}

	}
	fmt.Println(result)

	return result
}

//최상위 부모 노드 찾기 (재귀 함수 - 해당 조건을 충족 못시키면 찾을 때 까지 자기 자신 함수를 계속 호출)
func findparent(x int) int {
	if x == isnums[x] {
		return x
	} else {
		return findparent(isnums[x])
	}
}

func merge(a int, b int) bool {

	a = findparent(a)
	b = findparent(b)

	if a == b {
		return false
	} else {
		if a > b {
			isnums[a] = b
		} else {
			isnums[b] = a
		}
	}

	return true
}
*/

var lands []int

func solution(n int, costs [][]int) int {
	lands = make([]int, n)
	for i := range lands {
		lands[i] = i
	}

	sort.Slice(costs[:], func(i, j int) bool {
		for range costs[i] {
			if costs[i][2] == costs[j][2] {
				continue
			}
			return costs[i][2] < costs[j][2]
		}
		return false
	})

	var ret int
	for _, v := range costs {
		if merge(v[0], v[1]) {
			ret += v[2]
		}
	}
	return ret
}

func findparent(x int) int {
	if lands[x] == x {
		return x
	} else {
		return findparent(lands[x])
	}
}

func merge(a int, b int) bool {

	a = findparent(a)
	b = findparent(b)

	if a == b {
		return false
	} else {
		if a > b {
			lands[a] = b
		} else {
			lands[b] = a
		}
	}

	return true
}
