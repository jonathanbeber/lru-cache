package main

import (
	"log"
	"math/big"
	"os"
	"time"

	"github.com/jonathanbeber/lru-cache/lru"
)

func main() {
	cacheSize := 5
	value := 100000

	cacheLog := log.New(os.Stdout, "[LRUCache] ", log.LstdFlags)
	lruc := lru.New(cacheSize, factorial, cacheLog)

	logger := log.New(os.Stdout, "[Factorial] ", log.LstdFlags)
	logger.Printf("running factorial %d...\n", value)

	t1 := time.Now()
	lruc.Do(value)
	t2 := time.Now()
	diff := t2.Sub(t1)
	logger.Printf("without cache: time=%s\n", diff)

	t1c := time.Now()
	lruc.Do(value)
	t2c := time.Now()
	diffc := t2c.Sub(t1c)
	logger.Printf("with cache: time=%s\n", diffc)

	// removes cache
	for i := 1; i <= cacheSize; i++ {
		lruc.Do(i)
	}

	t1i := time.Now()
	lruc.Do(value)
	t2i := time.Now()
	diffi := t2i.Sub(t1i)
	logger.Printf("after invalidating cache: time=%s\n", diffi)

}

func factorial(x int) *big.Int {
	n := big.NewInt(int64(x))
	return recursiveFactorial(n)
}

func recursiveFactorial(x *big.Int) *big.Int {
	n := big.NewInt(1)
	if x.Cmp(big.NewInt(0)) == 0 {
		return n
	}
	return n.Mul(x, recursiveFactorial(n.Sub(x, n)))
}
