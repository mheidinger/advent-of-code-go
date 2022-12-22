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

type Tree struct {
	height  int
	visible bool
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
	forest := parseInput(input)
	last := len(forest) - 1

	// go through rows/columns from both sides, remember highest tree, mark trees higher as visible
	for x := 0; x < len(forest); x++ {
		highestTreeLeft := -1
		highestTreeRight := -1
		highestTreeTop := -1
		highestTreeBottom := -1
		for y := 0; y < len(forest[x]); y++ {
			if forest[x][y].height > highestTreeLeft {
				highestTreeLeft = forest[x][y].height
				forest[x][y].visible = true
			}
			if forest[x][last-y].height > highestTreeRight {
				highestTreeRight = forest[x][last-y].height
				forest[x][last-y].visible = true
			}
			if forest[y][x].height > highestTreeTop {
				highestTreeTop = forest[y][x].height
				forest[y][x].visible = true
			}
			if forest[last-y][x].height > highestTreeBottom {
				highestTreeBottom = forest[last-y][x].height
				forest[last-y][x].visible = true
			}
		}
	}

	numVisible := 0
	for x := 0; x < len(forest); x++ {
		for y := 0; y < len(forest[x]); y++ {
			if forest[x][y].visible {
				numVisible++
			}
		}
	}

	return numVisible
}

func getScenicScore(forest [][]*Tree, x, y int) int {
	treeHeight := forest[x][y].height
	size := len(forest)

	checkLeft := x - 1
	sightLeft := 0
	checkRight := x + 1
	sightRight := 0
	checkTop := y - 1
	sightTop := 0
	checkBottom := y + 1
	sightBottom := 0
	for it := 0; it < size; it++ {
		if checkLeft < 0 && checkRight >= size && checkTop < 0 && checkBottom >= size {
			break
		}

		if checkLeft >= 0 {
			sightLeft++
			if forest[checkLeft][y].height >= treeHeight {
				checkLeft = -1
			} else {
				checkLeft--
			}
		}
		if checkRight < size {
			sightRight++
			if forest[checkRight][y].height >= treeHeight {
				checkRight = size
			} else {
				checkRight++
			}
		}
		if checkTop >= 0 {
			sightTop++
			if forest[x][checkTop].height >= treeHeight {
				checkTop = -1
			} else {
				checkTop--
			}
		}
		if checkBottom < size {
			sightBottom++
			if forest[x][checkBottom].height >= treeHeight {
				checkBottom = size
			} else {
				checkBottom++
			}
		}
	}

	return sightLeft * sightRight * sightTop * sightBottom
}

func part2(input string) int {
	forest := parseInput(input)

	highestScore := 0
	for x := 0; x < len(forest); x++ {
		for y := 0; y < len(forest[x]); y++ {
			if score := getScenicScore(forest, x, y); score > highestScore {
				highestScore = score
			}
		}
	}
	return highestScore
}

func parseInput(input string) (ans [][]*Tree) {
	for _, line := range strings.Split(input, "\n") {
		rowNums := strings.Split(line, "")
		row := make([]*Tree, 0, len(rowNums))
		for _, num := range rowNums {
			row = append(row, &Tree{
				height:  cast.ToInt(num),
				visible: false,
			})
		}
		ans = append(ans, row)
	}
	return ans
}
