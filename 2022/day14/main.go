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

type Point struct {
	x int
	y int
}

func (point Point) Sub(point2 Point) Point {
	return Point{
		point.x - point2.x,
		point.y - point2.y,
	}
}

func (point Point) Add(point2 Point) Point {
	return Point{
		point.x + point2.x,
		point.y + point2.y,
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (point Point) Len() int {
	return Abs(point.x + point.y)
}

type Material string

const (
	MatAir   = "."
	MatRock  = "#"
	MatSand  = "o"
	MatSpawn = "+"
)

func getEmptyScene(low, high Point) [][]Material {
	rangeY := high.y - low.y

	scene := make([][]Material, 0, high.x+1)
	for x := 0; x <= high.x+3; x++ {
		row := make([]Material, 0, rangeY+1)
		for y := 0; y <= rangeY; y++ {
			row = append(row, MatAir)
		}
		scene = append(scene, row)
	}

	return scene
}

func insertRocks(scene [][]Material, lines [][]Point) {
	for _, line := range lines {
		drawLine(scene, line)
	}
}

func sign(in int) int {
	if in < 0 {
		return -1
	}
	return 1
}

func drawLine(scene [][]Material, line []Point) {
	for it := 0; it < len(line)-1; it++ {
		point1 := line[it]
		point2 := line[it+1]
		diff := point2.Sub(point1)
		for step := 0; step <= diff.Len(); step++ {
			if Abs(diff.x) > 0 {
				x := point1.x + (step * sign(diff.x))
				scene[x][point1.y] = MatRock
			} else {
				y := point1.y + (step * sign(diff.y))
				scene[point1.x][y] = MatRock
			}
		}
	}
}

func scaleLinesY(lines [][]Point, low Point) [][]Point {
	for x := range lines {
		for y := range lines[x] {
			lines[x][y] = lines[x][y].Sub(Point{0, low.y})
		}
	}

	return lines
}

func drawScene(scene [][]Material) {
	for _, row := range scene {
		for _, pixel := range row {
			fmt.Print(pixel)
		}
		fmt.Println()
	}
}

func simulateSand(scene [][]Material, spawn Point) (spawnedSand int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	for {
		sandPos := spawn
		for {
			if scene[sandPos.x+1][sandPos.y] == MatAir {
				sandPos = sandPos.Add(Point{1, 0})
			} else if scene[sandPos.x+1][sandPos.y-1] == MatAir {
				sandPos = sandPos.Add(Point{1, -1})
			} else if scene[sandPos.x+1][sandPos.y+1] == MatAir {
				sandPos = sandPos.Add(Point{1, 1})
			} else {
				break
			}
		}
		scene[sandPos.x][sandPos.y] = MatSand
		spawnedSand++
	}
}

func simulateSandPart2(scene [][]Material, spawn Point) (spawnedSand int) {
	for {
		sandPos := spawn
		for {
			if scene[sandPos.x+1][sandPos.y] == MatAir {
				sandPos = sandPos.Add(Point{1, 0})
			} else if scene[sandPos.x+1][sandPos.y-1] == MatAir {
				sandPos = sandPos.Add(Point{1, -1})
			} else if scene[sandPos.x+1][sandPos.y+1] == MatAir {
				sandPos = sandPos.Add(Point{1, 1})
			} else {
				break
			}
		}
		scene[sandPos.x][sandPos.y] = MatSand
		spawnedSand++
		if sandPos == spawn {
			break
		}
	}
	return spawnedSand
}

func part1(input string) int {
	lines, low, high := parseInput(input)
	lines = scaleLinesY(lines, low)

	scene := getEmptyScene(low, high)
	insertRocks(scene, lines)

	spawn := Point{0, 500}.Sub(Point{0, low.y})
	res := simulateSand(scene, spawn)
	drawScene(scene)
	return res
}

func part2(input string) int {
	lines, low, high := parseInput(input)
	margin := 300
	bottomLine := []Point{
		{high.x + 2, 0},
		{high.x + 2, high.y - low.y + margin*2},
	}

	low = low.Sub(Point{0, margin})
	high = high.Add(Point{0, margin})
	lines = scaleLinesY(lines, low)

	scene := getEmptyScene(low, high)
	insertRocks(scene, lines)
	// insert bottom row
	drawLine(scene, bottomLine)

	spawn := Point{0, 500}.Sub(Point{0, low.y})
	res := simulateSandPart2(scene, spawn)
	drawScene(scene)
	return res
}

func parseInput(input string) (ans [][]Point, low, high Point) {
	low = Point{99999, 99999}
	high = Point{0, 0}
	for _, line := range strings.Split(input, "\n") {
		lineSplit := strings.Split(line, " -> ")
		pointList := []Point{}
		for _, linePoint := range lineSplit {
			pointCoords := strings.Split(linePoint, ",")
			x := cast.ToInt(pointCoords[1])
			y := cast.ToInt(pointCoords[0])
			pointList = append(pointList, Point{x, y})

			if x < low.x {
				low.x = x
			}
			if y < low.y {
				low.y = y
			}
			if x > high.x {
				high.x = x
			}
			if y > high.y {
				high.y = y
			}
		}
		ans = append(ans, pointList)
	}
	return ans, low, high
}
