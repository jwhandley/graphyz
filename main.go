package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
	nodeRadius   = 5.0
)

type Graph struct {
	Nodes       []*Node `json:"nodes"`
	Edges       []*Edge `json:"links"`
	TotalDegree int
}

type Node struct {
	Name       string `json:"name"`
	Group      int    `json:"group"`
	degree     int
	isSelected bool
	pos        rl.Vector2
	vel        rl.Vector2
}

type Edge struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}

func (graph *Graph) applyForce(deltaTime float32, k float32) {
	center := rl.Vector2{
		X: screenWidth / 2,
		Y: screenHeight / 2,
	}
	for _, node := range graph.Nodes {
		delta := rl.Vector2Subtract(center, node.pos)
		node.vel = rl.Vector2Scale(delta, 0.1)
		node.vel = rl.Vector2Zero()
	}

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

	for _, edge := range graph.Edges {
		from := graph.Nodes[edge.Source]
		to := graph.Nodes[edge.Target]
		delta := rl.Vector2Subtract(from.pos, to.pos)
		dist := rl.Vector2Length(delta)

		if dist < 1e-2 {
			continue
		}
		l := float32(5.0)
		s := float32(math.Min(float64(from.degree), float64(to.degree)))
		dv := rl.Vector2Scale(rl.Vector2Normalize(delta), (dist-l)/s*float32(edge.Value))
		from.vel = rl.Vector2Subtract(from.vel, dv)
		to.vel = rl.Vector2Add(to.vel, dv)
	}

	for _, node := range graph.Nodes {
		node.pos = rl.Vector2Add(node.pos, rl.Vector2Scale(node.vel, deltaTime))
		node.pos = rl.Vector2Clamp(node.pos, rl.Vector2{X: -screenWidth, Y: -screenHeight}, rl.Vector2{X: screenWidth * 2, Y: screenHeight * 2})
	}
}

func main() {
	file, err := os.ReadFile("assets/les-mis.json")
	if err != nil {
		panic(err)
	}

	var graph Graph
	json.Unmarshal(file, &graph)
	colorMap := make(map[int]rl.Color, 0)
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
		graph.Nodes[edge.Source].degree += edge.Value
		graph.Nodes[edge.Target].degree += edge.Value
		graph.TotalDegree += edge.Value
	}
	k := float32(math.Sqrt(float64(screenWidth * screenHeight / graph.TotalDegree)))

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "graphyz")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	camera := new(rl.Camera2D)
	camera.Target = rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2}
	camera.Offset = rl.Vector2{X: screenWidth / 2, Y: screenHeight / 2}
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	anySelected := false
	for !rl.WindowShouldClose() {
		camera.Zoom += rl.GetMouseWheelMove() * 0.05
		if camera.Zoom > 3.0 {
			camera.Zoom = 3.0
		} else if camera.Zoom < 0.1 {
			camera.Zoom = 0.1
		}

		if rl.IsKeyPressed(rl.KeyR) {
			camera.Zoom = 1.0
			for _, node := range graph.Nodes {
				node.pos = rl.Vector2{
					X: float32(rand.Intn(screenWidth)),
					Y: float32(rand.Intn(screenHeight)),
				}
			}
		}

		graph.applyForce(rl.GetFrameTime(), k)
		mousePos := rl.GetMousePosition()
		mousePos.X = (mousePos.X-camera.Offset.X)/camera.Zoom + camera.Target.X
		mousePos.Y = (mousePos.Y-camera.Offset.Y)/camera.Zoom + camera.Target.Y
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(*camera)

		for _, edge := range graph.Edges {
			sourcePos := graph.Nodes[edge.Source].pos
			targetPos := graph.Nodes[edge.Target].pos
			thickness := float32(math.Sqrt(float64(edge.Value)))
			rl.DrawLineEx(sourcePos, targetPos, thickness, rl.LightGray)
		}
		for _, node := range graph.Nodes {
			dist := rl.Vector2Distance(mousePos, node.pos)
			radius := float32(math.Max(math.Sqrt(float64(node.degree)), 2))
			rl.DrawCircleV(node.pos, radius, colorMap[node.Group])
			if dist < radius {
				if rl.IsMouseButtonDown(0) && !anySelected {
					node.isSelected = true
					anySelected = true
				}
				//message := node.Name + ", Group " + strconv.Itoa(node.Group) + "\n Degree: " + strconv.Itoa(node.degree)
				message := fmt.Sprintf("%s, Group: %d\nDegree: %d", node.Name, node.Group, node.degree)
				rl.DrawText(message, int32(node.pos.X)+5, int32(node.pos.Y), 20, rl.Black)
				rl.DrawCircleV(node.pos, radius, rl.NewColor(80, 80, 80, 150))
			}

			if node.isSelected {
				if rl.IsMouseButtonDown(0) {
					node.pos = mousePos
				} else {
					node.isSelected = false
					anySelected = false
				}

			}
		}
		rl.EndMode2D()
		rl.DrawFPS(10, 10)
		zoomMessage := fmt.Sprintf("Zoom: %.2f", camera.Zoom)
		rl.DrawText(zoomMessage, screenWidth-100, 10, 20, rl.Black)
		rl.EndDrawing()
	}

}
