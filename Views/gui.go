package Views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

	garage := canvas.NewRectangle(theme.BackgroundColor())
	garage.Resize(fyne.NewSize(garageWidth, garageHeight))

	containerWithLayout := container.NewMax(garage, label)
	window.SetContent(containerWithLayout)
	window.Resize(fyne.NewSize(garageWidth, garageHeight+50))

	window.Show()

	go runSimulation(window)
}

func runSimulation(window fyne.Window) {
	for {
		if len(carGraphics) < maxCars {
			AddCarToGarage()
		}
		UpdateGarageView(window)
		time.Sleep(2 * time.Second)
	}
}

func AddCarToGarage() {
	car := canvas.NewRectangle(color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255,
	})
	car.Resize(fyne.NewSize(20, 20))
	car.Move(fyne.NewPos(float32(rand.Intn(garageWidth-20)), float32(rand.Intn(garageHeight-20)))
	carGraphics = append(carGraphics, car)

	fmt.Println("Auto agregado al estacionamiento.")
}

func UpdateGarageView(window fyne.Window) {
	contentObjects := make([]fyne.CanvasObject, 0)
	contentObjects = append(contentObjects, canvas.NewRectangle(theme.BackgroundColor()))

	for _, car := range carGraphics {
		contentObjects = append(contentObjects, car)
	}

	content := container.NewMax(contentObjects...)
	window.SetContent(content)
	window.Canvas().Refresh(content)

	fmt.Println("Actualización de la vista.")
}
