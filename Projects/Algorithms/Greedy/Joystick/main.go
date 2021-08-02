package main

/*
func main() {
	//name := "ABAAAAAAAAABA"

	name := "JAN"

	solution(name)
}

func solution(name string) int {
	var result int
	var left, right int
	var leftandright int
	var leftside, rightside int
	var index int
	//위 아래 이동 값
	for i := 0; i < len(name); i++ {
		if int(name[i]) < 78 {
			result += int(name[i]) - 65
		} else {
			result += 91 - int(name[i])
		}
	}

	for i := 0; i < len(name); i++ {
		if name[i] != 'A' {
			left = i
		}
	}

	for i := len(name) - 1; i > 0; i-- {
		if name[i] != 'A' {
			right = len(name) - i
		}
	}

	for i := 1; i < len(name); i++ {
		if name[i] != 'A' {
			leftside = 2 * i
			index = i
			break
		}
	}

	for i := len(name) - 1; i > index; i-- {
		if name[i] != 'A' {
			rightside = len(name) - i
		}
	}
	leftandright = leftside + rightside

	if left >= right {
		if right >= leftandright {
			result += leftandright
		} else {
			result += right
		}
	} else {
		if left >= leftandright {
			result += leftandright
		} else {
			result += left
		}
	}

	fmt.Println("result : ", result)

	return result
}
*/

var (
	word []string
)

func main() {
	//name := "ABAAAAAAAAABA"
	//name := "JAN"
	//name := "JEROEN"
	//name := "J"
	//name := "JAZ"
	//name := "AAAAAAAABAABA"
	name := "ABAAAAAAAAABB"

	solution(name)
}

func solution(name string) int {

	for _, v := range name {
		word = append(word, string(v))
	}

	leftandright := direction()
	updown := updown(name)

	return leftandright + updown
}

func direction() int {
	var l, r, lr int

	for i, v := range word {
		if i == 0 {
			continue
		}

		if v != "A" {
			r = i
		}
	}

	for i, v := range word {
		if i == 0 {
			continue
		}

		if v != "A" {
			l = len(word) - i
			break
		}
	}

	var limit, value int
	for i, v := range word {
		if i == 0 || v == "A" {
			continue
		}

		if limit >= 2 {
			break
		}

		limit++
		if limit == 1 {
			value += i * 2
		} else if limit == 2 {
			value += len(word) - i
		}
	}

	lr = value

	if l > r {
		if r > lr {
			return lr
		}
		return r
	}

	if l > lr {
		return lr
	}

	return l

}

func updown(name string) int {
	var ret int
	for _, v := range name {
		if int(v) < 78 {
			ret += int(v) - 65
		} else {
			ret += 91 - int(v)
		}
	}
	return ret
}
