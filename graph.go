package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Graph struct {
	Nodes       []*Node `json:"nodes"`
	Edges       []*Edge `json:"links"`
	TotalDegree float32
}

type Node struct {
	Name       string `json:"name"`
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

func (graph *Graph) applyForce(deltaTime float32, temperature float32) {
	if config.Gravity {
		center := rl.Vector2{
			X: float32(config.ScreenWidth) / 2,
			Y: float32(config.ScreenHeight) / 2,
		}
		for _, node := range graph.Nodes {
			delta := rl.Vector2Subtract(center, node.pos)
			node.vel = rl.Vector2Scale(delta, config.GravityStrength)
			node.acc = rl.Vector2Zero()
		}
	} else {
		for _, node := range graph.Nodes {
			node.vel = rl.Vector2Zero()
			node.acc = rl.Vector2Zero()
		}
	}

	if config.BarnesHut {
		rect := Rect{-float32(config.ScreenWidth), -float32(config.ScreenHeight), 2 * float32(config.ScreenWidth), 2 * float32(config.ScreenHeight)}
		qt := NewQuadTree(rect)

		for _, node := range graph.Nodes {
			qt.Insert(node)
		}
		qt.CalculateMasses()
		for _, node := range graph.Nodes {
			force := qt.CalculateForce(node, config.Theta)
			node.acc = rl.Vector2Add(node.acc, force)
		}
	} else {
		for i, node := range graph.Nodes {

			for j, other := range graph.Nodes {
				if i == j {
					continue
				}

				delta := rl.Vector2Subtract(node.pos, other.pos)
				dist := rl.Vector2LengthSqr(delta)
				if dist < 1e-2 {
					continue
				}
				var scale float32 = node.degree * other.degree
				force := rl.Vector2Scale(rl.Vector2Normalize(delta), 10*scale/dist)
				node.acc = rl.Vector2Add(node.acc, force)
			}

		}
	}

	for _, edge := range graph.Edges {
		from := graph.Nodes[edge.Source]
		to := graph.Nodes[edge.Target]
		delta := rl.Vector2Subtract(from.pos, to.pos)
		dist := rl.Vector2Length(delta)

		if dist < 1e-1 {
			continue
		}
		s := float32(math.Min(float64(from.radius), float64(to.radius)))
		var l float32 = from.radius + to.radius
		force := rl.Vector2Scale(rl.Vector2Normalize(delta), (dist-l)/s*edge.Value)
		from.acc = rl.Vector2Subtract(from.acc, force)
		to.acc = rl.Vector2Add(to.acc, force)

	}

	for _, node := range graph.Nodes {
		if !node.isSelected {
			node.vel = rl.Vector2Add(node.vel, node.acc)
			node.vel = rl.Vector2ClampValue(node.vel, -temperature, temperature)
			node.pos = rl.Vector2Add(node.pos, rl.Vector2Scale(node.vel, deltaTime))
			node.pos = rl.Vector2Clamp(node.pos, rl.NewVector2(-10*float32(config.ScreenWidth), -10*float32(config.ScreenHeight)), rl.NewVector2(10*float32(config.ScreenWidth), 10*float32(config.ScreenHeight)))
		}
	}
}

func ImportFromJson(filepath string) (*Graph, map[int]rl.Color, error) {
	var graph Graph
	colorMap := make(map[int]rl.Color, 0)
	file, err := os.ReadFile(filepath)
	if err != nil {
		return &graph, colorMap, err
	}

	json.Unmarshal(file, &graph)

	rand := rand.New(rand.NewSource(123))
	for _, node := range graph.Nodes {
		node.pos = rl.Vector2{
			X: float32(rand.Intn(int(config.ScreenWidth))),
			Y: float32(rand.Intn(int(config.ScreenHeight))),
		}
		if _, containsKey := colorMap[node.Group]; !containsKey {
			r := uint8(rand.Intn(255))
			g := uint8(rand.Intn(255))
			b := uint8(rand.Intn(255))
			colorMap[node.Group] = rl.NewColor(r, g, b, 255)
		}
	}
	for _, edge := range graph.Edges {
		if edge.Value == 0.0 {
			edge.Value = 1.0
		}

		graph.Nodes[edge.Source].degree += edge.Value
		graph.Nodes[edge.Target].degree += edge.Value
		graph.TotalDegree += edge.Value
	}

	for _, node := range graph.Nodes {
		node.radius = float32(math.Max(math.Sqrt(float64(node.degree)), 2))
	}

	return &graph, colorMap, nil
}
