import std/[sequtils, strutils, algorithm, tables]

# define our types and some labels
type
  CardLabel = enum
    Two = "2", Three = "3", Four = "4", Five = "5", Six = "6", Seven = "7",
    Eight = "8", Nine = "9", Ten = "T", Joker = "J", Queen = "Q", King = "K", Ace = "A"

  HandType = enum
    HighCard, OnePair, TwoPair, ThreeOfKind, FullHouse, FourOfKind, FiveOfKind

  Hand = object
    hand_type: HandType
    hand: seq[CardLabel]
    bet: int

# determine the hand type of a given hand
proc determineHandType(hand: seq[CardLabel], wild: bool = false): HandType =
  var handTable = initTable[CardLabel, int]()
  var duplicates: range[0..2] = 0
  var triplicate, fourOfKind, fiveOfKind = false
  var numJokers: range[0..5] = 0

  for card in hand:
    # don't consider jokers for any parts of hand types until later if wild.
    if card == Joker and wild:
      numJokers += 1
    else:
      let labelCount = handTable.getOrDefault(card)
      handTable[card] = labelCount + 1

  for (_, count) in handTable.pairs:
    case count:
      of 5: fiveOfKind = true
      of 4: fourOfKind = true
      of 3: triplicate = true
      of 2: duplicates += 1
      else: discard

  # Jokers are not wild or we have none.
  # this does not look particularly nice but it works! Should maybe refactor some day...
  if not wild or numJokers == 0:
    if fiveOfKind: return FiveOfKind
    if fourOfKind: return FourOfKind
    if duplicates != 0 and triplicate: return FullHouse
    if triplicate: return ThreeOfKind
    if duplicates > 1: return TwoPair
    if duplicates == 1: return OnePair
    return HighCard

  case numJokers:
    of 0: raiseAssert("Wait what?! Cover all cases but should never get here!")
    of 1:
      if fourOfKind: return FiveOfKind
      if triplicate: return FourOfKind
      if duplicates > 1: return FullHouse
      if duplicates == 1: return ThreeOfKind
      return OnePair
    of 2:
      if triplicate: return FiveOfKind
      if duplicates == 1: return FourOfKind
      else: return ThreeOfKind
    of 3:
      if duplicates != 0: return FiveOfKind
      else: return FourOfKind
    of 4..5: return FiveOfKind

# compare two hands
proc compareHands(a: Hand, b: Hand, wild: bool = false): int =
  for i in 0 ..< a.hand.len:
    let aOrd =
        if a.hand[i] == Joker and wild: -1
        else: ord(a.hand[i])
    let bOrd =
        if b.hand[i] == Joker and wild: -1
        else: ord(b.hand[i])

    if aOrd > bOrd: return 1
    elif bOrd > aOrd: return -1

# parse a hand from a string
func parseHand(line: string, wild: bool = false): Hand =
  let hand = line.split[0]
  var result: Hand
  result.bet = line.split[1].parseInt

  for card in hand:
    result.hand.add(parseEnum[CardLabel]($card))

  result.hand_type = determineHandType(result.hand, wild)
  return result

# calculate the sum of bets from a sequence of lines
func sumBetsFromLines(lines: seq[string], wild: bool = false): int =
  var handTypes = initTable[HandType, seq[Hand]]()

  for line in lines:
    let hand = line.parseHand(wild)
    handTypes.mgetOrPut(hand.hand_type, @[]).add(hand)

  var currentRank = 1
  var result = 0

  for handType in HandType.items:
    if handTypes.getOrDefault(handType).len > 0:
      var hands = handTypes[handType]
      hands.sort(proc(a, b: Hand): int = compareHands(a, b, wild))

      for hand in hands.mitems:
        result += currentRank * hand.bet
        currentRank += 1

  return result

try:
  let input = open("input.txt")
  defer: close(input)

  let lines = input.lines.toSeq

  echo("Part One: " & $sumBetsFromLines(lines))
  echo("Part Two: " & $sumBetsFromLines(lines, wild = true))
except IOError:
  echo("Error: Unable to read input file")

