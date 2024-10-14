package main

import (
	"fmt"
	"os"
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/raRelojesLogicos"
	"strconv"
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
	msgTypes := []ms.Message{raRelojesLogicos.Request{}, raRelojesLogicos.Reply{}, mm.Upgrade{}, mm.Barrier{}}
	msgs := ms.New(me, "../../ms/users.txt", msgTypes)

	requests := make(chan raRelojesLogicos.Request)
	replies := make(chan raRelojesLogicos.Reply)
	okBarrier := make(chan bool)

	go mm.ManageMsg(&msgs, file, requests, replies, okBarrier)

	// Crea RA

	raData := raRelojesLogicos.New(me, "../../ms/usersRAs.txt", "Writer")

	// Se comunica con la barrera
	msgs.Send(raRelojesLogicos.N+1, mm.Barrier{})
	<-okBarrier

	text := os.Args[1]

	for {
		raData.PreProtocol()
		mm.EscribirFichero(file, os.Args[1])
		for i := 1; i < raRelojesLogicos.N; i++ {
			if i != me {
				msgs.Send(me, mm.Upgrade{Text: text})
			}
		}
		raData.PostProtocol()
	}
}
