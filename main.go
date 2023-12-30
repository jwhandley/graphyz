package main

import (
	"fmt"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
	alphaTarget  = 1.0
	alphaDecay   = 0.01
	alphaInit    = float32(20.0)
	barnesHut    = true
	gravity      = true
)

func main() {
	graph, colorMap, err := ImportFromJson("assets/gove-rvs-materials/gove-rvs-materials/json/benchmark/block_2000.json")
	if err != nil {
		panic(err)
	}

	temperature := alphaInit

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
			temperature = alphaInit
			camera.Zoom = 1.0
			for _, node := range graph.Nodes {
				node.pos = rl.Vector2{
					X: float32(rand.Intn(screenWidth)),
					Y: float32(rand.Intn(screenHeight)),
				}
			}
		}

		graph.applyForce(rl.GetFrameTime(), temperature)
		temperature += (alphaTarget - temperature) * alphaDecay * rl.GetFrameTime()

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
				message := fmt.Sprintf("%s, Group: %d\nDegree: %.0f", node.Name, node.Group, node.degree)
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
		rl.DrawText(zoomMessage, screenWidth-110, 10, 20, rl.Black)
		rl.EndDrawing()
	}

}
