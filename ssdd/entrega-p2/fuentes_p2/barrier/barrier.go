/*
- AUTORES: Jorge Hernández Aznar (872838) y Paula Blasco Díaz (874320)
- ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
- Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
- FECHA: octubre de 2024
- FICHERO: barrier.go
- DESCRIPCIÓN: Implementación de una barrera para el problema de lectores escritores según el MessageSystem
*/

package main

import (
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/ra"
)

func main() {
	me := ra.N + 1
	messageTypes := []ms.Message{mm.Barrier{}}
	msgs := ms.New(me, "../ms/users.txt", messageTypes)
	for i := 1; i <= ra.N; i++ {
		_ = msgs.Receive()
	}
	for i := 1; i <= ra.N; i++ {
		msgs.Send(i, mm.Barrier{})
	}
}
