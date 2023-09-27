package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Markov(data []string) (map[string]map[string]int, func(float32) string) {
	startTime := time.Now()
	// Create a Markov chain from a dataset

	chain := make(map[string]map[string]int)

	allWords := strings.Join(data, " ")

	// Get two random unicode characters that aren't in the dataset,
	// These will be used for start and end markers
	randChar := func(not string) string {
		for {
			r := rune(rand.Intn(65535))
			if !strings.Contains(allWords, string(r)) && string(r) != not {
				return string(r)
			}
		}
	}

	start := randChar("\000")
	end := randChar(start)

	fmt.Println("Done random stuff in", time.Since(startTime))
	startTime = time.Now()

	startPart := start + " "
	endPart := " " + end

	wordsArray := strings.Split(startPart+strings.Join(data, endPart+" "+startPart)+endPart, " ")

	fmt.Println("Split strings in", time.Since(startTime))
	startTime = time.Now()

	for i, word := range wordsArray[:len(wordsArray)-1] {
		// Add a word to the Markov chain
		if chain[word] == nil {
			chain[word] = make(map[string]int)
		}
		chain[word][wordsArray[i+1]]++
	}

	fmt.Println("Added words in", time.Since(startTime))

	keys := make([]string, 0, len(chain[start]))
	for k := range chain[start] {
		keys = append(keys, k)
	}

	// Return the chain alongside a function to generate a tweet
	return chain, func(lengthMultiplier float32) string {
		startTime = time.Now()

		var words []string

		weightedRandom := func(list map[string]int) string {
			total := 0
			for _, v := range list {
				total += v
			}
			rand := rand.Intn(total)
			for k, v := range list {
				rand -= v
				if rand <= 0 {
					return k
				}
			}
			return ""
		}

		currentWord := weightedRandom(chain[start])

		for currentWord != end {
			words = append(words, currentWord)

			nextWordsTable := chain[currentWord]
			nextWords := make([]string, 0, len(nextWordsTable))
			for k := range nextWordsTable {
				nextWords = append(nextWords, k)
			}

			currentWord = weightedRandom(nextWordsTable)

			for currentWord == end {
				// generate random number, affected by lengthMultiplier
				rand := rand.Intn(int(100 * lengthMultiplier))
				if rand < 100 {
					break
				}
				currentWord = weightedRandom(nextWordsTable)
			}
		}

		fmt.Println("Generated tweet in", time.Since(startTime))

		return strings.Join(words, " ")
	}
}
