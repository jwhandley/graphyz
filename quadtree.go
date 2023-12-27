package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Quadtree struct {
	Center    rl.Vector2
	TotalMass float32
	Region    Rect
	Nodes     []Node
	Children  [4]*Quadtree
}

type Rect struct {
	X, Y, Width, Height float32
}

func NewQuadTree(boundary Rect) *Quadtree {
	qt := new(Quadtree)
	qt.Region = boundary
	qt.Nodes = make([]Node, 0)

	return qt
}
