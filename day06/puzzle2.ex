defmodule Puzzle do

  def update_entry(nil, val), do: {nil, [val]}
  def update_entry(l, val), do: {l, l ++ [val]}

  def build_edges([], o_map, i_map, allobjs), do: {o_map, i_map, MapSet.to_list(allobjs)}
  def build_edges([line | rest], o_map, i_map, allobjs) do
    [inner, outer] = String.split(line, ")")
    {_, updates_o_map} = Map.get_and_update(o_map, inner,
                 fn x -> update_entry(x, outer) end)
    updated_i_map = Map.put(i_map, outer, inner)
    new_all = MapSet.put(allobjs, inner) |> MapSet.put(outer)
    build_edges(rest, updates_o_map, updated_i_map, new_all)
  end

  def crawl("COM", _, acc), do: acc
  def crawl(obj, map, acc) do
    crawl(Map.get(map, obj), map, acc + 1)
  end

  def count_edges([], _, acc), do: acc
  def count_edges(["COM"], _, acc), do: acc
  def count_edges(["COM"|rest], map, acc), do: count_edges(rest, map, acc)
  def count_edges([obj|rest], map, acc) do
    IO.puts("counting #{obj}, current #{acc}")
    count_edges(rest, map, crawl(obj, map, 0) + acc)
  end

  def orbital_path("COM", _, lst), do: lst
  def orbital_path(start, map, lst) do
    next = Map.get(map, start)
    orbital_path(next, map, lst ++ [next])
  end

  def index([{x, idx}|_], x), do: idx
  def index([{_, _}|rest], x), do: index(rest, x)
end


input = IO.read(:stdio, :all)
  |> String.trim()
  |> String.split("\n")

{_out_edges, in_edges, _allobjs} = Puzzle.build_edges(input, %{}, %{}, MapSet.new())

#IO.puts("all #{inspect(allobjs)}")
#IO.puts("edges #{inspect(in_edges)}")

y_path = Puzzle.orbital_path("YOU", in_edges, ["YOU"])
s_path = Puzzle.orbital_path("SAN", in_edges, ["SAN"])

#IO.puts("YOU: #{inspect(y_path)}")
#IO.puts("SAN: #{inspect(s_path)}")

intersection = MapSet.intersection(MapSet.new(y_path), MapSet.new(s_path)) |> MapSet.to_list()
#IO.puts("int: #{inspect(intersection)}")

y_idx = Enum.with_index(y_path)
s_idx = Enum.with_index(s_path)

distances = Enum.map(intersection,
  fn x -> {x, Puzzle.index(y_idx, x) + Puzzle.index(s_idx, x)} end)
  |> Enum.sort(
      fn ({_,a}, {_,b}) -> a < b end)
IO.puts("dist: #{inspect(distances)}")

[{_ancestor,answer}|_] = distances
# need to subtract 2 as orgin bodies are objects
IO.puts("ANSWER: #{answer-2}")
