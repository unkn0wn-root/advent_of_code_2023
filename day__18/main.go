package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"strconv"
)

// represents a change in coordinates based on a direction or digit.
type delta map[string]image.Point

func readInputFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func compileRegex() *regexp.Regexp {
	return regexp.MustCompile(`(.) (.*) \(#(.*)(.)\)`)
}

func initializeDelta() delta {
	// initializes and returns the Delta map for direction mappings.
	return delta{
		"R": {1, 0}, "D": {0, 1}, "L": {-1, 0}, "U": {0, -1},
		"0": {1, 0}, "1": {0, 1}, "2": {-1, 0}, "3": {0, -1},
	}
}

// calculates the area based on the parsed input, delta, and parameters.
func calculateArea(input string, regex *regexp.Regexp, delta delta, directionIdx, lengthIdx, base int) int {
	currentPosition, totalArea := image.Point{0, 0}, 0
	for _, match := range regex.FindAllStringSubmatch(input, -1) {
		length, _ := strconv.ParseInt(match[lengthIdx], base, strconv.IntSize)
		newPosition := currentPosition.Add(delta[match[directionIdx]].Mul(int(length)))

		// calculate the area using the Shoelace formula
		totalArea += currentPosition.X*newPosition.Y - currentPosition.Y*newPosition.X + int(length)

		currentPosition = newPosition
	}

	return totalArea/2 + 1
}

func main() {
	input, err := readInputFile("input.txt")
	if err != nil {
		log.Fatal("Could not read the file. Error: ", err)
	}

	regex := compileRegex()
	delta := initializeDelta()

	fmt.Println("Part 1:", calculateArea(input, regex, delta, 1, 2, 10))
	fmt.Println("Part 2:", calculateArea(input, regex, delta, 4, 3, 16))
}
