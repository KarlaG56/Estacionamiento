package Views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

const (
	maxCars      = 20
	garageWidth  = 500
	garageHeight = 300
)

var appCtx = app.New()
var carGraphics []*canvas.Rectangle
var garageMutex sync.Mutex
var parkingSpacesSemaphore = make(chan struct{}, maxCars)
var accessGateSemaphore = make(chan struct{}, 1)

func ShowGUI() {
	window := appCtx.NewWindow("Estacionamiento")

	garage := createParkingGarage()

	containerWithLayout := container.New(layout.NewVBoxLayout(),
		garage,
	)

	window.SetContent(containerWithLayout)
	window.Resize(fyne.NewSize(garageWidth, garageHeight+50))

	window.ShowAndRun()

	go runSimulation()
}

func createParkingGarage() fyne.CanvasObject {
	parkingGarage := container.New(layout.NewHBoxLayout())

	leftSpaces := createParkingSpaces(5)
	parkingGarage.Add(leftSpaces)
	entranceImage := createEntranceImage()
	parkingGarage.Add(entranceImage)
	rightSpaces := createParkingSpaces(5)

	parkingGarage.Add(rightSpaces)

	return parkingGarage
}

func createParkingSpaces(count int) fyne.CanvasObject {
	parkingSpaces := container.New(layout.NewVBoxLayout())
	for i := 0; i < count; i++ {
		parkingSpace := canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Color verde para los lugares de estacionamiento
		parkingSpace.Resize(fyne.NewSize(50, 50))                                   // Tamaño de espacio de estacionamiento
		parkingSpaces.Add(parkingSpace)
	}
	return parkingSpaces
}

func createEntranceImage() fyne.CanvasObject {
	entranceImageData, _ := theme.FyneLogoResource().Content()
	entranceResource := fyne.NewStaticResource("entrance.png", entranceImageData)
	entranceImage := canvas.NewImageFromResource(entranceResource)
	entranceImage.Resize(fyne.NewSize(50, 50)) // Tamaño de la entrada/salida
	return entranceImage
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
	car.Resize(fyne.NewSize(20, 20))
	carGraphics = append(carGraphics, car)

	go moveCar(car)
}

func moveCar(car *canvas.Rectangle) {
	for {
		accessGateSemaphore <- struct{}{}
		<-accessGateSemaphore

		parkingSpacesSemaphore <- struct{}{}
		garageMutex.Lock()

		car.Move(fyne.NewPos(float32(rand.Intn(garageWidth-20)), float32(garageHeight/2-10)))

		garageMutex.Unlock()
		<-parkingSpacesSemaphore

		UpdateGarageView()

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
