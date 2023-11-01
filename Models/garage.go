package Models

import (
	"github.com/faiface/pixel"
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

type Auto struct {
	ID                        int
	Posicion                  pixel.Vec
	PosicionPrevia            pixel.Vec
	Carril                    int
	Estacionado               bool
	TiempoParaSalir           time.Time
	EnProcesoDeEntrada        bool
	EnProcesoDeTeleportacion  bool
	TiempoInicioTeleportacion time.Time
}

var CanalAutos chan Auto

func InicializarSistemaDeAutos() {
	CanalAutos = make(chan Auto)
	go GeneradorDeAutos()
}

func CrearNuevoAuto(id int) Auto {
	MutexAutos.Lock()
	defer MutexAutos.Unlock()
	nuevoAuto := Auto{
		ID:          id,
		Posicion:    pixel.V(0, 300),
		Carril:      -1,
		Estacionado: false,
	}
	ListaDeAutos = append(ListaDeAutos, nuevoAuto)
	return nuevoAuto
}

func AsignarTiempoDeSalida(auto *Auto) {
	rand.Seed(time.Now().UnixNano())
	tiempoDeSalida := time.Duration(rand.Intn(5)+1) * time.Second
	auto.TiempoParaSalir = time.Now().Add(tiempoDeSalida)
}

func ObtenerListaDeAutos() []Auto {
	MutexAutos.Lock()
	defer MutexAutos.Unlock()
	listaActualizadaDeAutos := make([]Auto, len(ListaDeAutos))
	copy(listaActualizadaDeAutos, ListaDeAutos)
	return listaActualizadaDeAutos
}

func GeneradorDeAutos() {
	id := 0
	for {
		id++
		nuevoAuto := CrearNuevoAuto(id)
		CanalAutos <- nuevoAuto
		time.Sleep(time.Millisecond * 500)
	}
}

func AsignarCarrilAlAuto(id int, carril int) {
	MutexAutos.Lock()
	defer MutexAutos.Unlock()
	for i := range ListaDeAutos {
		if ListaDeAutos[i].ID == id {
			ListaDeAutos[i].Carril = carril
			break
		}
	}
}

func ReiniciarPosicionDelAuto(id int) {
	MutexAutos.Lock()
	defer MutexAutos.Unlock()
	for i := range ListaDeAutos {
		if ListaDeAutos[i].ID == id {
			ListaDeAutos[i].Posicion = pixel.V(0, 300)
			break
		}
	}
}

func EncontrarPosicionDelAuto(id int) pixel.Vec {
	MutexAutos.Lock()
	defer MutexAutos.Unlock()
	for _, auto := range ListaDeAutos {
		if auto.ID == id {
			return auto.Posicion
		}
	}
	return pixel.Vec{}
}

func EstacionarAuto(auto *Auto, posX, posY float64) {
	auto.Posicion.X = posX
	auto.Posicion.Y = posY
	auto.Estacionado = true
	AsignarTiempoDeSalida(auto)
}

func EliminarAuto(indice int) {
	ListaDeAutos = append(ListaDeAutos[:indice], ListaDeAutos[indice+1:]...)
}

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

func LogicaDeMovimientoDeAutos() {
	for i := len(ListaDeAutos) - 1; i >= 0; i-- {
		if ListaDeAutos[i].Posicion.X < 100 && ListaDeAutos[i].Carril == -1 && !ListaDeAutos[i].EnProcesoDeEntrada {
			ListaDeAutos[i].Posicion.X += 10
			if ListaDeAutos[i].Posicion.X > 100 {
				ListaDeAutos[i].Posicion.X = 100
			}
		} else if ListaDeAutos[i].Carril != -1 && !ListaDeAutos[i].Estacionado {
			var destinoX, destinoY float64
			anchoCarril := 600.0 / 10
			if ListaDeAutos[i].Carril < 10 {
				destinoX = 100.0 + float64(ListaDeAutos[i].Carril)*anchoCarril + anchoCarril/2
				destinoY = 400 + (500-350)/2
			} else {
				destinoX = 100.0 + float64(ListaDeAutos[i].Carril-10)*anchoCarril + anchoCarril/2
				destinoY = 100 + (250-100)/2
			}
			EstacionarAuto(&ListaDeAutos[i], destinoX, destinoY)
		}
	}
	LogicaDeSalidaDelAuto()
}

func LogicaDeSalidaDelAuto() {
	for i := len(ListaDeAutos) - 1; i >= 0; i-- {
		if ListaDeAutos[i].Estacionado && time.Now().After(ListaDeAutos[i].TiempoParaSalir) && !ListaDeAutos[i].EnProcesoDeEntrada {
			if !ListaDeAutos[i].EnProcesoDeTeleportacion {
				ListaDeAutos[i].EnProcesoDeTeleportacion = true
				ListaDeAutos[i].TiempoInicioTeleportacion = time.Now()
				ListaDeAutos[i].Posicion = pixel.V(400, 300) // Posición de salida
			} else if time.Since(ListaDeAutos[i].TiempoInicioTeleportacion) >= 500*time.Millisecond {
				ActualizarEstadoCarril(ListaDeAutos[i].Carril, false)
				EliminarAuto(i)
			}
		}
	}
}

func ActualizarEstadoCarril(lane int, status bool) {
	MutexCarriles.Lock()
	defer MutexCarriles.Unlock()
	EstadoCarriles[lane] = status
}
