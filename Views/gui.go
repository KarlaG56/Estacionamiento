package Views

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math/rand"
	"time"
)

const (
	maxCars      = 20
	carWidth     = 40
	carHeight    = 20
	garageWidth  = 500
	garageHeight = 300
)

var cars []*car

type car struct {
	sprite *pixel.Sprite
	pos    pixel.Vec
	speed  float64
}

func drawParking(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	// Dibuja los límites del estacionamiento
	p1 := pixel.V(100, 460)
	p2 := pixel.V(100, 20)
	p3 := pixel.V(600, 20)
	p4 := pixel.V(600, 460)
	p5 := pixel.V(100, 300)
	p6 := pixel.V(100, 180)

	// Dibuja los límites exteriores del estacionamiento
	imd.Push(p1)
	imd.Push(p5)
	imd.Line(1)

	imd.Push(p6)
	imd.Push(p2)
	imd.Line(1)

	imd.Push(p2)
	imd.Push(p3)
	imd.Line(1)

	imd.Push(p3)
	imd.Push(p4)
	imd.Line(1)

	imd.Push(p4)
	imd.Push(p1)
	imd.Line(1)

	// Dibuja las plazas de estacionamiento
	width := 40.0
	height := 80.0
	space := 10.0

	// Dibuja las plazas horizontales
	for x := 105.0; x <= 595; x += width + space {
		p1 := pixel.V(x, 455)
		p2 := pixel.V(x, 455-height)
		p3 := pixel.V(x+width, 455-height)
		p4 := pixel.V(x+width, 455)

		imd.Push(p1)
		imd.Push(p2)
		imd.Line(1)

		imd.Push(p3)
		imd.Push(p4)
		imd.Line(1)

		imd.Push(p4)
		imd.Push(p1)
		imd.Line(1)
	}

	// Dibuja las plazas verticales
	for x := 105.0; x <= 595; x += width + space {
		p1 := pixel.V(x, 25)
		p2 := pixel.V(x, 25+height)
		p3 := pixel.V(x+width, 25+height)
		p4 := pixel.V(x+width, 25)

		imd.Push(p1)
		imd.Push(p2)
		imd.Line(1)

		imd.Push(p3)
		imd.Push(p4)
		imd.Line(1)

		imd.Push(p4)
		imd.Push(p1)
		imd.Line(1)
	}

	imd.Draw(win)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Estacionamiento",
		Bounds: pixel.R(0, 0, 640, 480),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	lastCarSpawn := time.Now()

	spriteSheet := createCarSprite(carWidth, carHeight)

	for !win.Closed() {
		if len(cars) < maxCars && time.Since(lastCarSpawn) > 2*time.Second {
			lastCarSpawn = time.Now()
			newCar := &car{
				sprite: spriteSheet,
				pos:    pixel.V(-40, 300), // Posición inicial desde el lado izquierdo
				speed:  rand.Float64()*200 + 100,
			}
			cars = append(cars, newCar)
		}

		win.Clear(colornames.White)
		drawParking(win)
		for _, car := range cars {
			if car.pos.X < 570 && car.pos.Y > 30 && car.pos.Y < 455 {
				car.sprite.Draw(win, pixel.IM.Moved(car.pos))
				car.pos = car.pos.Add(pixel.V(car.speed*1/60, 0))
			} else {
				car.speed = 0
			}
		}

		win.Update()
	}
}

func createCarSprite(width, height int) *pixel.Sprite {
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	transparent := color.RGBA{0, 0, 0, 0} // Color transparente
	red := colornames.Purple

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Define el área del carro como roja y el resto como transparente
			if x > 2 && x < width-3 && y > 2 && y < height-3 {
				rgba.Set(x, y, red)
			} else {
				rgba.Set(x, y, transparent)
			}
		}
	}

	pic := pixel.PictureDataFromImage(rgba)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	return sprite
}

func Show() {
	pixelgl.Run(run)
}
