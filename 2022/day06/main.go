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

func part1(input string) int {
	parsed := parseInput(input)

	result := -1
	for it := 3; it < len(parsed); it++ {
		lastFour := parsed[it-3 : it+1]
		if allUnique(lastFour) {
			result = it
			break
		}
	}

	return result + 1
}

func part2(input string) int {
	parsed := parseInput(input)

	result := -1
	for it := 13; it < len(parsed); it++ {
		lastThirteen := parsed[it-13 : it+1]
		if allUnique(lastThirteen) {
			result = it
			break
		}
	}

	return result + 1
}

func allUnique(input string) bool {
	for _, char := range input {
		if strings.Count(input, string(char)) > 1 {
			return false
		}
	}
	return true
}

func parseInput(input string) string {
	for _, line := range strings.Split(input, "\n") {
		return line
	}
	return ""
}
