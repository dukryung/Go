package main

import "fmt"

func main() {
	var begin = "1234567000"
	//var target = "cog"
	//var target = "hot"
	var target = "1234567899"
	//var words = []string{"hot", "dot", "dog", "lot", "log", "cog"}
	//var words = []string{"cog", "dot", "dog", "lot", "log", "hot"}
	//var words = []string{"hot", "dot", "dog", "lot", "log"}
	//var words = []string{"hot", "dot", "dog", "lot", "log"}
	var words = []string{"1234567800", "1234567890", "1234567899"}
	solution(begin, target, words)
}

func solution(begin string, target string, words []string) int {
	var result int
	var mincnt int
	var istarget bool

	for _, word := range words {
		if word == target {
			istarget = true
		}
	}

	if istarget == false {
		return 0
	}

	result = DFS(begin, target, words, mincnt)
	fmt.Println(result)
	return result

}

func DFS(begin string, target string, words []string, mincnt int) int {
	if begin == target {
		return mincnt
	}

	for i, word := range words {
		if word != "" {
			var samealphabetcnt = 0
			var sametargetalphacnt = 0

			for i := 0; i < len(word); i++ {
				if begin[i] == target[i] {
					sametargetalphacnt++
				}
				if begin[i] == word[i] {
					samealphabetcnt++
				}
			}

			if sametargetalphacnt == len(begin)-1 {
				fmt.Println("word", word)
				fmt.Println("mincnt : ", mincnt)
				return mincnt + 1
			}

			if samealphabetcnt == len(begin)-1 {
				words[i] = ""
				mincnt = DFS(word, target, words, mincnt+1)

				return mincnt

			}
		}
	}

	return 0
}
