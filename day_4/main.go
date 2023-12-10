package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func numSet(s string) map[int]struct{} {
	set := make(map[int]struct{})
	for _, numStr := range strings.Fields(s) {
		num := parseInt(numStr)
		set[num] = struct{}{}
	}

	return set
}

func interSizeCount(set1, set2 map[int]struct{}) int {
	count := 0
	for num := range set1 {
		if _, exists := set2[num]; exists {
			count++
		}
	}

	return count
}

func parseInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func readSolutionLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func solve(lines []string) (part1, part2 int) {
	counts := make(map[int]int)
	for l, line := range lines {
		winners, inHand, _ := strings.Cut(line[9:], "|")
		copies := interSizeCount(numSet(inHand), numSet(winners))

		part1 += int(math.Pow(2, float64(copies-1)))
		part2++

		card := l + 1
		counts[card]++
		count := counts[card]

		for x := 1; x <= copies; x++ {
			counts[card+x] += count
			part2 += count
		}

		delete(counts, card)
	}

	return part1, part2
}

func main() {
	lines, err := readSolutionLines("input.txt")
	if err != nil {
		fmt.Println("Error reading local solution file:", err)
		return
	}

	part1, part2 := solve(lines)

	fmt.Println("Part 1 count:", part1)
	fmt.Println("Part 2 count:", part2)
}
