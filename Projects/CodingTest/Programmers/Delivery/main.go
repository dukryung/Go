package main

func main() {

	solution(5, [][]int{{1, 2, 1}, {2, 3, 3}, {5, 2, 2}, {1, 4, 2}, {5, 3, 1}, {5, 4, 2}}, 3)
	//solution(6, [][]int{{1, 2, 1}, {1, 3, 2}, {2, 3, 2}, {3, 4, 3}, {3, 5, 2}, {3, 5, 3}, {5, 6, 1}}, 4)

}

func solution(N int, road [][]int, k int) int {
	length := make([]int, N+1)
	for i := range length {
		length[i] = k + 1
	}

	queue := make([]int, 0)
	queue = append(queue, 1)
	length[1] = 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, r := range road {
			if r[0] == cur || r[1] == cur {
				from, to := r[0], r[1]
				if r[1] == cur {
					from, to = r[1], r[0]
				}
				//시작 해당 노드 위치를 기준으로 걸리는 시간을 계산한 값과 다른 노드의 최소 걸리는 시간을 비교하는 로직
				if length[to] > length[from]+r[2] {
					length[to] = length[from] + r[2]
					queue = append(queue, to)
				}
			}
		}
	}
	cnt := 0
	for _, l := range length {
		if l <= k {
			cnt++
		}
	}

	return cnt
}
