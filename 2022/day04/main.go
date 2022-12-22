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

type Range struct {
	Lower  int
	Higher int
}

func (r Range) size() int {
	return r.Higher - r.Lower
}

type Pair struct {
	FirstRange  Range
	SecondRange Range
}

func (p Pair) fullyContained() bool {
	first := p.FirstRange
	second := p.SecondRange
	if p.SecondRange.size() > p.FirstRange.size() {
		second = p.FirstRange
		first = p.SecondRange
	}

	if second.Lower >= first.Lower && second.Higher <= first.Higher {
		return true
	}
	return false
}

func (p Pair) partlyContained() bool {
	first := p.FirstRange
	second := p.SecondRange
	if p.SecondRange.Lower < p.FirstRange.Lower {
		second = p.FirstRange
		first = p.SecondRange
	}

	if (first.Lower <= second.Lower && first.Higher >= second.Lower) ||
		(first.Lower <= second.Higher && first.Higher >= second.Higher) {
		return true
	}
	return false
}

func part1(input string) int {
	parsed := parseInput(input)

	fullyContained := 0
	for _, pair := range parsed {
		if pair.fullyContained() {
			fullyContained++
		}
	}

	return fullyContained
}

func part2(input string) int {
	parsed := parseInput(input)

	partlyContained := 0
	for _, pair := range parsed {
		if pair.partlyContained() {
			partlyContained++
		}
	}

	return partlyContained
}

func parseInput(input string) (ans []Pair) {
	for _, line := range strings.Split(input, "\n") {
		pairSplit := strings.Split(line, ",")
		rangeSplit1 := strings.Split(pairSplit[0], "-")
		rangeSplit2 := strings.Split(pairSplit[1], "-")
		pair := Pair{
			Range{cast.ToInt(rangeSplit1[0]), cast.ToInt(rangeSplit1[1])},
			Range{cast.ToInt(rangeSplit2[0]), cast.ToInt(rangeSplit2[1])},
		}
		ans = append(ans, pair)
	}
	return ans
}
