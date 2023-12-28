package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// represents a map of rune to int.
type Part map[rune]int

// a workflow rule.
type Rule struct {
	category, operator rune
	right              int
	consequence        string
}

// collection of rules.
type Workflow []Rule

func readLocalFile(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// parses a rule string and returns a Rule.
func parseRule(rule string) Rule {
	var r Rule
	if pos := strings.Index(rule, ":"); pos > -1 {
		r.category, r.operator, r.right = rune(rule[0]), rune(rule[1]), toInt(rule[2:pos])
		r.consequence = rule[pos+1:]
	} else {
		r.consequence = rule
	}

	return r
}

func toInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

// parses a part string and returns a Part.
func parsePart(line string) Part {
	var x, m, a, s int
	fmt.Sscanf(line, "{x=%d,m=%d,a=%d,s=%d}", &x, &m, &a, &s)
	return Part{'x': x, 'm': m, 'a': a, 's': s}
}

// parses the input and returns workflows and parts.
func parse(input []byte) (map[string]Workflow, []Part) {
	sections := strings.Split(string(input), "\n\n")

	// parse workflows
	workflows := make(map[string]Workflow)
	for _, line := range strings.Split(sections[0], "\n") {
		name, line := line[:strings.Index(line, "{")], line[strings.Index(line, "{")+1:len(line)-1]
		w := make(Workflow, 0, 4)
		for _, rule := range strings.Split(line, ",") {
			w = append(w, parseRule(rule))
		}

		workflows[name] = w
	}

	// parse parts
	parts := make([]Part, 0)
	for _, line := range strings.Split(sections[1], "\n") {
		parts = append(parts, parsePart(line))
	}

	return workflows, parts
}

// applies the given workflow on a part and returns the result.
func applyWorkflow(workflows map[string]Workflow, workflow string, part Part) bool {
	if workflow == "R" {
		return false
	} else if workflow == "A" {
		return true
	}

	for _, r := range workflows[workflow] {
		if evaluateRule(r, part) {
			return applyWorkflow(workflows, r.consequence, part)
		}
	}

	return false
}

// evaluates a rule against a part and returns the result.
func evaluateRule(r Rule, part Part) bool {
	switch r.operator {
	case '>':
		return part[r.category] > r.right
	case '<':
		return part[r.category] < r.right
	default:
		return true
	}
}

func pt1(workflows map[string]Workflow, parts []Part) int {
	sum := 0
	for _, p := range parts {
		if applyWorkflow(workflows, "in", p) {
			sum += sumPartValues(p)
		}
	}

	return sum
}

// returns the sum of values in a part.
func sumPartValues(p Part) int {
	sum := 0
	for _, v := range p {
		sum += v
	}

	return sum
}

// calculates the result for Part 2.
func count(workflows map[string]Workflow, workflow string, values map[rune][2]int) int {
	if workflow == "R" {
		return 0
	} else if workflow == "A" {
		return calculateProduct(values)
	}

	total := 0
	for _, r := range workflows[workflow] {
		v := values[r.category]
		tv, fv := getTrueFalseRanges(r, v)

		if tv[0] <= tv[1] {
			v2 := cloneMap(values)
			v2[r.category] = tv
			total += count(workflows, r.consequence, v2)
		}

		if fv[0] > fv[1] {
			break
		}

		values[r.category] = fv
	}

	return total
}

// returns true and false ranges for a rule.
func getTrueFalseRanges(r Rule, v [2]int) ([2]int, [2]int) {
	var tv, fv [2]int

	if r.operator == '<' {
		tv = [2]int{v[0], r.right - 1}
		fv = [2]int{r.right, v[1]}
	} else if r.operator == '>' {
		tv = [2]int{r.right + 1, v[1]}
		fv = [2]int{v[0], r.right}
	}

	return tv, fv
}

// calculates the product of ranges in values.
func calculateProduct(values map[rune][2]int) int {
	product := 1
	for _, v := range values {
		product *= (v[1] - v[0] + 1)
	}

	return product
}

// clones the values map.
func cloneMap(original map[rune][2]int) map[rune][2]int {
	c := make(map[rune][2]int, len(original))
	for k, v := range original {
		c[k] = v
	}

	return c
}

func main() {
	input, err := readLocalFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	workflows, parts := parse(input)

	part1 := pt1(workflows, parts)
	part2 := count(workflows, "in", map[rune][2]int{
		'x': {1, 4000},
		'm': {1, 4000},
		'a': {1, 4000},
		's': {1, 4000},
	})

	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}
