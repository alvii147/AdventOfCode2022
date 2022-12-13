package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func main() {
	// current cycle
	cycle := 0
	// register value
	X := 1
	// signal strength, i.e. the cycle number multiplied by the value of X
	signalStrength := 0
	// CRT writing position
	crtPosition := 0
	// length of screen
	screenLength := 40
	// length of sprite, must be odd
	spriteLength := 3
	// initialize row of pixels
	pixels := make([]rune, screenLength)
	// file with instructions
	filePath := "../instructions.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	// scan file line-by-line
	for scanner.Scan() {
		inputSplit := strings.Fields(scanner.Text())
		instruction := inputSplit[0]

		var V int
		var nCycles int

		// single cycle instruction where X is unchanged
		if instruction == "noop" {
			V = 0
			nCycles = 1
			// double cycle instruction where X is changed
		} else if instruction == "addx" {
			// convert V from string to int
			V, err = strconv.Atoi(inputSplit[1])
			if err != nil {
				fmt.Printf("failed to convert %s to integer\n", inputSplit[1])
				return
			}
			nCycles = 2
		} else {
			fmt.Printf("encountered unknown instruction %s\n", instruction)
			return
		}

		// iterate over cycles
		for i := 0; i < nCycles; i++ {
			// increment cycle count
			cycle++
			// update signal strength 20th cycle or a multiple of 40 cycles after that
			if cycle == 20 || (cycle-20)%40 == 0 {
				signalStrength += cycle * X
			}

			// check if pixel should be lit or dark
			if crtPosition >= X-(spriteLength/2) && crtPosition <= X+(spriteLength/2) {
				pixels[crtPosition] = '#'
			} else {
				pixels[crtPosition] = '.'
			}

			// increment CRT writing position
			crtPosition++
			// if reached end of pixels row
			if crtPosition == screenLength {
				// print row of pixels
				fmt.Println(string(pixels))
				// reset writing position and row
				crtPosition = 0
			}
		}

		// update register X
		X += V
	}

	fmt.Println(signalStrength)
}
