package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"

	"github.com/barkimedes/go-deepcopy"
	"github.com/davecgh/go-spew/spew"
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

type Number struct {
	InitialPos int
	Num        int
}

func findNumToMove(numbers []*Number, initialPos int) (int, *Number) {
	for pos, num := range numbers {
		if num.InitialPos == initialPos {
			return pos, num
		}
	}
	panic(fmt.Errorf("could not find initial pos %d", initialPos))
}

// 1, 2, -3, 4, 0, 3, -2

func moveNumber(numbers []*Number, num *Number, currentPos int) []*Number {
	removedList := append(numbers[:currentPos], numbers[currentPos+1:]...)
	newPos := (currentPos + num.Num) % len(removedList)
	if newPos < 0 {
		newPos = len(removedList) + newPos
	}

	copy := deepcopy.MustAnything(removedList).([]*Number)
	newList := append(copy[:newPos], num)
	return append(newList, removedList[newPos:]...)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func part1(input string) int {
	numbers := parseInput(input)

	spew.Println(numbers)
	for it := range numbers {
		pos, num := findNumToMove(numbers, it)
		numbers = moveNumber(numbers, num, pos)
		// spew.Println(numbers)
	}

	posZero := -1
	for it, num := range numbers {
		if num.Num == 0 {
			posZero = it
			break
		}
	}

	pos1 := (posZero + 1000) % len(numbers)
	pos2 := (posZero + 2000) % len(numbers)
	pos3 := (posZero + 3000) % len(numbers)
	fmt.Printf("%d %d %d\n", pos1, pos2, pos3)
	fmt.Printf("%d %d %d\n", numbers[pos1].Num, numbers[pos2].Num, numbers[pos3].Num)
	return numbers[pos1].Num + numbers[pos2].Num + numbers[pos3].Num
}

func part2(input string) int {
	numbers := parseInput(input)

	for _, num := range numbers {
		num.Num *= 811589153
	}
	spew.Println(numbers)
	for it := 0; it < 10; it++ {
		for it := range numbers {
			pos, num := findNumToMove(numbers, it)
			numbers = moveNumber(numbers, num, pos)
			// spew.Println(numbers)
		}
	}

	posZero := -1
	for it, num := range numbers {
		if num.Num == 0 {
			posZero = it
			break
		}
	}

	pos1 := (posZero + 1000) % len(numbers)
	pos2 := (posZero + 2000) % len(numbers)
	pos3 := (posZero + 3000) % len(numbers)
	fmt.Printf("%d %d %d\n", pos1, pos2, pos3)
	fmt.Printf("%d %d %d\n", numbers[pos1].Num, numbers[pos2].Num, numbers[pos3].Num)
	return numbers[pos1].Num + numbers[pos2].Num + numbers[pos3].Num
}

func parseInput(input string) (ans []*Number) {
	for it, line := range strings.Split(input, "\n") {
		num := &Number{
			InitialPos: it,
			Num:        cast.ToInt(line),
		}
		ans = append(ans, num)
	}
	return ans
}
