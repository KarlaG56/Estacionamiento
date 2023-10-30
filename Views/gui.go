package Views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"math/rand"
	"time"
	"fmt"
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

	// Asegurar el tamaño y ubicación de los autos
	go runSimulation(window)
}

func runSimulation(window fyne.Window) {
	for {
		if len(carGraphics) < maxCars {
			fyne.CurrentApp().Send(func() {
				AddCarToGarage()
			})
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
	car.Resize(fyne.NewSize(20, 20)) // Asegurar el tamaño visible
	car.Move(fyne.NewPos(float32(rand.Intn(garageWidth-20)), float32(rand.Intn(garageHeight-20))) // Asegurar la posición visible
	carGraphics = append(carGraphics, car)

	// Mensajes de registro para depurar
	fmt.Println("Auto agregado al estacionamiento.")
}

func UpdateGarageView(window fyne.Window) {
	contentObjects := make([]fyne.CanvasObject, 0) // Crear una nueva lista para evitar duplicaciones
	// Agregar el rectángulo del estacionamiento a la lista de objetos a mostrar
	contentObjects = append(contentObjects, canvas.NewRectangle(theme.BackgroundColor()))

	for _, car := range carGraphics {
		// Agregar cada auto a la lista de objetos a mostrar
		contentObjects = append(contentObjects, car)
	}

	// Actualizar la ventana con los objetos
	content := container.NewMax(contentObjects...)
	window.SetContent(content)
	window.Canvas().Refresh(content)

	// Mensajes de registro para depurar
	fmt.Println("Actualización de la vista.")
}
