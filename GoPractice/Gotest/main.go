package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 3)

	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("test")
			ch <- i
		}
	}()

	for i := 0; i < 100; i++ {
		var x int
		time.Sleep(time.Second * 5)
		x = <-ch
		fmt.Println("x :", x)
	}

}
