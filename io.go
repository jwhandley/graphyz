package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ImportFromJson(filepath string) (*Graph, map[int]rl.Color, error) {
	var graph Graph
	colorMap := make(map[int]rl.Color, 0)
	file, err := os.ReadFile(filepath)
	if err != nil {
		return &graph, colorMap, err
	}

	json.Unmarshal(file, &graph)
	var initialRadius float32 = 10.0
	initialAngle := float64(rl.Pi) * (3 - math.Sqrt(5))
	rand := rand.New(rand.NewSource(123))
	for i, node := range graph.Nodes {
		radius := initialRadius * float32(math.Sqrt(0.5+float64(i)))
		angle := float64(i) * initialAngle

		node.pos = rl.Vector2{
			X: radius*float32(math.Cos(angle)) + float32(config.ScreenWidth)/2,
			Y: radius*float32(math.Sin(angle)) + float32(config.ScreenHeight)/2,
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
	}

	for _, node := range graph.Nodes {
		node.radius = float32(math.Max(math.Sqrt(float64(node.degree)), 2))
	}

	return &graph, colorMap, nil
}
