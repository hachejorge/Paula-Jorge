package barrier

import (
	"practica2/ms"
	"practica2/ra"
)

type Barrier struct{}

func main() {
	me := ra.N + 1
	messageTypes := []ms.Message{Barrier{}}
	msgs := ms.New(me, "../ms/users.txt", messageTypes)
	for i := 1; i <= ra.N; i++ {
		_ = msgs.Receive()
	}
	for i := 1; i <= ra.N; i++ {
		msgs.Send(i, Barrier{})
	}
}
