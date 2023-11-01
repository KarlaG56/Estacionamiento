package Models

import (
	"math/rand"
	"sync"
	"time"
)

const (
	NumeroCarriles = 20
	AnchoCarril    = 150.0
)

var (
	EstadoCarriles               [NumeroCarriles]bool
	ListaDeAutos                 []Auto
	MutexAutos                   sync.Mutex
	MutexCarriles                sync.Mutex
	AutoEnProcesoDeEntradaSalida bool
)

func EsperarHastaPosicion(id int, posXObjetivo float64) {
	for {
		posAuto := EncontrarPosicionDelAuto(id)
		if posAuto.X >= posXObjetivo {
			break
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func BuscarCarrilDisponible() (int, bool) {
	MutexCarriles.Lock()
	defer MutexCarriles.Unlock()
	rand.Seed(time.Now().UnixNano())
	ordenCarriles := rand.Perm(NumeroCarriles)
	for _, carril := range ordenCarriles {
		if !EstadoCarriles[carril] {
			EstadoCarriles[carril] = true
			return carril, true
		}
	}
	return -1, false
}

func SeleccionarCarril(id int) {
	EsperarHastaPosicion(id, 100)
	carril, encontrado := BuscarCarrilDisponible()
	if !encontrado {
		ReiniciarPosicionDelAuto(id)
		return
	}
	AsignarCarrilAlAuto(id, carril)
}

func ActualizarEstadoCarril(lane int, status bool) {
	MutexCarriles.Lock()
	defer MutexCarriles.Unlock()
	EstadoCarriles[lane] = status
}
