package Scenes

import (
	"Simulador/Models"
	"Simulador/Views"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"time"
)

func Run() {

	Models.Init()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Parking Lot Simulation",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for car := range Models.CarChannel {
			Models.LaneMutex.Lock()
			for _, occupied := range Models.LaneStatus {
				if !occupied {
					break
				}
			}
			Models.LaneMutex.Unlock()

			go Models.Lane(car.ID)
		}
	}()

	for !win.Closed() {
		win.Clear(colornames.White)
		Views.DrawParkingLot(win, Models.GetCars())
		win.Update()
		Models.CarsMutex.Lock()
		Models.MoveCarsLogic()
		Models.CarsMutex.Unlock()

		time.Sleep(16 * time.Millisecond)
	}
}
