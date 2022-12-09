package main

import (
	"bufio"
	"fmt"
	"os"
)

// determine whether a single tree is visible from at least one side
func IsVisible(trees [][]int, i int, j int) bool {
	// trees grid dimensions
	l := len(trees)
	w := len(trees[0])

	visible := true
	// check for visibility from the left
	for k := 0; k < j; k++ {
		if trees[i][k] >= trees[i][j] {
			visible = false
			break
		}
	}

	if visible {
		return true
	}

	visible = true
	// check for visibility from the right
	for k := w - 1; k > j; k-- {
		if trees[i][k] >= trees[i][j] {
			visible = false
			break
		}
	}

	if visible {
		return true
	}

	visible = true
	// check for visibility from the top
	for k := 0; k < i; k++ {
		if trees[k][j] >= trees[i][j] {
			visible = false
			break
		}
	}

	if visible {
		return true
	}

	visible = true
	// check for visibility from the bottom
	for k := l - 1; k > i; k-- {
		if trees[k][j] >= trees[i][j] {
			visible = false
			break
		}
	}

	return visible
}

// get scenic score of a single tree
func GetScenicScore(trees [][]int, i int, j int) int {
	// trees grid dimensions
	l := len(trees)
	w := len(trees[0])

	scenicScore := 1
	viewingDistance := 0

	// get viewing distance on the right
	for k := j + 1; k < w; k++ {
		viewingDistance++
		if trees[i][k] >= trees[i][j] {
			break
		}
	}

	scenicScore *= viewingDistance
	viewingDistance = 0

	// get viewing distance on the left
	for k := j - 1; k >= 0; k-- {
		viewingDistance++
		if trees[i][k] >= trees[i][j] {
			break
		}
	}

	scenicScore *= viewingDistance
	viewingDistance = 0

	// get viewing distance on the bottom
	for k := i + 1; k < l; k++ {
		viewingDistance++
		if trees[k][j] >= trees[i][j] {
			break
		}
	}

	scenicScore *= viewingDistance
	viewingDistance = 0

	// get viewing distance on the top
	for k := i - 1; k >= 0; k-- {
		viewingDistance++
		if trees[k][j] >= trees[i][j] {
			break
		}
	}

	scenicScore *= viewingDistance

	return scenicScore
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
	// file with trees data
	filePath := "../trees.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	// ASCII value of '0' for converting runes to digits
	zeroASCII := int('0')
	// 2-dimensional trees grid
	trees := make([][]int, 0)

	// scan file line-by-line and construct trees grid
	lineIdx := 0
	for scanner.Scan() {
		input := scanner.Text()
		trees = append(trees, make([]int, len(input)))

		for charIdx := range input {
			trees[lineIdx][charIdx] = int(input[charIdx]) - zeroASCII
		}
		lineIdx += 1
	}

	// trees grid dimensions
	l := len(trees)
	w := len(trees[0])

	nonBorderingVisibleCount := 0
	maxScenicScore := 0

	// iterate over non-bordering trees
	for i := 1; i < l-1; i++ {
		for j := 1; j < w-1; j++ {
			// check if tree is visible and increment count
			if IsVisible(trees, i, j) {
				nonBorderingVisibleCount++
			}

			// get scenic score of tree and update max if appropriate
			scenicScore := GetScenicScore(trees, i, j)
			if maxScenicScore < scenicScore {
				maxScenicScore = scenicScore
			}
		}
	}

	// number of trees along the outer border
	// these trees are always visible
	borderingTreesCount := 2 * (l + w - 2)

	fmt.Println(nonBorderingVisibleCount + borderingTreesCount)
	fmt.Println(maxScenicScore)
}
