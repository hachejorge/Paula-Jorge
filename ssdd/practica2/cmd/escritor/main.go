package main

import (
	"fmt"
	"os"
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/ra"
	"strconv"
	"time"
)

// Como primer argumento se manda
func main() {
	fmt.Println("Iniciando lector " + os.Args[1])

	// Crea su propio fichero
	file := "fichero" + os.Args[1] + ".txt"
	_, err := os.Create(file)
	if err != nil {
		fmt.Println("Error al crear el fichero" + file)
		os.Exit(1)
	}

	me, _ := strconv.Atoi(os.Args[1])
	msgTypes := []ms.Message{mm.Reply{}, mm.Upgrade{}, mm.Barrier{}}
	msgs := ms.New(me, "../../ms/users.txt", msgTypes)
	fmt.Println("Se ha creado el msgs")

	okBarrier := make(chan bool)
	okUpgrade := make(chan bool)

	go mm.ManageMsg(&msgs, file, okBarrier, okUpgrade)

	// Crea RA

	raData := ra.New(me, "../../ms/usersRA.txt", "Writer")

	// Se comunica con la barrera
	msgs.Send(ra.N+1, mm.Barrier{})
	<-okBarrier
	fmt.Println("Se ha superado la barrera")

	text := os.Args[1]

	for {
		time.Sleep(2 * time.Second)
		fmt.Println("Quiero entrar a SC")
		raData.PreProtocol()
		fmt.Println("He entrado a SC")
		mm.EscribirFichero(file, os.Args[1])
		for i := 1; i <= ra.N; i++ {
			if i != me {
				msgs.Send(i, mm.Upgrade{Origin: me, Text: text})
			}
		}
		// Recibir N replies
		for i := 1; i < ra.N-1; i++ {
			<-okUpgrade
		}

		raData.PostProtocol()
		fmt.Println("He salido de SC")
	}
}
