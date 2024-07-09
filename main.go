package main

import (
	"example/cachd/cache"
	"fmt"
	"io"
	"math"
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
	//	fmt.Printf("Hits: %d\nMisses: %d\n", hits, misses)
	return float64(hits) / float64(hits+misses), nil

}

func showHitRatio(memory cache.ICache, size int, filename string) {
	// wordsize represents the size of the word in bytes
	hits := make([]float64, 5)
	for wordsize := 1; wordsize <= 16; wordsize *= 2 {
		memory.Init(size, wordsize)
		if magicnumber, err := calculateHitRatio(memory, filename); err != nil {
			fmt.Println(err)
		} else {
			hits[int(math.Log2(float64(wordsize)))] = magicnumber
		}
	}
	fmt.Printf("Hits for %v:\n", filename)
	for i := 0; i < 5; i++ {
		fmt.Printf("Wordsize %d:\t", int(math.Pow(2, float64(i))))
		fmt.Printf("Hit ratio is: %f\n", hits[i])
	}

}

func main() {
	// represents a 1KB cache
	size := 1024
	memorySA := cache.SACache{}
	memoryDM := cache.DMCache{}

	// tests memories for different wordsizes
	println("Memory Set Associative")
	showHitRatio(&memorySA, size, "referencia1.dat")
	println("-------------------------------------------------")
	showHitRatio(&memorySA, size, "referencia2.dat")
	println("\n\nMemory Direct Mapped")
	showHitRatio(&memoryDM, size, "referencia1.dat")
	println("-------------------------------------------------")
	showHitRatio(&memoryDM, size, "referencia2.dat")
}
