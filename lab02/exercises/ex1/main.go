package main

import "fmt"
import "strings"

func  WordCount (str string) map[string]int {
	wordList := strings.Split(str, " ")
	freq := map[string]int {}
	for _, v := range wordList {
		// fmt.Println(v)
		freq[v]++
		// fmt.Println(freq)

	}
	return freq
}

func main() {
	fmt.Println(WordCount("Ana are are mere"))
}
