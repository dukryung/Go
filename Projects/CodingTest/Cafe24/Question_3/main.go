package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var firstnum string
	scn := bufio.NewScanner(os.Stdin)
	scn.Scan()
	a, _ := strconv.Atoi(scn.Text())
	scn.Scan()
	b, _ := strconv.Atoi(scn.Text())

	if a < b {
		for i := a; i <= b; i++ {
			str := strconv.Itoa(i)
			firstnum += str
			fmt.Printf("%010s\n", firstnum)
		}
	} else if a > b {
		for i := a; i > 0; i-- {
			str := strconv.Itoa(i)
			firstnum += str
			fmt.Printf("%010s\n", firstnum)
		}
	}
}
