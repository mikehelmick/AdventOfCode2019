defmodule Puzzle do

  # Ensure that there is a pair
  def pair([x,x,_,_,_,_]), do: true
  def pair([_,x,x,_,_,_]), do: true
  def pair([_,_,x,x,_,_]), do: true
  def pair([_,_,_,x,x,_]), do: true
  def pair([_,_,_,_,x,x]), do: true
  def pair(_), do: false

  # Ensure ascenting
  def ascending([a,b,c,d,e,f]) when b >= a and c >= b and d >= c and e >= d and f >= e, do: true
  def ascending(_), do: false

  def solve(x, y, acc) when x > y, do: acc
  def solve(x, y, acc) do
    digits = Integer.digits(x)
    case {pair(digits), ascending(digits)} do
      {true, true} -> solve(x+1, y, acc ++ [x])
      _ -> solve(x+1, y, acc)
    end
  end
end

valid = Puzzle.solve(367479, 893698, [])
IO.puts("Answer #{length(valid)}")
IO.puts("Matches #{inspect(valid)}")
