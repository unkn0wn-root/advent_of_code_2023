type Result = {
    part1: number;
    part2: number;
};

// represents a camera lens with a label and focal length
class Lens {
  label: string;
  focalLength: number;

  constructor(label: string, focalLength: number) {
    this.label = label;
    this.focalLength = focalLength;
  }
}

function getLines(filePath: string): string[] {
  try {
    const content: string = Deno.readTextFileSync(filePath).trim();
    return content.split('\n').map(line => line.trim());
  } catch (error) {
    console.error(`Error reading file: ${error.message}`);
    return [];
  }
}

function calculateResults(): Result {
  const input: string[] = getLines("input.txt");
  const sequence: string[] = input[0].split(",");

  // Part 1: calculate hash sum for the sequence
  let hashSum: number = 0;
  sequence.forEach(step => hashSum += calculateHash(step));

  // Part 2: process sequence to calculate focusing power
  const boxes: Lens[][] = new Array(256).fill([]).map(() => []);

  for (const step of sequence) {
    // if the step contains "-", it's a removal operation
    if (step.includes("-")) {
      const data: string[] = step.split("-");
      const box: number = calculateHash(data[0]);
      removeFromBox(boxes, box, data[0]);
    } else {
      // if the step contains "=", it's an addition or update operation
      const data: string[] = step.split("=");
      const box: number = calculateHash(data[0]);
      const focal: number = parseInt(data[1], 10);
      changeBox(boxes, box, data[0], focal);
    }
  }

  let focusingPower: number = 0;
  for (let boxNum = 0; boxNum < boxes.length; boxNum++) {
    const box: Lens[] = boxes[boxNum];
    for (let i = 0; i < box.length; i++) {
      const lens: Lens = box[i];
      const power: number = (boxNum + 1) * (i + 1) * lens.focalLength;
      focusingPower += power;
    }
  }

  return { part1: hashSum, part2: focusingPower };
}

// adds or updates a lens in the specified box
function changeBox(boxes: Lens[][], box: number, label: string, focal: number): void {
  for (let i = 0; i < boxes[box].length; i++) {
    const lens: Lens = boxes[box][i];
    // if lens with the same label exists, update its focal length
    if (lens.label === label) {
      boxes[box][i].focalLength = focal;
      return;
    }
  }
  // if lens with the given label doesn't exist, add a new lens to the box
  boxes[box].push(new Lens(label, focal));
}

// removes a lens from the specified box based on its label
function removeFromBox(boxes: Lens[][], box: number, label: string): void {
  for (let i = 0; i < boxes[box].length; i++) {
    const lens: Lens = boxes[box][i];
    if (lens.label === label) {
      boxes[box].splice(i, 1);
      break;
    }
  }
}

// calculates a hash value for the given data string
function calculateHash(data: string): number {
  let result: number = 0;
  for (const char of data) {
    // update the result based on the ASCII value of each character
    result += char.charCodeAt(0);
    result *= 17;
    result %= 256;
  }
  return result;
}

function main(): void {
  const results: Result = calculateResults();
  console.log(`Part 1: ${results.part1}`);
  console.log(`Part 2: ${results.part2}`);
}

main();
