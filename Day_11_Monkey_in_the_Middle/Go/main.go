package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// regex pattern for parsing starting items
var ReStartingItems = regexp.MustCompile(`^\s*Starting items:\s*(.+)$`)

// regex pattern for parsing operation
var ReOperation = regexp.MustCompile(`^\s*Operation:\s*new\s*=\s*old\s*(\*|\+)\s*(.+)\s*$`)

// regex pattern for parsing divisor
var ReDivisor = regexp.MustCompile(`^\s*Test:\s*divisible\s*by\s*(\d+)$`)

// regex pattern for parsing monkey to throw to if divisible
var ReDivisibleThrowTo = regexp.MustCompile(`^\s*If\s*true:\s*throw\s*to\s*monkey\s*(\d+)$`)

// regex pattern for parsing monkey to throw to if indivisible
var ReIndivisibleThrowTo = regexp.MustCompile(`^\s*If\s*false:\s*throw\s*to\s*monkey\s*(\d+)$`)

// function type for basic math operators
type OperatorFunc func(a int, b int) int

// operator character to operator function map
var OperatorFuncMap = map[string]OperatorFunc{
	"+": func(a int, b int) int { return a + b },
	"-": func(a int, b int) int { return a - b },
	"*": func(a int, b int) int { return a * b },
	"/": func(a int, b int) int { return a / b },
}

// function type for operations on worry value
type OperationFunc func(old int) int

// struct representing monkey
type Monkey struct {
	// worry values for items
	Items []int
	// operation to perform on worry value
	Operation OperationFunc
	// divisor to check divisibility for on each item
	Divisor int
	// number of inspected items
	InspectionCount int
}

// inspect current items for monkey and throw to other monkeys
func (monkey *Monkey) Inspect(worryManager OperationFunc, monkeyThrowToDivisible *Monkey, monkeyThrowToIndivisible *Monkey) {
	for _, worryValue := range monkey.Items {
		// perform operation on item worry level
		newWorryValue := monkey.Operation(worryValue)
		// perform worry management operation on item worry level
		newWorryValue = worryManager(newWorryValue)

		// check if divisible by divisor and through to another monkey
		if newWorryValue%monkey.Divisor == 0 {
			monkeyThrowToDivisible.Items = append(monkeyThrowToDivisible.Items, newWorryValue)
		} else {
			monkeyThrowToIndivisible.Items = append(monkeyThrowToIndivisible.Items, newWorryValue)
		}

		// increment inspection count for monkey
		monkey.InspectionCount++
	}

	// clear items for monkey
	monkey.Items = monkey.Items[:0]
}

// compute greatest common divisor of two integers
func gcd(a int, b int) int {
	if a == 0 {
		return b
	}

	if b == 0 {
		return a
	}

	return gcd(b, a%b)
}

// compute greatest common divisor of slice of integers
func GCD(a []int) int {
	v := a[0]
	for i := 1; i < len(a); i++ {
		v = gcd(v, a[i])
	}

	return v
}

// compute lowest common multiple of two integers
func lcm(a int, b int) int {
	return (a * b) / gcd(a, b)
}

// compute lowest common multiple of slice of integers
func LCM(a []int) int {
	v := a[0]
	for i := 1; i < len(a); i++ {
		v = lcm(v, a[i])
	}

	return v
}

// insert integer into sorted array and shift elements to the right
func InsertSorted(a []int, v int) {
	// don't insert if new integer is less than lowest value
	n := len(a)
	if v < a[n-1] {
		return
	}

	// perform linear scan to find where new integer should be inserted
	insert_idx := -1
	for i := range a {
		if a[i] > v {
			continue
		}

		insert_idx = i
		break
	}

	// shift slice to the right after insertion index
	copy(a[insert_idx+1:n], a[insert_idx:n-1])
	// insert new integer
	a[insert_idx] = v
}

