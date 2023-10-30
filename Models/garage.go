package Models

import "sync"

type Garage struct {
	CarsInGarage int
	MaxCars      int
	Capacity     chan struct{}
	GarageMutex  *sync.Mutex
}

func NewGarage(maxCars int) *Garage {
	return &Garage{
		CarsInGarage: 0,
		MaxCars:      maxCars,
		Capacity:     make(chan struct{}, maxCars),
		GarageMutex:  &sync.Mutex{},
	}
}

func (g *Garage) AddCar() {
	g.GarageMutex.Lock()
	defer g.GarageMutex.Unlock()
	g.CarsInGarage++
	g.Capacity <- struct{}{}
}

func (g *Garage) RemoveCar() {
	<-g.Capacity
	g.GarageMutex.Lock()
	defer g.GarageMutex.Unlock()
	g.CarsInGarage--
}
