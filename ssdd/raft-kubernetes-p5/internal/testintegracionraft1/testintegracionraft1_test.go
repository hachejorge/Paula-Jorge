package testintegracionraft1

import (
	"fmt"
	"math/rand"
	"raft/internal/comun/check"

	//"log"
	//"crypto/rand"
	//	"os"
	//	"path/filepath"

	"raft/internal/comun/rpctimeout"
	"raft/internal/despliegue"
	"raft/internal/raft"
	"strconv"
	"testing"
	"time"
)

const (
	//hosts
	// MAQUINA1 = "127.0.0.1"
	// MAQUINA2 = "127.0.0.1"
	// MAQUINA3 = "127.0.0.1"
	MAQUINA1 = "192.168.3.6"
	MAQUINA2 = "192.168.3.7"
	MAQUINA3 = "192.168.3.8"

	//puertos
	PUERTOREPLICA1 = "31139"
	PUERTOREPLICA2 = "31138"
	PUERTOREPLICA3 = "31137"

	//nodos replicas
	REPLICA1 = MAQUINA1 + ":" + PUERTOREPLICA1
	REPLICA2 = MAQUINA2 + ":" + PUERTOREPLICA2
	REPLICA3 = MAQUINA3 + ":" + PUERTOREPLICA3

	// paquete main de ejecutables relativos a PATH previo
	EXECREPLICA = "./main"

	// comandos completo a ejecutar en máquinas remota con ssh. Ejemplo :
	// 				cd $HOME/raft; go run cmd/srvraft/main.go 127.0.0.1:29001

	// Ubicar, en esta constante, nombre de fichero de vuestra clave privada local
	// emparejada con la clave pública en authorized_keys de máquinas remotas

	PRIVKEYFILE = "id_rsa"
)

// PATH de los ejecutables de modulo golang de servicio Raft
// var PATH string = filepath.Join(os.Getenv("HOME"), "tmp", "p5", "raft")
var PATH string = "/misc/alumnos/sd/sd2425/a872838/practica3/cmd/srvraft"

// go run cmd/srvraft/main.go 0 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003
var EXECREPLICACMD string = "cd " + PATH + "; " + EXECREPLICA

// TEST primer rango
func TestPrimerasPruebas(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
		3,
		[]string{REPLICA1, REPLICA2, REPLICA3},
		[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Test1 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T1:soloArranqueYparada",
		func(t *testing.T) { cfg.soloArranqueYparadaTest1(t) })

	// Test2 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T2:ElegirPrimerLider",
		func(t *testing.T) { cfg.elegirPrimerLiderTest2(t) })

	// Test3: tenemos el primer primario correcto
	t.Run("T3:FalloAnteriorElegirNuevoLider",
		func(t *testing.T) { cfg.falloAnteriorElegirNuevoLiderTest3(t) })

	// Test4: Tres operaciones comprometidas en configuración estable
	t.Run("T4:tresOperacionesComprometidasEstable",
		func(t *testing.T) { cfg.tresOperacionesComprometidasEstable(t) })
}

