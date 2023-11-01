package main

import (
	"Simulador/Scenes"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	Scenes.Run()
}

func main() {
	pixelgl.Run(run)
}
