package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const probAlive = 5
const (
	surviveMin = 2
	surviveMax = 3
	reproduceAt = 3
)

const sleepTime = 200 * time.Millisecond

type board [][]bool

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: life BOARD_SIZE ITERATIONS")
		return
	}
	boardSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid size %s.\n", os.Args[1])
		return
	}
	iterations, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid iterations %s.\n", os.Args[2])
		return
	}
	fmt.Printf("Board size is %d.\n", boardSize)
	board := makeBoard(boardSize)
	board.populateRandomly()
	board.run(iterations)
}

func makeBoard(size int) board {
	b := make(board, size)
	for i := range b {
		b[i] = make([]bool, size)
	}
	return b
}

func (b board) populateRandomly() {
	rand.Seed(int64(os.Getpid()))
	for i, row := range b {
		for j := range row {
			if 0 == rand.Intn(probAlive) {
				b[i][j] = true
			}
		}
	}
}

func (b board) run(iterations int) {
	printBoard(b)
	for round := 0; round < iterations; round ++ {
		time.Sleep(sleepTime)
		b = b.update()
		printBoard(b)
	}
}

func printBoard(b board) {
	for _, row := range b {
		for _, state := range row {
			output := "-"
			if state {
				output = "*"
			}
			fmt.Print(output)
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func (b board) update() board {
	newBoard := makeBoard(len(b))
	for rNum, row := range b {
		for cNum := range row {
			currentState := row[cNum]
			count := b.countNeighbors(rNum, cNum)
			newBoard[rNum][cNum] = determineState(currentState, count)
		}
	}
	return newBoard
}

func (b board) countNeighbors(rNum int, cNum int) (count uint) {
	for r := -1; r <= 1; r++ {
		for c := -1; c <= 1; c++ {
			if r == 0 && c == 0 {
				continue // don't count yourself
			}
			if b.isAlive(rNum+r, cNum+c) {
				count++
			}
		}
	}
	return
}

func (b board) isAlive(rNum int, cNum int) bool {
	if rNum < 0 || rNum >= len(b) {
		return false
	}
	if cNum < 0 || cNum >= len(b) {
		return false
	}
	return b[rNum][cNum]
}

func determineState(current bool, count uint) bool {
	return (count == reproduceAt) || (current && count >= surviveMin && count <= surviveMax)
}
