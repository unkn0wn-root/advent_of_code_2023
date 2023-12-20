package main

import (
	"fmt"
	"image"
	"os"
	"strings"
)

var (
	right = image.Point{1, 0}
	left  = image.Point{-1, 0}
	up    = image.Point{0, -1}
	down  = image.Point{0, 1}

	reflectors = map[rune]map[image.Point]image.Point{
		'/':  {up: right, right: up, down: left, left: down},
		'\\': {up: left, right: down, left: up, down: right},
	}

	forks = map[rune][]image.Point{
		'-': {right, left},
		'|': {up, down},
	}
)

type Particle struct {
	position  image.Point
	direction image.Point
}

type FloorMap map[image.Point]rune

func readFromFile(filename string) (FloorMap, int, int) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return parseFloorMap(string(content))
}

// convert the input string into a FloorMap, and returns the map, width, and height.
func parseFloorMap(input string) (FloorMap, int, int) {
	floorMap := make(FloorMap)
	var width, height int
	for y, line := range strings.Split(strings.TrimSpace(input), "\n") {
		width = len(line) - 1
		height = y
		for x, symbol := range strings.TrimSpace(line) {
			floorMap[image.Point{x, y}] = symbol
		}
	}

	return floorMap, width, height
}

// simulate the movement of particles on the floor map and returns the number of visited positions.
func traverseFloor(floorMap FloorMap, width, height int, particle Particle) int {
	visited := map[image.Point][]image.Point{}
	particles := []Particle{particle}

	for len(particles) > 0 {
		particle, particles = particles[0], particles[1:]

		// check if the particle is out of bounds.
		if particle.position.X > width || particle.position.X < 0 || particle.position.Y > height || particle.position.Y < 0 {
			continue
		}

		// check if the particle has visited the current position in the same direction before.
		if v, ok := visited[particle.position]; ok && containsDirection(v, particle.direction) {
			continue
		}

		// mark the current position as visited in the specified direction.
		visited[particle.position] = append(visited[particle.position], particle.direction)

		// check if there is a reflector at the current position.
		if v, ok := reflectors[floorMap[particle.position]][particle.direction]; ok {
			particle.direction = v
		}

		// check if there is a fork at the current position.
		if fork, ok := forks[floorMap[particle.position]]; ok {
			// create a new particle with the second direction of the fork.
			particles = append(particles, Particle{position: particle.position.Add(fork[1]), direction: fork[1]})
			// change the direction of the current particle to the first direction of the fork.
			particle.direction = fork[0]
		}

		// move the particle to the next position.
		particle.position = particle.position.Add(particle.direction)
		particles = append(particles, particle)
	}

	return len(visited)
}

// simulate the movement of a particle starting from the top-left corner and returns the number of visited positions.
func partOne(floorMap FloorMap, width int, height int) int {
	return traverseFloor(floorMap, width, height, Particle{position: image.Point{0, 0}, direction: right})
}

// simulate the movement of particles from different starting positions and returns the maximum number of visited positions.
func partTwo(floorMap FloorMap, width int, height int) int {
	maxCoverage := 0

	for x := 0; x < width; x++ {
		maxCoverage = max(
			maxCoverage, traverseFloor(
				floorMap, width, height, Particle{position: image.Point{x, 0}, direction: down},
			))

		maxCoverage = max(
			maxCoverage, traverseFloor(
				floorMap, width, height, Particle{position: image.Point{x, height}, direction: up},
			))
	}

	for y := 0; y < height; y++ {
		maxCoverage = max(
			maxCoverage, traverseFloor(
				floorMap, width, height, Particle{position: image.Point{0, y}, direction: right},
			))

		maxCoverage = max(
			maxCoverage, traverseFloor(
				floorMap, width, height, Particle{position: image.Point{width, y}, direction: left},
			))
	}

	return maxCoverage
}

func main() {
	floorMap, width, height := readFromFile("input.txt")
	fmt.Println("Part 1:", partOne(floorMap, width, height))
	fmt.Println("Part 2:", partTwo(floorMap, width, height))
}

// check if a direction is present in a slice of directions.
func containsDirection(directions []image.Point, dir image.Point) bool {
	for _, d := range directions {
		if d == dir {
			return true
		}
	}

	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
