package main

import (
	//"errors"
	//"fmt"
	//"log"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
	"strconv"
	"strings"
	//"time"
)

func main() {
	almacen := make(map[string]string)

	dns := "raft-service.default.svc.cluster.local:6000"
	// obtener entero de indice de este nodo
	meStrg := os.Args[1]
	nombre := strings.Split(meStrg, "-")[0]

	me, err := strconv.Atoi(strings.Split(meStrg, "-")[1])
	check.CheckError(err, "Main, mal numero entero de indice de nodo:")

	var direcciones []string
	for i := 0; i < 3; i++ {
		nodo := nombre + "-" + strconv.Itoa(i) + "." + dns
		direcciones = append(direcciones, nodo)
	}

	var nodos []rpctimeout.HostPort
	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range direcciones {
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	canalAplicarOperacion := make(chan raft.AplicaOperacion, 1000)

	// *** Depuración ***
	fmt.Println("Depuración: Argumentos para NuevoNodo:")
	fmt.Printf("  nodos: %+v\n", nodos)
	fmt.Printf("  me: %d\n", me)
	fmt.Printf("  direcciones: %+v\n", direcciones)

	// Parte Servidor
	nr := raft.NuevoNodo(nodos, me, canalAplicarOperacion)
	rpc.Register(nr)

	go aplicarOperacion(almacen, canalAplicarOperacion)

	fmt.Println("Réplica escucha en:", me, "de", direcciones)

	l, err := net.Listen("tcp", direcciones[me])
	check.CheckError(err, "Main listen error:")

	for {
		rpc.Accept(l)
	}
}

// aplicaOperacion recibe las operaciones sometidas del nodo raft y almacena o devuelve el valor correspondiente a la clave
func aplicarOperacion(almacen map[string]string, canal chan raft.AplicaOperacion) {
	for {
		op := <-canal
		if op.Operacion.Operacion == "leer" {
			op.Operacion.Valor = almacen[op.Operacion.Clave]
		} else if op.Operacion.Operacion == "escribir" {
			almacen[op.Operacion.Clave] = op.Operacion.Valor
			op.Operacion.Valor = "ESCRITO CORRECTAMENTE"
		}
		canal <- op
	}
}
