package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Solution struct{}

func (s *Solution) PartOne(input string) interface{} {
	return s.Solve(input, `\d`)
}

func (s *Solution) PartTwo(input string) interface{} {
	return s.Solve(input, `\d|one|two|three|four|five|six|seven|eight|nine`)
}

func readInputText(name string) string {
	data, err := os.ReadFile(name)

	if err != nil {
		panic(err)
	}

	return string(data)
}

func (s *Solution) Solve(input string, rx string) int {
	lines := strings.Split(strings.Trim(input, " "), "\n")
	sum := 0

	for _, line := range lines {
		first := regexp.MustCompile(rx).FindString(line)
		last := regexp.MustCompile(rx).FindString(line)

		match := regexp.MustCompile(rx).FindAllString(line, -1)

		if len(match) > 0 {
			last = match[len(match)-1]
		}

		sum += s.ParseMatch(first)*10 + s.ParseMatch(last)
	}

	return sum
}

func (s *Solution) ParseMatch(st string) int {
	switch st {
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	default:
		d, _ := strconv.Atoi(st)
		return d
	}
}

func main() {
	input := readInputText("input.txt")
	solver := &Solution{}
	fmt.Println("Part One:", solver.PartOne(input))
	fmt.Println("Part Two:", solver.PartTwo(input))
}
