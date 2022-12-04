package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

// check if one range contains another
func CheckRangeContainment(a1 int, a2 int, b1 int, b2 int) bool {
	return (a1 <= b1 && a2 >= b2) || (b1 <= a1 && b2 >= a2)
}

// check if one range overlaps another
func CheckRangeOverlap(a1 int, a2 int, b1 int, b2 int) bool {
	return (a1 <= b2 && a2 >= b1) || (b1 <= a2 && b2 >= a1)
}

func main() {
	// file with camp assignments input
	filePath := "../camp_assignments.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	containsCount := 0
	overlapsCount := 0
	// regex for parsing lines into camp assignments
	r := regexp.MustCompile(`^(\d+)-(\d+),(\d+)-(\d+)$`)
	// scan file line-by-line
	for scanner.Scan() {
		inputLine := scanner.Text()
		// parse input line
		matches := r.FindAllStringSubmatch(inputLine, -1)
		if len(matches) < 1 {
			fmt.Println("failed to parse input line", inputLine)
			return
		}

		// check that 4 values are found representing 2 ranges
		if len(matches[0][1:]) != 4 {
			fmt.Println("expected 4 input values representing 2 ranges, got", len(matches[0][1:]), "values")
		}

		// convert parsed strings to integers
		assignment := make([]int, 4)
		for i := 0; i < 4; i++ {
			j, err := strconv.Atoi(matches[0][i+1])
			if err != nil {
				fmt.Println("failed to convert", j, "to int", err.Error())
				return
			}
			assignment[i] = j
		}

		// check for range containment and increment
		if CheckRangeContainment(
			assignment[0],
			assignment[1],
			assignment[2],
			assignment[3],
		) {
			containsCount += 1
		}

		// check for range overlap and increment
		if CheckRangeOverlap(
			assignment[0],
			assignment[1],
			assignment[2],
			assignment[3],
		) {
			overlapsCount += 1
		}
	}

	fmt.Println(containsCount)
	fmt.Println(overlapsCount)
}
