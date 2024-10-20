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
