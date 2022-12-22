package main

import (
	_ "embed"
	"flag"
	"fmt"
	"regexp"
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
		ans := part1(input, 2000000)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input, 4000000)
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

func (point Point) ManhattanDistance(point2 Point) int {
	return Abs(point.x-point2.x) + Abs(point.y-point2.y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Sensor struct {
	pos        Point
	beacon     Point
	distBeacon int
}

func isCovered(sensors []*Sensor, point Point, beaconRet bool) bool {
	covered := false
	for _, sens := range sensors {
		if sens.beacon == point {
			return beaconRet
		}

		if sens.pos.ManhattanDistance(point) <= sens.distBeacon {
			covered = true
		}
	}
	return covered
}

func isCoveredPart2(sensors []*Sensor, point Point) (bool, int) {
	for _, sens := range sensors {
		pointDistance := sens.pos.ManhattanDistance(point)
		if pointDistance <= sens.distBeacon {
			skipY := Abs(sens.distBeacon - pointDistance)
			return true, skipY
		}
	}
	return false, 0
}

func part1(input string, y int) int {
	sensors, minX, maxX := parseInput(input)

	covered := 0
	for it := minX - 50; it <= maxX+50; it++ {
		point := Point{it, y}
		if isCovered(sensors, point, false) {
			covered++
		}
	}

	return covered
}

func part2(input string, searchMax int) int {
	sensors, _, _ := parseInput(input)

	for x := 0; x <= searchMax; x++ {
		for y := 0; y <= searchMax; y++ {
			point := Point{x, y}
			covered, skipY := isCoveredPart2(sensors, point)
			if !covered {
				return x*4000000 + y
			}
			y += skipY
			// fmt.Printf("checked point %v, skip %d\n", point, skipY)
		}
	}

	return -1
}

var reg = regexp.MustCompile(`.*x=(-?\d+), y=(-?\d+).*x=(-?\d+), y=(-?\d+)`)

func parseInput(input string) (ans []*Sensor, minX, maxX int) {
	maxDistBeacon := 0
	minX = 99999999999999999
	maxX = 0
	for _, line := range strings.Split(input, "\n") {
		matches := reg.FindStringSubmatch(line)
		sens := &Sensor{
			pos:    Point{cast.ToInt(matches[1]), cast.ToInt(matches[2])},
			beacon: Point{cast.ToInt(matches[3]), cast.ToInt(matches[4])},
		}
		sens.distBeacon = sens.pos.ManhattanDistance(sens.beacon)
		ans = append(ans, sens)

		if Abs(sens.distBeacon) > maxDistBeacon {
			maxDistBeacon = Abs(sens.distBeacon)
		}
		if sens.pos.x < minX {
			minX = sens.pos.x
		} else if sens.pos.x > maxX {
			maxX = sens.pos.x
		}
	}
	return ans, minX - maxDistBeacon, maxX + maxDistBeacon
}
