package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Queue struct {
	queue []*Pos
}
type Pos struct {
	x int
	y int
}

var maze = make([][]int, 5)
var n, m int
var dx = []int{-1, 1, 0, 0}
var dy = []int{0, 0, -1, 1}

func (q *Queue) Enqueue(pos *Pos) {
	q.queue = append(q.queue, pos)
}

func (q *Queue) Dequeue() (*Pos, error) {
	if len(q.queue) == 0 {
		return nil, nil
	}

	pos, poslist := q.queue[0], q.queue[1:]
	q.queue = poslist
	return pos, nil
}

func main() {

	//맵 지정
	maze = [][]int{{1, 0, 1, 0, 1, 0}, {1, 1, 1, 1, 1, 1}, {0, 0, 0, 0, 0, 1}, {1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1}}

	//맵 길이 지정 받기
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	slice := strings.Split(text, " ")
	nums := make([]int, 2)
	for i, num := range slice {
		nums[i], _ = strconv.Atoi(num)
		n = nums[0]
		m = nums[1]
	}

	//미로 찾기(횟수 확인)
	result, err := BFS(0, 0)
	if err != nil {
		fmt.Println("[ERR] internal error")
	}
	fmt.Println("result : ", result)
	fmt.Println(maze)

}

func BFS(x int, y int) (int, error) {
	q := &Queue{}
	q.Enqueue(&Pos{x, y})
	for {

		if len(q.queue) == 0 {
			fmt.Println("quene end")
			break
		}

		pos, err := q.Dequeue()
		if err != nil {
			fmt.Println("deque err")
			return 0, nil
		}

		x = pos.x
		y = pos.y

		for i := range dy {
			nx := x + dx[i]
			ny := y + dy[i]

			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}

			if maze[nx][ny] == 0 {
				continue
			}

			if maze[nx][ny] == 1 {
				maze[nx][ny] = maze[x][y] + 1

				q.Enqueue(&Pos{nx, ny})
			}
		}
	}

	return maze[n-1][m-1], nil

}
