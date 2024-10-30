/*
- AUTORES: Jorge Hernández Aznar (872838) y Paula Blasco Díaz (874320)
- ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
- Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
- FECHA: octubre de 2024
- FICHERO: msgManager.go
- DESCRIPCIÓN: Implementación de un handler de mensajes entre escritores y lectores
*/

package mm

import (
	"fmt"
	"os"
	"practica2/ms"
)

type Barrier struct{}

type Upgrade struct {
	Origin int
	Text   string
}

type Reply struct {
	Pid int
}

// Función que escribe al final del fichero file la cadena text
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

// Función que gestiona los mensajes recibidos por escritores y lectores
func ManageMsg(msgs *ms.MessageSystem, file string, okBarrier chan bool, okUpgrade chan bool) {
	for {
		switch msg := (msgs.Receive()).(type) {
		case Barrier:
			// Da el ok a la barrera
			okBarrier <- true
		case Upgrade:
			// Si recibe una update un proceso escritor modifica el fichero
			EscribirFichero(file, msg.Text)
			// Manda el ack de que ha actualizado el fichero
			msgs.Send(msg.Origin, Reply{})
			fmt.Println("Upgrade recibida y reply enviada a ", msg.Origin)
		case Reply:
			// Recibe las replies de las actualizaciones del fichero
			fmt.Println("Reply a la upgrade recibida de ", msg.Pid)
			okUpgrade <- true
		}
	}
}
