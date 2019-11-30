defmodule Puzzle do

  def print_lines([]) do
    IO.puts("--")
  end

  def print_lines([x | rest]) do
    IO.puts(x)
    print_lines(rest)
  end

end


input = IO.read(:stdio, :all)
  |> String.split("\n")

Puzzle.print_lines(input)
