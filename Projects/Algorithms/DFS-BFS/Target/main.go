package main

import "fmt"

func main() {

	var numbers = []int{1, 1, 1, 1, 1}
	var target = 3
	answer := solution(numbers, target)

	fmt.Println("answer : ", answer)
}

/*
func solution(numbers []int, target int) int {

	currentnum := numbers[0]
	var answer int
	answer += dfs(currentnum, 1, numbers, target)
	answer += dfs(-currentnum, 1, numbers, target)

	return answer
}

func dfs(prev int, index int, numbers []int, target int) int {

	if index == len(numbers) {
		if prev == target {
			return 1
		}
		return 0
	}

	cur1 := prev + numbers[index]
	cur2 := prev - numbers[index]

	var answer int

	answer += dfs(cur1, index+1, numbers, target)
	answer += dfs(cur2, index+1, numbers, target)

	return answer
}
*/

func solution(numbers []int, target int) int {
	var answer int

	curnum := numbers[0]
	answer += DFS(curnum, 1, numbers, target)
	answer += DFS(-curnum, 1, numbers, target)

	return answer
}
func DFS(prevnum int, index int, numbers []int, target int) int {
	var answer int

	if index == len(numbers) {
		if prevnum == target {
			return 1
		}
		return 0
	}

	firstcurnum := prevnum + numbers[index]
	secondcurnum := prevnum - numbers[index]

	answer += DFS(firstcurnum, index+1, numbers, target)
	answer += DFS(secondcurnum, index+1, numbers, target)

	return answer
}
