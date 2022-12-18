package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var RegexParser = regexp.MustCompile(
	`^Valve\s+(\D{2})\s+has\s+flow\s+rate\s*=\s*(\d+)\s*;\s*tunnels?\s+leads?\s+to\s+valves?\s+([\D,\s]+)$`,
)

type Valve struct {
	Name              string
	FlowRate          int
	Neighbours        []*Valve
	ShortestDistances map[*Valve]int
	ShortestPaths     map[*Valve][]*Valve
}

type ValveQueueItem struct {
	Distance int
	Path     []*Valve
	Index    int
	V        *Valve
}

type ValveQueue []*ValveQueueItem

func (vq ValveQueue) Len() int {
	return len(vq)
}

func (vq ValveQueue) Less(i, j int) bool {
	return vq[i].Distance < vq[j].Distance
}

func (vq ValveQueue) Swap(i, j int) {
	vq[i], vq[j] = vq[j], vq[i]
	vq[i].Index = i
	vq[j].Index = j
}

func (vq *ValveQueue) Push(x any) {
	n := len(*vq)
	item := x.(*ValveQueueItem)
	item.Index = n
	*vq = append(*vq, item)
}

func (vq *ValveQueue) Pop() any {
	oldVQ := *vq
	n := len(oldVQ)
	item := oldVQ[n-1]
	oldVQ[n-1] = nil
	item.Index = -1
	*vq = oldVQ[0 : n-1]

	return item
}

func (valve *Valve) Dijkstra() {
	if valve.ShortestDistances != nil && valve.ShortestPaths != nil {
		return
	}

	visited := make(map[*Valve]bool)
	visited[valve] = true

	distances := make(map[*Valve]int)
	paths := make(map[*Valve][]*Valve)

	valveQueue := make(ValveQueue, 1)
	valveQueue[0] = &ValveQueueItem{
		Distance: 0,
		Path:     make([]*Valve, 0),
		Index:    0,
		V:        valve,
	}
	heap.Init(&valveQueue)

	for valveQueue.Len() > 0 {
		item := heap.Pop(&valveQueue).(*ValveQueueItem)
		d, ok := distances[item.V]
		if !ok || item.Distance < d {
			distances[item.V] = item.Distance
			paths[item.V] = item.Path
		}

		for _, neighbourValve := range item.V.Neighbours {
			_, isVisited := visited[neighbourValve]
			if isVisited {
				continue
			}

			visited[neighbourValve] = true
			heap.Push(
				&valveQueue,
				&ValveQueueItem{
					Distance: item.Distance + 1,
					Path:     append(item.Path, neighbourValve),
					V:        neighbourValve,
				},
			)
		}
	}

	valve.ShortestDistances = distances
	valve.ShortestPaths = paths
}

type ValveSet map[*Valve]bool

func NewValveSet() ValveSet {
	vs := make(ValveSet)

	return vs
}

func (vs ValveSet) Add(valve *Valve) {
	vs[valve] = true
}

func (vs ValveSet) Remove(valve *Valve) error {
	_, ok := vs[valve]
	if !ok {
		return errors.New("Couldn't find valve to remove")
	}

	delete(vs, valve)

	return nil
}

func (vs ValveSet) Contains(valve *Valve) bool {
	_, ok := vs[valve]

	return ok
}

// insert integer into sorted array and shift elements to the right
func InsertSorted(rewards []int, valves []*Valve, reward int, valve *Valve) {
	// don't insert if new integer is less than lowest value
	n := len(rewards)
	if reward < rewards[n-1] {
		return
	}

	// perform linear scan to find where new integer should be inserted
	insert_idx := -1
	for i := range rewards {
		if rewards[i] > reward {
			continue
		}

		insert_idx = i
		break
	}

	// shift slice to the right after insertion index
	copy(rewards[insert_idx+1:n], rewards[insert_idx:n-1])
	copy(valves[insert_idx+1:n], valves[insert_idx:n-1])
	// insert new integer
	rewards[insert_idx] = reward
	valves[insert_idx] = valve
}

