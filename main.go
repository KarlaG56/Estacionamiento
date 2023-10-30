package main

import (
	"Simulador/Models"
	"Simulador/Scenes"
	"Simulador/Views"
)

const maxCars = 20

func main() {
	garage := Models.NewGarage(maxCars)

	go Scenes.RunSimulation(garage)
	Views.ShowGUI()

	// Simulated indefinite loop to keep the program running
	select {}
}
