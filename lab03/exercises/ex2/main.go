package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var mu sync.Mutex
var v map[int]int

func sum(s int) {
	mu.Lock()
	
	if (s % 2 == 0) {
		v[0] = v[0] + s
	}

	mu.Unlock()

	wg.Done()
}
func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}
	v = map[int]int{}
	for _, v := range s  {
		wg.Add(1)
		go sum(v)
	}

	wg.Wait()
	fmt.Println(v[0])
}