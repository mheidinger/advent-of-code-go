package main

import (
	_ "embed"
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/mheidinger/advent-of-code-go/cast"
	"github.com/mheidinger/advent-of-code-go/util"
	"golang.org/x/exp/slices"
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

type Valve struct {
	id        string
	flowRate  int
	tunnels   []string
	distances map[string]int

	// for path finding
	visited  bool
	distance int
}

// START Pathfindin

func setValveDistances(valves map[string]*Valve, startValve *Valve) {
	calcValveDistances(valves, startValve)
	for _, valve := range valves {
		startValve.distances[valve.id] = valve.distance
	}
}

func calcValveDistances(valves map[string]*Valve, startValve *Valve) {
	for _, valve := range valves {
		valve.distance = 9999999
		if valve.id == startValve.id {
			valve.distance = 0
		}
		valve.visited = false
	}
	current := startValve

	for current != nil {
		neighbours := current.tunnels
		for _, neighbour := range neighbours {
			if current.distance+1 < valves[neighbour].distance {
				valves[neighbour].distance = current.distance + 1
			}
		}

		current.visited = true
		current = getLowestValve(valves)
	}
}

func getLowestValve(valves map[string]*Valve) *Valve {
	var minimum *Valve
	for _, valve := range valves {
		if (minimum == nil || valve.distance < minimum.distance) && !valve.visited {
			minimum = valve
		}
	}
	return minimum
}

// END Pathfinding

func getNextValve(valves map[string]*Valve, currValve string, targetValves []string, timeLeft int) (string, int, int) {
	calcValveDistances(valves, valves[currValve])

	pressureRelief := 0
	var valve *Valve
	for _, targetValveID := range targetValves {
		targetValve := valves[targetValveID]
		valvePressureRelief := targetValve.flowRate * (timeLeft - targetValve.distance - 1)
		if valvePressureRelief > pressureRelief {
			pressureRelief = valvePressureRelief
			valve = targetValve
		}
	}
	return valve.id, valve.distance, pressureRelief
}

type Path struct {
	valves         []string
	pressureRelief int
}

func getMaxPressureRelief(valves map[string]*Valve, currentValve *Valve, targetValves []*Valve, timeLeft int) []*Path {
	paths := []*Path{}

	if len(targetValves) == 0 {
		paths = append(paths, &Path{
			valves:         []string{},
			pressureRelief: 0,
		})
		return paths
	}

	for it, targetValve := range targetValves {
		dist := currentValve.distances[targetValve.id]
		newTimeLeft := timeLeft - dist - 1
		if newTimeLeft <= 0 {
			paths = append(paths, &Path{
				valves:         []string{},
				pressureRelief: 0,
			})
			continue
		}
		pressureRelief := targetValve.flowRate * newTimeLeft
		reducedTargetValves := slices.Delete(slices.Clone(targetValves), it, it+1)
		nestedPaths := getMaxPressureRelief(valves, targetValve, reducedTargetValves, newTimeLeft)
		for _, nestedPath := range nestedPaths {
			valves := []string{targetValve.id}
			paths = append(paths, &Path{
				valves:         append(valves, nestedPath.valves...),
				pressureRelief: nestedPath.pressureRelief + pressureRelief,
			})
		}
	}
	return paths
}

func part1(input string) int {
	valves := parseInput(input)

	targetValves := []*Valve{}
	for _, valve := range valves {
		setValveDistances(valves, valve)
		if valve.flowRate > 0 {
			targetValves = append(targetValves, valve)
		}
	}

	paths := getMaxPressureRelief(valves, valves["AA"], targetValves, 30)

	pressureRelief := 0
	for _, path := range paths {
		if path.pressureRelief > pressureRelief {
			pressureRelief = path.pressureRelief
		}
	}

	return pressureRelief
}

func nextProduct(a []*Path, r int) func() []*Path {
	p := make([]*Path, r)
	x := make([]int, len(p))
	return func() []*Path {
		p := p[:len(x)]
		for i, xi := range x {
			p[i] = a[xi]
		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(a) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return p
	}
}

func getPathKey(valves []string) string {
	keyBuilder := strings.Builder{}
	for _, valve := range valves {
		keyBuilder.WriteString(valve)
	}
	return keyBuilder.String()
}

func getInvertTargetValves(valves []*Valve, used []string) []*Valve {
	clone := slices.Clone(valves)
	for _, usedValve := range used {
		for it, valve := range clone {
			if valve.id == usedValve {
				clone = slices.Delete(clone, it, it+1)
				break
			}
		}
	}
	return clone
}

func part2(input string) int {
	valves := parseInput(input)

	targetValves := []*Valve{}
	for _, valve := range valves {
		setValveDistances(valves, valve)
		if valve.flowRate > 0 {
			targetValves = append(targetValves, valve)
		}
	}

	paths := getMaxPressureRelief(valves, valves["AA"], targetValves, 26)

	pathMap := make(map[string]*Path, len(paths))
	for _, path := range paths {
		key := getPathKey(path.valves)
		pathMap[key] = path
	}

	pressureRelief := 0
	for {
		fmt.Println(len(pathMap))
		var path *Path
		var key string
		for selKey, selPath := range pathMap {
			path = selPath
			key = selKey
			break
		}
		fmt.Println(path.valves)
		invertValves := getInvertTargetValves(targetValves, path.valves)
		fmt.Println(invertValves)
		complementPaths := getMaxPressureRelief(valves, valves["AA"], invertValves, 26)
		fmt.Println(len(complementPaths))
		for _, compPath := range complementPaths {
			if path.pressureRelief+compPath.pressureRelief > pressureRelief {
				pressureRelief = path.pressureRelief + compPath.pressureRelief
			}
			// delete(pathMap, getPathKey(compPath.valves))
		}
		delete(pathMap, key)
		if len(pathMap) == 0 {
			break
		}
	}

	return pressureRelief
}

var reg = regexp.MustCompile(`Valve (..) has flow rate=(\d+); tunnels? leads? to valves? (.*)`)

func parseInput(input string) (ans map[string]*Valve) {
	ans = map[string]*Valve{}
	for _, line := range strings.Split(input, "\n") {
		matches := reg.FindStringSubmatch(line)
		valve := &Valve{
			id:        matches[1],
			flowRate:  cast.ToInt(matches[2]),
			tunnels:   strings.Split(matches[3], ", "),
			distances: map[string]int{},
		}
		ans[valve.id] = valve
	}
	return ans
}
