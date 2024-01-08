package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
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
	Debug           bool    `yaml:"Debug"`
}

var config Config
var temperature float32
var mutex sync.Mutex

const EPSILON = 1e-2

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

func updatePhysics(graph *Graph, numSteps int) {
	targetTime := time.Millisecond * 16
	var frameTime float32 = 0.016 / float32(numSteps)

	for {
		startTime := time.Now()
		totalTime := time.Duration(0)

		for totalTime <= targetTime {
			graph.ApplyForce(frameTime)
			elapsedTime := time.Since(startTime)
			totalTime += elapsedTime

			// Update frameTime and temperature for the next iteration
			frameTime = float32(elapsedTime.Seconds())
			temperature += (config.AlphaTarget - temperature) * config.AlphaDecay * frameTime
		}

		// Sleep for the remaining time of the target time, if any
		remainingTime := targetTime - totalTime
		if remainingTime > 0 {
			time.Sleep(remainingTime)
		}
	}
}

func main() {
	if config.Debug {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	path := os.Args[1]
	graph, colorMap, err := ImportFromJson(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	temperature = config.AlphaInit
	go updatePhysics(graph, 8)

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	windowTitle := fmt.Sprintf("graphyz - %s", strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
	rl.InitWindow(config.ScreenWidth, config.ScreenHeight, windowTitle)
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := new(rl.Camera2D)
	camera.Target = rl.Vector2{X: float32(config.ScreenWidth) / 2, Y: float32(config.ScreenHeight) / 2}
	camera.Offset = rl.Vector2{X: float32(config.ScreenWidth) / 2, Y: float32(config.ScreenHeight) / 2}
	camera.Rotation = 0.0
	camera.Zoom = 1.0

	anySelected := false
	panMode := false
	for !rl.WindowShouldClose() {
		camera.Zoom += rl.GetMouseWheelMove() * 0.05
		if camera.Zoom > 3.0 {
			camera.Zoom = 3.0
		} else if camera.Zoom < 0.1 {
			camera.Zoom = 0.1
		}

		if rl.IsKeyPressed(rl.KeyR) {
			temperature = config.AlphaInit
			camera.Zoom = 1.0
			graph.resetPosition()
		}

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !anySelected {
			panMode = true
		}

		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			panMode = false
		}

		if panMode {
			camera.Offset = rl.Vector2Add(camera.Offset, rl.GetMouseDelta())
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
			rl.DrawCircleV(node.pos, node.radius, colorMap[node.Group])
		}

		for _, node := range graph.Nodes {
			dist := rl.Vector2Distance(mousePos, node.pos)
			if dist < node.radius {
				if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !anySelected {
					node.isSelected = true
					anySelected = true
					panMode = false
				}
				var name string
				if len(node.Name) > 0 {
					name = node.Name
				} else if len(node.Label) > 0 {
					name = node.Label
				}
				message := fmt.Sprintf("%s, Group: %d\nDegree: %.0f", name, node.Group, node.degree)
				rl.DrawText(message, int32(mousePos.X)+5, int32(mousePos.Y), 20, rl.Black)
				rl.DrawCircleV(node.pos, node.radius, rl.NewColor(80, 80, 80, 150))
			}

			if node.isSelected {
				if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
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
		rl.DrawText(zoomMessage, config.ScreenWidth-105, 10, 20, rl.Black)

		numNodes := fmt.Sprintf("Number of nodes: %d", len(graph.Nodes))
		rl.DrawText(numNodes, 10, config.ScreenHeight-45, 20, rl.Black)
		numEdges := fmt.Sprintf("Number of edges: %d", len(graph.Edges))
		rl.DrawText(numEdges, 10, config.ScreenHeight-25, 20, rl.Black)
		rl.EndDrawing()
	}

}
