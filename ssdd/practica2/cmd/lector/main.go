package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"practica2/ms"
	mm "practica2/msgManager"
	"practica2/ra"
	"strconv"
	"time"
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
	msgTypes := []ms.Message{ra.Request{}, ra.Reply{}, mm.Upgrade{}, mm.Barrier{}}
	msgs := ms.New(me, "../../ms/users.txt", msgTypes)
	fmt.Println("Se ha creado el msgs")

	okBarrier := make(chan bool)
	okUpgrade := make(chan bool)

	go mm.ManageMsg(&msgs, file, okBarrier, okUpgrade)

	// Crea RA

	raData := ra.New(me, "../../ms/usersRA.txt", "Reader")

	// Se comunica con la barrera
	msgs.Send(ra.N+1, mm.Barrier{})
	<-okBarrier
	fmt.Println("Se ha superado la barrera")

	for {
		time.Sleep(2 * time.Second)
		fmt.Println("Quiero entrar a SC")
		raData.PreProtocol()
		fmt.Println("He entrado a SC")
		LeerFichero(file)
		raData.PostProtocol()
		fmt.Println("He salido de SC")

	}
}
