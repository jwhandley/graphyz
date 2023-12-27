# Graphyz

A simple program for real-time force directed graph layout simulation. Loosely based on the Fruchterman-Reingold algorithm, d3-force, and force atlas 2.
Work in progress.

## Example

![](examples/graphyz-example.png)

## References
- Force-directed graph layouts: https://en.wikipedia.org/wiki/Force-directed_graph_drawing
- d3-force: https://github.com/d3/d3-force
- Force atlas 2: https://journals.plos.org/plosone/article?id=10.1371/journal.pone.0098679

## To-do
- Implement an optimization such as the Barnes-Hut algorithm for the many body forces (needed for scaling beyond a few hundred nodes)
- Allow users to use their own graph data
- Add support for additional formats
