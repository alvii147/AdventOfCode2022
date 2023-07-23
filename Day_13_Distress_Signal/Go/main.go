package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// struct representing packet data
// includes either a single value or a slice of packet data
type Packet struct {
	// packet data value
	// if packet includes slice of packet data, this must be -1
	Value int
	// slice of packet data
	// should be nil if value is non-negative
	Elements []Packet
}

// compare two packets
// returns 1 if receiver packet is smaller
// returns -1 if comparison packet is smaller
// returns 0 if packets are equal
func (packet *Packet) Compare(cmpPacket *Packet) int {
	// default return value
	// this is what is returned if one of packet values run out
	defaultComparison := -1
	// shortest length of elements
	maxIter := len(cmpPacket.Elements)

	if len(packet.Elements) < len(cmpPacket.Elements) {
		defaultComparison = 1
		maxIter = len(packet.Elements)
	} else if len(packet.Elements) == len(cmpPacket.Elements) {
		defaultComparison = 0
	}

	// case where both packets are single values
	if packet.Value >= 0 && cmpPacket.Value >= 0 {
		if packet.Value < cmpPacket.Value {
			return 1
		} else if packet.Value == cmpPacket.Value {
			return 0
		}

		return -1
	}

	// case where receiver packet is single value
	// and comparison packet is list of elements
	if packet.Value >= 0 && cmpPacket.Value < 0 {
		// convert receiver packet to list of elements
		packetElement := make([]Packet, 1)
		packetElement[0] = Packet{
			Value: packet.Value,
		}
		newPacket := Packet{
			Value:    -1,
			Elements: packetElement,
		}

		// make recursive call
		return newPacket.Compare(cmpPacket)
	}

	// case where receiver packet is list of elements
	// and comparison packet is single value
	if packet.Value < 0 && cmpPacket.Value >= 0 {
		// convert comparison packet to list of elements
		packetElement := make([]Packet, 1)
		packetElement[0] = Packet{
			Value: cmpPacket.Value,
		}
		newCmpPacket := Packet{
			Value:    -1,
			Elements: packetElement,
		}

		// make recursive call
		return packet.Compare(&newCmpPacket)
	}

	// case where both packets are lists of elements
	// iterate over each element and make recursive call
	for i := 0; i < maxIter; i++ {
		packetElement := packet.Elements[i]
		cmpPacketElement := cmpPacket.Elements[i]

		// recursively compare packet elements
		comparison := packetElement.Compare(&cmpPacketElement)
		// keep comparing if equal
		if comparison == 0 {
			continue
		}

		// return comparison if not equal
		return comparison
	}

	return defaultComparison
}

// parse line of packet data
func ParsePacketData(packetText string) (*Packet, error) {
	// make sure packet text begins with "[" and ends with "]"
	if len(packetText) < 2 || packetText[0] != '[' || packetText[len(packetText)-1] != ']' {
		return nil, errors.New("packetText must begin with \"[\" and end with \"]\"")
	}

	// remove "[" and "]"
	packetText = packetText[1 : len(packetText)-1]

	// indices of commas
	// this are only commas on the top level
	commaIndices := make([]int, 0)
	// stack size of unclosed brackets
	bracketsStack := 0

	// iterate over characters in packet text
	for index, char := range packetText {
		if char == '[' {
			// if character is open bracket, increase size of stack
			bracketsStack++
		} else if char == ']' {
			// if character is close bracket, decrease size of stack
			bracketsStack--
		} else if char == ',' && bracketsStack == 0 {
			// if character is comma and stack is empty, store index of comma
			commaIndices = append(commaIndices, index)
		}
	}

	// slice for holding packet sections split by commas
	packetSectionsTexts := make([]string, 1)
	if len(commaIndices) == 0 {
		// if no commas are found, use whole packet text as packet section
		packetSectionsTexts[0] = packetText
	} else {
		// get first packet section
		packetSectionsTexts[0] = packetText[:commaIndices[0]]

		// iterate over comma indices and collect packet sections
		for i := 1; i < len(commaIndices); i++ {
			packetSectionsTexts = append(packetSectionsTexts, packetText[commaIndices[i-1]+1:commaIndices[i]])
		}

		// get last packet section
		packetSectionsTexts = append(packetSectionsTexts, packetText[commaIndices[len(commaIndices)-1]+1:])
	}

	// root packet
	packet := Packet{
		Value: -1,
	}
	// iterate over packet sections
	for _, packetSectionText := range packetSectionsTexts {
		// skip if packet section is empty
		if len(packetSectionText) == 0 {
			continue
		}

		// attempt to convert packet section to integer
		num, err := strconv.Atoi(packetSectionText)
		if err != nil {
			// if integer conversion fails, process it as a whole packet section using recursion
			packetSection, err := ParsePacketData(packetSectionText)
			if err != nil {
				return nil, fmt.Errorf("Unable to parse \"%s\", %s", packetSectionText, err)
			}

			// add parsed packet section to root packet
			packet.Elements = append(packet.Elements, *packetSection)
		} else {
			// add packet with parsed integer to root packet
			packet.Elements = append(
				packet.Elements,
				Packet{
					Value: num,
				},
			)
		}
	}

	return &packet, nil
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
	// file with signal packets
	filePath := "../signals.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Printf("Failed to get file scanner, %s\n", err)
		return
	}

	// slice of packets
	packets := make([]Packet, 0)
	// current packet index
	packetIndex := 0
	// sum of indices in right order
	rightOrderIndicesSum := 0

	// scan file line-by-line
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty lines
		if len(line) == 0 {
			continue
		}

		// parse packet text
		packet, err := ParsePacketData(line)
		if err != nil {
			fmt.Printf("Error parsing line \"%s\", %s", line, err)
			return
		}
		// store packet in slice
		packets = append(packets, *packet)
		// increment packet index
		packetIndex++

		// if packet pair full, compare packets
		if packetIndex%2 == 0 {
			comparison := packets[len(packets)-2].Compare(&packets[len(packets)-1])
			// if packets in right order, add index of pair to sum of indices
			if comparison == 1 {
				rightOrderIndicesSum += packetIndex / 2
			}
		}
	}

	fmt.Println(rightOrderIndicesSum)

	// sort packets
	sort.Slice(packets, func(i int, j int) bool {
		return packets[i].Compare(&packets[j]) > 0
	})

	// divider packet texts
	divider1PacketText := "[[2]]"
	divider2PacketText := "[[6]]"

	// parse divider packets
	divider1Packet, err := ParsePacketData(divider1PacketText)
	divider2Packet, err := ParsePacketData(divider2PacketText)

	// indices of dividers after sorting
	divider1SortedIndex := -1
	divider2SortedIndex := -1

	// iterate over packets and find sorted indices for dividers
	for index, packet := range packets {
		if divider1SortedIndex < 0 && divider1Packet.Compare(&packet) > 0 {
			// add one to index since indexing begins at 1
			divider1SortedIndex = index + 1
		}

		if divider2SortedIndex < 0 && divider2Packet.Compare(&packet) > 0 {
			// add one to index since indexing begins at 1
			// add another one to index since divider 1 has already been placed
			divider2SortedIndex = index + 2
		}
	}

	fmt.Println(divider1SortedIndex * divider2SortedIndex)
}
