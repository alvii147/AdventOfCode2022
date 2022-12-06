package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// regular expression for parsing a line of crates
// expects format "[X] [Y]     [Z]"
var cratesParserRegex *regexp.Regexp = regexp.MustCompile(`(?:(?:\[(\S)\]\s?)|(?:(\s{4})))`)

// regular expression for parsing move
// expects format "move x from y to z"
var moveParserRegex *regexp.Regexp = regexp.MustCompile(`^\s*move\s*(\d+)\s*from\s*(\d*)\s*to\s*(\d*)\s*$`)

// parse single line of crates using regular expressions
func ParseLineOfCrates(cratesLine string, nStacks int) []string {
	// get regex matches and groups
	matches := cratesParserRegex.FindAllStringSubmatch(cratesLine, nStacks)
	// initialize crates
	crates := make([]string, nStacks)
	for i := 0; i < nStacks; i++ {
		// no crates are represented by empty strings
		crates[i] = ""
	}

	// store crates from regex matches
	for i := 0; i < len(matches); i++ {
		crates[i] = strings.TrimSpace(matches[i][1])
	}

	return crates
}

// parse move from single move input line
func ParseMove(moveLine string) (int, int, int, error) {
	// get regex match and groups
	match := moveParserRegex.FindStringSubmatch(moveLine)

	// get number of crates to move
	nCrates, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse %s", match[1])
	}

	// get index of source stack
	fromIdx, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse %s", match[2])
	}
	// decrement once to make zero-indexed
	fromIdx--

	// get index of destination stack
	toIdx, err := strconv.Atoi(match[3])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse %s", match[3])
	}
	// decrement once to make zero-indexed
	toIdx--

	return nCrates, fromIdx, toIdx, nil
}

// reverse generic slice
func ReverseSlice[T any](s []T) {
	// swap using left and right iterators
	for l, r := 0, len(s)-1; l < r; l, r = l+1, r-1 {
		s[l], s[r] = s[r], s[l]
	}
}

// transfer multiple crates from one stack to another
// retainOrder indicates whether or not to move crates in original or reverse order
func TransferCrates(fromStack []string, toStack []string, nCrates int, retainOrder bool) ([]string, []string) {
	// create copy of source stack
	fromStackCopy := make([]string, len(fromStack))
	for i := 0; i < len(fromStack); i++ {
		fromStackCopy[i] = fromStack[i]
	}
	// create copy of destination stack
	toStackCopy := make([]string, len(toStack))
	for i := 0; i < len(toStack); i++ {
		toStackCopy[i] = toStack[i]
	}

	// get stack partition to move
	stackToMove := fromStackCopy[len(fromStackCopy)-nCrates:]
	// if retainOrder is true, move in order
	// if retainOrder is false, move in reverse order, much like stack data structure
	if !retainOrder {
		ReverseSlice(stackToMove)
	}

	// delete last nCrates crates from source stack
	fromStackCopy = fromStackCopy[:len(fromStackCopy)-nCrates]
	// copy over to destination stack
	toStackCopy = append(toStackCopy, stackToMove...)

	return fromStackCopy, toStackCopy
}

// get line-by-line file scanner
func GetFileScanner(filePath string) (*bufio.Scanner, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s, %s", filePath, err.Error())
	}

	scanner := bufio.NewScanner(f)

	return scanner, nil
}

func main() {
	// file with input crates stacks data
	filePath := "../crates.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	// lines of file input
	inputLines := make([]string, 0)
	// line iterator index
	lineIdx := 0
	// line index of blank line separator between stacks and moves data
	dividerLineIdx := -1
	// scan file line-by-line
	for scanner.Scan() {
		line := scanner.Text()
		// store file input line
		inputLines = append(inputLines, line)

		// store divider line index if blank line found
		if strings.TrimSpace(line) == "" {
			dividerLineIdx = lineIdx
		}

		// increment line index
		lineIdx++
	}

	if dividerLineIdx < 0 {
		fmt.Println("failed to find division between stacks and moves data in input file")
		return
	}

	// get number of stacks
	nStacks := len(strings.Fields(inputLines[dividerLineIdx-1]))

	// initialize nest slice representing stacks
	stacks1 := make([][]string, nStacks)
	for i := 0; i < nStacks; i++ {
		stacks1[i] = make([]string, 0)
	}

	// iterate over crate lines from bottom up
	for i := dividerLineIdx - 2; i >= 0; i-- {
		// parse and store crates line in stacks
		crates := ParseLineOfCrates(inputLines[i], nStacks)
		for j, crate := range crates {
			if crate == "" {
				continue
			}

			stacks1[j] = append(stacks1[j], crate)
		}
	}

	// create copy of stacks for part 2 of puzzle
	stacks2 := make([][]string, nStacks)
	for i := 0; i < nStacks; i++ {
		stacks2[i] = make([]string, len(stacks1[i]))
		for j := 0; j < len(stacks1[i]); j++ {
			stacks2[i][j] = stacks1[i][j]
		}
	}

	// iterate over moves
	for i := dividerLineIdx + 1; i < len(inputLines); i++ {
		// parse current move
		nCrates, fromIdx, toIdx, err := ParseMove(inputLines[i])
		if err != nil {
			fmt.Println("failed to parse move", inputLines[i], err.Error())
			return
		}

		// transfer appropriate crates between stacks
		stacks1[fromIdx], stacks1[toIdx] = TransferCrates(stacks1[fromIdx], stacks1[toIdx], nCrates, false)
		stacks2[fromIdx], stacks2[toIdx] = TransferCrates(stacks2[fromIdx], stacks2[toIdx], nCrates, true)
	}

	// gather top crates in each stack
	topCrates1 := make([]string, nStacks)
	topCrates2 := make([]string, nStacks)
	for i := 0; i < nStacks; i++ {
		topCrates1[i] = stacks1[i][len(stacks1[i])-1]
		topCrates2[i] += stacks2[i][len(stacks2[i])-1]
	}

	// concatenate top crates
	fmt.Println(strings.Join(topCrates1, ""))
	fmt.Println(strings.Join(topCrates2, ""))
}
