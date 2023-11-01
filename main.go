package main

import (
	"Simulador/Models"
	"Simulador/Scenes"
	"Simulador/Views"
)

func main() {
	// Crear un garaje con capacidad para 20 vehículos
	garage := Models.NewGarage(20)

	// Iniciar la simulación en segundo plano
	go Scenes.RunSimulation(garage)

	// Iniciar la interfaz gráfica
	Views.Show()
}
