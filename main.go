package main

import (
	"example/cachd/cache"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

// This function opens and reads the contents of a file, retrieves a 5-character line, extracts the 4-digit memory address
// present in each line of the file and then tries to convert the addread read into a 16-bit integer.
// If all is done successfully stores the number of hits and misses in the cache passed as parameter of the addresses
// contained within the file.
// Finally, the function returns the percentage of hits or an error if anything goes wrong in the process.
func calculateHitRatio(memory cache.ICache, filename string) (float64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()         //Closes the file after the function ends
	line := make([]byte, 5) //Gets the first 5 characters of the line (4 for the address and 1 for the '\n' character)
	var hits, misses int
	for {
		n, err := f.Read(line)
		//If the end of the file is reached, the loop is broken
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n >= 4 {
			//Tries to convert the 4-digit address into a 64-bit integer
			address, err := strconv.ParseUint(string(line[:4]), 10, 16)
			if err != nil {
				return 0, err
			}
			//Looks up the address in the cache and increments the hits or misses counter based on the return value
			if memory.Lookup(uint16(address)) {
				hits++
			} else {
				misses++
			}
		}
	}
	return float64(hits) / float64(hits+misses), nil

}

func showHitRatio(memory cache.ICache, size int, filename string) {
	//Wordsize in bytes
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
		fmt.Printf("Hit ratio is: %.2f%%\n", hits[i]*100)
	}

}

/*
For the purposes of this experiment, the actual information stored in the memory isn't important
so the functions signatures and behaviors aren't fully designed and implemented to deal with the
treatment of such data
*/
func main() {
	//Represents a 1KB cache
	memorySize := 1024

	//Creates the two types of cache
	memorySA := cache.SACache{}
	memoryDM := cache.DMCache{}

	//T'ests memories for different wordsizes
	println("Memory Set Associative")
	showHitRatio(&memorySA, memorySize, "referencia1.dat")
	println("-------------------------------------------------")
	showHitRatio(&memorySA, memorySize, "referencia2.dat")
	println("\n\nMemory Direct Mapped")
	showHitRatio(&memoryDM, memorySize, "referencia1.dat")
	println("-------------------------------------------------")
	showHitRatio(&memoryDM, memorySize, "referencia2.dat")
}
