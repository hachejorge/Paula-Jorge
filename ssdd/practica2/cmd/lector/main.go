package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/raRelojesLogicos"
	"strconv"
)

func LeerFichero(file string) string {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error al leer el fichero " + file)
	}
	return string(text)
}

// Como primer argumento se manda
func main() {
	fmt.Println("Iniciando lector " + os.Args[1])

	// Crea su propio fichero
	file := "fichero" + os.Args[1] + ".txt"
	_, err := os.Create(file)
	if err != nil {
		fmt.Println("Error al crear el fichero")
		os.Exit(1)
	}

	me, _ := strconv.Atoi(os.Args[1])
	msgTypes := []ms.Message{raRelojesLogicos.Request{}, raRelojesLogicos.Reply{}, mm.Upgrade{}, mm.Barrier{}}
	msgs := ms.New(me, "../../ms/users.txt", msgTypes)

	requests := make(chan raRelojesLogicos.Request) //
	replies := make(chan raRelojesLogicos.Reply)    //
	okBarrier := make(chan bool)

	go mm.ManageMsg(&msgs, file, requests, replies, okBarrier)

	// Crea RA

	raData := raRelojesLogicos.New(me, "../../ms/usersRa.txt", "Reader")

	// Se comunica con la barrera
	msgs.Send(raRelojesLogicos.N+1, mm.Barrier{})
	<-okBarrier

	for {
		raData.PreProtocol()
		LeerFichero(file)
		raData.PostProtocol()
	}
}
