defmodule Puzzle do

  def update_entry(nil, val), do: {nil, [val]}
  def update_entry(l, val), do: {l, l ++ [val]}

  # builds incoming and outgoing edges representing the orbital graph
  def build_edges([], o_map, i_map, allobjs), do: {o_map, i_map, MapSet.to_list(allobjs)}
  def build_edges([line | rest], o_map, i_map, allobjs) do
    IO.puts("line #{line}")
    [inner, outer] = String.split(line, ")")
    {_, updates_o_map} = Map.get_and_update(o_map, inner,
                 fn x -> update_entry(x, outer) end)
    updated_i_map = Map.put(i_map, outer, inner)
    new_all = MapSet.put(allobjs, inner) |> MapSet.put(outer)
    build_edges(rest, updates_o_map, updated_i_map, new_all)
  end

  # Counts edges from a body to COM
  def crawl("COM", _, acc), do: acc
  def crawl(obj, map, acc) do
    IO.puts("crawl #{obj}")
    crawl(Map.get(map, obj), map, acc + 1)
  end

  # Counts all edges in the system
  def count_edges([], _, acc), do: acc
  def count_edges(["COM"], _, acc), do: acc
  def count_edges(["COM"|rest], map, acc), do: count_edges(rest, map, acc)
  def count_edges([obj|rest], map, acc) do
    IO.puts("counting #{obj}, current #{acc}")
    count_edges(rest, map, crawl(obj, map, 0) + acc)
  end
end

input = IO.read(:stdio, :all)
  |> String.trim()
  |> String.split("\n")

# just built in and out edges in case out was needed for part 2
{_out_edges, in_edges, allobjs} = Puzzle.build_edges(input, %{}, %{}, MapSet.new())

IO.puts("all #{inspect(allobjs)}")
IO.puts("edges #{inspect(in_edges)}")

count = Puzzle.count_edges(allobjs, in_edges, 0)
IO.puts("answer #{count}")