func main() {
	// file with monkeys input data
	filePath := "../monkeys.txt"

	// read input file into bytes slice
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read file %s, %s\n", filePath, err.Error())
		return
	}

	// get file contents string
	fileString := string(fileBytes)

	// unparsed input for each monkey
	monkeyInputs := strings.Split(fileString, "\n\n")
	// number of monkeys
	nMonkeys := len(monkeyInputs)
	// initialize slice of monkeys
	monkeys := make([]*Monkey, nMonkeys)
	// initialize slice of maps with indexes of monkeys to throw to
	throwToMap := make([]map[bool]int, nMonkeys)
	// initialize slice of divisors
	divisors := make([]int, nMonkeys)

	// iterate over monkey input
	for i, monkeyInput := range monkeyInputs {
		// create monkey struct
		monkey := &Monkey{}
		// set inspection count to zero
		monkey.InspectionCount = 0

		inputLines := strings.Split(monkeyInput, "\n")

		// parse and set starting items for monkey
		result := ReStartingItems.FindStringSubmatch(inputLines[1])
		itemsStrings := strings.Split((result[1]), ",")

		monkey.Items = make([]int, len(itemsStrings))
		for i := range itemsStrings {
			monkey.Items[i], err = strconv.Atoi(strings.TrimSpace(itemsStrings[i]))
			if err != nil {
				fmt.Printf("failed to parse %s into integer, %s\n", itemsStrings[i], err.Error())
				return
			}
		}

		// parse operation for monkey
		result = ReOperation.FindStringSubmatch(inputLines[2])
		operator := result[1]
		operand := result[2]

		// get operator function based on operator symbol
		operatorFunc, ok := OperatorFuncMap[operator]
		if !ok {
			fmt.Printf("unknown operator %s encountered\n", operator)
			return
		}

		// if operand is self
		if operand == "old" {
			// set operation function with self
			monkey.Operation = func(old int) int {
				return operatorFunc(old, old)
			}
			// otherwise expect operand to be integer
		} else {
			operatorValue, err := strconv.Atoi(operand)
			if err != nil {
				fmt.Printf("failed to parse %s into integer\n", operand)
				return
			}

			// set operation function with parsed operand
			monkey.Operation = func(old int) int {
				return operatorFunc(old, operatorValue)
			}
		}

		// parse and set divisor
		result = ReDivisor.FindStringSubmatch(inputLines[3])
		divisor, err := strconv.Atoi(result[1])
		if err != nil {
			fmt.Printf("failed to parse %s into integer\n", result[1])
			return
		}
		monkey.Divisor = divisor
		divisors[i] = divisor

		// set monkey in slice of monkeys
		monkeys[i] = monkey

		// initialize throw to map
		throwToMap[i] = make(map[bool]int)

		// parse and set monkey index to throw to if divisible
		result = ReDivisibleThrowTo.FindStringSubmatch(inputLines[4])
		divisibleThrowTo, err := strconv.Atoi(result[1])
		if err != nil {
			fmt.Printf("failed to parse %s into integer\n", result[1])
			return
		}
		throwToMap[i][true] = divisibleThrowTo

		// parse and set monkey index to throw to if indivisible
		result = ReIndivisibleThrowTo.FindStringSubmatch(inputLines[5])
		indivisibleThrowTo, err := strconv.Atoi(result[1])
		if err != nil {
			fmt.Printf("failed to parse %s into integer\n", result[1])
			return
		}
		throwToMap[i][false] = indivisibleThrowTo
	}

	// get lcm of divisors
	divisorsLCM := LCM(divisors)
	// worryManager := func(old int) int {
	// 	return old / 3
	// }
	worryManager := func(old int) int {
		return old % divisorsLCM
	}

	// perform rounds of inspections
	// nRounds := 20
	nRounds := 10000
	for i := 0; i < nRounds; i++ {
		for j := range monkeys {
			monkeys[j].Inspect(
				worryManager,
				monkeys[throwToMap[j][true]],
				monkeys[throwToMap[j][false]],
			)
		}
	}

	// get monkey business value from top-2 inspection counts
	inspectionCounts := make([]int, 2)
	for j := range monkeys {
		InsertSorted(inspectionCounts, monkeys[j].InspectionCount)
	}

	fmt.Println(inspectionCounts[0] * inspectionCounts[1])
}
