package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// struct for maintaining top k integers
type TopKNums struct {
	K    int
	Nums []int
}

// initialize top k integers
func (tk *TopKNums) Init(k int) {
	tk.K = k
	tk.Nums = make([]int, tk.K)
	// initialize integers slice to zero
	for i := range tk.Nums {
		tk.Nums[i] = 0
	}
}

// insert new integer
func (tk *TopKNums) Insert(x int) {
	// return early if new integer is less than lowest value
	if x < tk.Nums[tk.K-1] {
		return
	}

	// perform linear scan to find where new integer should be inserted
	insert_idx := -1
	for i := range tk.Nums {
		if tk.Nums[i] > x {
			continue
		}

		insert_idx = i
		break
	}

	// shift slice to the right after insertion index
	copy(tk.Nums[insert_idx+1:tk.K], tk.Nums[insert_idx:tk.K-1])
	// insert new integer
	tk.Nums[insert_idx] = x
}

// compute top k sum
func (tk *TopKNums) Sum() int {
	sum := 0
	for i := range tk.Nums {
		sum += tk.Nums[i]
	}

	return sum
}

func main() {
	// file with input calories data
	filePath := "../calories.txt"
	// open and read file
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("failed to open file", filePath, err.Error())
		return
	}

	// store top 3 calories
	k := 3
	topKCalories := TopKNums{}
	topKCalories.Init(k)
	// calorie count for current elf
	currentElfCalories := 0

	// scan input line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// check if blank line
		if len(strings.TrimSpace(line)) == 0 {
			// insert calorie count
			topKCalories.Insert(currentElfCalories)
			// set calorie count to zero for next elf
			currentElfCalories = 0
			continue
		}

		// convert calories string to int
		calories, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("failed to convert", line, "to integer", err.Error())
			return
		}

		// increment elf calorie count
		currentElfCalories += calories
	}

	// insert calorie count
	topKCalories.Insert(currentElfCalories)
	// compute and print sum of top 3 calories
	topKCaloriesSum := topKCalories.Sum()
	fmt.Println(topKCaloriesSum)
}
