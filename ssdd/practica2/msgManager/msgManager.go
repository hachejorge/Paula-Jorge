package mm

import (
	"fmt"
	"os"
	"practica2/barrier"
	"practica2/ms"
	"practica2/ra"
)

type Upgrade struct {
	Text string
}

func EscribirFichero(file string, text string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error al abrir el fichero" + file)
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println("Error al escribir en el fichero" + file)
		os.Exit(1)
	}
}

func ManageMsg(msgs *ms.MessageSystem, file string, request chan ra.Request, reply chan ra.Reply, okBarrier chan bool) {
	for {
		switch msg := (msgs.Receive()).(type) {
		case ra.Request: //
			request <- msg
		case ra.Reply: //
			reply <- msg
		case barrier.Barrier:
			okBarrier <- true
		case Upgrade:
			EscribirFichero(file, msg.Text)
		}

	}
}
