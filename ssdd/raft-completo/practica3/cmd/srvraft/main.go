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
	//"time"
)

func main() {
	almacen := make(map[string]string)
	// obtener entero de indice de este nodo
	me, err := strconv.Atoi(os.Args[1])
	check.CheckError(err, "Main, mal numero entero de indice de nodo:")

	var nodos []rpctimeout.HostPort
	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range os.Args[2:] {
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	canalAplicarOperacion := make(chan raft.AplicaOperacion, 1000)

	// Parte Servidor
	nr := raft.NuevoNodo(nodos, me, canalAplicarOperacion)
	rpc.Register(nr)

	go aplicarOperacion(almacen, canalAplicarOperacion)

	fmt.Println("Replica escucha en :", me, " de ", os.Args[2:])

	l, err := net.Listen("tcp", os.Args[2:][me])
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
