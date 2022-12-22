package main

import (
	_ "embed"
	"flag"
	"fmt"
	"strings"
	"unicode"

	"github.com/mheidinger/advent-of-code-go/data-structures/set"
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

type Rucksack struct {
	Comp1 []rune
	Comp2 []rune
}

func (ruck Rucksack) getCommon() int {
	for _, rune1 := range ruck.Comp1 {
		for _, rune2 := range ruck.Comp2 {
			if rune1 == rune2 {
				return toInt(rune1)
			}
		}
	}
	panic(fmt.Errorf("common not found"))
}

func (ruck Rucksack) getCommonThree(ruck2, ruck3 Rucksack) int {
	ruck1Set := set.NewIntSet([]int{})
	for _, rune1 := range ruck.Comp1 {
		ruck1Set.Add(toInt(rune1))
	}
	ruck2Set := set.NewIntSet([]int{})
	for _, rune2 := range ruck2.Comp1 {
		ruck2Set.Add(toInt(rune2))
	}
	ruck3Set := set.NewIntSet([]int{})
	for _, rune3 := range ruck3.Comp1 {
		ruck3Set.Add(toInt(rune3))
	}
	totalMap := map[int]int{}
	for _, num1 := range ruck1Set.Keys() {
		totalMap[num1] = 1
	}
	for _, num1 := range ruck2Set.Keys() {
		if _, ok := totalMap[num1]; ok {
			totalMap[num1] = 2
		}
	}
	for _, num1 := range ruck3Set.Keys() {
		if cur, ok := totalMap[num1]; ok && cur == 2 {
			return num1
		}
	}
	panic(fmt.Errorf("nothing found"))
}

func part1(input string) int {
	parsed := parseInput(input)

	total := 0
	for _, sack := range parsed {
		total += sack.getCommon()
	}

	return total
}

func part2(input string) int {
	parsed := parseInputAll(input)

	total := 0
	for it := 0; it < len(parsed); {
		total += parsed[it].getCommonThree(parsed[it+1], parsed[it+2])
		it += 3
	}

	return total
}

func toInt(char rune) int {
	if unicode.IsLower(char) { // lowercase
		return int(char) - 96
	}
	// uppercase
	return int(char) - 64 + 26
}

func parseInputAll(input string) (ans []Rucksack) {
	for _, line := range strings.Split(input, "\n") {
		sack := Rucksack{[]rune{}, []rune{}}
		for _, char := range []rune(line) {
			sack.Comp1 = append(sack.Comp1, char)
		}
		ans = append(ans, sack)
	}
	return ans
}

func parseInput(input string) (ans []Rucksack) {
	for _, line := range strings.Split(input, "\n") {
		sack := Rucksack{[]rune{}, []rune{}}
		comp1 := line[:len(line)/2]
		for _, char := range []rune(comp1) {
			sack.Comp1 = append(sack.Comp1, char)
		}
		comp2 := line[len(line)/2:]
		for _, char := range []rune(comp2) {
			sack.Comp2 = append(sack.Comp2, char)
		}

		if len(sack.Comp1) != len(sack.Comp2) {
			panic(fmt.Errorf("length of compartments not equal"))
		}

		ans = append(ans, sack)
	}
	return ans
}
