import std/[sequtils, strutils, strmisc, tables]

type Day12[T] = tuple[part1: T, part2: T]

type
  Record = object
    springs: string
    groups: seq[int]

type FileReadError* = object of IOError

# read file from local dir. Throws FileReadError if file not found
proc readInput(filename: string): seq[string] {.raises: [FileReadError]} =
    try:
        return filename.readFile().strip().splitLines()
    except IOError as err:
        raise FileReadError.newException("Could not read input file", err)

# calculate possible solutions for a given record
proc possibleSolutions(record: Record, memo: var Table[Record, int]): int =
  if memo.hasKey(record):
    return memo[record]

  # base case: if there are no groups, check if all springs are '#'
  if record.groups.len < 1:
    if record.springs.allIt(it != '#'):
      memo[record] = 1
      return 1
    else:
      memo[record] = 0
      return 0

  # check if the length of springs is valid
  if record.springs.len < record.groups.foldl(a+b) + record.groups.len - 1:
    memo[record] = 0
    return 0

  # if the first spring is '.', calculate solutions recursively
  if record.springs[0] == '.':
    let solutions = possibleSolutions(
      Record(
        springs: record.springs[1..^1],
        groups: record.groups
      ), memo)
    memo[record] = solutions
    return solutions

  var solutions: int
  let
    cur = record.groups[0]
    allNonOperational = record.springs[0..cur-1].allIt(it != '.')
    endd = min(cur + 1, record.springs.len)

  # check if all required springs are non-operational, and then calculate solutions
  if allNonOperational and
    ((record.springs.len() > cur and record.springs[cur] != '#') or
      record.springs.len <= cur):
        solutions = possibleSolutions(
          Record(
            springs:record.springs[endd..^1],
            groups: record.groups[1..^1]
          ),
          memo)

  # if the first spring is '?', calculate with and without replacing it
  if record.springs[0] == '?':
    solutions += possibleSolutions(
        Record(springs: record.springs[1..^1], groups: record.groups),
        memo,
      )
  memo[record] = solutions
  return solutions

# repeat a string with a separator
proc unfold(s: string, sep: string): string =
  for i in 1..5:
    result.add s
    if i < 5:
      result.add sep

# solve part 1 or part 2 for a given record
proc solvePart(record: Record, memo: var Table[Record, int]): int =
  return record.possibleSolutions(memo)

# solve both parts here
proc solve(lines: seq[string]): Day12[int] =
  var memo = initTable[Record, int]()
  var part1, part2: int

  for line in lines:
    let (springs, _, right) = line.partition(" ")
    let pattern = right.split(',').mapIt(it.parseInt)

    # part 1
    let solution1 = Record(springs: springs, groups: pattern)
    part1 += solvePart(solution1, memo)

    # part 2
    let solution2 = Record(
      springs: springs.unfold("?"),
      groups: right.unfold(",").split(',').mapIt(it.parseInt)
    )
    part2 += solvePart(solution2, memo)

  return (part1, part2)

# main enter point
try:
    let filename = "input.txt"
    let (part1, part2) = solve(readInput(filename))
    echo "Part 1 (sum counts): ", part1
    echo "Part 2 (arrangement counts): ", part2
except CatchableError as err:
    echo "Error: ", err.msg
    echo "Stack trace: ", err.getStackTrace()
