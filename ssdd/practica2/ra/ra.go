/*
- AUTOR: Rafael Tolosana Calasanz
- ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
- Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
- FECHA: septiembre de 2021
- FICHERO: ricart-agrawala.go
- DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
*/
package ra

import (
	"fmt"
	"practica2/ms"
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec/vclock"
)

const (
	N = 3
)

type Request struct {
	Clock vclock.VClock
	Pid   int
	Op    string
}

type Reply struct {
	Pid int
}

type Pair struct {
	Op1 string
	Op2 string
}

type RASharedDB struct {
	// OurSeqNum int
	// HigSeqNum int
	Me        int
	Op        string
	OutRepCnt int
	ReqCS     bool
	Exclusion map[Pair]bool
	RepDefd   []bool
	Ms        *ms.MessageSystem
	Done      chan bool
	Chrep     chan bool  //channel replies
	Mutex     sync.Mutex // mutex para proteger concurrencia sobre las variables
	// TODO: completar
	//logger    *govec.GoLog
	VClock    vclock.VClock
	VClockMax vclock.VClock
}

func New(me int, usersFile string, op string) *RASharedDB {
	messageTypes := []ms.Message{Request{}, Reply{}}
	msgs := ms.New(me, usersFile, messageTypes)

	//logger := govec.InitGoVector(strconv.Itoa(me), "LogFile", govec.GetDefaultConfig())

	vClock := vclock.New()
	for i := 0; i < N; i++ {
		vClock.Set(strconv.Itoa(i), 0)
	}

	vClockMax := vclock.New()
	for i := 0; i < N; i++ {
		vClockMax.Set(strconv.Itoa(i), 0)
	}

	ra := RASharedDB{me, op, 0, false, make(map[Pair]bool), make([]bool, N), &msgs, make(chan bool), make(chan bool), sync.Mutex{}, vClock, vClockMax}

	ra.Exclusion[Pair{"Reader", "Reader"}] = false
	ra.Exclusion[Pair{"Reader", "Writer"}] = true
	ra.Exclusion[Pair{"Writer", "Reader"}] = true
	ra.Exclusion[Pair{"Writer", "Writer"}] = true

	go func() {
		for {
			select {
			case <-ra.Done:
				return
			default:
				switch msg := (ra.Ms.Receive()).(type) {
				case Request:

					ra.VClockMax.Tick(strconv.Itoa(ra.Me - 1))

					ra.VClockMax.Merge(msg.Clock)

					ra.Mutex.Lock()

					fmt.Print("Mi reloj que solicita entrar ")
					ra.VClock.PrintVC()
					fmt.Println("VS")
					fmt.Print("Otro reloj que solicita entrar de ", msg.Pid)
					msg.Clock.PrintVC()

					deferIt := ra.ReqCS && happensBefore(ra.VClock, msg.Clock, ra.Me, msg.Pid) && ra.Exclusion[Pair{ra.Op, msg.Op}]

					if deferIt {
						ra.RepDefd[msg.Pid-1] = true
						fmt.Println("Peticion a ", msg.Pid, "diferida")
					} else {
						ra.Ms.Send(msg.Pid, Reply{ra.Me})
						fmt.Println("Enviada reply a ", msg.Pid)
					}

					ra.Mutex.Unlock()

				case Reply:
					if ra.ReqCS {
						fmt.Println("Recibido reply de ", msg.Pid)
						ra.OutRepCnt = ra.OutRepCnt - 1
						if ra.OutRepCnt == 0 {
							ra.Chrep <- true
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
	ra.VClock.Tick(strconv.Itoa(ra.Me - 1))
	ra.VClockMax = ra.VClock.Copy()

	fmt.Print("Mi reloj que solicita entrar ")
	ra.VClockMax.PrintVC()

	ra.OutRepCnt = N - 1

	for i := 1; i <= N; i++ {
		if i != ra.Me {
			//vClockSend.PrintVC();
			ra.Ms.Send(i, Request{Clock: ra.VClock, Pid: ra.Me, Op: ra.Op})
			fmt.Println("Request enviada para acceder a SC a", i)
		}
	}
	ra.Mutex.Unlock()

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

func happensBefore(a vclock.VClock, b vclock.VClock, pid_a int, pid_b int) bool {
	if a.Compare(b, vclock.Descendant) {
		return true
	} else if a.Compare(b, vclock.Concurrent) {
		return pid_a < pid_b
	} else {
		return false
	}
}
