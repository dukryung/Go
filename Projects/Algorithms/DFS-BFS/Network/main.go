package main

import (
	"fmt"
)

var arr []int

func main() {
	n := 3
	//n := 4
	var testarr = [][]int{{1, 1, 0}, {1, 1, 0}, {0, 0, 1}}
	//var testarr = [][]int{{1, 0, 0, 1}, {0, 1, 1, 0}, {0, 1, 1, 0}, {1, 1, 0, 1}}
	solution(n, testarr)
}

/*
type Queue struct {
	queue []*Pos
}
type Pos struct {
	x int
	y int
}

func (q *Queue) EnQueue(pos *Pos) {
	q.queue = append(q.queue, pos)
}

func (q *Queue) DeQueue() *Pos {
	if len(q.queue) == 0 {
		return nil
	}
	pos, poslist := q.queue[0], q.queue[1:]
	q.queue = poslist
	return pos
}

func solution(n int, computers [][]int) int {

	var answer int
	arr = make([]int, n)

	for i := 0; i < n; i++ {
		if arr[i] == 1 {
			continue
		}

		bfs(i, i, computers)
		answer++
	}
	fmt.Println("answer : ", answer)

	return answer
}

func bfs(x int, y int, computers [][]int) {
	q := &Queue{}

	q.EnQueue(&Pos{x, y})

	for {

		if len(q.queue) == 0 {
			break
		}

		pos := q.DeQueue()

		computers[pos.x][pos.y] = 0

		// x의 값과  같은 경우
		for i := 0; i < len(computers); i++ {
			if computers[pos.x][i] == 1 {
				q.EnQueue(&Pos{pos.x, i})
				computers[pos.x][i] = 0
				arr[pos.x] = 1
				arr[i] = 1
			}
		}
		// y의 값과 같은 경우
		for i := 0; i < len(computers); i++ {
			if computers[i][pos.y] == 1 {
				q.EnQueue(&Pos{i, pos.y})
				computers[i][pos.y] = 0
				arr[pos.y] = 1
				arr[i] = 1
			}
		}
	}

	return

}
*/

var (
	networks  [][]int
	Queue     []Pos
	isVisited []bool
)

type Pos struct {
	x int
	y int
}

func EnQueue(position Pos) {
	Queue = append(Queue, position)
}

func DeQueue() Pos {
	pos := Queue[0]
	Queue = Queue[1:]

	return pos
}

func solution(n int, computers [][]int) int {
	var ret int
	isVisited = make([]bool, len(computers))
	networks = make([][]int, len(computers))
	copy(networks, computers)

	for i := range isVisited {
		if isVisited[i] {
			continue
		}
		BFS(i, i)
		isVisited[i] = true
		fmt.Println("isVisited : ", isVisited)

		ret++
	}
	fmt.Println("ret : ", ret)
	return ret
}

func BFS(x int, y int) {
	var pos Pos
	pos.x = x
	pos.y = y
	EnQueue(pos)

	for {
		if len(Queue) == 0 {
			break
		}
		fmt.Println("-----------1")
		pos = DeQueue()
		fmt.Println("pos :", pos)
		networks[pos.x][pos.y] = 0
		for i := 0; i < len(networks); i++ {

			if networks[pos.x][i] == 1 {
				isVisited[i] = true
				networks[pos.x][i] = 0
				pos.y = i
				EnQueue(pos)
			}
		}

		for i := 0; i < len(networks); i++ {

			if networks[i][pos.y] == 1 {
				isVisited[i] = true
				networks[i][pos.y] = 0
				pos.x = i
				EnQueue(pos)
			}
		}
		fmt.Println(networks)
		//time.Sleep(time.Second * 10)
	}
}
