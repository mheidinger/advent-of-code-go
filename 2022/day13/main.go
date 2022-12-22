package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sort"
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

type Num struct {
	parent    *Num
	plain     int
	arr       []*Num
	isDivider bool
}

// returns bool for inOrder and bool for whether to continue the checks
func inOrder(num1 *Num, num2 *Num) (bool, bool) {
	fmt.Printf("Compare %+v with %+v\n", num1, num2)
	if num1.arr == nil && num2.arr == nil {
		return num1.plain < num2.plain, num1.plain == num2.plain
	} else if num1.arr != nil && num2.arr != nil {
		for it, nestedNum := range num1.arr {
			if it >= len(num2.arr) {
				return false, false
			}
			order, cont := inOrder(nestedNum, num2.arr[it])
			if !cont {
				return order, cont
			}
		}
		if len(num1.arr) < len(num2.arr) {
			return true, false
		} else if len(num1.arr) == len(num2.arr) {
			return true, true
		} else {
			return false, false
		}
	} else if num1.arr == nil && num2.arr != nil {
		arrNum := &Num{arr: []*Num{{plain: num1.plain}}}
		return inOrder(arrNum, num2)
	} else if num1.arr != nil && num2.arr == nil {
		arrNum := &Num{arr: []*Num{{plain: num2.plain}}}
		return inOrder(num1, arrNum)
	}

	panic(fmt.Errorf("Unknown situation: %v == %v", num1, num2))
}

func part1(input string) int {
	pairs := parseInput(input)

	sumInOrderIt := 0
	for it, pair := range pairs {
		if order, _ := inOrder(pair[0], pair[1]); order {
			sumInOrderIt += it + 1
			fmt.Printf("Pair %d is in order\n\n", it+1)
		} else {
			fmt.Printf("Pair %d is not in order\n\n", it+1)
		}
	}

	return sumInOrderIt
}

func part2(input string) int {
	packets := parseInputPart2(input)
	packets = append(packets, &Num{
		isDivider: true,
		arr:       []*Num{{plain: 2}},
	})
	packets = append(packets, &Num{
		isDivider: true,
		arr:       []*Num{{plain: 6}},
	})

	sort.Slice(packets, func(i, j int) bool {
		order, _ := inOrder(packets[i], packets[j])
		return order
	})

	dividerProd := 1
	for it, packet := range packets {
		if packet.isDivider {
			dividerProd *= it + 1
		}
	}

	return dividerProd
}

func parseInput(input string) (ans [][2]*Num) {
	lineSplit := strings.Split(input, "\n")
	for it := 0; it < len(lineSplit); it += 3 {
		num1 := parseLine(lineSplit[it])
		num2 := parseLine(lineSplit[it+1])
		ans = append(ans, [2]*Num{num1, num2})
	}
	return ans
}

func parseInputPart2(input string) (ans []*Num) {
	for _, line := range strings.Split(input, "\n") {
		if line != "" {
			ans = append(ans, parseLine(line))
		}
	}
	return ans
}

func parseLine(line string) *Num {
	// Remove outermost array as we'll start with a "hardcoded" root
	line = line[1 : len(line)-1]

	root := &Num{arr: []*Num{}}

	// Parser Variables
	current := root
	numCache := ""

	for _, char := range line {
		switch char {
		case ',':
			if numCache != "" {
				current.arr = append(current.arr, &Num{parent: current, plain: cast.ToInt(numCache)})
				numCache = ""
			}
		case '[':
			newNum := &Num{parent: current, arr: []*Num{}}
			current.arr = append(current.arr, newNum)
			current = newNum
		case ']':
			if numCache != "" {
				current.arr = append(current.arr, &Num{parent: current, plain: cast.ToInt(numCache)})
				numCache = ""
			}
			current = current.parent
		default:
			// Write numbers into a cache, on a comma add whole number into current root
			// Needed to support numbers with multiple digits
			numCache += string(char)
		}
	}

	// We remove the outermost parantheses, if the last char is a number we therefore won't save it
	if numCache != "" {
		root.arr = append(root.arr, &Num{parent: nil, plain: cast.ToInt(numCache)})
	}

	return root
}
