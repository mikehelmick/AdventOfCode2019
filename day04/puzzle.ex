defmodule Puzzle do

  def to_list(num) do
    [div(num,             100000),
     div(rem(num, 100000), 10000),
     div(rem(num, 10000),   1000),
     div(rem(num, 1000),     100),
     div(rem(num, 100),       10),
     rem(num, 10),]
  end

  def rule1([x,x,_,_,_,_]), do: true
  def rule1([_,x,x,_,_,_]), do: true
  def rule1([_,_,x,x,_,_]), do: true
  def rule1([_,_,_,x,x,_]), do: true
  def rule1([_,_,_,_,x,x]), do: true
  def rule1(_), do: false

  def rule2([a,b,c,d,e,f]) when b >= a and c >= b and d >= c and e >= d and f >= e, do: true
  def rule2(_), do: false

  def solve(x, y, acc) when x > y, do: acc
  def solve(x, y, acc) do
    digits = to_list(x)
    case {rule1(digits), rule2(digits)} do
      {true, true} -> solve(x+1, y, acc ++ [x])
      _ -> solve(x+1, y, acc)
    end
  end
end


valid = Puzzle.solve(367479, 893698, [])
IO.puts("Answer #{length(valid)}")
IO.puts("Matches #{valid}")