func InformedSearchRecursive(
	currentValve *Valve,
	openedValves ValveSet,
	unopenedValves ValveSet,
	minutesLeft int,
	pressureReleased int,
) int {
	if minutesLeft < 1 {
		return pressureReleased
	}

	if len(unopenedValves) < 1 {
		for valve := range openedValves {
			pressureReleased += (valve.FlowRate * minutesLeft)
		}

		return pressureReleased
	}

	currentValve.Dijkstra()
	distances, paths := currentValve.ShortestDistances, currentValve.ShortestPaths
	rewards := make(map[*Valve]int)
	for valve := range unopenedValves {
		openingCost := distances[valve] + 1
		openingReward := (minutesLeft - openingCost) * valve.FlowRate
		rewards[valve] = openingReward
	}

	if len(rewards) < 1 {
		for valve := range openedValves {
			pressureReleased += (valve.FlowRate * minutesLeft)
		}

		return pressureReleased
	}

	k := 2
	topValves := make([]*Valve, k)
	topValveRewards := make([]int, k)
	for valve, reward := range rewards {
		parentValve := valve
		if valve != currentValve {
			parentValve = paths[valve][0]
		}

		InsertSorted(topValveRewards, topValves, reward, parentValve)
	}

	// tvn := make([]string, k)
	// for i := range tvn {
	// 	if topValves[i] != nil {
	// 		tvn[i] = topValves[i].Name
	// 	}
	// }
	// fmt.Println(currentValve.Name)
	// fmt.Println(tvn)

	for valve := range openedValves {
		pressureReleased += valve.FlowRate
	}

	maxPressureReleased := 0
	for _, nextValve := range topValves {
		if nextValve == nil {
			continue
		}

		if nextValve == currentValve {
			openedValves.Add(currentValve)
			unopenedValves.Remove(currentValve)
		}

		p := InformedSearchRecursive(nextValve, openedValves, unopenedValves, minutesLeft-1, pressureReleased)
		if p > maxPressureReleased {
			maxPressureReleased = p
		}

		if nextValve == currentValve {
			openedValves.Remove(currentValve)
			unopenedValves.Add(currentValve)
		}
	}

	return maxPressureReleased
}

func InformedSearch(startingValve *Valve, valvesMap map[string]*Valve, totalMinutes int) int {
	openedValves := NewValveSet()
	unopenedValves := NewValveSet()
	for _, valve := range valvesMap {
		unopenedValves.Add(valve)
	}

	return InformedSearchRecursive(startingValve, openedValves, unopenedValves, totalMinutes, 0)
}

func DepthFirstSearch(startingValve *Valve, valvesMap map[string]*Valve, totalMinutes int) int {
	openedValves := NewValveSet()
	unopenedValves := NewValveSet()
	for _, valve := range valvesMap {
		unopenedValves.Add(valve)
	}

	return DepthFirstSearchRecursive(startingValve, openedValves, unopenedValves, totalMinutes, 0)
}

func DepthFirstSearchRecursive(
	currentValve *Valve,
	openedValves ValveSet,
	unopenedValves ValveSet,
	minutesLeft int,
	pressureReleased int,
) int {
	depth := 30 - minutesLeft
	for i := 0; i < depth; i++ {
		fmt.Printf(" ")
	}
	fmt.Println(currentValve.Name, len(unopenedValves))
	if minutesLeft < 1 {
		return pressureReleased
	}

	if len(unopenedValves) < 1 {
		for valve := range openedValves {
			pressureReleased += (valve.FlowRate * minutesLeft)
		}

		return pressureReleased
	}

	for valve := range openedValves {
		pressureReleased += valve.FlowRate
	}

	maxPressureReleased := 0

	if unopenedValves.Contains(currentValve) {
		openedValves.Add(currentValve)
		unopenedValves.Remove(currentValve)
		p := DepthFirstSearchRecursive(currentValve, openedValves, unopenedValves, minutesLeft-1, pressureReleased)
		if p > maxPressureReleased {
			maxPressureReleased = p
		}
		openedValves.Remove(currentValve)
		unopenedValves.Add(currentValve)
	}

	for _, valve := range currentValve.Neighbours {
		p := DepthFirstSearchRecursive(valve, openedValves, unopenedValves, minutesLeft-1, pressureReleased)
		if p > maxPressureReleased {
			maxPressureReleased = p
		}
	}

	return maxPressureReleased
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
	// file with valves data
	filePath := "../valves.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	valvesMap := make(map[string]*Valve)
	tempNeighbouringValvesMap := make(map[string][]string)

	for scanner.Scan() {
		line := scanner.Text()
		result := RegexParser.FindStringSubmatch(line)

		name := strings.TrimSpace(result[1])
		flowRate, err := strconv.Atoi(result[2])
		if err != nil {
			fmt.Printf("failed to parse %s into integer\n", result[1])
			return
		}

		valve := &Valve{
			Name:     name,
			FlowRate: flowRate,
		}
		valvesMap[name] = valve

		neighbours := strings.Split(result[3], ",")
		for i := range neighbours {
			neighbours[i] = strings.TrimSpace(neighbours[i])
		}
		tempNeighbouringValvesMap[name] = neighbours
	}

	for name, valve := range valvesMap {
		neighbours := tempNeighbouringValvesMap[name]
		valve.Neighbours = make([]*Valve, len(neighbours))
		for i, neighbour := range neighbours {
			valve.Neighbours[i] = valvesMap[neighbour]
		}
	}

	fmt.Println(InformedSearch(valvesMap["AA"], valvesMap, 30))
	// fmt.Println(DepthFirstSearch(valvesMap["AA"], valvesMap, 30))
}
