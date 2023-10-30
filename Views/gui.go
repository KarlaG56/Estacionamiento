package Views

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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

	garage := canvas.NewRectangle(theme.BackgroundColor())
	garage.Resize(fyne.NewSize(garageWidth, garageHeight))

	containerWithLayout := container.NewMax(garage, label)
	window.SetContent(containerWithLayout)
	window.Resize(fyne.NewSize(garageWidth, garageHeight+50))

	window.Show()

	// Comenzar a agregar autos al estacionamiento
	go runSimulation()
}

func runSimulation() {
	for {
		if len(carGraphics) < maxCars {
			AddCarToGarage()
		}
		UpdateGarageView()
		time.Sleep(2 * time.Second)
	}
}

func AddCarToGarage() {
	car := canvas.NewRectangle(color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255})
	car.Resize(fyne.NewSize(20, 20))
	car.Move(fyne.NewPos(rand.Intn(garageWidth-20), rand.Intn(garageHeight-20)))
	carGraphics = append(carGraphics, car)
}

func UpdateGarageView() {
	contentObjects := []fyne.CanvasObject{canvas.NewRectangle(theme.BackgroundColor())}

	for _, car := range carGraphics {
		contentObjects = append(contentObjects, car)
	}

	window := appCtx.Driver().AllWindows()[0]
	window.SetContent(container.NewMax(contentObjects...))
	window.Canvas().Refresh(window.Content())
}
