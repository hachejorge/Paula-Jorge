package main

import (
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/raRelojesLogicos"
)

func main() {
	me := raRelojesLogicos.N + 1
	messageTypes := []ms.Message{mm.Barrier{}}
	msgs := ms.New(me, "../ms/users.txt", messageTypes)
	for i := 1; i <= raRelojesLogicos.N; i++ {
		_ = msgs.Receive()
	}
	for i := 1; i <= raRelojesLogicos.N; i++ {
		msgs.Send(i, mm.Barrier{})
	}
}
