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
	if len(os.Args) > 1 && os.Args[1] == "serialise" {
		Serialise()
		return
	}

	fmt.Println("Starting server...")

	startTime := time.Now()

	chains := []string{"chain", "chain0", "chain1"}
	gens := make(map[string]func(float32) string)
	var wg sync.WaitGroup

	// load chains
	for _, v := range chains {
		go func(v string) {
			chain := Deserialise(v)
			gens[v] = FromChain(chain)
			fmt.Println("Loaded", v)
			wg.Done()
		}(v)
	}

	// wait for all chains to load
	wg.Add(len(chains))
	wg.Wait()

	fmt.Println("Loaded all chains in", time.Since(startTime))

	r := gin.Default()
	r.GET("/:chain", func(c *gin.Context) {
		chainName := c.Param("chain")
		gen, chainFound := gens[chainName]
		if !chainFound {
			c.String(404, "Chain not found")
			return
		}
		c.String(200, gen(1))
	})
	r.GET("/:chain/:multi", func(c *gin.Context) {
		chainName := c.Param("chain")
		multi, ok := strconv.ParseFloat(c.Param("multi"), 32)
		if ok != nil {
			multi = 1
		}
		gen, chainFound := gens[chainName]
		if !chainFound {
			c.String(404, "Chain not found")
			return
		}
		c.String(200, gen(float32(multi)))
	})
	r.Run()
}
