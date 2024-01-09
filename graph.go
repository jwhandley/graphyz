package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Body interface {
	size() float32
	position() rl.Vector2
}

type Graph struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"links"`
}

type Node struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	Group      int    `json:"group"`
	degree     float32
	isSelected bool
	radius     float32
	pos        rl.Vector2
	vel        rl.Vector2
	acc        rl.Vector2
}

type Edge struct {
	Source int     `json:"source"`
	Target int     `json:"target"`
	Value  float32 `json:"value"`
}

func (graph *Graph) resetPosition() {
	var initialRadius float32 = 10.0
	initialAngle := float64(rl.Pi) * (3 - math.Sqrt(5))
	for i, node := range graph.Nodes {
		radius := initialRadius * float32(math.Sqrt(0.5+float64(i)))
		angle := float64(i) * initialAngle

		node.pos = rl.Vector2{
			X: radius*float32(math.Cos(angle)) + float32(config.ScreenWidth)/2,
			Y: radius*float32(math.Sin(angle)) + float32(config.ScreenHeight)/2,
		}
	}
}

func (graph *Graph) ApplyForce(deltaTime float32, qt *QuadTree) {
	graph.resetAcceleration()
	if config.Gravity {
		graph.gravityForce()
	}

	graph.attractionForce()

	if config.BarnesHut {
		graph.repulsionBarnesHut(qt)
	} else {
		graph.repulsionNaive()
	}

	graph.updatePositions(deltaTime)
}

func (graph *Graph) updatePositions(deltaTime float32) {
	for _, node := range graph.Nodes {
		if !node.isSelected {
			node.vel = rl.Vector2Add(node.vel, node.acc)
			node.vel = rl.Vector2Scale(node.vel, 1-config.VelocityDecay)
			node.vel = rl.Vector2ClampValue(node.vel, -100, 100)
			node.pos = rl.Vector2Add(node.pos, rl.Vector2Scale(node.vel, deltaTime))
			node.pos = rl.Vector2Clamp(node.pos, rl.NewVector2(-10*float32(config.ScreenWidth), -10*float32(config.ScreenHeight)), rl.NewVector2(10*float32(config.ScreenWidth), 10*float32(config.ScreenHeight)))
		}
	}
}

func (graph *Graph) resetAcceleration() {
	for _, node := range graph.Nodes {
		node.acc = rl.Vector2Zero()
	}
}

func (graph *Graph) gravityForce() {
	center := rl.Vector2{
		X: float32(config.ScreenWidth) / 2,
		Y: float32(config.ScreenHeight) / 2,
	}
	for _, node := range graph.Nodes {
		delta := rl.Vector2Subtract(center, node.pos)
		force := rl.Vector2Scale(delta, config.GravityStrength*node.size()*temperature)
		node.acc = rl.Vector2Add(node.acc, force)
	}
}

func (graph *Graph) attractionForce() {
	for _, edge := range graph.Edges {
		from := graph.Nodes[edge.Source]
		to := graph.Nodes[edge.Target]
		force := calculateAttractionForce(from, to, edge.Value)
		from.acc = rl.Vector2Subtract(from.acc, force)
		to.acc = rl.Vector2Add(to.acc, force)

	}
}

func (graph *Graph) repulsionBarnesHut(qt *QuadTree) {
	qt.Clear()

	for _, node := range graph.Nodes {
		qt.Insert(node)
	}
	qt.CalculateMasses()
	for _, node := range graph.Nodes {
		force := qt.CalculateForce(node, config.Theta)
		node.acc = rl.Vector2Add(node.acc, force)
	}
}

func (graph *Graph) repulsionNaive() {
	for i, node := range graph.Nodes {

		for j, other := range graph.Nodes {
			if i == j {
				continue
			}

			force := calculateRepulsionForce(node, other)
			node.acc = rl.Vector2Add(node.acc, force)
		}

	}
}

func (node *Node) size() float32 {
	return node.degree
}

func (node *Node) position() rl.Vector2 {
	return node.pos
}

func calculateRepulsionForce(b1 Body, b2 Body) rl.Vector2 {
	delta := rl.Vector2Subtract(b1.position(), b2.position())
	dist := rl.Vector2LengthSqr(delta)
	if dist*dist < b1.size()*b2.size() {
		dist = b1.size() * b2.size()
	}
	scale := b1.size() * b2.size() * temperature
	force := rl.Vector2Scale(rl.Vector2Normalize(delta), 10*scale/dist)
	return force
}

func calculateAttractionForce(from *Node, to *Node, weight float32) rl.Vector2 {
	delta := rl.Vector2Subtract(from.pos, to.pos)
	dist := rl.Vector2Length(delta)

	if dist < EPSILON {
		dist = EPSILON
	}
	s := float32(math.Min(float64(from.radius), float64(to.radius)))
	var l float32 = from.radius + to.radius
	return rl.Vector2Scale(rl.Vector2Normalize(delta), (dist-l)/s*weight*temperature)
}
