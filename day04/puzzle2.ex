defmodule Puzzle do

  def ascending([a,b,c,d,e,f]) when b >= a and c >= b and d >= c and e >= d and f >= e, do: true
  def ascending(_), do: false

  def short_pair([a,a,b,_,_,_]) when a != b, do: true
  def short_pair([a,b,b,c,_,_]) when a != b and b != c, do: true
  def short_pair([_,a,b,b,c,_]) when a != b and b != c, do: true
  def short_pair([_,_,a,b,b,c]) when a != b and b != c, do: true
  def short_pair([_,_,_,a,b,b]) when a != b, do: true
  def short_pair(_), do: false

  def solve(x, y, acc) when x > y, do: acc
  def solve(x, y, acc) do
    digits = Integer.digits(x)
    case {short_pair(digits), ascending(digits)} do
      {true, true} -> solve(x+1, y, acc ++ [x])
      _ -> solve(x+1, y, acc)
    end
  end
end

valid = Puzzle.solve(367479, 893698, [])
IO.puts("Answer #{length(valid)}")
IO.puts("Matches #{inspect(valid)}")
