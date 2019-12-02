defmodule PuzzleTwo do

  def solve_one(val, acc) do
    fuel = Integer.floor_div(val,3)-2
    case fuel do
      x when x <= 0 -> acc
      x -> solve_one(x, x + acc)
    end
  end

  def solve([], acc) do
    acc
  end
  def solve([x | rest], acc) do
    case Integer.parse(x) do
      :error -> solve(rest, acc)
      {val, _} -> solve(rest, solve_one(val,0) + acc)
    end
  end
end

input = IO.read(:stdio, :all)
  |> String.split("\n")

solution = PuzzleTwo.solve(input, 0)
IO.puts("Solution: #{solution}")
