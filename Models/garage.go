package Models

import (
	"github.com/faiface/pixel"
	"math/rand"
	"sync"
	"time"
)

const (
	numLanes  = 20
	LaneWidth = 150.0
)

var (
	LaneStatus           [numLanes]bool
	Cars                 []Car
	CarsMutex            sync.Mutex
	LaneMutex            sync.Mutex
	CarEnteringOrExiting bool
)

type Car struct {
	ID                int
	Position          pixel.Vec
	PreviousPosition  pixel.Vec
	Lane              int
	Parked            bool
	ExitTime          time.Time
	IsEntering        bool
	Teleporting       bool
	TeleportStartTime time.Time
}

var CarChannel chan Car

func Init() {
	CarChannel = make(chan Car)
	go CarGenerator()
}

func CreateCar(id int) Car {
	CarsMutex.Lock()
	defer CarsMutex.Unlock()
	Car := Car{
		ID:       id,
		Position: pixel.V(0, 300),
		Lane:     -1,
		Parked:   false,
	}
	Cars = append(Cars, Car)
	return Car
}

func SetExitTime(car *Car) {
	rand.Seed(time.Now().UnixNano())
	exitIn := time.Duration(rand.Intn(5)+1) * time.Second
	car.ExitTime = time.Now().Add(exitIn)
}

func GetCars() []Car {
	return Cars
}

func AssignLaneToCar(id int, lane int) {
	CarsMutex.Lock()
	defer CarsMutex.Unlock()
	for i := range Cars {
		if Cars[i].ID == id {
			Cars[i].Lane = lane
		}
	}
}

func ResetCarPosition(id int) {
	CarsMutex.Lock()
	defer CarsMutex.Unlock()
	for i := range Cars {
		if Cars[i].ID == id {
			Cars[i].Position = pixel.V(0, 300)
		}
	}
}

func FindCarPosition(id int) pixel.Vec {
	CarsMutex.Lock()
	defer CarsMutex.Unlock()
	for _, car := range Cars {
		if car.ID == id {
			return car.Position
		}
	}
	return pixel.Vec{}
}

func ParkCar(car *Car, targetX, targetY float64) {
	car.Position.X = targetX
	car.Position.Y = targetY
	car.Parked = true
	SetExitTime(car)
}

func removeCar(index int) {
	Cars = append(Cars[:index], Cars[index+1:]...)
}

func CarGenerator() {
	id := 0
	for {
		id++
		car := CreateCar(id)
		CarChannel <- car
		time.Sleep(time.Millisecond * 500)
	}
}

func WaitForPosition(id int, targetX float64) {
	for {
		carPos := FindCarPosition(id)
		if carPos.X >= targetX {
			break
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func FindAvailableLane() (int, bool) {
	LaneMutex.Lock()
	defer LaneMutex.Unlock()
	rand.Seed(time.Now().UnixNano())
	lanes := rand.Perm(numLanes)
	for _, l := range lanes {
		if !LaneStatus[l] {
			LaneStatus[l] = true
			return l, true
		}
	}
	return -1, false
}

func Lane(id int) {
	CreateCar(id)
	WaitForPosition(id, 100)
	lane, foundLane := FindAvailableLane()
	if !foundLane {
		ResetCarPosition(id)
		return
	}
	AssignLaneToCar(id, lane)
}

func MoveCarsLogic() {
	for i := len(Cars) - 1; i >= 0; i-- {
		if Cars[i].Position.X < 100 && Cars[i].Lane == -1 && !Cars[i].IsEntering {
			Cars[i].Position.X += 10
			if Cars[i].Position.X > 100 {
				Cars[i].Position.X = 100
			}
		} else if Cars[i].Lane != -1 && !Cars[i].Parked {
			var targetX, targetY float64
			laneWidth := 600.0 / 10
			if Cars[i].Lane < 10 {
				targetX = 100.0 + float64(Cars[i].Lane)*laneWidth + laneWidth/2
				targetY = 400 + (500-350)/2
			} else {
				targetX = 100.0 + float64(Cars[i].Lane-10)*laneWidth + laneWidth/2
				targetY = 100 + (250-100)/2
			}
			ParkCar(&Cars[i], targetX, targetY)
		}
	}
	ExitCarLogic()
}

func ExitCarLogic() {
	for i := len(Cars) - 1; i >= 0; i-- {
		if Cars[i].Parked && time.Now().After(Cars[i].ExitTime) && !Cars[i].IsEntering {
			if !Cars[i].Teleporting {
				Cars[i].Teleporting = true
				Cars[i].TeleportStartTime = time.Now()
				Cars[i].Position.X = 400
				Cars[i].Position.Y = 300
			} else if time.Since(Cars[i].TeleportStartTime) >= time.Millisecond*500 {
				updateLaneStatus(Cars[i].Lane, false)
				removeCar(i)
			}
		}
	}
}

func updateLaneStatus(lane int, status bool) {
	LaneMutex.Lock()
	defer LaneMutex.Unlock()
	LaneStatus[lane] = status
}
