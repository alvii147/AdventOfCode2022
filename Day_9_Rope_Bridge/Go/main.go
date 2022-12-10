package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// get sign of integer
func Sign(x int) int {
	if x == 0 {
		return 0
	}

	return x / int(math.Abs(float64(x)))
}

// get movement direction for following knot given pair of knot coordinates
func FollowKnot(x2 int, y2 int, x1 int, y1 int) (int, int) {
	delx, dely := x2-x1, y2-y1

	// don't move if one step or less away in both dimensions
	if math.Abs(float64(delx)) <= 1 && math.Abs(float64(dely)) <= 1 {
		return 0, 0
	}

	return Sign(delx), Sign(dely)
}

// struct representing coordinates
type Coordinates struct {
	X int
	Y int
}

// struct representing knot in rope
type Knot struct {
	// coordinates of location of knot
	Location Coordinates
	// map of coordinates visited by knot
	Visited map[Coordinates]bool
}

// create new knot starting coordinates
func NewKnot(x int, y int) *Knot {
	knot := &Knot{
		Location: Coordinates{
			X: x,
			Y: y,
		},
		Visited: make(map[Coordinates]bool),
	}

	// add starting coordinates to visited
	knot.Visited[Coordinates{
		X: x,
		Y: y,
	}] = true

	return knot
}

// struct representing rope
type Rope struct {
	// length of rope
	Length int
	// slice of knot pointers
	Knots []*Knot
}

// create new rope of given length
func NewRope(length int) *Rope {
	rope := &Rope{
		Length: length,
		Knots:  make([]*Knot, length),
	}

	// create rope knots
	for i := 0; i < length; i++ {
		rope.Knots[i] = NewKnot(0, 0)
	}

	return rope
}

// move first knot by movement direction and update following knots
func (rope *Rope) Move(delx int, dely int) {
	// following knot coordinates
	var c1 Coordinates
	// preceding knot coordinates
	var c2 Coordinates

	// iterate over knots
	for i := 0; i < rope.Length; i++ {
		// coordinates of current knot
		c1 = rope.Knots[i].Location

		// if not head knot
		if i != 0 {
			// get movement direction based on preceding knot
			delx, dely = FollowKnot(c2.X, c2.Y, c1.X, c1.Y)
		}

		// update knot coordinates
		c1.X += delx
		c1.Y += dely
		rope.Knots[i].Location = c1

		// update visited coordinates for knot
		rope.Knots[i].Visited[c1] = true

		// set current knot coordinates as preceding knot coordinates
		c2 = c1
	}
}

// convert direction character into movement direction
func GetDirection(direction string) (int, int, error) {
	delx, dely := 0, 0

	if strings.ToUpper(direction) == "U" {
		dely = 1
	} else if strings.ToUpper(direction) == "D" {
		dely = -1
	} else if strings.ToUpper(direction) == "R" {
		delx = 1
	} else if strings.ToUpper(direction) == "L" {
		delx = -1
	} else {
		return 0, 0, fmt.Errorf("encountered invalid direction %s", direction)
	}

	return delx, dely, nil
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
	ropeLength := 10
	rope := NewRope(ropeLength)

	// file with rope motions data
	filePath := "../rope_motions.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	// scan file line-by-line
	for scanner.Scan() {
		inputSplit := strings.Fields(scanner.Text())

		// get movement directions and number of steps
		direction := inputSplit[0]
		nSteps, err := strconv.Atoi(inputSplit[1])
		if err != nil {
			fmt.Printf("failed to convert %s into integer, %s\n", inputSplit[1], err.Error())
			return
		}

		// get movement direction
		delx, dely, err := GetDirection(direction)
		if err != nil {
			fmt.Printf("failed to get direction, %s\n", err.Error())
			return
		}

		// move rope step-by-step
		for i := 0; i < nSteps; i++ {
			rope.Move(delx, dely)
		}
	}

	// print number of visited coordinates for tail knot
	fmt.Println(len(rope.Knots[ropeLength-1].Visited))
}
