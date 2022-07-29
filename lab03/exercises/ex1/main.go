package main

import (
	"fmt"
	// "sync"
)

// var wg sync.WaitGroup
// var mu sync.Mutex

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		if (v % 2 == 1) {
			sum = sum + v
		}
	}
	c <- sum
}
func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	c := make(chan int)

	go sum(s[:len(s) / 4], c)
	go sum(s[len(s) / 4 : len(s) / 2], c)
	go sum(s[len(s) / 2 : 3 * len(s) / 4], c)
	go sum(s[3 * len(s) / 4 :], c)

	x, y, z, t := <-c, <-c, <-c, <-c 

	fmt.Println(x + y + z + t) 

	
}