// TEST primer rango
func TestAcuerdosConFallos(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
		3,
		[]string{REPLICA1, REPLICA2, REPLICA3},
		[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Test5: Se consigue acuerdo a pesar de desconexiones de seguidor
	t.Run("T5:AcuerdoAPesarDeDesconexionesDeSeguidor ",
		func(t *testing.T) { cfg.AcuerdoApesarDeSeguidor(t) })

	t.Run("T5:SinAcuerdoPorFallos ",
		func(t *testing.T) { cfg.SinAcuerdoPorFallos(t) })

	t.Run("T5:SometerConcurrentementeOperaciones ",
		func(t *testing.T) { cfg.SometerConcurrentementeOperaciones(t) })

}

// ---------------------------------------------------------------------
//
// Canal de resultados de ejecución de comandos ssh remotos
type canalResultados chan string

func (cr canalResultados) stop() {
	close(cr)

	// Leer las salidas obtenidos de los comandos ssh ejecutados
	for s := range cr {
		fmt.Println(s)
	}
}

// ---------------------------------------------------------------------
// Operativa en configuracion de despliegue y pruebas asociadas
type configDespliegue struct {
	t           *testing.T
	conectados  []bool
	numReplicas int
	nodosRaft   []rpctimeout.HostPort
	cr          canalResultados
	lider       int
}

// Crear una configuracion de despliegue
func makeCfgDespliegue(t *testing.T, n int, nodosraft []string,
	conectados []bool) *configDespliegue {
	cfg := &configDespliegue{}
	cfg.t = t
	cfg.conectados = conectados
	cfg.numReplicas = n
	cfg.nodosRaft = rpctimeout.StringArrayToHostPortArray(nodosraft)
	cfg.cr = make(canalResultados, 2000)
	cfg.lider = -1

	return cfg
}

func (cfg *configDespliegue) stop() {
	//cfg.stopDistributedProcesses()

	time.Sleep(50 * time.Millisecond)

	cfg.cr.stop()
}

// --------------------------------------------------------------------------
// FUNCIONES DE SUBTESTS

// Se pone en marcha una replica ?? - 3 NODOS RAFT
func (cfg *configDespliegue) soloArranqueYparadaTest1(t *testing.T) {
	t.Skip("SKIPPED soloArranqueYparadaTest1")

	fmt.Println(t.Name(), ".....................")

	cfg.t = t // Actualizar la estructura de datos de tests para errores

	//	cfg.stopDistributedProcesses()

	// Poner en marcha replicas en remoto con un tiempo de espera incluido
	cfg.startDistributedProcesses()
	fmt.Println("Procesos levantados")

	//	fmt.Println("Procesos iniciados")
	//time.Sleep(2000 * time.Millisecond)

	// Comprobar estado replica 0
	cfg.comprobarEstadoRemoto(0, 0, false, -1)

	// Comprobar estado replica 1
	cfg.comprobarEstadoRemoto(1, 0, false, -1)

	// Comprobar estado replica 2
	cfg.comprobarEstadoRemoto(2, 0, false, -1)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// Primer lider en marcha - 3 NODOS RAFT
func (cfg *configDespliegue) elegirPrimerLiderTest2(t *testing.T) {
	t.Skip("SKIPPED ElegirPrimerLiderTest2")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()
	time.Sleep(6000 * time.Millisecond)

	cfg.estadoNodo(0)
	cfg.estadoNodo(1)
	cfg.estadoNodo(2)

	// Se ha elegido lider ?
	fmt.Printf("Probando lider en curso\n")
	cfg.pruebaUnLider(3)

	// Parar réplicas alamcenamiento en remoto
	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// Fallo de un primer lider y reeleccion de uno nuevo - 3 NODOS RAFT
func (cfg *configDespliegue) falloAnteriorElegirNuevoLiderTest3(t *testing.T) {
	t.Skip("SKIPPED FalloAnteriorElegirNuevoLiderTest3")
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	lider := cfg.pruebaUnLider(3)
	cfg.estadoNodo(lider)

	// Desconectar lider
	fmt.Println("Desconectamos el lider")
	cfg.stopDistributedProcess(lider)

	//time.Sleep(6000 * time.Millisecond)
	fmt.Println("Relanzamos todos los procesos")
	cfg.restartDistributedProcesses()

	time.Sleep(6000 * time.Millisecond)

	fmt.Printf("Comprobar nuevo lider\n")
	lider = cfg.pruebaUnLider(3)
	cfg.estadoNodo(lider)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses() //parametros

	fmt.Println(".............", t.Name(), "Superado")

}

// 3 operaciones comprometidas con situacion estable y sin fallos - 3 NODOS RAFT
func (cfg *configDespliegue) tresOperacionesComprometidasEstable(t *testing.T) {
	t.Skip("SKIPPED tresOperacionesComprometidasEstable")
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	lider := cfg.pruebaUnLider(3)
	cfg.estadoNodo(lider)

	cfg.checkSometerOperation(0, "escribir", "a", "primero", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(1, "escribir", "b", "segundo", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(2, "leer", "a", "", "primero")

	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// Se consigue acuerdo a pesar de desconexiones de seguidor -- 3 NODOS RAFT
func (cfg *configDespliegue) AcuerdoApesarDeSeguidor(t *testing.T) {
	//t.Skip("SKIPPED AcuerdoApesarDeSeguidor")
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	lider := cfg.pruebaUnLider(3)
	cfg.estadoNodo(lider)

	// Comprometer una entrada
	cfg.checkSometerOperation(0, "escribir", "a", "primero", "ESCRITO CORRECTAMENTE")

	//  Obtener un lider y, a continuación desconectar una de los nodos Raft
	cfg.stopDistributedProcess(randomExcluding(lider))

	// Comprobar varios acuerdos con una réplica desconectada
	cfg.checkSometerOperation(1, "escribir", "b", "segundo", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(2, "leer", "a", "", "primero")
	cfg.checkSometerOperation(3, "leer", "b", "", "segundo")

	// reconectar nodo Raft previamente desconectado y comprobar varios acuerdos
	cfg.restartDistributedProcesses()

	cfg.checkSometerOperation(4, "escribir", "c", "tercero", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(5, "escribir", "d", "cuarto", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(6, "leer", "d", "", "cuarto")

	cfg.checkLoggers(6)

	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// NO se consigue acuerdo al desconectarse mayoría de seguidores -- 3 NODOS RAFT
func (cfg *configDespliegue) SinAcuerdoPorFallos(t *testing.T) {
	//t.Skip("SKIPPED SinAcuerdoPorFallos")
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	lider := cfg.pruebaUnLider(3)
	cfg.estadoNodo(lider)

	// Comprometer una entrada
	cfg.checkSometerOperation(0, "escribir", "a", "primero", "ESCRITO CORRECTAMENTE")

	//  Obtener un lider y, a continuación desconectar los followers
	cfg.stopFollowers()

	// Comprobar varios acuerdos con ningún lider conectado
	cfg.checkSometerOperationFail("escribir", "b", "segundo")
	cfg.checkSometerOperationFail("leer", "a", "")
	cfg.checkSometerOperationFail("leer", "b", "")

	// reconectar los nodo Raft previamente desconectados y comprobar varios acuerdos
	cfg.restartDistributedProcesses()

	cfg.checkSometerOperation(4, "escribir", "c", "tercero", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(5, "escribir", "d", "cuarto", "ESCRITO CORRECTAMENTE")
	cfg.checkSometerOperation(6, "leer", "d", "", "cuarto")

	cfg.checkLoggers(6)

	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// Se somete 5 operaciones de forma concurrente -- 3 NODOS RAFT
func (cfg *configDespliegue) SometerConcurrentementeOperaciones(t *testing.T) {
	//t.Skip("SKIPPED SometerConcurrentementeOperaciones")
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	// Obtener un lider y, a continuación someter una operacion
	lider := cfg.pruebaUnLider(3)
	cfg.estadoNodo(lider)
	cfg.checkSometerOperation(0, "escribir", "a", "primero", "ESCRITO CORRECTAMENTE")

	// Someter 5  operaciones concurrentes
	go cfg.someterOperation("escribir", "b", "segundo A")
	go cfg.someterOperation("escribir", "b", "segundo B")
	go cfg.someterOperation("leer", "a", "")
	go cfg.someterOperation("escribir", "c", "tercero")
	go cfg.someterOperation("leer", "b", "")

	time.Sleep(6000 * time.Millisecond)

	// Comprobar estados de nodos Raft, sobre todo
	// el avance del mandato en curso e indice de registro de cada uno
	// que debe ser identico entre ellos
	cfg.checkLoggers(5)

	fmt.Println(".............", t.Name(), "Superado")

}

// --------------------------------------------------------------------------
// FUNCIONES DE APOYO
// Comprobar que hay un solo lider
// probar varias veces si se necesitan reelecciones
func (cfg *configDespliegue) pruebaUnLider(numreplicas int) int {
	for iters := 0; iters < 10; iters++ {
		time.Sleep(500 * time.Millisecond)
		mapaLideres := make(map[int][]int)
		for i := 0; i < numreplicas; i++ {
			if cfg.conectados[i] {
				if _, mandato, eslider, _ := cfg.obtenerEstadoRemoto(i); eslider {
					mapaLideres[mandato] = append(mapaLideres[mandato], i)
				}
			}
		}

		ultimoMandatoConLider := -1
		for mandato, lideres := range mapaLideres {
			if len(lideres) > 1 {
				cfg.t.Fatalf("mandato %d tiene %d (>1) lideres",
					mandato, len(lideres))
			}
			if mandato > ultimoMandatoConLider {
				ultimoMandatoConLider = mandato
			}
		}

		if len(mapaLideres) != 0 {
			cfg.lider = mapaLideres[ultimoMandatoConLider][0]
			return cfg.lider // Termina

		}
	}
	cfg.t.Fatalf("un lider esperado, ninguno obtenido")

	return -1 // Termina
}

func (cfg *configDespliegue) obtenerEstadoRemoto(
	indiceNodo int) (int, int, bool, int) {
	var reply raft.EstadoRemoto
	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoNodo",
		raft.Vacio{}, &reply, 10*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")

	return reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider
}

func (cfg *configDespliegue) obtenerEstadoLogger(indiceNodo int) (int, int) {
	var reply raft.EstadoLogger
	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoLogger",
		raft.Vacio{}, &reply, 10*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ObtenerEstadoLogger")

	return reply.Indice, reply.Mandato
}

// start  gestor de vistas; mapa de replicas y maquinas donde ubicarlos;
// y lista clientes (host:puerto)
func (cfg *configDespliegue) startDistributedProcesses() {
	//cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)

	for i, endPoint := range cfg.nodosRaft {
		despliegue.ExecMutipleHosts(EXECREPLICACMD+
			" "+strconv.Itoa(i)+" "+
			rpctimeout.HostPortArrayToString(cfg.nodosRaft),
			[]string{endPoint.Host()}, cfg.cr, PRIVKEYFILE)

		// dar tiempo para se establezcan las replicas
		//time.Sleep(500 * time.Millisecond)
	}

	// aproximadamente 500 ms para cada arranque por ssh en portatil
	time.Sleep(2500 * time.Millisecond)
}

func (cfg *configDespliegue) restartDistributedProcesses() {
	for i, endPoint := range cfg.nodosRaft {
		if !cfg.conectados[i] {
			despliegue.ExecMutipleHosts(EXECREPLICACMD+
				" "+strconv.Itoa(i)+" "+
				rpctimeout.HostPortArrayToString(cfg.nodosRaft),
				[]string{endPoint.Host()}, cfg.cr, PRIVKEYFILE)

			// dar tiempo para se establezcan las replicas
			//time.Sleep(500 * time.Millisecond)
			cfg.conectados[i] = true
			fmt.Printf("Nodo %d ha sido relanzado\n", i)
		}
	}
	time.Sleep(2500 * time.Millisecond)
}

func (cfg *configDespliegue) stopDistributedProcesses() {
	var reply raft.Vacio

	for _, endPoint := range cfg.nodosRaft {
		err := endPoint.CallTimeout("NodoRaft.ParaNodo",
			raft.Vacio{}, &reply, 10*time.Millisecond)
		check.CheckError(err, "Error en llamada RPC Para nodo")
	}
}

// Comprobar estado remoto de un nodo con respecto a un estado prefijado
func (cfg *configDespliegue) comprobarEstadoRemoto(idNodoDeseado int,
	mandatoDeseado int, esLiderDeseado bool, IdLiderDeseado int) {
	idNodo, mandato, esLider, idLider := cfg.obtenerEstadoRemoto(idNodoDeseado)

	//cfg.t.Log("Estado replica 0: ", idNodo, mandato, esLider, idLider, "\n")

	// Muestra el estado actual de la réplica remota con los valores reales recibidos
	fmt.Printf("Estado real del nodo     - ID: %d, Mandato: %d, Es Líder: %t, ID del Líder: %d\n",
		idNodo, mandato, esLider, idLider)

	// Muestra el estado esperado para comparación
	fmt.Printf("Estado esperado del nodo - ID: %d, Mandato: %d, Es Líder: %t, ID del Líder: %d\n",
		idNodoDeseado, mandatoDeseado, esLiderDeseado, IdLiderDeseado)

	if idNodo != idNodoDeseado || mandato != mandatoDeseado ||
		esLider != esLiderDeseado || idLider != IdLiderDeseado {
		cfg.t.Fatalf("Estado incorrecto en replica %d en subtest %s",
			idNodoDeseado, cfg.t.Name())
	}

}

// Comprobar estado remoto de un nodo con respecto a un estado prefijado
func (cfg *configDespliegue) estadoNodo(idNodoDeseado int) {
	idNodo, mandato, esLider, idLider := cfg.obtenerEstadoRemoto(idNodoDeseado)

	//cfg.t.Log("Estado replica 0: ", idNodo, mandato, esLider, idLider, "\n")

	// Muestra el estado actual de la réplica remota con los valores reales recibidos
	fmt.Printf("Estado del nodo - ID: %d, Mandato: %d, Es Líder: %t, ID del Líder: %d\n",
		idNodo, mandato, esLider, idLider)
}

// stopDistributedProcess detiene el nodo con identificador id
func (cfg *configDespliegue) stopDistributedProcess(id int) {
	var reply raft.Vacio

	err := cfg.nodosRaft[id].CallTimeout("NodoRaft.ParaNodo", raft.Vacio{}, &reply, 10*time.Millisecond)
	check.CheckError(err, "Error en la llamada RPC ParaNodo")
	cfg.conectados[id] = false
	fmt.Printf("Se desconecta el nodo %d\n", id)

}

// someterOperation manda la operación y devuelve el resultado obtenido de la operación
func (cfg *configDespliegue) someterOperation(op string, clave string, valor string) (int, bool, int, int, string) {
	var reply raft.ResultadoRemoto
	err := cfg.nodosRaft[cfg.lider].CallTimeout("NodoRaft.SometerOperacionNodo", raft.TipoOperacion{op, clave, valor}, &reply, 10000*time.Millisecond)
	check.CheckError(err, "Error en la llamada RPC SometerOpertacionRaft")
	cfg.lider = reply.IdLider
	return reply.IdLider, reply.EsLider, reply.Mandato, reply.IndiceRegistro, reply.ValorADevolver
}

// checkSometerOperation somete una operación y queda a la espera de la respuesta para comprobar si el valor y indice obtenidos coinciden con los esperados
func (cfg *configDespliegue) checkSometerOperation(indiceEsperado int, op string, clave string, valor string, valorEsperado string) {
	_, _, _, indiceObtenido, valorObtenido := cfg.someterOperation(op, clave, valor)

	if indiceObtenido != indiceEsperado || valorEsperado != valorObtenido {
		cfg.t.Fatalf("Operación no comprometida correctamente indiceObtenido %d, esperado %d, valorObtenido %s, esperado %s", indiceObtenido, indiceEsperado, valorObtenido, valorEsperado)
	}
	fmt.Printf("Operación comprometida correctamente indiceObtenido %d, esperado %d, valorObtenido %s, esperado %s\n", indiceObtenido, indiceEsperado, valorObtenido, valorEsperado)

}

// checkSometerOperationFail intenta comprometer una operación dada y si no logra un consenso, se da por válido, en caso contrario el comportamiento no sería el esperado
func (cfg *configDespliegue) checkSometerOperationFail(op string, clave string, valor string) {
	var reply raft.ResultadoRemoto
	err := cfg.nodosRaft[cfg.lider].CallTimeout("NodoRaft.SometerOperacionNodo", raft.TipoOperacion{op, clave, valor}, &reply, 5000*time.Millisecond)
	if err == nil {
		cfg.t.Fatalf("Se ha sometido una operación que no debería en el test %s", cfg.t.Name())
	} else {
		fmt.Printf("No se ha logrado un consenso para someter\n")
	}
}

// randomExcluding genera un número aleatorio entre 0 y 2 excluyendo el número pasado como parámetro.
func randomExcluding(exclude int) int {
	if exclude < 0 || exclude > 2 {
		panic("El número excluido debe estar entre 0 y 2")
	}
	for {
		num := rand.Intn(3) // Genera un número aleatorio entre 0 y 2
		if num != exclude {
			return num
		}
	}
}

// stopFollowers detiene todos los nodos con rol distinto a LEADER, generalmente followers
func (cfg *configDespliegue) stopFollowers() {
	var reply raft.Vacio

	for nodo, endPoint := range cfg.nodosRaft {
		if nodo != cfg.lider {
			err := endPoint.CallTimeout("NodoRaft.ParaNodo",
				raft.Vacio{}, &reply, 10*time.Millisecond)
			check.CheckError(err, "Error en llamada RPC Para nodo")
			cfg.conectados[nodo] = false
		}
	}
	fmt.Printf("Se desconectan los followers\n")
}

// checkLogger comprueba de todos los nodos tengan commit el último índice index y sean para todos igual tanto el mandatos como los últimos commits
func (cfg *configDespliegue) checkLoggers(index int) {
	indices := make([]int, cfg.numReplicas)
	mandatos := make([]int, cfg.numReplicas)
	for i, _ := range cfg.nodosRaft {
		indices[i], mandatos[i] = cfg.obtenerEstadoLogger(i)
		fmt.Printf("Nodo %d, logger cuyo último commit es el indice %d y del mandato %d \n", i, indices[i], mandatos[i])
	}

	if indices[0] != index {
		cfg.t.Fatalf("Obtenido indice no esperado en el subtest %s", cfg.t.Name())
	}
	for i := 1; i < cfg.numReplicas; i++ {
		if indices[0] != indices[i] || mandatos[0] != mandatos[i] {
			cfg.t.Fatalf("No coincide el indice o mandato en los logs de los nodos en el subtest %s", cfg.t.Name())
		}
	}
}
