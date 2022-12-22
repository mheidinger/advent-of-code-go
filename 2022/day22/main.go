package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/mheidinger/advent-of-code-go/cast"
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
	DirRight = iota
	DirDown
	DirLeft
	DirUp
)

type Position struct {
	X         int
	Y         int
	Direction int
}

func execStep(board [][]rune, step string, pos Position) Position {
	switch step {
	case "R":
		pos.Direction = (pos.Direction + 1) % 4
	case "L":
		pos.Direction = (pos.Direction - 1) % 4
		if pos.Direction < 0 {
			pos.Direction = 4 + pos.Direction
		}
	default:
		numSteps := cast.ToInt(step)
		pos = walkSteps(board, numSteps, pos)
	}

	return pos
}

func walkSteps(board [][]rune, numSteps int, pos Position) Position {
	for it := 0; it < numSteps; it++ {
		newPos := getNewPos(board, pos)
		if board[newPos.Y][newPos.X] == '#' {
			break
		}
		pos = newPos
	}
	return pos
}

func getNewPos(board [][]rune, pos Position) Position {
	newPos := pos
	switch pos.Direction {
	case DirRight:
		newPos.X++
	case DirDown:
		newPos.Y++
	case DirLeft:
		newPos.X--
	case DirUp:
		newPos.Y--
	}
	if newPos.Y < 0 || newPos.Y >= len(board) ||
		newPos.X < 0 || newPos.X >= len(board[newPos.Y]) ||
		board[newPos.Y][newPos.X] == ' ' {
		// wrap around
		newPos = getWrapAroundPos(board, pos)
	}

	return newPos
}

func getWrapAroundPos(board [][]rune, pos Position) Position {
	newPos := pos
	for {
		if pos.Y < 0 || pos.Y >= len(board) ||
			pos.X < 0 || pos.X >= len(board[pos.Y]) ||
			board[pos.Y][pos.X] == ' ' {
			break
		}
		newPos = pos

		switch pos.Direction {
		case DirRight:
			pos.X -= 1
		case DirDown:
			pos.Y -= 1
		case DirLeft:
			pos.X += 1
		case DirUp:
			pos.Y += 1
		}
	}
	return newPos
}

func part1(input string) int {
	board, path := parseInput(input)

	pos := Position{
		X:         strings.Index(string(board[0]), "."),
		Y:         0,
		Direction: DirRight,
	}

	numCache := ""
	for _, step := range path {
		if step >= '0' && step <= '9' {
			numCache += string(step)
			continue
		}

		if len(numCache) > 0 {
			pos = execStep(board, numCache, pos)
			numCache = ""
		}
		pos = execStep(board, string(step), pos)
	}

	return ((pos.Y + 1) * 1000) + ((pos.X + 1) * 4) + pos.Direction
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (ans [][]rune, path string) {
	lines := strings.Split(input, "\n")
	boardStr := lines[:len(lines)-2]
	board := make([][]rune, 0, len(boardStr))
	for _, str := range boardStr {
		board = append(board, []rune(str))
	}
	return board, lines[len(lines)-1]
}
