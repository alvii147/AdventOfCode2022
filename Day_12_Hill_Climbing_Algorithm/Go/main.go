package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"unicode/utf8"
)

// struct representing row and column coordinates
type Coordinates struct {
	// row index
	I int
	// column index
	J int
}

// priority queue item
type QueueItem struct {
	// distance from source to item
	Distance int64
	// index of item in sorted queue
	Index int
	// coordinates of item
	Coor Coordinates
}

// priority queue for breadth first traversal
type Queue []*QueueItem

// get length of priority queue
func (queue Queue) Len() int {
	return len(queue)
}

// comparision priority queue items by index
func (queue Queue) Less(i, j int) bool {
	return queue[i].Distance < queue[j].Distance
}

// swap priority queue items by index
func (queue Queue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
	queue[i].Index = i
	queue[j].Index = j
}

// push new item into priority queue
func (queue *Queue) Push(x any) {
	n := len(*queue)
	queueItem := x.(*QueueItem)
	queueItem.Index = n
	*queue = append(*queue, queueItem)
}

// pop smallest item from priority queue
// in the context of this problem, it is the coordinates with shortest distance
func (queue *Queue) Pop() any {
	oldQueue := *queue
	n := len(oldQueue)
	queueItem := oldQueue[n-1]
	oldQueue[n-1] = nil
	queueItem.Index = -1
	*queue = oldQueue[0 : n-1]

	return queueItem
}

// struct representing elevation grid
type Grid struct {
	// 2d slice representing elevation
	Elevation [][]int
	// length of grid
	L int
	// width of grid
	W int
	// starting coordinates
	Start Coordinates
	// ending coordinates
	End Coordinates
	// slice of coordinates with the lowest elevation
	Lowest []Coordinates
}

// get elevation at given coordinates
func (grid *Grid) E(c Coordinates) int {
	return grid.Elevation[c.I][c.J]
}

// get neighbours of given coordinates
// neighbours are the set of coordinates from which the traverser can move to the given coordinates
func (grid *Grid) GetNeighbours(c Coordinates, visited map[Coordinates]bool) []Coordinates {
	// possible candidates for neighbours
	candidates := []Coordinates{
		// top
		{I: c.I - 1, J: c.J},
		// bottom
		{I: c.I + 1, J: c.J},
		// left
		{I: c.I, J: c.J - 1},
		// right
		{I: c.I, J: c.J + 1},
	}

	// slice of valid neighbours
	neighbours := []Coordinates{}

	for _, candidate := range candidates {
		// not a valid neighbour if indices out of range
		if candidate.I < 0 || candidate.I >= grid.L || candidate.J < 0 || candidate.J >= grid.W {
			continue
		}

		// not a valid neighbour if already visited
		_, ok := visited[candidate]
		if ok {
			continue
		}

		// not a valid neighbour if neighbour's elevation is higher than one level above that of current coordinates
		if grid.E(c)-grid.E(candidate) > 1 {
			continue
		}

		// add to slice of neighbours
		neighbours = append(neighbours, candidate)
	}

	return neighbours
}

// perform Dijkstra's algorithm to compute distances from source to all destinations
// in the current context, the source is the ending coordinates
func (grid *Grid) Dijkstra(source Coordinates) map[Coordinates]int64 {
	// map of visited coordinates
	visited := make(map[Coordinates]bool)
	// set source as visited
	visited[source] = true

	// map of current distances to coordinates
	distances := make(map[Coordinates]int64)
	// priority queue to store coordinates
	queue := make(Queue, 1)
	// add source coordinates to priority queue
	queue[0] = &QueueItem{
		Coor:     source,
		Distance: 0,
		Index:    0,
	}
	heap.Init(&queue)

	// loop until queue is empty
	for queue.Len() > 0 {
		// get next queue item
		queueItem := heap.Pop(&queue).(*QueueItem)
		// if queue item distance is shorter than currently stored distance, update stored distance
		currDistance, ok := distances[queueItem.Coor]
		if !ok || queueItem.Distance < currDistance {
			distances[queueItem.Coor] = queueItem.Distance
		}

		// get neighbouring coordinates
		neighbours := grid.GetNeighbours(queueItem.Coor, visited)
		for _, neighbour := range neighbours {
			// set neighbour as visisted
			visited[neighbour] = true
			// push neighbour to priority queue
			heap.Push(
				&queue,
				&QueueItem{
					Distance: queueItem.Distance + 1,
					Coor:     neighbour,
				},
			)
		}
	}

	return distances
}

// get elevation from character, and whether it's a starting/ending point
func GetElevation(c rune) (int, bool, bool) {
	isStart := false
	isEnd := false

	// 'S' represents starting point and has the same elevation as 'a'
	if c == 'S' {
		isStart = true
		c = 'a'
		// 'E' represents ending point and has the same elevation as 'z'
	} else if c == 'E' {
		isEnd = true
		c = 'z'
	}

	return int(c) - int('a'), isStart, isEnd
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
	// file with elevation data
	filePath := "../elevation.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	// create grid object
	grid := Grid{}
	// initialize 2d slice of elevation
	grid.Elevation = make([][]int, 0)
	// initialize starting, ending, and lowest coordinates
	grid.Start = Coordinates{}
	grid.End = Coordinates{}
	grid.Lowest = make([]Coordinates, 0)

	// row index
	i := 0
	// scan file line-by-line
	for scanner.Scan() {
		rowText := scanner.Text()
		rowLen := utf8.RuneCountInString(rowText)
		// declare slice of row
		rowElevation := make([]int, rowLen)
		// set grid width
		grid.W = rowLen

		for j, c := range rowText {
			// get and set elevation value
			elevation, isStart, isEnd := GetElevation(c)
			rowElevation[j] = elevation

			// set starting/endpoint points
			if isStart {
				grid.Start.I = i
				grid.Start.J = j
			} else if isEnd {
				grid.End.I = i
				grid.End.J = j
			}

			// add to lowest elevations slice
			if elevation == 0 {
				grid.Lowest = append(grid.Lowest, Coordinates{I: i, J: j})
			}
		}

		// add row of elevations to grid
		grid.Elevation = append(grid.Elevation, rowElevation)
		// increment row index
		i++
	}

	// set grid length
	grid.L = i
	// compute distances using Dijkstra's algorithm
	distances := grid.Dijkstra(grid.End)
	// print distance from starting point
	fmt.Println(distances[grid.Start])

	// iterate over computed distances and find shortest distance among lowest elevation coordinates
	shortestDistance := math.MaxInt64
	for _, start := range grid.Lowest {
		distance, ok := distances[start]
		if ok && distance < int64(shortestDistance) {
			shortestDistance = int(distance)
		}
	}
	fmt.Println(shortestDistance)
}
