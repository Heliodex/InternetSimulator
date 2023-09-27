package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
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

	_, gen := Markov(tweets)
	_, gen0 := Markov(s0tweets)
	_, gen1 := Markov(s1tweets)


	fmt.Println("Done! Starting server...")

	r := gin.Default()
	r.GET("/:multi", func(c *gin.Context) {
		multi, ok := strconv.ParseFloat(c.Param("multi"), 32)
		if ok != nil {
			multi = 1
		}
		c.String(200, gen(float32(multi)))
	})
	r.GET("/:multi/:s", func(c *gin.Context) {
		multi, ok := strconv.ParseFloat(c.Param("multi"), 32)
		if ok != nil {
			multi = 1
		}
		sentiment := c.Param("s")
		if sentiment == "1" {
			c.String(200, gen1(float32(multi)))
		} else {
			c.String(200, gen0(float32(multi)))
		}
	})
	r.Run()
}
