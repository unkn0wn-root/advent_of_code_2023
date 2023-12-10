package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var square = [][]int{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

func main() {
	inputFile := readLocalInput()
	fmt.Println(solutionPartOne(inputFile))
	fmt.Println(solutionPartTwo(inputFile))
}

// read and split input file for each processing line instead of doing it in each func
func readLocalInput() []string {
	filePath := "input.txt"
	content, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	return strings.Split(string(content), "\n")
}

func solutionPartOne(inputString []string) int {
	sum := 0
	for y, line := range inputString {
		line = strings.TrimSpace(line) + "."

		var number string
		valid := false

		for x := 0; x < len(line); x++ {
			if isDigit(line[x]) {
				number += string(line[x])
				if checkAdjacent(inputString, x, y, isSymbol) {
					valid = true
				}
			} else {
				if valid && number != "" {
					n, _ := strconv.Atoi(number)
					sum += n
				}

				number = ""
				valid = false
			}
		}
	}

	return sum
}

func solutionPartTwo(inputString []string) int {
	sum := 0
	for y, line := range inputString {
		line = strings.TrimSpace(line) + "."

		for x := 0; x < len(line); x++ {
			if line[x] == '*' {
				n := calculateGearRatio(inputString, x, y)
				sum += n
			}
		}
	}

	return sum
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isSymbol(c byte) bool {
	return !isDigit(c) && c != '.'
}

func checkAdjacent(lines []string, x, y int, checkFunc func(byte) bool) bool {
	for _, d := range square {
		dx, dy := d[0], d[1]

		if isValidCoordinate(x+dx, y+dy, lines) && checkFunc(lines[y+dy][x+dx]) {
			return true
		}
	}

	return false
}

func isValidCoordinate(x, y int, lines []string) bool {
	return y >= 0 && y < len(lines) && x >= 0 && x < len(lines[y])
}

func calculateGearRatio(lines []string, x, y int) int {
	n, n1 := 0, 0
	for _, d := range square {
		dx, dy := d[0], d[1]

		if isValidCoordinate(x+dx, y+dy, lines) {
			n = extractNumber(lines[y+dy], x+dx)
			if n > 0 {
				if n1 == 0 {
					n1 = n
					continue
				}

				if n1 != n {
					return n * n1
				}
			}
		}
	}

	return 0
}

func extractNumber(s string, x int) int {
	if x < 0 || x >= len(s) || !isDigit(s[x]) {
		return -1
	}

	number := string(s[x])
	for i := 1; x+i < len(s) && isDigit(s[x+i]); i++ {
		number += string(s[x+i])
	}

	for i := 1; x-i >= 0 && isDigit(s[x-i]); i++ {
		number = string(s[x-i]) + number
	}

	n, _ := strconv.Atoi(number)

	return n
}
