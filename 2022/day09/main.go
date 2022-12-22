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

type Command struct {
	direction string
	steps     int
}

type Position struct {
	x int
	y int
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func follow(head, tail Position) Position {
	diffX := head.x - tail.x
	diffY := head.y - tail.y
	if Abs(diffX) <= 1 && Abs(diffY) <= 1 {
		return tail
	}

	// linear diff
	if diffY == 0 && diffX > 0 {
		tail.x++
	} else if diffY == 0 && diffX < 0 {
		tail.x--
	} else if diffX == 0 && diffY > 0 {
		tail.y++
	} else if diffX == 0 && diffY < 0 {
		tail.y--
	}
	if diffX == 0 || diffY == 0 {
		return tail
	}

	// diagonal diff, even
	if Abs(diffX) == 2 && Abs(diffY) == 2 {
		if diffX > 0 {
			tail.x++
		} else {
			tail.x--
		}
		if diffY > 0 {
			tail.y++
		} else {
			tail.y--
		}
		return tail
	}

	// diagonal diff, uneven
	if Abs(diffY) >= 2 && diffY > 0 {
		tail.x = head.x
		tail.y++
	} else if Abs(diffY) >= 2 && diffY < 0 {
		tail.x = head.x
		tail.y--
	} else if Abs(diffX) >= 2 && diffX > 0 {
		tail.y = head.y
		tail.x++
	} else if Abs(diffX) >= 2 && diffX < 0 {
		tail.y = head.y
		tail.x--
	}

	return tail
}

func part1(input string) int {
	commands := parseInput(input)

	visitMap := map[string]bool{} // keys are x:y
	head := Position{0, 0}
	tail := Position{0, 0}

	for _, cmd := range commands {
		for it := 0; it < cmd.steps; it++ {
			switch cmd.direction {
			case "R":
				head.x++
			case "L":
				head.x--
			case "U":
				head.y++
			case "D":
				head.y--
			default:
				panic(fmt.Errorf("unknown direction: %s", cmd.direction))
			}

			tail = follow(head, tail)
			visitMap[fmt.Sprintf("%d:%d", tail.x, tail.y)] = true
			fmt.Printf("head: %d:%d, tail: %d:%d\n", head.x, head.y, tail.x, tail.y)
		}
		fmt.Println("------")
	}

	return len(visitMap)
}

func part2(input string) int {
	commands := parseInput(input)

	visitMap := map[string]bool{} // keys are x:y
	rope := make([]Position, 10)

	for _, cmd := range commands {
		for it := 0; it < cmd.steps; it++ {
			switch cmd.direction {
			case "R":
				rope[0].x++
			case "L":
				rope[0].x--
			case "U":
				rope[0].y++
			case "D":
				rope[0].y--
			default:
				panic(fmt.Errorf("unknown direction: %s", cmd.direction))
			}

			for knotPos := range rope {
				rope[knotPos+1] = follow(rope[knotPos], rope[knotPos+1])
				if knotPos == len(rope)-2 {
					visitMap[fmt.Sprintf("%d:%d", rope[knotPos+1].x, rope[knotPos+1].y)] = true
					break
				}
			}

			// for _, knot := range rope {
			// fmt.Printf("%03d:%03d  ", knot.x, knot.y)
			// }
			// fmt.Println()
		}
		// fmt.Println("----------")
	}

	return len(visitMap)
}

func parseInput(input string) (ans []Command) {
	for _, line := range strings.Split(input, "\n") {
		lineParts := strings.Split(line, " ")
		ans = append(ans, Command{
			direction: lineParts[0],
			steps:     cast.ToInt(lineParts[1]),
		})
	}
	return ans
}
