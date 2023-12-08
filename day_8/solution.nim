import std/[tables, strutils, strscans, sugar, math]

var pattern: string
var map = initTable[string, tuple[left, right: string]]()

type FileReadError* = object of IOError

# read file from local dir. Throws FileReadError if file not found
proc readInput(filename: string): seq[string] {.raises: [FileReadError]} =
    try:
        return filename.readFile().splitLines()
    except IOError as err:
        raise FileReadError.newException("Could not read input file", err)

proc solvePath(start: string, endCond: proc (s: string): bool): int =
    var nSteps = 0
    var current = start

    while not endCond(current):
        let currMap = map[current]

        case pattern[nSteps mod pattern.len]:
            of 'R': current = currMap.right
            of 'L': current = currMap.left
            else: discard

        nSteps += 1

    return nSteps

# parse each line of input and build a map of portals
proc parseInput(lines: seq[string]): seq[string] =
    var p2Start: seq[string]

    for i in 2 ..< lines.high:
        let line = lines[i]
        var key: string
        var lValue, rValue: string

        if scanf(line, "$+ = ($+, $+)", key, lValue, rValue):
            map[key] = (lValue, rValue)
            if key[2] == 'A':
                p2Start.add(key)

    return p2Start

# solve it using lcm. Didn't know about chinese remainder theorem which is pretty cool btw.
proc solve(lines: seq[string]): (int, int) =
    pattern = lines[0]
    let p2Start = parseInput(lines)
    let solvePart1 = solvePath("AAA", s => s.cmp("ZZZ") == 0)

    var paths: seq[int]
    for start in p2Start:
        var p = solvePath(start, e => e[2] == 'Z')
        paths.add(p)

    let solvePart2 = paths.lcm()
    return (solvePart1, solvePart2)

try:
    let filename = "input.txt"
    let (part1, part2) = solve(readInput(filename))
    echo "Part 1 (steps to reach ZZZ): ", part1
    echo "Part 2 (nodes ends with Z): ", part2
except CatchableError as err:
    echo "Error: ", err.msg
    echo "Stack trace: ", err.getStackTrace()
