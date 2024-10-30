/*
- AUTORES: Jorge Hernández Aznar (872838) y Paula Blasco Díaz (874320)
- ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
- Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
- FECHA: octubre de 2024
- FICHERO: ricart-agrawala.go
- DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado para lectores-escritores en Go con relojes vectoriales
*/
package ra

import (
	"fmt"
	"practica2/ms"
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec/vclock"
)

// Constante de número de procesos
const (
	N = 3
)

// Elementos de la Request
type Request struct {
	Clock vclock.VClock
	Pid   int
	Op    string
}

// Elementos de la Reply
type Reply struct {
	Pid int
}

// Pareja de valores para la exclusión
type Pair struct {
	Op1 string
	Op2 string
}

type RASharedDB struct {
	Me        int               // Identificador del proceso
	Op        string            // Tipo de proceso, "Reader" o "Writer"
	OutRepCnt int               // Candidad de ack's para entrar a SC recibidos
	ReqCS     bool              // Intención de entrar a SC
	Exclusion map[Pair]bool     // Mapa con las reglas de exclusión
	RepDefd   []bool            // Vector con las replies a procesos diferidas
	Ms        *ms.MessageSystem // Módulo actor para el middleware
	Done      chan bool         // channel para finalizar el algoritmo
	Chrep     chan bool         // channel replies
	Mutex     sync.Mutex        // mutex para proteger concurrencia sobre las variables
	VClock    vclock.VClock     // Reloj vectorial que solicita entrar a SC (OurSequenceNumber)
	VClockMax vclock.VClock     // Reloj vectorial máximo registrado por el proceso (MaxSequecuenceNumber)
}

func New(me int, usersFile string, op string) *RASharedDB {
	// Se inicializa el ms
	messageTypes := []ms.Message{Request{}, Reply{}}
	msgs := ms.New(me, usersFile, messageTypes)

	// Inicializa los relojes vectoriales
	vClock := vclock.New()
	vClockMax := vclock.New()
	for i := 0; i < N; i++ {
		vClock.Set(strconv.Itoa(i), 0)
		vClockMax.Set(strconv.Itoa(i), 0)

	}

	// Crea la estructura de datos
	ra := RASharedDB{me, op, 0, false, make(map[Pair]bool), make([]bool, N), &msgs, make(chan bool), make(chan bool), sync.Mutex{}, vClock, vClockMax}

	// Valores de la matriz de exclusión mutua
	ra.Exclusion[Pair{"Reader", "Reader"}] = false
	ra.Exclusion[Pair{"Reader", "Writer"}] = true
	ra.Exclusion[Pair{"Writer", "Reader"}] = true
	ra.Exclusion[Pair{"Writer", "Writer"}] = true

	// Proceso que se encarga de recibir las replies y las requests
	go func() {
		for {
			select {
			case <-ra.Done:
				return
			default:
				switch msg := (ra.Ms.Receive()).(type) {
				// Recibe request
				case Request:

					ra.Mutex.Lock()
					// Aumenta el reloj máximo
					ra.VClockMax.Tick(strconv.Itoa(ra.Me - 1))
					// Se queda con los valores máximos de mezclarlo con el reloj recibido
					ra.VClockMax.Merge(msg.Clock)

					fmt.Print("Mi reloj que solicita entrar ")
					ra.VClock.PrintVC()
					fmt.Println("VS")
					fmt.Print("Otro reloj que solicita entrar de ", msg.Pid)
					msg.Clock.PrintVC()

					// Se difiere si se solicita entrar a SC, si tiene reloj anterior y tiene que haber exclusión mutua porque hay algún proceso "Writer"
					deferIt := ra.ReqCS && happensBefore(ra.VClock, msg.Clock, ra.Me, msg.Pid) && ra.Exclusion[Pair{ra.Op, msg.Op}]
					ra.Mutex.Unlock()

					// Si se difiere se modifica el vector con los procesos diferidos
					if deferIt {
						ra.RepDefd[msg.Pid-1] = true
						fmt.Println("Request de ", msg.Pid, "diferida")
					} else { // Si no se manda la concesión a la petición
						ra.Ms.Send(msg.Pid, Reply{ra.Me})
						fmt.Println("Enviada reply a ", msg.Pid)
					}

				// Rebibe reply
				case Reply:
					if ra.ReqCS {
						fmt.Println("Recibido reply de ", msg.Pid)
						// Se aumentan las respuesta recibidas
						ra.OutRepCnt = ra.OutRepCnt - 1
						if ra.OutRepCnt == 0 {
							ra.Chrep <- true // Si llegan todas las replies necesarias se concede acceso a SC
						}
					}

				}
			}
		}
	}()

	return &ra
}

// Pre: Verdad
// Post: Realiza  el  PreProtocol  para el  algoritmo de
//
//	Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol() {
	ra.Mutex.Lock()
	ra.ReqCS = true
	// Se aumenta el reloj máximo y se copia siendo este el que solicita la entrada a SC
	ra.VClockMax.Tick(strconv.Itoa(ra.Me - 1))
	ra.VClock = ra.VClockMax.Copy()

	fmt.Print("Mi reloj que solicita entrar ")
	ra.VClockMax.PrintVC()
	ra.Mutex.Unlock()

	// Se reinician la cantidad de replies recibidads
	ra.OutRepCnt = N - 1

	for i := 1; i <= N; i++ {
		if i != ra.Me {
			// Mando todas las peticiones
			ra.Ms.Send(i, Request{Clock: ra.VClock, Pid: ra.Me, Op: ra.Op})
			fmt.Println("Request enviada para acceder a SC a", i)
		}
	}

	<-ra.Chrep
}

// Pre: Verdad
// Post: Realiza  el  PostProtocol  para el  algoritmo de
//
//	Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol() {
	ra.Mutex.Lock()
	ra.ReqCS = false
	for i, defered := range ra.RepDefd {
		if defered {
			ra.Ms.Send(i+1, Reply{})
			ra.RepDefd[i] = false
			fmt.Println("Enviada confirmación para acceder a SC a", i+1)
		}
	}
	ra.Mutex.Unlock()
}

func (ra *RASharedDB) Stop() {
	ra.Ms.Stop()
	ra.Done <- true
}

// Función que verifica que a sea anterior a b y en caso de ser concurrente distingue según el menor pid de ambos relojes
func happensBefore(a vclock.VClock, b vclock.VClock, pid_a int, pid_b int) bool {
	if a.Compare(b, vclock.Descendant) {
		return true
	} else if a.Compare(b, vclock.Concurrent) {
		return pid_a < pid_b
	} else {
		return false
	}
}
