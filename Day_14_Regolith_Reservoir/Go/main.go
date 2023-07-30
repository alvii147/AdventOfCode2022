package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

const (
	// integer representing rock path
	ROCK = 1
	// integer representing sand
	SAND = 2
)

// struct representing coordinates in cave
type Coordinates struct {
	X int
	Y int
}

// map representing cave
// this maps coordinates to an integer representing its contents
type Cave map[Coordinates]int

// get range of cave
// this returns the maximum and minimum xy values of cave
// minimum value of y is set to zero
func (c Cave) GetRange() (Coordinates, Coordinates) {
	// set top left coordinates
	topLeft := Coordinates{
		X: math.MaxInt,
		Y: 0,
	}
	// set bottom right coordinates
	bottomRight := Coordinates{
		X: 0,
		Y: 0,
	}

	// iterate over occupied coordinates and update ranges
	for coordinates := range c {
		if coordinates.X < topLeft.X {
			topLeft.X = coordinates.X
		}

		if coordinates.X > bottomRight.X {
			bottomRight.X = coordinates.X
		}

		if coordinates.Y > bottomRight.Y {
			bottomRight.Y = coordinates.Y
		}
	}

	return topLeft, bottomRight
}

// add rock path between coordinates
func (c Cave) AddRockPath(from Coordinates, to Coordinates) error {
	// get change in xy
	delX := to.X - from.X
	delY := to.Y - from.Y

	// make sure path is horizontal or vertical
	if delX > 0 && delY > 0 {
		return errors.New("Path must be horizontal or vertical")
	}

	// get step in x direction
	xStep := 1
	if delX < 0 {
		xStep = -1
	}

	// get step in y direction
	yStep := 1
	if delY < 0 {
		yStep = -1
	}

	// add rock path
	for x := from.X; x != to.X+xStep; x += xStep {
		for y := from.Y; y != to.Y+yStep; y += yStep {
			c[Coordinates{X: x, Y: y}] = ROCK
		}
	}

	return nil
}

// drop sand unit into cave with no bottom from given x coordinate
// stop when sand unit falls beyond lowest piece of rock
// returns true if sand unit rests
// returns false if sand unit falls beyond lowest piece of rock
func (c Cave) DropSandBottomless(x int) bool {
	// get bottom right to get lowest piece of rock
	_, bottomRight := c.GetRange()
	// set up current coordinates to given x coordinate
	current := Coordinates{
		X: x,
		Y: -1,
	}

	for {
		// stop and return false if sand unit falls beyond lowest piece of rock
		if current.Y > bottomRight.Y {
			return false
		}

		// slice of next possible coordinates for sand unit
		nexts := []Coordinates{
			// directly below
			{
				X: current.X,
				Y: current.Y + 1,
			},
			// below, to the left
			{
				X: current.X - 1,
				Y: current.Y + 1,
			},
			// below, to the right
			{
				X: current.X + 1,
				Y: current.Y + 1,
			},
		}

		// whether or not next available sand unit coordinates are found
		nextSandPositionFound := false
		// iterate over next possible coordinates
		for _, next := range nexts {
			_, ok := c[next]
			// if coordinates are empty, break and continue dropping
			if !ok {
				current = next
				nextSandPositionFound = true
				break
			}
		}

		// if next sand unit coordinates found, continue dropping
		if nextSandPositionFound {
			continue
		}

		// break if all three positions checked
		// but no next available coordinates found
		break
	}

	// place sand unit in current coordinates
	c[current] = SAND

	return true
}

// drop sand unit into cave with bottom from given x coordinate
// stop when sand unit rests
// returns true if sand unit is dropped
// returns false if no more sand units can be dropped
func (c Cave) DropSandWithBottom(x int, bottomY int) bool {
	// if point where sand pours from is occupied, return false
	_, ok := c[Coordinates{
		X: x,
		Y: 0,
	}]
	if ok {
		return false
	}

	// set up current coordinates to given x coordinate
	current := Coordinates{
		X: x,
		Y: -1,
	}

	for {
		// stop dropping if bottom floor reached
		if current.Y >= bottomY-1 {
			break
		}

		// slice of next possible coordinates for sand unit
		nexts := []Coordinates{
			// directly below
			{
				X: current.X,
				Y: current.Y + 1,
			},
			// below, to the left
			{
				X: current.X - 1,
				Y: current.Y + 1,
			},
			// below, to the right
			{
				X: current.X + 1,
				Y: current.Y + 1,
			},
		}

		// whether or not next available sand unit coordinates are found
		nextSandPositionFound := false
		// iterate over next possible coordinates
		for _, next := range nexts {
			_, ok := c[next]
			// if coordinates are empty, break and continue dropping
			if !ok {
				current = next
				nextSandPositionFound = true
				break
			}
		}

		// if next sand unit coordinates found, continue dropping
		if nextSandPositionFound {
			continue
		}

		// break if all three positions checked
		// but no next available coordinates found
		break
	}

	// place sand unit in current coordinates
	c[current] = SAND

	return true
}

