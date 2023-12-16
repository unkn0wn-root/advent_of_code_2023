const rock: number = 'O'.charCodeAt(0);
const empty: number = '.'.charCodeAt(0);

const directions: number[][] = [
  [0, -1], // N
  [-1, 0], // W
  [0, 1],  // S
  [1, 0],  // E
];

// move rocks in the specified direction
function tiltRocks(grid: number[][], dir: number[]): void {
  const width: number = grid[0].length;
  const height: number = grid.length;

  // determine the starting point and step for iteration
  let startCol: number = 0, colStep: number = 1;
  if (dir[0] === 1) {
    startCol = width - 1;
    colStep = -1;
  }

  let startRow: number = 0, rowStep: number = 1;
  if (dir[1] === 1) {
    startRow = height - 1;
    rowStep = -1;
  }

  // iterate through the grid and tilt rocks
  for (let row: number = startRow; row < height && row >= 0; row += rowStep) {
    for (let col: number = startCol; col < width && col >= 0; col += colStep) {
      if (grid[row][col] !== rock) {
        continue;
      }

      // move the rock in the specified direction until an obstacle is encountered
      let c: number = col, r: number = row;
      while (
        c + dir[0] >= 0 && c + dir[0] < width &&
        r + dir[1] >= 0 && r + dir[1] < height &&
        grid[r + dir[1]][c + dir[0]] === empty
      ) {
        c += dir[0];
        r += dir[1];
      }

      // place the rock at the new position
      if (r !== row || c !== col) {
        grid[r][c] = rock;
        grid[row][col] = empty;
      }
    }
  }
}

// calculates the weight of rocks in the grid
function calculateWeight(grid: number[][]): number {
  const width: number = grid[0].length;
  const height: number = grid.length;
  let weight: number = 0;

  for (let row: number = 0; row < height; row++) {
    for (let col: number = 0; col < width; col++) {
      if (grid[row][col] === rock) {
        weight += height - row;
      }
    }
  }

  return weight;
}

// checks if two grids are equal
function areGridsEqual(a: number[][], b: number[][]): boolean {
  for (let row: number = 0; row < a.length; row++) {
    if (!a[row].every((value, col) => value === b[row][col])) {
      return false;
    }
  }

  return true;
}

// creates a deep copy of the grid
function copyGrid(a: number[][]): number[][] {
  return a.map(row => row.slice());
}

// applies the tilt cycle to the grid
function tiltCycle(grid: number[][]): void {
  for (const dir of directions) {
    tiltRocks(grid, dir);
  }
}

// reads input from a file and returns the grid
function readInputFromFile(filename: string): number[][] {
  const content: string = Deno.readTextFileSync(filename).trim();
  return content.split('\n').map(line => line.split('').map(char => char.charCodeAt(0)));
}

// applies Floyd's Tortoise and Hare algorithm to the grid
function floyd(func: (grid: number[][]) => void, x0: number[][]): [number, number[][]] {
  const hare: number[][] = copyGrid(x0);
  const tortoise: number[][] = copyGrid(x0);

  func(tortoise);
  func(hare);
  func(hare);

  while (!areGridsEqual(tortoise, hare)) {
    func(tortoise);
    func(hare);
    func(hare);
  }

  func(hare);

  // move hare forward until the start of the next cycle
  // tortoise remains in place
  let cycleLength: number = 1;
  while (!areGridsEqual(tortoise, hare)) {
    func(hare);
    cycleLength++;
  }

  // return the length of the cycle and the current state of the grid
  return [cycleLength, hare];
}

function main(): void {
  try {
    const grid: number[][] = readInputFromFile("input.txt");

    // Part 1
    const grid1: number[][] = copyGrid(grid);
    tiltRocks(grid1, directions[0]);
    const pt1: number = calculateWeight(grid1);

    // Part 2
    const cycle = (g: number[][]): void => {
      tiltCycle(g);
    };
    const [cycleLength, state] = floyd(cycle, grid);
    const cyclesRemaining: number = 1e9 % cycleLength;

    for (let i: number = 0; i < cyclesRemaining; i++) {
      tiltCycle(state);
    }

    const pt2: number = calculateWeight(state);

    console.log("Part 1:", pt1);
    console.log("Part 2:", pt2);
  } catch (error) {
    console.error(error.message);
  }
}

main();
