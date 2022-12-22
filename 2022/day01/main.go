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

func part1(input string) int {
	parsed := parseInput(input)

	highestCal := 0
	for _, elf := range parsed {
		elfCal := 0
		for _, cal := range elf {
			elfCal += cal
		}
		if elfCal > highestCal {
			highestCal = elfCal
		}
	}

	return highestCal
}

func part2(input string) int {
	parsed := parseInput(input)

	highestCal1 := -1
	highestCal2 := -2
	highestCal3 := -3
	for _, elf := range parsed {
		elfCal := 0
		for _, cal := range elf {
			elfCal += cal
		}
		if elfCal > highestCal1 && highestCal1 < highestCal2 && highestCal1 < highestCal3 {
			highestCal1 = elfCal
			continue
		}
		if elfCal > highestCal2 && highestCal2 < highestCal3 && highestCal2 < highestCal1 {
			highestCal2 = elfCal
			continue
		}
		if elfCal > highestCal3 && highestCal3 < highestCal2 && highestCal3 < highestCal1 {
			highestCal3 = elfCal
			continue
		}
	}

	return highestCal1 + highestCal2 + highestCal3
}

func parseInput(input string) (ans [][]int) {
	currentElf := []int{}
	for _, line := range strings.Split(input, "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			ans = append(ans, currentElf)
			currentElf = []int{}
			continue
		}
		currentElf = append(currentElf, cast.ToInt(line))
	}
	return ans
}
