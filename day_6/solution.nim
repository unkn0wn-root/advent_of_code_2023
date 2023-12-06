import std/[math, sequtils, strformat, strutils]

# remove all whitespace from a string and convert it to list of ints
proc stripAndParse(s: string): seq[int] = s.splitWhitespace.mapIt(it.parseInt)

# remove all whitespace from and join it back together to be able to calculate max diff ways
proc stripAllSpace(s: string): string = s.splitWhitespace.join("")

# read the whole file into mem. Don't care about performance here since it's small txt file
proc readFromLocalFile(filename: string): string = return readFile(filename)

# func to solve the quadratic equation
proc solveQuadratic(time, dist: int): int =
  let
    tm = time.toFloat
    ds = dist.toFloat
    temp = sqrt(tm * tm - ds * 4)
    eps = 1e-12
  (floor((tm + temp) / 2 - eps) - ceil((tm - temp) / 2 + eps)).toInt + 1

# part 1 and 2 are almost the same, except for the input parsing
proc part1(data: string): int =
  var
    time, dist: seq[int]
  for line in data.splitLines():
    # There will always be Time and Distance so no need to check for anything else. Same in part2
    if line.startsWith("Time: "): time = stripAndParse(line[5..^1])
    else: dist = stripAndParse(line[9..^1])

  prod((0..time.high)
    .toSeq
    .mapIt(solveQuadratic(time[it], dist[it])))

proc part2(data: string): int =
  var
    time, dist: int
  for line in data.splitLines():
    if line.startsWith("Time: "): time = stripAllSpace(line[5..^1]).parseInt()
    else: dist = stripAllSpace(line[9..^1]).parseInt()

  solveQuadratic(time, dist)

# @toDo - better error handling for reading file and parsing input
let data = readFromLocalFile("input.txt").strip()

echo &"Part 1 - (each race): { part1(data) }"
echo &"Part 2 - (one race): { part2(data) }"
