# Graphyz

A simple program for real-time force directed graph layout simulation written with Go and Raylib. Loosely based on the Fruchterman-Reingold algorithm, d3-force, and force atlas 2. Uses the Barnes-Hut algorithm to optimize force calculations.

Work in progress.

## Example

![](examples/graphyz-example.png)

## References
- Force-directed graph layouts: https://en.wikipedia.org/wiki/Force-directed_graph_drawing
- d3-force: https://github.com/d3/d3-force
- Force atlas 2: https://journals.plos.org/plosone/article?id=10.1371/journal.pone.0098679
- Go bindings for Raylib: https://github.com/gen2brain/raylib-go 
- Fruchterman-Reingold paper: https://www.mathe2.uni-bayreuth.de/axel/papers/reingold:graph_drawing_by_force_directed_placement.pdf

## To-do
- Use GPU instancing to speed up rendering for large graphs
- Allow users to use their own graph data
- Add support for additional graph serialization formats
