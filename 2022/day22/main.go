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
	switch newPos.Direction {
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
		newPos = getWrapAroundPosCube(board, pos)
	}

	return newPos
}

func translate(min1, anchor2, val int, invert bool) int {
	diff := val - min1
	if !invert {
		return anchor2 + diff
	}
	return anchor2 - diff
}

func getWrapAroundPosCube(board [][]rune, pos Position) Position {
	if pos.Y == 0 && pos.X < 100 {
		// Edge 4a
		pos.Direction = DirRight
		pos.Y = translate(50, 150, pos.X, false)
		pos.X = 0
	} else if pos.Y == 0 && pos.X >= 100 {
		// Edge 5a
		pos.Direction = DirUp
		pos.X = translate(100, 0, pos.X, false)
		pos.Y = 199
	} else if pos.Y == 49 && pos.X >= 100 {
		// Edge 3b
		pos.Direction = DirLeft
		pos.Y = translate(100, 50, pos.X, false)
		pos.X = 99
	} else if pos.Y == 100 && pos.X < 50 {
		// Edge 6a
		pos.Direction = DirRight
		pos.Y = translate(0, 50, pos.X, false)
		pos.X = 50
	} else if pos.Y == 149 && pos.X >= 50 {
		// Edge 1a
		pos.Direction = DirLeft
		pos.Y = translate(50, 150, pos.X, false)
		pos.X = 49
	} else if pos.Y == 199 && pos.X < 50 {
		// Edge 5b
		pos.Direction = DirDown
		pos.X = translate(0, 100, pos.X, false)
		pos.Y = 0
	} else if pos.X == 50 && pos.Y < 50 {
		// Edge 7b
		pos.Direction = DirRight
		pos.Y = translate(0, 149, pos.Y, true)
		pos.X = 0
	} else if pos.X == 149 && pos.Y < 50 {
		// Edge 2b
		pos.Direction = DirLeft
		pos.Y = translate(0, 149, pos.Y, true)
		pos.X = 99
	} else if pos.X == 50 && pos.Y >= 50 && pos.Y < 100 {
		// Edge 6b
		pos.Direction = DirDown
		pos.X = translate(50, 0, pos.Y, false)
		pos.Y = 100
	} else if pos.X == 99 && pos.Y >= 50 && pos.Y < 100 {
		// Edge 3a
		pos.Direction = DirUp
		pos.X = translate(50, 100, pos.Y, false)
		pos.Y = 49
	} else if pos.X == 0 && pos.Y >= 100 && pos.Y < 150 {
		// Edge 7a
		pos.Direction = DirRight
		pos.Y = translate(100, 49, pos.Y, true)
		pos.X = 50
	} else if pos.X == 99 && pos.Y >= 100 && pos.Y < 150 {
		// Edge 2a
		pos.Direction = DirLeft
		pos.Y = translate(100, 49, pos.Y, true)
		pos.X = 149
	} else if pos.X == 0 && pos.Y >= 150 {
		// Edge 4b
		pos.Direction = DirDown
		pos.X = translate(150, 50, pos.Y, false)
		pos.Y = 0
	} else if pos.X == 49 && pos.Y >= 150 {
		// Edge 1b
		pos.Direction = DirUp
		pos.X = translate(150, 50, pos.Y, false)
		pos.Y = 149
	} else {
		panic(fmt.Errorf("unknown edge: %v", pos))
	}
	if pos.X < 0 || pos.Y < 0 {
		panic(fmt.Errorf("X or Y is negative: %v", pos))
	}
	return pos
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
	board, path := parseInput(input)

	pos := Position{
		X:         strings.Index(string(board[0]), "."),
		Y:         0,
		Direction: DirRight,
	}

	numCache := ""
	for it, step := range path {
		if it == 1111 {
			// fmt.Println("bla")
		}

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
	if len(numCache) > 0 {
		pos = execStep(board, numCache, pos)
	}

	return ((pos.Y + 1) * 1000) + ((pos.X + 1) * 4) + pos.Direction
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
