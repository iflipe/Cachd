package main

import (
	"example/cachd/cache"
	"fmt"
	"io"
	"os"
	"strconv"
)

func calculateHitRatio(memory cache.ICache, filename string) (float64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	address := make([]byte, 5)
	var hits, misses int
	for {
		n, err := f.Read(address)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n >= 4 {
			add, err := strconv.ParseUint(string(address[:4]), 10, 16)
			if err != nil {
				return 0, err
			}
			if memory.Lookup(uint16(add)) {
				hits++
			} else {
				misses++
			}
		}
	}
	fmt.Printf("Hits: %d\nMisses: %d\n", hits, misses)
	return float64(hits) / float64(hits+misses), nil

}

func main() {
	// represents a 1KB cache
	size := 1024

	testCache := cache.DMCache{}
	// wordsize represents the size of the word in bytes
	for wordsize := 1; wordsize <= 16; wordsize *= 2 {
		testCache.Init(size, wordsize)
		fmt.Printf("----------- For wordsize: %d --------------\n", wordsize)
		if magicnumber, err := calculateHitRatio(&testCache, "referencia1.dat"); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Hit ratio is: %f\n", magicnumber)
		}
		if magicnumber, err := calculateHitRatio(&testCache, "referencia2.dat"); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Hit ratio is: %f\n", magicnumber)
		}
	}
}
