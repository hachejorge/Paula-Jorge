// /*
//   - AUTOR: Rafael Tolosana Calasanz
//   - ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
//   - Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
//   - FECHA: septiembre de 2021
//   - FICHERO: ricart-agrawala.go
//   - DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
//     */
package ra

// import (
// 	"practica2/ms"
// 	"sync"

// 	"github.com/DistributedClocks/GoVector/govec"
// )

// const (
// 	N = 4
// )

// type Request struct {
// 	Clock govec.VectorClock
// 	Pid   int
// }

// type Reply struct{}

// type RASharedDB struct {
// 	// OurSeqNum int
// 	// HigSeqNum int
// 	OutRepCnt int
// 	ReqCS     bool
// 	RepDefd   []bool
// 	ms        *MessageSystem
// 	done      chan bool
// 	chrep     chan bool  //channel replies
// 	Mutex     sync.Mutex // mutex para proteger concurrencia sobre las variables
// 	// TODO: completar
// 	logger *govec.GoLog
// 	vClock govec.VectorClock
// }

// func New(me int, usersFile string) *RASharedDB {
// 	messageTypes := []Message{Request{}, Reply{}}
// 	msgs := ms.New(me, usersFile, messageTypes)
// 	logger := govec.InitGoVector(msgs.me, "LogFile", govec.GetDefaultConfig())

// 	vClock := govec.NewVectorClock()

// 	ra := RASharedDB{0, false, []bool{}, &msgs, make(chan bool), make(chan bool), sync.Mutex{}, logger, vClock}
// 	// TODO completar

// 	go func() {
// 		for {
// 			select {
// 			case <-ra.done:
// 				return
// 			default:
// 				switch msg := (ra.ms.Receive()).(type) {
// 				case Request:
// 					ra.Mutex.Lock()

// 					// deferIt := ra.ReqCS && (msg.Clock > ra.OurSeqNum || (msg.Clock == ra.OurSeqNum && msg.Pid > ra.ms.me))
// 					deferIt := ra.ReqCS && (ra.vClock.Compare(msg.Clock, Descendant) == 1 || (ra.vClock.Compare(msg.Clock, Equal) == 1 && msg.Pid > ra.ms.me))

// 					ra.vClock.Merge(msg.Clock)

// 					ra.Mutex.Unlock()

// 					if deferIt {
// 						RepDefd[msg.Pid-1]
// 					} else {
// 						ra.vClock.Tick(ra.ms.me)
// 						ra.ms.Send(msg.Pid, Reply{})
// 					}

// 				case Reply:

// 					ra.vClock.Tick(ra.ms.me)

// 					if ra.ReqCs {
// 						ra.OutRepCnt = ra.OutRepCnt - 1
// 						if ra.OutRepCnt == 0 {
// 							ra.chrep <- true
// 						}
// 					}

// 				}
// 			}
// 		}
// 	}()

// 	return &ra
// }

// // Pre: Verdad
// // Post: Realiza  el  PreProtocol  para el  algoritmo de
// //
// //	Ricart-Agrawala Generalizado
// func (ra *RASharedDB) PreProtocol() {
// 	ra.Mutex.Lock()
// 	ra.ReqCS = true
// 	ra.OurSeqNum = HigSeqNum + 1
// 	ra.Mutex.Unlock()

// 	ra.OutRepCnt = ra.ms.peers.len() - 1

// 	for i := 1; i <= ra.ms.peers.len(); i++ {
// 		if i != ra.ms.me {
// 			ra.ms.Send(i, Request{ra.OurSeqNum, ra.ms.me})
// 		}
// 	}
// 	<-ra.chrep
// }

// // Pre: Verdad
// // Post: Realiza  el  PostProtocol  para el  algoritmo de
// //
// //	Ricart-Agrawala Generalizado
// func (ra *RASharedDB) PostProtocol() {
// 	ra.ReqCS = false
// 	for i, defered := range ra.RepDefd {
// 		if defered == true {
// 			ra.RepDefd[i-1] = false
// 			ra.ms.Send(i, Reply{})
// 		}
// 	}
// }

// func (ra *RASharedDB) Stop() {
// 	ra.ms.Stop()
// 	ra.done <- true
// }
