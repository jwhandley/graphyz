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
