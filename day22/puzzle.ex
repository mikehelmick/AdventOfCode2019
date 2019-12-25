defmodule Puzzle do
  require Integer

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

  def part1(input) do
    deck = generateList(0, 10007, [])
    IO.puts("#{inspect(deck)}")
    IO.puts("Length #{length(deck)} first: #{List.first(deck)} last: #{List.last(deck)}")

    deck = Puzzle.process(input, deck)

    ans = Puzzle.find(deck, 2019, 0)
    IO.puts("Answer #{ans}")
  end

  # code for part 2

  def  pow(n, k, mod), do: pow(n, k, mod, 1)
  defp pow(_, 0, mod, acc), do: rem(acc, mod)
  defp pow(n, k, mod, acc), do: pow(n, k - 1, mod, rem(n * acc, mod))

  def expand(place, target, _a, _c, acc) when place>target, do: acc
  def expand(1, target, a, c, acc) do
    expand(2, target, a, c, acc ++ [Integer.mod(a, c)])
  end
  def expand(place, target, a, c, acc) do
    last = List.last(acc)
    val = Integer.mod(last * last, c)
    expand(place + 1, target, a, c, acc ++ [Integer.mod(val, c)])
  end

  def combine_powers([], _, acc), do: acc
  def combine_powers([0|rest], [_|rest_powers], acc), do: acc
  def combine_powers([1|rest], [power|rest_powers], acc) do
    combine_powers(rest, rest_powers, acc * power)
  end

  def fastPow(a, b, c) do
    # b as list of zeros and ones
    IO.puts("Calc: #{a}^#{b} % #{c}")
    b_bin = String.graphemes(Integer.to_string(b, 2))
      |> Enum.map(fn x -> Integer.parse(x, 10) end)
      |> Enum.map(fn {x,_} -> x end)
      |> Enum.reverse()
    powers = expand(1, length(b_bin), a, c, [])
    IO.puts("Effective powers: #{inspect(powers)}")
    Integer.mod(combine_powers(b_bin, powers, 1), c)
  end

  #def inverse(n, cards) do
  #  Integer.mod(pow(n, cards-2), cards)
  #end

  def extended_gcd(a, b) do
    {last_remainder, last_x} = extended_gcd(abs(a), abs(b), 1, 0, 0, 1)
    {last_remainder, last_x * (if a < 0, do: -1, else: 1)}
  end

  defp extended_gcd(last_remainder, 0, last_x, _, _, _), do: {last_remainder, last_x}
  defp extended_gcd(last_remainder, remainder, last_x, x, last_y, y) do
    quotient   = div(last_remainder, remainder)
    remainder2 = rem(last_remainder, remainder)
    extended_gcd(remainder, remainder2, x, last_x - quotient*x, y, last_y - quotient*y)
  end

  def inverse(e, et) do
    {g, x} = extended_gcd(e, et)
    if g != 1, do: raise "The maths are broken!"
    rem(x+et, et)
  end

  def getPos(offset, increment, i, cards) do
    Integer.mod(offset + i * increment, cards)
  end

  def calcDeal(command, cards, inc, off) do
    case String.starts_with?(command, "deal into new stack") do
      true ->
        inc = inc * -1
        inc = Integer.mod(inc, cards)
        off = off + inc
        off = Integer.mod(off, cards)
        {inc, off}
      false -> {inc, off}
    end
  end

  def calcCut(command, cards, inc, off) do
    case String.starts_with?(command, "cut ") do
      true ->
        [_,n_str] = String.split(command)
        {num, _} = Integer.parse(n_str)
        off = off + num * inc
        off = Integer.mod(off, cards)
        {inc, off}
      false -> {inc, off}
    end
  end

  def calcDealIncrement(command, cards, inc, off) do
    case String.starts_with?(command, "deal with increment ") do
      true ->
        [_,_,_,n_str] = String.split(command)
        {num, _} = Integer.parse(n_str)
        inc = inc * inverse(num, cards)
        inc = Integer.mod(inc, cards)
        {inc, off}
      false -> {inc, off}
    end
  end

  def getIncrementOffset([], _cards, inc, off), do: {inc, off}
  def getIncrementOffset([line|rest], cards, inc, off) do
    {inc, off} = calcDeal(line, cards, inc, off)
    {inc, off} = calcCut(line, cards, inc, off)
    {inc, off} = calcDealIncrement(line, cards, inc, off)
    IO.puts("#{line} -> {#{inc}, #{off}}")
    getIncrementOffset(rest, cards, inc, off)
  end

  # Based on the writeup at
  # https://www.reddit.com/r/adventofcode/comments/ee0rqi/2019_day_22_solutions/fbnkaju/
  def part2(input) do
    size = 119315717514047
    rounds = 101741582076661

    {increment, offset} = getIncrementOffset(input, size, 1, 0)
    IO.puts("Inc: #{increment} Off: #{offset}")


    f_increment = fastPow(increment, rounds, size)
    f_offset = offset * (1-increment) * Integer.mod(inverse(1-increment, size), size)
    f_offset = Integer.mod(f_offset, size)

    IO.puts("Answer: #{getPos(f_offset, f_increment, 2020, size)}")
  end
end

input = IO.read(:stdio, :all)
  |> String.trim()
  |> String.split("\n")

Puzzle.part1(input)
Puzzle.part2(input)
