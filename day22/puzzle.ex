defmodule Puzzle do

  def generateList(x, y, acc) when x == y, do: acc
  def generateList(x, y, acc) do
    generateList(x+1, y, acc ++ [x])
  end

  def dealNewStack(command, deck) do
    case String.starts_with?(command, "deal into new stack") do
      true -> Enum.reverse(deck)
      false -> deck
    end
  end

  def cutNcards(command, deck) do
    case String.starts_with?(command, "cut ") do
      true ->
        [_,n_str] = String.split(command)
        {num, _} = Integer.parse(n_str)
        case num do
          x when x > 0 ->
            cut_cards = Enum.slice(deck, 0, x)
            rest = Enum.slice(deck, x..-1)
            rest ++ cut_cards
          x when x < 0 ->
            rest = Enum.slice(deck, 0..x-1)
            cut_cards = Enum.slice(deck, x..-1)
            cut_cards ++ rest
        end
      false -> deck
    end
  end

  def doIncrementDeal(map, [], _, _, _), do: map
  def doIncrementDeal(map, [card|rest], size, inc, pos) do
    pos = Integer.mod(pos, size)
    doIncrementDeal(Map.put(map, pos, card), rest, size, inc, pos+inc)
  end

  def dealWithIncrement(command, deck) do
    case String.starts_with?(command, "deal with increment ") do
      true ->
        [_,_,_,n_str] = String.split(command)
        {num, _} = Integer.parse(n_str)
        len = length(deck)
        map = doIncrementDeal(Map.new(), deck, len, num, 0)
        Enum.map(Enum.sort(Map.keys(map)), fn x -> Map.get(map, x) end)
      false -> deck
    end
  end

  def process([], deck), do: deck
  def process([command|rest], deck) do
    deck = dealNewStack(command, deck)
    deck = cutNcards(command, deck)
    deck = dealWithIncrement(command, deck)
    process(rest, deck)
  end

  def find([], _, _), do: nil
  def find([x|_], target, pos) when x == target, do: pos
  def find([_|rest], target, pos), do: find(rest, target, pos+1)

  def part1(deck, input) do
    IO.puts("#{inspect(deck)}")
    IO.puts("Length #{length(deck)} first: #{List.first(deck)} last: #{List.last(deck)}")

    deck = Puzzle.process(input, deck)

    ans = Puzzle.find(deck, 2019, 0)
    IO.puts("Answer #{ans}")
  end

  def doRound(deck, _, 0, _), do: deck
  def doRound(deck, input, times, origin) do
    if Integer.mod(times,100) do
      IO.puts("--> #{times}")
    end
    deck = Puzzle.process(input, deck)

    hash = :crypto.hash(:md5, deck) |> Base.encode16()
    if String.equivalent?(hash, origin) do
      IO.puts("Back to original with #{times} ")
    end

    doRound(deck, input, times-1, origin)
  end

  def part2(deck, input) do
    origin = :crypto.hash(:md5, deck) |> Base.encode16()
    IO.puts("Hash of original order: #{origin}")
    deck = doRound(deck, input, 101741582076661, origin)

    ans = Enum.at(deck, 2020)
    IO.puts("Answer: #{ans}")
  end
end

input = IO.read(:stdio, :all)
  |> String.trim()
  |> String.split("\n")

# part 1 -> 10007
# part 2 -> 119315717514047
deck = Puzzle.generateList(0, 119315717514047, [])

# Puzzle.part1(deck, input)
Puzzle.part2(deck, input)
