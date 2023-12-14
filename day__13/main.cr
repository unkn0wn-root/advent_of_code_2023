def read_input(file_path)
  begin
    File.read(file_path)
  rescue ex : Exception
    puts "Error reading the input file: #{ex.message}"
    exit(1)
  end
end

# calculates the difference between two rows
def row_diff(row1, row2)
  row1.zip(row2).count { |e1, e2| e1 != e2 }
end

# finds the reflection row
def find_ref_row(grid, required_smudges)
  max_mid = -1
  (0...grid.size - 1).each do |middle|
    up = middle
    down = up + 1

    is_reflection = true
    smudges = 0

    while up >= 0 && down <= grid.size - 1
      diff = row_diff(grid[up], grid[down])
      smudges += diff
      if smudges > required_smudges
        is_reflection = false
        break
      end
      up -= 1
      down += 1
    end

    max_mid = middle if smudges == required_smudges && is_reflection && (up <= 0 || down >= grid.size - 1)
  end
  max_mid + 1
end

# parses the input into grids
def parse(input)
  input.chunk do |line|
    /\A\s*\z/ !~ line || nil
  end.map do |_, lines|
    lines.map(&.strip).map { |line| line.chars }
  end
end

# finds the reflection column
def find_ref_column(grid, required_smudges)
  find_ref_row(grid.transpose, required_smudges)
end

# calculates the score for a grid
def score(grid, required_smudges)
  ref_row = find_ref_row(grid, required_smudges)
  ref_column = find_ref_column(grid, required_smudges)
  100 * ref_row + ref_column
end

def main
  file_path = "input.txt"
  input = read_input(file_path)
  grids = parse(input)

  part1 = grids.map { |grid| score(grid, 0) }.sum
  part2 = grids.map { |grid| score(grid, 1) }.sum

  puts "Part 1: #{part1}"
  puts "Part 2: #{part2}"
end

# start main func
main
