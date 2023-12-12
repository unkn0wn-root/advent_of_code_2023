package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"strings"
)

func readLocalFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(content)), nil
}

func main() {
	fileContent, err := readLocalFile("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	getDistances := func(expand int) (totalDistance int) {
		galaxies := []image.Point{}
		disy := 0

		for y, row := range fileContent {
			if !strings.Contains(row, "#") {
				disy += expand - 1
			}

			disx := 0
			for x, char := range row {
				col := ""
				for _, s := range fileContent {
					col += string(s[x])
				}

				if !strings.Contains(col, "#") {
					disx += expand - 1
				}

				if char == '#' {
					for _, g := range galaxies {
						totalDistance += int(math.Abs(float64(x+disx-g.X)) + math.Abs(float64(y+disy-g.Y)))
					}
					galaxies = append(galaxies, image.Point{x + disx, y + disy})
				}
			}
		}

		return totalDistance
	}

	fmt.Println("Part 1 (sum lengths):", getDistances(2))
	fmt.Println("Part 2 (sum lengths):", getDistances(1000000))
}
