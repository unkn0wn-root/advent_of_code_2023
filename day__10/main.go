package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// position represents the x, y coordinates in the grid
type position struct {
	x int
	y int
}

// zoo represents the visual representation and count for part 2
type zoo struct {
	v     string
	count int
}

func main() {
	lines, file, err := readInputFile("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	// find the path history and count for part 1
	history, count := findCount(lines)

	// find the area for part 2
	numberOfInsideElements, visual := findVisualArea(history, lines)

	fmt.Println("Part 1: ", count)
	fmt.Println("Part 2: ", numberOfInsideElements)

	// print visual representation
	for i := 0; i < len(lines); i++ {
		if _, ok := visual[i]; ok {
			fmt.Println(visual[i].v, visual[i].count)
		}
	}
}

func readInputFile(filename string) ([]string, *os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return lines, file, nil
}

// findCount returns the path history and count for part 1
func findCount(input []string) ([]position, int) {
	history := []position{}

	// find the starting position 'S'
	x, y := findSChar(input)
	history = append(history, position{x: x, y: y})

	// get the starting position after 'S'
	x, y = getStartingPosition(x, y, input)

	// traverse the path until reaching the starting position 'S' again
	for input[x][y] != 'S' {
		lastP := history[len(history)-1]
		x1, y1 := findNextChar(x, y, input, lastP.x, lastP.y)
		history = append(history, position{x: x, y: y})
		x, y = x1, y1
	}

	return history, len(history) / 2
}

// getStartingPosition finds the next position after 'S'
func getStartingPosition(x, y int, input []string) (int, int) {
	if y != (len(input[x])-1) && (input[x][y+1] == '-' || input[x][y+1] == 'J' || input[x][y+1] == '7') {
		y++
		return x, y
	}

	if x != (len(input)-1) && (input[x+1][y] == '|' || input[x+1][y] == 'J' || input[x+1][y] == 'L') {
		x++
		return x, y
	}

	if y != 0 && (input[x][y-1] == 'F' || input[x][y-1] == '-' || input[x][y-1] == 'L') {
		y--
		return x, y
	}

	if x != 0 && (input[x-1][y] == '|' || input[x-1][y] == 'L' || input[x-1][y] == '7') {
		x--
		return x, y
	}

	return x, y
}

// findNextChar finds the next position based on the current position and the grid layout
func findNextChar(x, y int, input []string, lastx, lasty int) (int, int) {
	if input[x][y] == '-' {
		if y != 0 && lasty == (y-1) {
			y++
			return x, y
		}

		if y != (len(input[x])-1) && lasty == (y+1) {
			y--
			return x, y
		}
	}

	if input[x][y] == 'J' {
		if x != 0 && y != 0 && lastx == (x-1) {
			y--
			return x, y
		}

		if x != 0 && y != 0 && lasty == (y-1) {
			x--
			return x, y
		}
	}

	if input[x][y] == '|' {
		if x != 0 && x != len(input)-1 && lastx == (x-1) {
			x++
			return x, y
		}

		if x != 0 && x != len(input)-1 && lastx == (x+1) {
			x--
			return x, y
		}
	}

	if input[x][y] == 'L' {
		if x != 0 && y != len(input[x])-1 && lastx == (x-1) {
			y++
			return x, y
		}

		if x != 0 && y != len(input[x])-1 && lasty == (y+1) {
			x--
			return x, y
		}
	}

	if input[x][y] == '7' {
		if y != 0 && x != (len(input)-1) && lasty == (y-1) {
			x++
			return x, y
		}

		if x != (len(input)-1) && y != 0 && lastx == (x+1) {
			y--
			return x, y
		}
	}

	if input[x][y] == 'F' {
		if x != (len(input)-1) && y != (len(input[x])-1) && lasty == (y+1) {
			x++
			return x, y
		}

		if x != (len(input)-1) && y != (len(input[x])-1) && lastx == (x+1) {
			y++
			return x, y
		}
	}

	return x, y
}

// findVisualArea calculates the visual representation and count for part 2
func findVisualArea(path []position, input []string) (int, map[int]zoo) {
	replaceWith := map[string]string{
		"J": "┘", "L": "└", "7": "┐", "F": "┌", "|": "│", "-": "─",
	}
	mapPosition := make(map[int][]int)
	resultMap := make(map[int]zoo)
	sum := 0

	for _, p := range path {
		mapPosition[p.x] = append(mapPosition[p.x], p.y)
	}

	for k, v := range mapPosition {
		a := strings.Split(input[k], "")
		for _, j := range v {
			if a[j] != "S" {
				a[j] = replaceWith[a[j]]
			} else {
				a[j] = replaceWith[string(replaceSChar(path[0], path[1], path[len(path)-1], input))]
			}
		}
		// clean edges
		left, right := false, false
		for i, j := range a {
			if i == (len(a) - 1 - i) {
				break
			}

			r := a[len(a)-i-1]

			if j == "┘" || j == "└" || j == "┐" || j == "┌" || j == "│" || j == "─" {
				left = true
			}

			if r == "┘" || r == "└" || r == "┐" || r == "┌" || r == "│" || r == "─" {
				right = true
			}

			if !left {
				a[i] = " "
			}

			if !right {
				a[len(a)-1-i] = " "
			}
		}
		count := 0
		isInside := 0
		last := "-"

		// traverse the cleaned input line and count the number of inside elements '|'
		for i, char := range a {
			if char == "│" {
				isInside++
				continue
			}

			if char == "─" {
				continue
			}

			if last == "-" && (char == "┘" || char == "└" || char == "┐" || char == "┌") {
				if char == "┘" || char == "└" || char == "┐" || char == "┌" {
					last = char
					continue
				}
			} else if last != "-" && (char == "┘" || char == "└" || char == "┐" || char == "┌") {
				if last == "└" && char == "┐" {
					isInside++
				}

				if last == "┌" && char == "┘" {
					isInside++
				}

				last = "-"
				continue
			}

			if isInside%2 == 0 {
				a[i] = " "
			}

			if isInside%2 != 0 {
				a[i] = "\033[0;32m█\033[0m"
				count++
			}
		}
		resultMap[k] = zoo{
			v:     strings.Join(a, ""),
			count: count,
		}
		sum += count
	}

	return sum, resultMap
}

// replaceSChar replaces 'S' with other characters and checks if it reaches the starting position
func replaceSChar(s, first, last position, input []string) byte {
	for _, char := range "-|JL7F" {
		duplicateInput := input
		a := strings.Split(duplicateInput[s.x], "")
		a[s.y] = string(char)
		duplicateInput[s.x] = strings.Join(a, "")
		x, y := findNextChar(s.x, s.y, duplicateInput, last.x, last.y)
		if x == first.x && y == first.y {
			return byte(char)
		}
	}

	return 'S'
}

// findSChar finds the position of 'S' in the grid
func findSChar(input []string) (int, int) {
	for i, line := range input {
		for j, char := range line {
			if char == 'S' {
				return i, j
			}
		}
	}

	return 0, 0
}
