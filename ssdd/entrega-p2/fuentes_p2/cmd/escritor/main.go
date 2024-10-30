/*
- AUTORES: Jorge Hernández Aznar (872838) y Paula Blasco Díaz (874320)
- ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
- Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
- FECHA: octubre de 2024
- FICHERO: main.go
- DESCRIPCIÓN: Implementación del escritor para el problema lectores-escritores mediante el algoritmo de ra
*/

package main

import (
	"fmt"
	"os"
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/ra"
	"strconv"
)

// Como primer argumento se manda el pid artificial del proceso
func main() {
	fmt.Println("Iniciando escritor " + os.Args[1])

	// Crea su propio fichero según su ID
	file := "fichero" + os.Args[1] + ".txt"
	_, err := os.Create(file)
	if err != nil {
		fmt.Println("Error al crear el fichero" + file)
		os.Exit(1)
	}

	// Crea el msgSystem con tipos de mensaje Barrier, Upgrade y Reply para actualizar los ficheros
	me, _ := strconv.Atoi(os.Args[1])
	msgTypes := []ms.Message{mm.Reply{}, mm.Upgrade{}, mm.Barrier{}}
	msgs := ms.New(me, "../../ms/users.txt", msgTypes)

	okBarrier := make(chan bool)
	okUpgrade := make(chan bool)

	// Proceso que gestiona la comunicación de mensajes con la barrera, y otros lectores y escritores
	go mm.ManageMsg(&msgs, file, okBarrier, okUpgrade)

	// Crea RA
	raData := ra.New(me, "../../ms/usersRA.txt", "Writer")

	// Se comunica con la barrera
	msgs.Send(ra.N+1, mm.Barrier{})
	<-okBarrier
	fmt.Println("Se ha superado la barrera")

	// Texto a escribir en el fichero, su PID
	text := os.Args[1]

	for {
		fmt.Println("Quiero entrar a SC")
		raData.PreProtocol()
		fmt.Println("He entrado a SC")
		// Modifico el fichero
		mm.EscribirFichero(file, os.Args[1])
		for i := 1; i <= ra.N; i++ {
			if i != me {
				// Mando el texto a actualizar al resto de procesos
				msgs.Send(i, mm.Upgrade{Origin: me, Text: text})
				fmt.Println("Enviada upgrade a ", i)
			}
		}
		// Recibir N replies confirmando la upgrade
		for i := 1; i < ra.N; i++ {
			<-okUpgrade
		}
		fmt.Println("Respuestas al upgrade completadas")
		raData.PostProtocol()
		fmt.Println("He salido de SC")
	}
}
