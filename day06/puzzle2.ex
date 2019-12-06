defmodule Puzzle do

  def update_entry(nil, val), do: {nil, [val]}
  def update_entry(l, val), do: {l, l ++ [val]}

  # Incoming and outgoing edges, from a stream of edges.
  def build_edges([], o_map, i_map, allobjs), do: {o_map, i_map, MapSet.to_list(allobjs)}
  def build_edges([line | rest], o_map, i_map, allobjs) do
    [inner, outer] = String.split(line, ")")
    {_, updates_o_map} = Map.get_and_update(o_map, inner,
                 fn x -> update_entry(x, outer) end)
    updated_i_map = Map.put(i_map, outer, inner)
    new_all = MapSet.put(allobjs, inner) |> MapSet.put(outer)
    build_edges(rest, updates_o_map, updated_i_map, new_all)
  end

  # Get the path from an object back to COM as a list.
  def orbital_path("COM", _, lst), do: lst
  def orbital_path(start, map, lst) when is_list(lst) do
    next = Map.get(map, start)
    orbital_path(next, map, lst ++ [next])
  end
end


input = IO.read(:stdio, :all)
  |> String.trim()
  |> String.split("\n")

{_out_edges, in_edges, _allobjs} = Puzzle.build_edges(input, %{}, %{}, MapSet.new())

#IO.puts("all #{inspect(allobjs)}")
#IO.puts("edges #{inspect(in_edges)}")

# Get the orbital path for each item.
y_path = Puzzle.orbital_path("YOU", in_edges, ["YOU"])
s_path = Puzzle.orbital_path("SAN", in_edges, ["SAN"])

#IO.puts("YOU: #{inspect(y_path)}")
#IO.puts("SAN: #{inspect(s_path)}")

# Find the intersection. First item in the intersection should be the LCA.
intersection = MapSet.intersection(MapSet.new(y_path), MapSet.new(s_path)) |> MapSet.to_list()
#IO.puts("int: #{inspect(intersection)}")

# For each item in the intersection set, calculate the distance from target to ancestor
# Sort by second item in tuple so first is answer.
distances = Enum.map(intersection,
  fn x -> {x, Enum.find_index(y_path, fn a -> a == x end) + Enum.find_index(s_path, fn a -> a == x end) - 2} end)
  |> Enum.sort(
      fn ({_,a}, {_,b}) -> a < b end)
IO.puts("dist: #{inspect(distances)}")

[{_ancestor,answer}|_] = distances
# need to subtract 2 as orgin bodies are objects
IO.puts("ANSWER: #{answer}")