// get string of pixels displaying cave contents
func (c Cave) String(bottomless bool, bottomY int) string {
	// get top left and bottom right coordinates
	topLeft, bottomRight := c.GetRange()
	// set lowest point to be bottom floor
	if !bottomless {
		bottomRight.Y = bottomY
	}
	// cave width and height
	caveWidth, caveHeight := bottomRight.X-topLeft.X+2, bottomRight.Y-topLeft.Y+1
	// slice to use for storage of cave pixels
	cavePixels := make([]rune, caveWidth*caveHeight)
	// rune representing empty cave pixel
	emptyPixel := '.'
	// runes representing cave pixels with rocks and sand
	pixelDisplayMap := map[int]rune{
		ROCK: '#',
		SAND: 'o',
	}

	// iterate over pixels
	for y := topLeft.Y; y < bottomRight.Y+1; y++ {
		for x := topLeft.X; x < bottomRight.X+1; x++ {
			// check contents of coordinates
			pixelValue, ok := c[Coordinates{
				X: x,
				Y: y,
			}]
			var pixelRune rune
			// get pixel display rune
			if !ok {
				pixelRune = emptyPixel
			} else {
				pixelRune = pixelDisplayMap[pixelValue]
			}

			cavePixels[((y-topLeft.Y)*(bottomRight.X-topLeft.X+2))+(x-topLeft.X)] = pixelRune
		}

		// add newline rune
		cavePixels[((y-topLeft.Y+1)*(bottomRight.X-topLeft.X+2))-1] = '\n'
	}

	// set bottom floor to rocks
	if !bottomless {
		for x := topLeft.X; x < bottomRight.X+1; x++ {
			cavePixels[((bottomRight.Y-topLeft.Y)*(bottomRight.X-topLeft.X+2))+(x-topLeft.X)] = '#'
		}
	}

	return string(cavePixels)
}

// get line-by-line file scanner
func GetFileScanner(filePath string) (*bufio.Scanner, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file %s, %s", filePath, err.Error())
	}

	scanner := bufio.NewScanner(f)

	return scanner, nil
}

func main() {
	// file with rock paths
	filePath := "../rocks.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Printf("Failed to get file scanner, %s\n", err)
		return
	}

	// create cave maps
	caveBottomless := Cave{}
	caveWithBottom := Cave{}
	// regex for right arrow delimiter
	reRightArrowDelimiter := regexp.MustCompile("\\s*->\\s*")
	// regex for comma delimiter
	reCommaDelimiter := regexp.MustCompile("\\s*,\\s*")

	// scan file line-by-line
	for scanner.Scan() {
		line := scanner.Text()
		// slice of coordinates representing rock paths
		coordinatesRockPath := make([]Coordinates, 0)
		// split by right arrow
		coordinatesArrowSplit := reRightArrowDelimiter.Split(line, -1)

		for _, coordinatesString := range coordinatesArrowSplit {
			// split by commma
			coordinatesCommaSplit := reCommaDelimiter.Split(coordinatesString, -1)
			// make sure exactly two coordinate values are obtained
			if len(coordinatesCommaSplit) != 2 {
				fmt.Printf("Failed 2 coordinate values when parsing \"%s\", got %d\n", coordinatesString, len(coordinatesCommaSplit))
				return
			}

			// convert to integer to get x coordinate
			xCoordinate, err := strconv.Atoi(coordinatesCommaSplit[0])
			if err != nil {
				fmt.Printf("Failed to convert \"%s\" to integer, %s\n", coordinatesCommaSplit[0], err)
				return
			}

			// convert to integer to get y coordinate
			yCoordinate, err := strconv.Atoi(coordinatesCommaSplit[1])
			if err != nil {
				fmt.Printf("Failed to convert \"%s\" to integer, %s\n", coordinatesCommaSplit[1], err)
				return
			}

			// add coordinates to rock path
			coordinatesRockPath = append(
				coordinatesRockPath,
				Coordinates{
					X: xCoordinate,
					Y: yCoordinate,
				},
			)
		}

		// add rock path to cave
		for i := 1; i < len(coordinatesRockPath); i++ {
			fromCoordinates := coordinatesRockPath[i-1]
			toCoordinates := coordinatesRockPath[i]
			caveBottomless.AddRockPath(fromCoordinates, toCoordinates)
			caveWithBottom.AddRockPath(fromCoordinates, toCoordinates)
		}
	}

	// number of sand units dropped
	numSandUnitsDropped := 0

	for {
		// drop sand unit
		rested := caveBottomless.DropSandBottomless(500)
		// break if dropped sand unit does not rest
		if !rested {
			break
		}

		numSandUnitsDropped++
	}

	// print cave after dropping sand
	// fmt.Println(caveBottomless.String(true, 0))
	// print number of sand units dropped
	fmt.Println(numSandUnitsDropped)

	// number of sand units dropped
	numSandUnitsDropped = 0
	// get bottom right point
	_, bottomRight := caveWithBottom.GetRange()
	// get bottom floor y coordinate
	bottomY := bottomRight.Y + 2

	for {
		// drop sand unit
		rested := caveWithBottom.DropSandWithBottom(500, bottomY)
		// break if dropped sand unit does not rest
		if !rested {
			break
		}

		numSandUnitsDropped++
	}

	// print cave after dropping sand
	// fmt.Println(caveWithBottom.String(false, bottomY))
	// print number of sand units dropped
	fmt.Println(numSandUnitsDropped)
}
