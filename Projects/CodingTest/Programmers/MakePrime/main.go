package main

import (
	"fmt"
	"math"
	"sort"
)

const (
	Prime = iota
	NotPrime
)

var (
	numbers []int
	isPrime []int
)

func main() {

	result := solution([]int{1, 2, 7, 6, 4})
	fmt.Println("result : ", result)
}

func solution(nums []int) int {
	numbers = ReverseNums(nums)
	isPrime = MakePrime()

	return solve()
}

func ReverseNums(nums []int) []int {
	reversenums := make([]int, len(nums))
	copy(reversenums, nums)
	sort.Sort(sort.Reverse(sort.IntSlice(reversenums)))

	return reversenums
}

func MakePrime() []int {
	var upper int
	upper = 0
	for i := 0; i < 3; i++ {
		upper += numbers[i]
	}

	ret := make([]int, upper+1)

	ret[0] = NotPrime
	ret[1] = NotPrime
	ret[2] = Prime

	for i := 4; i <= upper; i += 2 {
		ret[i] = NotPrime
	}

	sqrt := int(math.Sqrt((float64(upper))))

	for i := 3; i <= sqrt; i++ {
		if ret[i] == Prime {
			for k := i * i; k <= upper; k += i {
				ret[k] = NotPrime
			}
		}
	}

	return ret
}

func solve() int {
	var ret int
	var sumnumbers []int
	for i := 0; i < len(numbers)-2; i++ {
		for j := i + 1; j < len(numbers)-1; j++ {
			for k := j + 1; k < len(numbers); k++ {
				result := numbers[i] + numbers[j] + numbers[k]
				sumnumbers = append(sumnumbers, result)
			}
		}
	}

	for _, value := range sumnumbers {
		if isPrime[value] == Prime {
			ret++
		}
	}
	return ret
}
