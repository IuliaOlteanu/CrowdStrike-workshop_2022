package main

import "fmt"
import "strings"

func noVowels (str string) map[string]int {
	wordList := strings.Split(str, " ")
	var vowels = [] rune{'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U'}
	freq := map[string]int{}
	for _, v := range wordList {
		
		for i := range v {
			for j := range vowels {
				if rune (v[i]) == vowels[j] {
					// fmt.Println(v[i], vowels[j])
					freq[v]++
				}
			}
		}
		//return freq
	}
	return freq
}
	

func main() {
	fmt.Println(noVowels("Ana are mere"))
}
