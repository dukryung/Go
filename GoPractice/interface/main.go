package main

import (
	"fmt"
	"strconv"
)

type StructA struct {
	val int
}

func (a *StructA) AAA(num int) int {
	a.val = num
	return a.val
}

func (a *StructA) BBB(num int) string {
	num = 5
	return strconv.Itoa(num)
}

func main() {
	sa := &StructA{}

	sa.AAA(5)
	fmt.Println(sa)

}
