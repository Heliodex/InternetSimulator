package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func loadTweets() [][]interface{} {
	file, err := os.ReadFile("../data/tweets.json")
	if err != nil {
		log.Fatalln(err)
	}

	var jsonData [][]interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	return jsonData
}

func Serialise() {
	// tweet[0] is sentiment (0 or 1), tweet[1] is text
	tweetSets := make(map[string][]string)

	for _, tweet := range loadTweets() {
		stringTweet := tweet[1].(string)
		if tweet[0] == float64(0) {
			tweetSets["chain0"] = append(tweetSets["chain0"], stringTweet)
		} else {
			tweetSets["chain1"] = append(tweetSets["chain1"], stringTweet)
		}
		tweetSets["chain"] = append(tweetSets["chain"], stringTweet)
	}
	fmt.Println("Loaded tweets file")

	fmt.Println("Loaded", len(tweetSets["chain"]), "tweets")
	fmt.Println("Loaded", len(tweetSets["chain0"]), "s0tweets")
	fmt.Println("Loaded", len(tweetSets["chain1"]), "s1tweets")

	fmt.Println("Generating Markov chains. This may take a while...")

	var wg sync.WaitGroup

	// Load chains
	for k, v := range tweetSets {
		go func(k string, v []string) {
			chain := Markov(v)

			f, err := os.Create("../data/" + k + ".gob")
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
		}(k, v)
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
