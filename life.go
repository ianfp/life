// A Go implementation of Conway's Game of Life.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	// The first argument is the size of the board.
	argBoardSize = 1
	// The second argument is the number of rounds to play the game.
	argIterations = 2
)

// The chance that a given cell starts the game alive is 1 / `probAlive`.
const probAlive = 5
const (
	// Fewer neighbors than this and a cell dies of lonliness.
	surviveMin = 2
	// More neighbors than this and a cell dies of overcrowding.
	surviveMax = 3
	// If a dead cell has this many living neighbors, it comes back to life.
	reproduceAt = 3
)

// How long to sleep between rounds.
const sleepTime = 200 * time.Millisecond

// The game board: true means alive, false means dead.
type board [][]bool

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: life BOARD_SIZE ITERATIONS")
		return
	}
	boardSize, err := strconv.Atoi(os.Args[argBoardSize])
	if err != nil {
		fmt.Printf("Invalid size %s.\n", os.Args[argBoardSize])
		return
	}
	iterations, err := strconv.Atoi(os.Args[argIterations])
	if err != nil {
		fmt.Printf("Invalid iterations %s.\n", os.Args[argIterations])
		return
	}
	fmt.Printf("Board size is %d.\n", boardSize)
	board := makeBoard(boardSize)
	board.populateRandomly()
	board.run(iterations)
}

// Make a new game board.
func makeBoard(size int) board {
	b := make(board, size)
	for i := range b {
		b[i] = make([]bool, size)
	}
	return b
}

// Populate a new game board randomly according to the `probAlive` setting.
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

// Run the game.
func (b board) run(iterations int) {
	printBoard(b)
	for round := 0; round < iterations; round++ {
		time.Sleep(sleepTime)
		b = b.update()
		printBoard(b)
	}
}

// Print the current state of the board to STDOUT.
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

// Returns a new board that represents the state of the game on the next round.
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

// Returns the number of living neighbors a given cell has.
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

// Whether a cell is alive in the current state of the board.
func (b board) isAlive(rNum int, cNum int) bool {
	if rNum < 0 || rNum >= len(b) {
		return false
	}
	if cNum < 0 || cNum >= len(b) {
		return false
	}
	return b[rNum][cNum]
}

// Given the state of a cell and the number of living neighbors,
// return whether it survives the current round.
func determineState(current bool, count uint) bool {
	return (count == reproduceAt) || (current && count >= surviveMin && count <= surviveMax)
}
