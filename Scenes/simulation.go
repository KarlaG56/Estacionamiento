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

	Models.InicializarSistemaDeAutos()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Parking Lot Simulation",
		Bounds: pixel.R(0, 0, 800, 600),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for car := range Models.CanalAutos {
			Models.MutexCarriles.Lock()
			for _, occupied := range Models.EstadoCarriles {
				if !occupied {
					break
				}
			}
			Models.MutexCarriles.Unlock()

			go Models.SeleccionarCarril(car.ID)
		}
	}()

	for !win.Closed() {
		win.Clear(colornames.White)
		Views.DrawParkingLot(win, Models.ObtenerListaDeAutos())
		win.Update()
		Models.MutexAutos.Lock()
		Models.LogicaDeMovimientoDeAutos()
		Models.MutexAutos.Unlock()

		time.Sleep(16 * time.Millisecond)
	}
}
