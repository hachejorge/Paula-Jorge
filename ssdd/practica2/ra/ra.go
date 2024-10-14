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
	"practica2/ms"
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vclock"
)

const (
	N = 3
)

type Request struct {
	Clock govec.VectorClock
	Pid   int
	Op    string
}

type Reply struct{}

type Pair struct {
	Op1 string
	Op2 string
}

type RASharedDB struct {
	// OurSeqNum int
	// HigSeqNum int
	me        int
	Op        string
	OutRepCnt int
	ReqCS     bool
	Exclusion map[Pair]bool
	RepDefd   []bool
	ms        *ms.MessageSystem
	done      chan bool
	chrep     chan bool  //channel replies
	Mutex     sync.Mutex // mutex para proteger concurrencia sobre las variables
	// TODO: completar
	logger *govec.GoLog
	vClock govec.VectorClock
}

func New(me int, usersFile string, op string) *RASharedDB {
	messageTypes := []ms.Message{Request{}, Reply{}}
	msgs := ms.New(me, usersFile, messageTypes)

	logger := govec.InitGoVector(strconv.Itoa(me), "LogFile", govec.GetDefaultConfig())

	vClock := govec.NewVectorClock()

	ra := RASharedDB{me, op, 0, false, make(map[Pair]bool), []bool{}, &msgs, make(chan bool), make(chan bool), sync.Mutex{}, logger, vClock}

	ra.Exclusion[Pair{"Reader", "Reader"}] = false
	ra.Exclusion[Pair{"Reader", "Writer"}] = true
	ra.Exclusion[Pair{"Writer", "Reader"}] = true
	ra.Exclusion[Pair{"Writer", "Writer"}] = true

	go func() {
		for {
			select {
			case <-ra.done:
				return
			default:
				switch msg := (ra.ms.Receive()).(type) {
				case Request:
					ra.Mutex.Lock()

					ra.vClock.Tick(ra.me)
					//deferIt := ra.ReqCS && (ra.vClock.Compare(msg.Clock, govec.Descendant) == 1 || (ra.vClock.Compare(msg.Clock, govec.Equal) == 1 && msg.Pid > ra.me))
					deferIt := ra.ReqCS && happensBefore(ra.vClock, msg.Clock, ra.me, msg.Pid) && ra.Exclusion[Pair{ra.Op, msg.Op}]

					ra.vClock.Merge(msg.Clock)

					ra.Mutex.Unlock()

					if deferIt {
						ra.RepDefd[msg.Pid-1] = true
					} else {
						ra.ms.Send(msg.Pid, Reply{})
					}

				case Reply:

					if ra.ReqCS {
						ra.OutRepCnt = ra.OutRepCnt - 1
						if ra.OutRepCnt == 0 {
							ra.chrep <- true
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
	ra.vClock.Tick(ra.me - 1)
	ra.Mutex.Unlock()

	ra.OutRepCnt = N - 1

	for i := 1; i <= N; i++ {
		if i != ra.me {
			ra.ms.Send(i, Request{Clock: ra.vClock, Pid: ra.me, Op: ra.Op})
		}
	}
	<-ra.chrep
}

// Pre: Verdad
// Post: Realiza  el  PostProtocol  para el  algoritmo de
//
//	Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol() {
	ra.ReqCS = false
	for i, defered := range ra.RepDefd {
		if defered == true {
			ra.RepDefd[i] = false
			ra.ms.Send(i+1, Reply{})
		}
	}
}

func (ra *RASharedDB) Stop() {
	ra.ms.Stop()
	ra.done <- true
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
