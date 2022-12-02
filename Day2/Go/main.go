package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// get rock paper scissors score given user's play and opponents's play
// 1 => rock, 2 => paper, 3 => scissors
// score is a sum of user's play and the outcome score
func RockPaperScissorsScore(opponentPlays int, IPlay int) int {
	// scores based on outcome
	outcomeScores := map[string]int{
		"loss": 0,
		"draw": 3,
		"win":  6,
	}

	// difference between user's and opponent's play
	diff := IPlay - opponentPlays
	// if plays are same, it's a draw
	if diff == 0 {
		return IPlay + outcomeScores["draw"]
	}

	// if user is 1 ahead or 2 behind, user wins
	if diff == 1 || diff == -2 {
		return IPlay + outcomeScores["win"]
	}

	// otherwise (if opponent is 1 ahead or 2 behind) opponent wins
	return IPlay + outcomeScores["loss"]
}

// what to play as user given opponent's play and wanted outcome
func WhatToPlay(opponentPlays int, outcome string) int {
	// outcome to how many steps to shift in order to get user's play
	outcomeModShiftMap := map[string]int{
		// X => loss for user, shift by 2 steps, i.e if opponent plays rock(1), play scissors(rock + 2)
		"X": 2,
		// Y => draw, shift by 0 steps, i.e if opponent plays rock(1), play paper(rock + 0)
		"Y": 0,
		// X => win for user, shift by 1 step, i.e if opponent plays rock(1), play paper(rock + 1)
		"Z": 1,
	}

	// compute user's play
	IPlay := ((opponentPlays + outcomeModShiftMap[outcome] - 1) % 3) + 1

	return IPlay
}

// convert play strings to play values
// this assumes opponent plays one of A, B, & C
// and assumes user plays one of X, Y, & Z
func StringToPlayValue(opponentPlaysString string, IPlayString string) (int, int) {
	asciiLowercaseA := 97
	asciiLowercaseX := 120
	// convert play string to lowercase, convert to ascii value, then subtract appropriate ascii value
	opponentPlays := int([]rune(strings.ToLower(opponentPlaysString))[0]) - asciiLowercaseA + 1
	iPlay := int([]rune(strings.ToLower(IPlayString))[0]) - asciiLowercaseX + 1

	return opponentPlays, iPlay
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
	// file with input rock paper scissors data
	filePath := "../rockpaperscissors.txt"

	scanner, err := GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}

	totalScore := 0
	// scan file line by line
	for scanner.Scan() {
		strategy := strings.Fields(scanner.Text())
		// convert play strings to play values
		opponentPlays, IPlay := StringToPlayValue(strategy[0], strategy[1])
		// compute score for round
		totalScore += RockPaperScissorsScore(opponentPlays, IPlay)
	}

	fmt.Println(totalScore)

	totalScore = 0
	// scan file line by line
	scanner, err = GetFileScanner(filePath)
	if err != nil {
		fmt.Println("failed to get file scanner", err.Error())
		return
	}
	for scanner.Scan() {
		strategy := strings.Fields(scanner.Text())
		// convert play strings to play values
		opponentPlays, _ := StringToPlayValue(strategy[0], strategy[1])
		// figure out what user should play
		IPlay := WhatToPlay(opponentPlays, strategy[1])
		// compute score for round
		totalScore += RockPaperScissorsScore(opponentPlays, IPlay)
	}

	fmt.Println(totalScore)
}
