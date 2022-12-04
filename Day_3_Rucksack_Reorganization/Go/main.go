package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

// get line-by-line file scanner
func GetFileScanner(filePath string) (*bufio.Scanner, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s, %s", filePath, err.Error())
	}

	scanner := bufio.NewScanner(f)

	return scanner, nil
}

// get priority of item based on ASCII value
func GetItemPriority(item string) int {
	itemASCII := int([]rune(item)[0])
	aASCII := int('a')
	AASCII := int('A')
	zASCII := int('z')

	// when item is lowercase, priority is given by how many steps ahead of 'a' the item is
	itemPriority := itemASCII - aASCII + 1
	// when item is uppercase
	if itemASCII < aASCII {
		// priority is given by how many steps ahead of 'A' the item is
		// plus the number of letters in the alphabet (i.e. 'z' - 'a' + 1)
		itemPriority = itemASCII - AASCII + 1 + zASCII - aASCII + 1
	}

	return itemPriority
}

// identify duplicate item between two given strings
func GetDuplicateItems(items1 string, items2 string) string {
	// hash map of elements in items1
	items1Map := make(map[string]bool)
	// hash map of duplicate elements
	duplicateItemsMap := make(map[string]bool)

	// loop over items1 and store all items
	for i := 0; i < len(items1); i++ {
		item := string(items1[i])
		items1Map[item] = true
	}

	// loop over items2 and store any duplicate items
	for i := 0; i < len(items2); i++ {
		item := string(items2[i])
		_, duplicate := items1Map[item]
		if duplicate {
			duplicateItemsMap[item] = true
		}
	}

	// create string of duplicate items
	duplicateItems := ""
	for item := range duplicateItemsMap {
		duplicateItems += item
	}

	return duplicateItems
}

func main() {
	// file with input rucksack contents
	filePath := "../rucksacks.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	sumOfPriorities := 0
	// scan file line-by-line
	for scanner.Scan() {
		rucksack := scanner.Text()
		halflen := utf8.RuneCountInString(rucksack) / 2
		// split rucksack items halfway to get compartment contents
		compartment1, compartment2 := rucksack[:halflen], rucksack[halflen:]

		// get duplicate items between two compartments
		duplicateItems := GetDuplicateItems(compartment1, compartment2)
		// return early if not exactly one duplicate is found
		duplicatesCount := utf8.RuneCountInString(duplicateItems)
		if duplicatesCount != 1 {
			fmt.Println(duplicatesCount, "duplicates found in rucksack", rucksack, "expected exactly 1")
			return
		}

		// update sum of item priorities
		sumOfPriorities += GetItemPriority(string(duplicateItems[0]))
	}

	fmt.Println(sumOfPriorities)

	scanner, err = GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	sumOfPriorities = 0
	// scan file n lines at a time
	n := 3
	for true {
		// slice of n rucksacks
		nRucksacks := make([]string, n)
		// condition for breaking out of outer loop
		breakOuterLoop := false

		// loop over and store n rucksack items
		for i := 0; i < n; i++ {
			// break out of inner and outer loops if no more lines are read
			if !scanner.Scan() {
				breakOuterLoop = true
				break
			}

			nRucksacks[i] = scanner.Text()
		}

		if breakOuterLoop {
			break
		}

		duplicateItems := nRucksacks[0]
		// loop over rucksack items and gather duplicate items
		for i := 1; i < n; i++ {
			duplicateItems = GetDuplicateItems(duplicateItems, nRucksacks[i])
		}

		// return early if not exactly one duplicate is found
		duplicatesCount := utf8.RuneCountInString(duplicateItems)
		if duplicatesCount != 1 {
			fmt.Println(duplicatesCount, "duplicates found in rucksacks", nRucksacks, "expected exactly 1")
			return
		}

		sumOfPriorities += GetItemPriority(string(duplicateItems[0]))
	}

	// update sum of item priorities
	fmt.Println(sumOfPriorities)
}
