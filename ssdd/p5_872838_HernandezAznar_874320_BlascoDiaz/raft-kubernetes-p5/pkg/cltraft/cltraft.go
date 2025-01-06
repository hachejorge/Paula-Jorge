package main

import (
	"fmt"
	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
	"strconv"
	"time"
)

func main() {
	dns := "raft-service.default.svc.cluster.local:6000"
	name := "raft"
	var direcciones []string
	for i := 0; i < 3; i++ {
		nodo := name + "-" + strconv.Itoa(i) + "." + dns
		direcciones = append(direcciones, nodo)
	}
	var nodos []rpctimeout.HostPort
	for _, endPoint := range direcciones {
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	time.Sleep(10 * time.Second)

	// Almacén interno para las claves y valores
	almacen := make(map[string]int)

	// Variables para controlar la suma de los valores
	sumaTotal := 0
	alphabet := []rune("abcdefghijklmnopqrstuvwxyz") // Letras del alfabeto

	// Bucle infinito para realizar operaciones continuamente
	for {
		// Seleccionar la clave según el alfabeto cíclico
		clave := string(alphabet[sumaTotal%len(alphabet)]) // Usamos la sumaTotal para ciclo

		// El valor es la suma de todos los valores previos
		valor := sumaTotal + 1

		// Operación de escritura
		operacion1 := raft.TipoOperacion{"escribir", clave, strconv.Itoa(valor)}
		fmt.Printf("Operación de escritura: clave = %s, valor = %d\n", clave, valor)

		var reply raft.ResultadoRemoto
		err := nodos[0].CallTimeout("NodoRaft.SometerOperacionNodo", operacion1, &reply, 5000*time.Millisecond)
		check.CheckError(err, "SometerOperacion")

		// Asegurarse de que la operación se realice solo cuando haya un líder
		for reply.IdLider == -1 {
			fmt.Printf("Operación 1 sometida a 0\n")
			err = nodos[0].CallTimeout("NodoRaft.SometerOperacionNodo", operacion1, &reply, 5000*time.Millisecond)
			check.CheckError(err, "SometerOperacion")
		}
		if !reply.EsLider {
			fmt.Printf("Operación 1 sometida a %d\n", reply.IdLider)
			err = nodos[reply.IdLider].CallTimeout("NodoRaft.SometerOperacionNodo", operacion1, &reply, 5000*time.Millisecond)
			check.CheckError(err, "SometerOperacion")
		}

		// Almacenar el valor en el "almacén" interno
		almacen[clave] = valor
		//fmt.Printf("Valor almacenado internamente: %s = %d\n", clave, valor)

		// Operación de lectura
		operacion2 := raft.TipoOperacion{"leer", clave, ""}
		fmt.Printf("Operación de lectura: clave = %s\n", clave)

		err = nodos[reply.IdLider].CallTimeout("NodoRaft.SometerOperacionNodo", operacion2, &reply, 5000*time.Millisecond)
		check.CheckError(err, "SometerOperacion")

		// Verificar que el valor leído es el esperado
		if reply.ValorADevolver == strconv.Itoa(almacen[clave]) {
			fmt.Printf("Valor correcto leído: %s = %s\n", clave, reply.ValorADevolver)
		} else {
			fmt.Printf("¡Error! Valor leído no coincide con el valor esperado para la clave %s. Esperado: %d, Leído: %s\n", clave, almacen[clave], reply.ValorADevolver)
		}

		// Incrementar la sumaTotal para la siguiente clave y valor
		sumaTotal++

		// Pausa entre operaciones (puedes ajustarlo según la frecuencia de las operaciones que desees)
		time.Sleep(2 * time.Second)
	}
}
