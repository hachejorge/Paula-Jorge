/*
* AUTOR: Rafael Tolosana Calasanz y Unai Arronategui
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2022
* FICHERO: server-draft.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes a la práctica 1
 */
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"practica1/com"
	"strings"
)

func readEndpoints(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var endpoints []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			endpoints = append(endpoints, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func processRequest(conn net.Conn, worker string) {
	// Recibe del cliente la petición
	var request com.Request
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&request)
	com.CheckError(err)
	// Crea conexión con el worker y envia petición
	conn_w, err := net.Dial("tcp", worker)
	com.CheckError(err)
	encoder := gob.NewEncoder(conn_w)
	err = encoder.Encode(request)
	com.CheckError(err)
	// Recibe respuesta del worker
	var reply com.Reply
	decoder = gob.NewDecoder(conn_w)
	err = decoder.Decode(&reply)
	com.CheckError(err)
	conn_w.Close()
	// Envia respuesta al cliente y cierra la conexión
	encoder = gob.NewEncoder(conn)
	encoder.Encode(&reply)
	conn.Close()
}

func poolWorkers(workersChan chan net.Conn, worker string) {
	for {
		// Recibe una petición y la procesa
		conn := <-workersChan
		processRequest(conn, worker)
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Println("Error: endpoint missing: go run server.go ip:port fileWorkers")
		os.Exit(1)
	}
	endpoint := args[1]
	listener, err := net.Listen("tcp", endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	workers, err := readEndpoints(os.Args[2])
	if err != nil {
		fmt.Println("Error reading endpoints:", err)
		return
	}

	workersChan := make(chan net.Conn)

	// Se lanzan los workers
	for _, ep := range workers {
		split := strings.Split(ep, ":")
		ip := split[0]

		// Construir el comando SSH para ejecutar el worker en la máquina remota
		cmd := exec.Command("ssh", ip, "go", "run", "/misc/alumnos/sd/sd2425/a872838/practica1/cmd/worker/worker.go", ep)

		// Ejecutar el comando SSH
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error running SSH command on %s: %v\n", ep, err)
			continue
		}
	}

	log.Println("***** Listening for new connection in endpoint ", endpoint)

	// Se crea el pool de workers
	for _, w := range workers {
		go poolWorkers(workersChan, w)
	}

	for {
		// Aceptas peticiones
		conn, err := listener.Accept()
		com.CheckError(err)
		workersChan <- conn
	}
}
