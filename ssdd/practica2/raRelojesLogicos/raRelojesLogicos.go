/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
 */
package raRelojesLogicos

import (
	"practica2/ms"
	"sync"
)

const (
	N = 3
)

type Request struct {
	Clock int
	Pid   int
	Op    string
}

type Reply struct{}

type Pair struct {
	Op1 string
	Op2 string
}

type RASharedDB struct {
	Me        int    // Pid del proceso
	Op        string // Tipo de proceso "Reader" o "Writer"
	OurSeqNum int
	HigSeqNum int
	OutRepCnt int
	ReqCS     bool
	Exclusion map[Pair]bool
	RepDefd   []bool
	ms        *ms.MessageSystem
	done      chan bool
	chrep     chan bool  //channel replies
	Mutex     sync.Mutex // mutex para proteger concurrencia sobre las variables
	// TODO: completar
}

func New(me int, usersFile string, op string) *RASharedDB {
	messageTypes := []ms.Message{Request{}, Reply{}}
	msgs := ms.New(me, usersFile, messageTypes)
	ra := RASharedDB{me, op, 0, 0, 0, false, make(map[Pair]bool), []bool{}, &msgs, make(chan bool), make(chan bool), sync.Mutex{}}
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
					if ra.HigSeqNum < msg.Clock {
						ra.HigSeqNum = msg.Clock
					}
					deferIt := ra.ReqCS && (msg.Clock > ra.OurSeqNum || (msg.Clock == ra.OurSeqNum && msg.Pid > ra.Me)) && ra.Exclusion[Pair{ra.Op, msg.Op}]
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
	ra.OurSeqNum = ra.HigSeqNum + 1
	ra.Mutex.Unlock()

	ra.OutRepCnt = N - 1

	for i := 1; i <= N; i++ {
		if i != ra.Me {
			ra.ms.Send(i, Request{ra.OurSeqNum, ra.Me, ra.Op})
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
		if defered {
			ra.RepDefd[i-1] = false
			ra.ms.Send(i, Reply{})
		}
	}
}

func (ra *RASharedDB) Stop() {
	ra.ms.Stop()
	ra.done <- true
}
