package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInputFile(filePath string) (*bufio.Scanner, *os.File, error) {
	if file, err := os.Open(filePath); err != nil {
		return nil, nil, err
	} else {
		return bufio.NewScanner(file), file, nil
	}
}

func solve(scanner *bufio.Scanner) (int, int) {
	sumPart1, sumPart2 := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		nums := stringToIntSlice(strings.Split(line, " "))
		lastVals := getLastGenerationValues(nums)
		firstVal, lastVal := calculateSums(lastVals)
		sumPart1 += lastVal
		sumPart2 += firstVal
	}

	return sumPart1, sumPart2
}

func getLastGenerationValues(nums []int) [][]int {
	lastVals := [][]int{{nums[0], nums[len(nums)-1]}}
	for {
		newNums := make([]int, 0, len(nums)-1)
		allZero := true
		for i := 0; i < len(nums)-1; i++ {
			diff := nums[i+1] - nums[i]
			if diff != 0 {
				allZero = false
			}

			newNums = append(newNums, diff)
		}

		nums = newNums
		lastVals = append(lastVals, []int{nums[0], nums[len(nums)-1]})

		if allZero {
			break
		}
	}

	return lastVals
}

func calculateSums(lastVals [][]int) (int, int) {
	firstVal, lastVal := 0, 0
	for i := len(lastVals) - 1; i >= 0; i-- {
		previousVals := lastVals[i]
		firstVal = previousVals[0] - firstVal
		lastVal += previousVals[1]
	}

	return firstVal, lastVal
}

func stringToIntSlice(strs []string) []int {
	nums := make([]int, len(strs))
	for i, str := range strs {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)
			return nil
		}

		nums[i] = num
	}

	return nums
}

func main() {
	inputFilename := "input.txt"
	scanner, file, err := readInputFile(inputFilename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	defer file.Close()

	sum1, sum2 := solve(scanner)

	fmt.Println("Part 1 (next value, each):", sum1)
	fmt.Println("Part 2 (previous value, each):", sum2)
}
