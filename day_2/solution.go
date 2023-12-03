package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	data := strings.TrimSpace(readLocalInput())
	parsed := parseData(data)
	firstPart := solutionPart1(parsed)
	secPart := solutionPart2(parsed)

	fmt.Println("Part 1 count:", firstPart)
	fmt.Println("Part 2 count:", secPart)
}

type subset struct {
	Red   int
	Green int
	Blue  int
}

type gameData struct {
	Game    int
	Subsets []subset
}

func readLocalInput() string {
	filePath := "input.txt"
	content, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func parseData(data string) []gameData {
	lines := strings.Split(data, "\n")
	parsed := make([]gameData, len(lines))

	for i, line := range lines {
		sets := strings.Split(line, ":")
		gameNumber, _ := strconv.Atoi(strings.Fields(sets[0])[1])
		subsets := make([]subset, 0)

		for _, sub := range strings.Split(sets[1], ";") {
			oneSubset := make(map[string]int)
			for _, c := range strings.Split(sub, ",") {
				c = strings.TrimSpace(c)
				parts := strings.Split(c, " ")
				number, _ := strconv.Atoi(parts[0])
				color := parts[1]
				oneSubset[color] = number
			}

			subsets = append(subsets, subset{
				Red:   oneSubset["red"],
				Green: oneSubset["green"],
				Blue:  oneSubset["blue"],
			})
		}

		parsed[i] = gameData{
			Game:    gameNumber,
			Subsets: subsets,
		}
	}

	return parsed
}

func solutionPart1(parsed []gameData) int {
	sum := 0
	for _, it := range parsed {
		ok := true

		for _, curr := range it.Subsets {
			ok = ok && curr.Red <= 12 && curr.Green <= 13 && curr.Blue <= 14
		}

		if ok {
			sum += it.Game
		}
	}

	return sum
}

func solutionPart2(parsed []gameData) int {
	sum := 0
	for _, game := range parsed {
		max := subset{Red: 0, Green: 0, Blue: 0}
		for _, curr := range game.Subsets {
			if curr.Red > max.Red {
				max.Red = curr.Red
			}

			if curr.Green > max.Green {
				max.Green = curr.Green
			}

			if curr.Blue > max.Blue {
				max.Blue = curr.Blue
			}
		}

		sum += max.Red * max.Green * max.Blue
	}

	return sum
}
