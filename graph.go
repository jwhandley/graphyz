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
	pos        rl.Vector2
	vel        rl.Vector2
}

type Edge struct {
	Source int     `json:"source"`
	Target int     `json:"target"`
	Value  float32 `json:"value"`
}

func (graph *Graph) applyForce(deltaTime float32, temperature float32) {
	if gravity {
		center := rl.Vector2{
			X: screenWidth / 2,
			Y: screenHeight / 2,
		}
		for _, node := range graph.Nodes {
			delta := rl.Vector2Subtract(center, node.pos)
			node.vel = rl.Vector2Scale(delta, gravityStrength)
		}
	} else {
		for _, node := range graph.Nodes {
			node.vel = rl.Vector2Zero()
		}
	}

	if barnesHut {
		rect := Rect{-screenWidth, -screenHeight, 2 * screenWidth, 2 * screenHeight}
		qt := NewQuadTree(rect)
		for _, node := range graph.Nodes {
			qt.Insert(node)
		}
		qt.CalculateMasses()
		for _, node := range graph.Nodes {
			force := qt.CalculateForce(node, theta)
			node.vel = rl.Vector2Add(node.vel, force)
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
				scale := float32(node.degree * other.degree)
				dv := rl.Vector2Scale(rl.Vector2Normalize(delta), 10*scale/dist)
				node.vel = rl.Vector2Add(node.vel, dv)
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
		s := float32(math.Min(float64(from.degree), float64(to.degree)))
		l := float32(5)
		dv := rl.Vector2Scale(rl.Vector2Normalize(delta), (dist-l)/s*float32(edge.Value))
		from.vel = rl.Vector2Subtract(from.vel, dv)
		to.vel = rl.Vector2Add(to.vel, dv)
	}

	for _, node := range graph.Nodes {
		node.vel = rl.Vector2Clamp(node.vel, rl.NewVector2(-temperature, -temperature), rl.NewVector2(temperature, temperature))
		node.pos = rl.Vector2Add(node.pos, rl.Vector2Scale(node.vel, deltaTime))
		node.pos = rl.Vector2Clamp(node.pos, rl.NewVector2(-10*screenWidth, -10*screenHeight), rl.NewVector2(10*screenWidth, 10*screenHeight))
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
			X: float32(rand.Intn(screenWidth)),
			Y: float32(rand.Intn(screenHeight)),
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

	return &graph, colorMap, nil
}
