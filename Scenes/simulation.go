package Scenes

import (
	"Simulador/Models"
	"fmt"
	"math/rand"
	"time"

	_ "Simulador/Models"
)

const (
	poisonInterval = 1
	maxGarage      = 100
)

func RunSimulation(garage *Models.Garage) {
	for {
		carArrives(garage)
		time.Sleep(poisonInterval * time.Second)
	}
}

func carArrives(garage *Models.Garage) {
	carID := rand.Intn(maxGarage)
	fmt.Printf("Car %d arrives.\n", carID)

	if garage.CarsInGarage < garage.MaxCars {
		garage.AddCar()
		go func(id int) {
			if parkCar(id) {
				garage.RemoveCar()
			}
		}(carID)
	} else {
		fmt.Printf("Car %d is blocked - Parking is full.\n", carID)
	}
}

func parkCar(carID int) bool {
	// Simulating checking for available parking spots
	for i := 1; i <= maxGarage; i++ {
		fmt.Printf("Car %d is looking for a parking spot...\n", carID)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		// Check if a parking spot is available
		if i == carID {
			fmt.Printf("Car %d parked successfully.\n", carID)
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
			fmt.Printf("Car %d leaves.\n", carID)
			return true
		}
	}
	fmt.Printf("Car %d couldn't find a parking spot.\n", carID)
	return false
}
