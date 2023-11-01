package Views

import (
	"Simulador/Models"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
)

var (
	background *pixel.Sprite
	bgPicture  pixel.Picture
	carSprites map[int]*pixel.Sprite
)

func LoadCarImages() {
	// initialize the map
	carSprites = make(map[int]*pixel.Sprite)
	// iterate over car types or IDs
	for _, car := range Models.GetCars() {
		imgPath := "Assets/auto.png" // define how you get each car's image path
		file, err := os.Open(imgPath)
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(file)
		if err != nil {
			panic(err)
		}
		file.Close()
		pic := pixel.PictureDataFromImage(img)
		carSprites[car.ID] = pixel.NewSprite(pic, pic.Bounds())
	}
}

func loadBackground() {
	file, err := os.Open("Assets/background.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bgPicture = pixel.PictureDataFromImage(img)
	background = pixel.NewSprite(bgPicture, bgPicture.Bounds())
}

func DrawParkingLot(win *pixelgl.Window, cars []Models.Car) {
	if background == nil {
		loadBackground()
	}

	background.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	imd.Push(pixel.V(100, 500), pixel.V(700, 500))
	imd.Line(2)
	imd.Push(pixel.V(100, 100), pixel.V(700, 100))
	imd.Line(2)

	imd.Push(pixel.V(700, 100), pixel.V(700, 500))
	imd.Line(2)

	parkingWidth := 600.0
	laneWidth := parkingWidth / 10

	for i := 0.0; i < 10.0; i++ {
		xOffset := 100.0 + i*laneWidth
		imd.Push(pixel.V(xOffset, 500), pixel.V(xOffset, 350))
		imd.Line(2)

		imd.Push(pixel.V(xOffset, 250), pixel.V(xOffset, 100))
		imd.Line(2)
	}

	carWidth := laneWidth / 4
	carHeight := laneWidth / 4

	LoadCarImages()
	for _, car := range cars {
		sprite := carSprites[car.ID]
		if sprite != nil {
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(car.Position))
		} else {
			imd.Color = colornames.Red
			imd.Push(pixel.V(car.Position.X-carWidth/2, car.Position.Y-carHeight/2))
			imd.Push(pixel.V(car.Position.X+carWidth/2, car.Position.Y-carHeight/2))
			imd.Push(pixel.V(car.Position.X+carWidth/2, car.Position.Y+carHeight/2))
			imd.Push(pixel.V(car.Position.X-carWidth/2, car.Position.Y+carHeight/2))
			imd.Polygon(0)
		}
	}

	imd.Draw(win)
}
