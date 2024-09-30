/*
* AUTOR: Rafael Tolosana Calasanz y Unai Arronategui
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*                       Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2022
* FICHERO: server-draft.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*                               correspondientes a la práctica 1
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
	"os/user"
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
	fmt.Println("Request recibida del cliente...")
	// Crea conexión con el worker y envia petición
	conn_w, err := net.Dial("tcp", worker)
	com.CheckError(err)
	encoder := gob.NewEncoder(conn_w)
	err = encoder.Encode(request)
	com.CheckError(err)
	fmt.Println("Request enviada al worker...")
	// Recibe respuesta del worker
	var reply com.Reply
	decoder = gob.NewDecoder(conn_w)
	err = decoder.Decode(&reply)
	com.CheckError(err)
	conn_w.Close()
	fmt.Println("Reply recibida del worker...")
	// Envia respuesta al cliente y cierrala conexión
	encoder = gob.NewEncoder(conn)
	encoder.Encode(&reply)
	conn.Close()
	fmt.Println("Reply enviada al cliente...")
}
func poolWorkers(workersChan chan net.Conn, worker string) {
	for {
		// Recibe una petición y la procesa
		conn := <-workersChan
		fmt.Println("Procesando request...")
		processRequest(conn, worker)
	}
}

func main() {
	args := os.Args
	if len(args) != 3 {
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

	fmt.Println("Lanzando los workers...")
	// Se lanzan los workers
	for _, ep := range workers {
		split := strings.Split(ep, ":")
		ip := split[0]

		// Construir el comando SSH para ejecutar el worker en la máquina remota
		usuario, err := user.Current()
		settingWorkers := fmt.Sprintf("ssh %s@%s 'cd /misc/alumnos/sd/sd2425/a872838/practica1/cmd/worker && go run worker.go %s'", usuario.Username, ip, ep)
		fmt.Println(settingWorkers)
		cmd := exec.Command("bash", "-c", settingWorkers)

		// Ejecutar el comando SSH
		err = cmd.Start()
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
		fmt.Println("Conexión aceptada, delegando al gestor de workers...")
		workersChan <- conn
	}
}
