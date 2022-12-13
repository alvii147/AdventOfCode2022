package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// map for maintaining a hash of occurrences of runes in a string
type RuneCounterMap map[rune]int

// add a rune to counter and return the number of newly added runes
func (m RuneCounterMap) addToCounterMap(k rune) int {
	// get current count value for given key
	// default is 0 if key does not exist
	count, ok := m[k]
	if !ok {
		count = 0
	}

	// increment count
	m[k] = count + 1

	// determine if a unique rune has been added
	addedUnique := 0
	if count == 0 {
		addedUnique = 1
	}

	return addedUnique
}

// remove a rune to counter and return the number of newly remove runes
func (m RuneCounterMap) removeFromCounterMap(k rune) (int, error) {
	// get current count value for given key
	// return error if key does not exist
	count, ok := m[k]
	if !ok {
		return 0, errors.New("cannot remove non-existent key")
	}

	// decrement count
	m[k] = count - 1

	// determine if a unique rune has been removed
	removedUnique := 0
	if count == 1 {
		removedUnique = 1
	}

	return removedUnique, nil
}

func main() {
	// file with input signal
	filePath := "../signal.txt"
	// read input file into bytes slice
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("failed to read file", filePath, err.Error())
		return
	}

	// get runes slice from file contents
	signal := []rune(strings.TrimSpace(string(fileBytes)))
	// initialize rune counter map
	counter := make(RuneCounterMap, 0)
	// current count of unique runes
	uniqueCount := 0
	// number of unique runes to look for
	windowSize := 14
	// rune iterator index
	i := 0

	// add first runes to counter map to fill window size
	for i < windowSize {
		r := signal[i]
		uniqueCount += counter.addToCounterMap(r)
		i++
	}

	// slide window along runes
	for i < len(signal)+1 {
		// end loop if desired number of unique consecutive runes found
		if uniqueCount == windowSize {
			fmt.Println(i)
			break
		}

		// throw error if unable to find unique consecutive runes
		if i >= len(signal) {
			fmt.Println("unable to find", windowSize, "unique consecutive runes")
			return
		}

		// add rune to the right of the window
		r := signal[i]
		uniqueCount += counter.addToCounterMap(r)

		// remove rune from the left of the window
		r = signal[i-windowSize]
		removedUnique, err := counter.removeFromCounterMap(r)
		if err != nil {
			fmt.Println("error removing", string(r), "from counter map", err.Error())
			return
		}
		uniqueCount -= removedUnique

		i++
	}
}
