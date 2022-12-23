package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/mheidinger/advent-of-code-go/util"
)

//go:embed input.txt
var input string

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

	if part == 1 {
		ans := part1(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

const (
	DirNorth = iota
	DirSouth
	DirWest
	DirEast
)

type Elf struct {
	x     int
	y     int
	nextX int
	nextY int
}

func (elf *Elf) SetNextPos(board [][]*Elf, startDirection int) bool {
	elfNW := board[elf.y-1][elf.x-1]
	elfN := board[elf.y-1][elf.x]
	elfNE := board[elf.y-1][elf.x+1]
	elfE := board[elf.y][elf.x+1]
	elfSE := board[elf.y+1][elf.x+1]
	elfS := board[elf.y+1][elf.x]
	elfSW := board[elf.y+1][elf.x-1]
	elfW := board[elf.y][elf.x-1]

	if elfNW == nil && elfN == nil && elfNE == nil && elfE == nil &&
		elfSE == nil && elfS == nil && elfSW == nil && elfW == nil {
		return false
	}

	direction := startDirection
	for it := 0; it < 4; it++ {
		switch direction {
		case DirNorth:
			if elfNW == nil &&
				elfN == nil &&
				elfNE == nil {
				elf.nextX = elf.x
				elf.nextY = elf.y - 1
				return true
			}
		case DirSouth:
			if elfSW == nil &&
				elfS == nil &&
				elfSE == nil {
				elf.nextX = elf.x
				elf.nextY = elf.y + 1
				return true
			}
		case DirWest:
			if elfNW == nil &&
				elfW == nil &&
				elfSW == nil {
				elf.nextX = elf.x - 1
				elf.nextY = elf.y
				return true
			}
		case DirEast:
			if elfNE == nil &&
				elfE == nil &&
				elfSE == nil {
				elf.nextX = elf.x + 1
				elf.nextY = elf.y
				return true
			}
		}
		direction = (direction + 1) % 4
	}
	return false
}

func (elf *Elf) Move(board [][]*Elf) {
	if elf.nextX == 0 {
		return
	}

	board[elf.y][elf.x] = nil
	board[elf.nextY][elf.nextX] = elf
	elf.x = elf.nextX
	elf.y = elf.nextY
	elf.nextX = 0
	elf.nextY = 0
}

func Simulate(board [][]*Elf, elves []*Elf, startDirection int) bool {
	elfMoved := false
	for _, elf := range elves {
		if elf.SetNextPos(board, startDirection) {
			elfMoved = true
		}
	}

	for _, elf1 := range elves {
		checkX := elf1.nextX
		checkY := elf1.nextY
		for _, elf2 := range elves {
			if elf1 != elf2 && checkX == elf2.nextX && checkY == elf2.nextY {
				elf1.nextX = 0
				elf1.nextY = 0
				elf2.nextX = 0
				elf2.nextY = 0
			}
		}
	}

	for _, elf := range elves {
		elf.Move(board)
	}

	return elfMoved
}

func getFreeFields(board [][]*Elf) int {
	minX, maxX, minY, maxY := 9999, 0, 9999, 0
	for y, row := range board {
		for x, elf := range row {
			if elf == nil {
				continue
			}
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	numFree := 0
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if board[y][x] == nil {
				numFree++
				// fmt.Print(".")
			} else {
				// fmt.Print("#")
			}
		}
		// fmt.Println()
	}
	return numFree
}

func part1(input string) int {
	board, elves := parseInput(input)

	startDirection := DirNorth
	for it := 0; it < 10; it++ {
		Simulate(board, elves, startDirection)
		startDirection = (startDirection + 1) % 4
	}

	return getFreeFields(board)
}

func part2(input string) int {
	board, elves := parseInput(input)

	numRounds := 0
	startDirection := DirNorth
	for {
		numRounds++
		if !Simulate(board, elves, startDirection) {
			break
		}
		startDirection = (startDirection + 1) % 4
	}

	return numRounds
}

func parseInput(input string) (ans [][]*Elf, elves []*Elf) {
	ans = make([][]*Elf, 400)
	for it := range ans {
		ans[it] = make([]*Elf, 400)
	}
	elves = make([]*Elf, 0)

	for y, line := range strings.Split(input, "\n") {
		for x, char := range line {
			if char == '#' {
				elf := &Elf{x: 200 + x, y: 200 + y}
				ans[elf.y][elf.x] = elf
				elves = append(elves, elf)
			}
		}
	}
	return ans, elves
}
