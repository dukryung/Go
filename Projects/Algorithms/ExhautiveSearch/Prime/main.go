package main

import "fmt"

func main() {
	var teststr string
	teststr = "1234"
	for _, c := range teststr {
		y := int(c - '0')
		x := int('0' - '0')

		fmt.Println("y : ", y)
		fmt.Println("x : ", x)
	}
}
