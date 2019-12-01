defmodule Puzzle do

  def solve_one(val) do
    Integer.floor_div(val,3)-2
  end

  def solve([], acc) do
    acc
  end

  def solve([x | rest], acc) do
    case Integer.parse(x) do
      :error -> solve(rest, acc)
      {val, _} -> solve(rest, solve_one(val) + acc)
    end
  end

end


input = IO.read(:stdio, :all)
  |> String.split("\n")

solution = Puzzle.solve(input, 0)
IO.puts("Solution: #{solution}")
