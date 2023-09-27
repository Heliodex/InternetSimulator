package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"encoding/gob"
)

func Serialise() {
	file, err := os.ReadFile("../data/tweets.json")
	if err != nil {
		log.Fatalln(err)
	}

	var jsonData [][]interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Loaded tweets file")

	// tweet[0] is sentiment (0 or 1), tweet[1] is text
	var s0tweets []string
	var s1tweets []string
	var tweets []string

	for _, tweet := range jsonData {
		stringTweet := tweet[1].(string)
		if tweet[0] == float64(0) {
			s0tweets = append(s0tweets, stringTweet)
		} else {
			s1tweets = append(s1tweets, stringTweet)
		}
		tweets = append(tweets, stringTweet)
	}

	fmt.Println("Loaded", len(tweets), "tweets")
	fmt.Println("Loaded", len(s0tweets), "s0tweets")
	fmt.Println("Loaded", len(s1tweets), "s1tweets")
	fmt.Println("Generating Markov chains. This may take a while...")

	var wg sync.WaitGroup

	// Load chains
	for _, v := range []string{"chain", "chain0", "chain1"} {
		go func(v string) {
			chain := Markov(s0tweets)

			f, err := os.Create("../data/" + v + ".gob")
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()

			enc := gob.NewEncoder(f)
			err = enc.Encode(chain)
			if err != nil {
				log.Fatalln(err)
			}
			wg.Done()
		}(v)
	}

	// Wait for all chains to load
	wg.Add(3)
	wg.Wait()
}

func Deserialise(chainName string) map[string]map[string]int {
	f, err := os.Open("../data/" + chainName + ".gob")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var chain map[string]map[string]int
	err = dec.Decode(&chain)
	if err != nil {
		log.Fatalln(err)
	}

	return chain
}
