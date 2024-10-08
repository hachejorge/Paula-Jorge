/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
 */
package ra

import (
	"ms"
	"sync"
)

type Request struct {
	Clock int
	Pid   int
}

type Reply struct{}

type RASharedDB struct {
	OurSeqNum int
	HigSeqNum int
	OutRepCnt int
	ReqCS     bool
	RepDefd   []bool
	ms        *MessageSystem
	done      chan bool
	chrep     chan bool  //channel replies
	Mutex     sync.Mutex // mutex para proteger concurrencia sobre las variables
	// TODO: completar
}

func New(me int, usersFile string) *RASharedDB {
	messageTypes := []Message{Request{}, Reply{}}
	msgs := ms.New(me, usersFile, messageTypes)
	ra := RASharedDB{0, 0, 0, false, []bool{}, &msgs, make(chan bool), make(chan bool), &sync.Mutex{}}
	// TODO completar

	go func() {
		for {
			select {
			case <-ra.done:
				return
			default:
				switch msg := (ra.ms.Receive()).(type) {
				case Request:
					ra.Mutex.Lock()
					ra.HigSeqNum = cmd.Max(ra.HigSeqNum, msg.Clock)
					deferIt := ra.ReqCS && (msg.Clock > ra.OurSeqNum || (msg.Clock == ra.OurSeqNum && msg.Pid > ra.ms.me))
					ra.Mutex.Unlock()

					if deferIt {
						RepDefd[msg.Pid-1]
					} else {
						ra.ms.Send(msg.Pid, Reply{})
					}

				case Reply:
					if ra.ReqCs {
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
	ra.OurSeqNum = HigSeqNum + 1
	ra.Mutex.Signal()

	ra.OutRepCnt = ra.ms.peers.len() - 1

	for i := 1; i <= ra.ms.peers.len(); i++ {
		if i != ra.ms.me {
			ra.ms.Send(i, Request{ra.OurSeqNum, ra.ms.me})
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
			ra.RepDefd[i-1] = false
			ra.ms.Send(i, Reply{})
		}
	}
}

func (ra *RASharedDB) Stop() {
	ra.ms.Stop()
	ra.done <- true
}
