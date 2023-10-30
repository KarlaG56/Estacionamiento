package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math/rand"
	"time"
)

const (
	maxCars      = 20
	garageWidth  = 500
	garageHeight = 300
)

var appCtx = app.New()
var carGraphics []*canvas.Rectangle

func ShowGUI() {
	window := appCtx.NewWindow("Estacionamiento")
	label := widget.NewLabel("Estacionamiento")

	garage := createParkingGarage()

	entrance := createEntrance()

	containerWithLayout := container.NewBorder(nil, entrance, nil, nil, container.NewMax(garage, label))

	window.SetContent(containerWithLayout)
	window.Resize(fyne.NewSize(garageWidth, garageHeight+50))

	window.ShowAndRun()

	go runSimulation()
}

func createParkingGarage() fyne.CanvasObject {
	parkingGarage := container.NewGridWithRows(5)
	for i := 0; i < 5; i++ {
		row := container.NewGridWithColumns(4)
		for j := 0; j < 4; j++ {
			parkingSpace := canvas.NewRectangle(color.RGBA{R: 192, G: 192, B: 192, A: 255}) // Color gris
			parkingSpace.Resize(fyne.NewSize(50, 50))                                       // Tamaño de espacio de estacionamiento
			row.AddObject(parkingSpace)
		}
		parkingGarage.AddObject(row)
	}
	return parkingGarage
}

func createEntrance() fyne.CanvasObject {
	entranceButton := widget.NewButton("Entrada", func() {
		// Acción al hacer clic en la entrada
	})
	return entranceButton
}

func runSimulation() {
	for {
		select {
		case <-time.After(time.Second):
			if len(carGraphics) < maxCars {
				go AddCarToGarage()
			}
		}
	}
}

func AddCarToGarage() {
	car := canvas.NewRectangle(color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255,
	})
	car.Resize(fyne.NewSize(5, 5))                                               // Cambia el tamaño del rectángulo a 5x5 (ajusta según tus preferencias)
	car.Move(fyne.NewPos(rand.Intn(garageWidth-20), rand.Intn(garageHeight-20))) // Ajusta el tamaño máximo
	carGraphics = append(carGraphics, car)

	go moveCar(car)
}

func moveCar(car *canvas.Rectangle) {
	for {
		// Modificar la posición del carro en el estacionamiento
		car.Move(fyne.NewPos(rand.Intn(garageWidth-20), rand.Intn(garageHeight-20)))

		// Refrescar la vista para actualizar la posición del carro
		UpdateGarageView()

		// Simular tiempo de ocupación del cajón
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	}
}

func UpdateGarageView() {
	contentObjects := []fyne.CanvasObject{canvas.NewRectangle(color.RGBA{R: 192, G: 192, B: 192, A: 255})}

	for _, car := range carGraphics {
		contentObjects = append(contentObjects, car)
	}

	window := appCtx.Driver().AllWindows()[0]
	window.SetContent(container.NewMax(contentObjects...))
	window.Canvas().Refresh(window.Content())
}
