package main

import (
	"encoding/json"
	"math"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Matplotlib tab20c colormap
var Colors = [20]rl.Color{
	{49, 130, 189, 255},
	{107, 174, 214, 255},
	{158, 202, 225, 255},
	{198, 219, 239, 255},
	{230, 85, 13, 255},
	{253, 141, 60, 255},
	{253, 174, 107, 255},
	{253, 208, 162, 255},
	{49, 163, 84, 255},
	{116, 196, 118, 255},
	{161, 217, 155, 255},
	{199, 233, 192, 255},
	{117, 107, 177, 255},
	{158, 154, 200, 255},
	{188, 189, 220, 255},
	{218, 218, 235, 255},
	{99, 99, 99, 255},
	{150, 150, 150, 255},
	{189, 189, 189, 255},
	{217, 217, 217, 255},
}

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
	for i, node := range graph.Nodes {
		radius := initialRadius * float32(math.Sqrt(0.5+float64(i)))
		angle := float64(i) * initialAngle

		node.pos = rl.Vector2{
			X: radius*float32(math.Cos(angle)) + float32(config.ScreenWidth)/2,
			Y: radius*float32(math.Sin(angle)) + float32(config.ScreenHeight)/2,
		}
		if _, containsKey := colorMap[node.Group]; !containsKey {
			colorMap[node.Group] = Colors[node.Group%len(Colors)]
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
