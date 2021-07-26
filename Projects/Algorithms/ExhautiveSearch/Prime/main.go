package main

import (
	"fmt"
	"math"
	"sort"
)

const (
	Prime = iota
	NotPrime
	CheckedPrime
)

var (
	nums    []int
	picked  []bool
	isPrime []int
)

func main() {
	result := solution("1234")
	fmt.Println("result : ", result)
}

func solution(numbers string) int {
	nums = StrtoIntSlice(numbers)
	picked = make([]bool, len(nums))
	isPrime = makePrimes()

	return solve(0)
}

func StrtoIntSlice(str string) []int {
	ret := make([]int, len(str))

	for i, v := range str {
		ret[i] = int(v - '0')
	}

	sort.Ints(ret)
	return ret
}

func makePrimes() []int {
	reverseNums := make([]int, len(nums))
	copy(reverseNums, nums)
	sort.Sort(sort.Reverse(sort.IntSlice(reverseNums)))

	upper := 0

	for _, v := range reverseNums {
		upper *= 10
		upper += v
	}

	ret := make([]int, upper+1)

	ret[0] = NotPrime
	ret[1] = NotPrime
	ret[2] = Prime

	for i := 4; i <= upper; i += 2 {
		ret[i] = NotPrime
	}

	sqrt := int(math.Sqrt(float64(upper)))

	for i := 3; i <= sqrt; i++ {
		if ret[i] == Prime {
			for j := i * i; j <= upper; j += i {
				ret[j] = NotPrime
			}
		}
	}

	return ret
}

func solve(number int) int {
	ret := 0
	for i := range nums {
		if picked[i] {
			continue
		}

		picked[i] = true

		next := number*10 + nums[i]
		if isPrime[next] == Prime {
			ret++
			isPrime[next] = CheckedPrime
		}

		ret += solve(next)
		picked[i] = false
	}

	return ret
}
