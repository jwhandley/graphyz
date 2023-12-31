package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ScreenWidth     int32   `yaml:"ScreenWidth"`
	ScreenHeight    int32   `yaml:"ScreenHeight"`
	AlphaTarget     float32 `yaml:"AlphaTarget"`
	AlphaDecay      float32 `yaml:"AlphaDecay"`
	AlphaInit       float32 `yaml:"AlphaInit"`
	BarnesHut       bool    `yaml:"BarnesHut"`
	Capacity        int     `yaml:"Capacity"`
	Gravity         bool    `yaml:"Gravity"`
	Theta           float32 `yaml:"Theta"`
	GravityStrength float32 `yaml:"Grav"`
}

var config Config

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

var mutex = &sync.Mutex{}

func updatePhysics(graph *Graph) {
	const deltaTime = time.Millisecond * 16
	temperature := config.AlphaInit
	for {
		temperature += (config.AlphaTarget - temperature) * config.AlphaDecay * 0.016
		mutex.Lock()
		graph.applyForce(0.016, temperature)
		mutex.Unlock()
		time.Sleep(deltaTime)
	}
}

func main() {
	path := os.Args[1]
	graph, colorMap, err := ImportFromJson(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	go updatePhysics(graph)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, "graphyz")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := new(rl.Camera2D)
	camera.Target = rl.Vector2{X: float32(config.ScreenWidth) / 2, Y: float32(config.ScreenHeight) / 2}
	camera.Offset = rl.Vector2{X: float32(config.ScreenWidth) / 2, Y: float32(config.ScreenHeight) / 2}
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
					X: float32(rand.Intn(int(config.ScreenWidth))),
					Y: float32(rand.Intn(int(config.ScreenHeight))),
				}
			}
		}

		mousePos := rl.GetMousePosition()
		mousePos.X = (mousePos.X-camera.Offset.X)/camera.Zoom + camera.Target.X
		mousePos.Y = (mousePos.Y-camera.Offset.Y)/camera.Zoom + camera.Target.Y
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.BeginMode2D(*camera)

		mutex.Lock()
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
				rl.DrawText(message, int32(mousePos.X)+5, int32(mousePos.Y), 20, rl.Black)
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
		mutex.Unlock()
		rl.EndMode2D()
		rl.DrawFPS(10, 10)
		zoomMessage := fmt.Sprintf("Zoom: %.2f", camera.Zoom)
		rl.DrawText(zoomMessage, config.ScreenWidth-110, 10, 20, rl.Black)
		rl.EndDrawing()
	}

}
