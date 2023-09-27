package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args)  > 1 && os.Args[1] == "serialise" {
		Serialise()
		return
	}

	fmt.Println("Starting server...")

	startTime := time.Now()

	chains := make(map[string]map[string]map[string]int)
	
	var wg sync.WaitGroup

	// load chains
	for _, v := range []string{"chain", "chain0", "chain1"} {
		go func(v string) {
			chain := Deserialise(v)
			chains[v] = chain
			fmt.Println("Loaded", v)
			wg.Done()
		}(v)
	}

	// wait for all chains to load
	wg.Add(3)
	wg.Wait()

	fmt.Println("Loaded all chains in", time.Since(startTime))

	chain, chain0, chain1 := chains["chain"], chains["chain0"], chains["chain1"]
	
	gen := FromChain(chain)
	gen0 := FromChain(chain0)
	gen1 := FromChain(chain1)

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
