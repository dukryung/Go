package main

func remove(s []int, i int) []int {
	switch {
	case len(s)-1 == i:
		return append(s[:len(s)-1])
	case i == 0:
		return append(s[1:])
	default:
		return append(s[:i], s[i+1:]...)

	}
}

func Solution(n int, lost []int, reserve []int) int {
	// 여벌 체육복을 가져왔으나 도난 당한 경우는 제외한다.
	for p := 0; p < len(lost); p++ {
		for q := 0; q < len(reserve); q++ {
			if reserve[q] == lost[p] {
				lost = remove(lost, p)
				p--
				reserve = remove(reserve, q)
				q--
				break
			}
		}
	}

	can := n - len(lost) // 체육수업을 들을 수 있는 학생 수

	// 체육복을 잃어버린 학생이 앞, 또는 뒤 번호에 체육복을 빌려준다.
	for p := 0; p < len(lost); p++ {
		for q := 0; q < len(reserve); q++ {
			if reserve[q] == lost[p]-1 || reserve[q] == lost[p]+1 {
				can++
				lost = remove(lost, p)
				p--
				reserve = remove(reserve, q)
				q--
				break
			}
		}
	}
	return can
}
