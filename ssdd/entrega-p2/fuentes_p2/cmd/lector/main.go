/*
- AUTORES: Jorge Hernández Aznar (872838) y Paula Blasco Díaz (874320)
- ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
- Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
- FECHA: octubre de 2024
- FICHERO: main.go
- DESCRIPCIÓN: Implementación del lector para el problema lectores-escritores mediante el algoritmo de ra
*/

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

// Como primer argumento se manda el pid artificial del proceso
func main() {
	fmt.Println("Iniciando lector " + os.Args[1])

	// Crea su propio fichero según su PID
	file := "fichero" + os.Args[1] + ".txt"
	_, err := os.Create(file)
	if err != nil {
		fmt.Println("Error al crear el fichero")
		os.Exit(1)
	}

	// Crea el msgSystem con tipos de mensaje Barrier, Upgrade y Reply para actualizar los ficheros
	me, _ := strconv.Atoi(os.Args[1])
	msgTypes := []ms.Message{mm.Reply{}, mm.Upgrade{}, mm.Barrier{}}
	msgs := ms.New(me, "../../ms/users.txt", msgTypes)
	fmt.Println("Se ha creado el msgs")

	okBarrier := make(chan bool)
	okUpgrade := make(chan bool)

	// Proceso que gestiona la comunicación de mensajes con la barrera, y otros lectores y escritores
	go mm.ManageMsg(&msgs, file, okBarrier, okUpgrade)

	// Crea RA indicando su PID y el tipo de proceso
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
		textFile := LeerFichero(file)
		fmt.Println(textFile)
		raData.PostProtocol()
		fmt.Println("He salido de SC")

	}
}
