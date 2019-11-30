defmodule Puzzle do

  def solve([], acc) do
    acc
  end

  def solve([x | rest], acc) do
    case Integer.parse(x) do
      :error -> solve(rest, acc)
      {val, _} -> solve(rest, val + acc)
    end
  end

end


input = IO.read(:stdio, :all)
  |> String.split("\n")

solution = Puzzle.solve(input, 0)
IO.puts("Solution: #{solution}")
