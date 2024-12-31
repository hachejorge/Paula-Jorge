// Módulo de implementación del algoritmo de raft
// Autores: Jorge Hernández 872838 y Paula Blaco 874320

package raft

//
// API
// ===
// Este es el API que vuestra implementación debe exportar
//
// nodoRaft = NuevoNodo(...)
//   Crear un nuevo servidor del grupo de elección.
//
// nodoRaft.Para()
//   Solicitar la parado de un servidor
//
// nodo.ObtenerEstado() (yo, mandato, esLider)
//   Solicitar a un nodo de elección por "yo", su mandato en curso,
//   y si piensa que es el msmo el lider
//
// nodoRaft.SometerOperacion(operacion interface()) (indice, mandato, esLider)

// type AplicaOperacion

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	//"crypto/rand"
	"sync"
	"time"

	//"net/rpc"

	"raft/internal/comun/rpctimeout"
)

const (
	// Constante para fijar valor entero no inicializado
	IntNOINICIALIZADO = -1

	//  false deshabilita por completo los logs de depuracion
	// Aseguraros de poner kEnableDebugLogs a false antes de la entrega
	kEnableDebugLogs = true

	// Poner a true para logear a stdout en lugar de a fichero
	kLogToStdout = false

	// Cambiar esto para salida de logs en un directorio diferente
	kLogOutputDir = "../../logs_raft/"
)

type TipoOperacion struct {
	Operacion string // La operaciones posibles son "leer" y "escribir"
	Clave     string // Clave en donde se escribirá o leerá el valor
	Valor     string // en el caso de la lectura Valor = ""
}

// A medida que el nodo Raft conoce las operaciones de las  entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados
type AplicaOperacion struct {
	Indice    int // en la entrada de registro
	Operacion TipoOperacion
}

// Tipo de dato Go que representa un solo nodo (réplica) de raft
type NodoRaft struct {
	Mux sync.Mutex // Mutex para proteger acceso a estado compartido

	// Host:Port de todos los nodos (réplicas) Raft, en mismo orden
	Nodos   []rpctimeout.HostPort
	Yo      int // indice de este nodos en campo array "nodos"
	IdLider int // Identificador del líder de los nodos, -1 si es desconocido

	// Utilización opcional de este logger para depuración
	// Cada nodo Raft tiene su propio registro de trazas (logs)
	Logger *log.Logger

	// Canales para comunicarse entre funciones
	Leader    chan bool // Indica que el nodo se convierte en líder
	Follower  chan bool // Indica que el nodo se convierte en seguidor
	Heartbeat chan bool // Avisa sobre la recepción del heartbeat por parte del líder

	ApplyOp    chan AplicaOperacion // Canal para comunicar la aplicación de las operaciones
	OpCommited chan string          // Canal para comunicar el valor obtenido al ejecutar la operación

	CurrentTerm int // Periodo actual
	VotedFor    int // A quién he votado

	VotesReceived int    // Número de votos recibidos
	Rol           string // Rol del nodo, ("LEADER", "FOLLOWER", "CANDIDATE")

	CommitIndex int // Índice de la última entrada comprometida (commited)
	LastAplied  int // Índice de la última entrada aplicada

	NextIndex  []int // Índice para cada servidor del log que tenemos que enviarle
	MatchIndex []int // Índice del log máximo que este replicado en esa máquina
	// Nodo 0     Lider    [1,1,2]  -> MatchIndex [x,2,1]
	// Nodo 1	  Follower [1,1]
	// Nodo 2     Follower [1]

	Logs []Entry // Vector con todas las entries cometidas

	ACKsCommit int // contador de respuestas recibidas para cometer una operación
}

type Entry struct {
	Index int
	Term  int
	Op    TipoOperacion
}

// Creacion de un nuevo nodo de eleccion
//
// Tabla de <Direccion IP:puerto> de cada nodo incluido a si mismo.
//
// <Direccion IP:puerto> de este nodo esta en nodos[yo]
//
// Todos los arrays nodos[] de los nodos tienen el mismo orden

// canalAplicar es un canal donde, en la practica 5, se recogerán las
// operaciones a aplicar a la máquina de estados. Se puede asumir que
// este canal se consumira de forma continúa.
//
// NuevoNodo() debe devolver resultado rápido, por lo que se deberían
// poner en marcha Gorutinas para trabajos de larga duracion
func NuevoNodo(nodos []rpctimeout.HostPort, yo int,
	canalAplicarOperacion chan AplicaOperacion) *NodoRaft {

	nr := &NodoRaft{}
	nr.Nodos = nodos
	nr.Yo = yo
	nr.IdLider = IntNOINICIALIZADO
	nr.Rol = "FOLLOWER"
	nr.CurrentTerm = 0
	nr.VotedFor = -1
	nr.VotesReceived = 0
	nr.ACKsCommit = 0
	nr.CommitIndex = -1
	nr.LastAplied = -1

	nr.NextIndex = make([]int, len(nr.Nodos))
	nr.MatchIndex = make([]int, len(nr.Nodos))

	// Inicialmente el siguiente índice a mandar es 0 y coincidimos en el -1
	for i := range nr.Nodos {
		nr.NextIndex[i] = 0
		nr.MatchIndex[i] = -1
	}

	nr.Logs = []Entry{}

	nr.ApplyOp = canalAplicarOperacion
	nr.Follower = make(chan bool)
	nr.Leader = make(chan bool)
	nr.Heartbeat = make(chan bool)
	nr.OpCommited = make(chan string)

	if kEnableDebugLogs {
		nombreNodo := nodos[yo].Host() + "_" + nodos[yo].Port()
		logPrefix := fmt.Sprintf("%s", nombreNodo)

		fmt.Println("LogPrefix: ", logPrefix)

		if kLogToStdout {
			nr.Logger = log.New(os.Stdout, nombreNodo+" -->> ",
				log.Lmicroseconds|log.Lshortfile)
		} else {
			err := os.MkdirAll(kLogOutputDir, os.ModePerm)
			if err != nil {
				panic(err.Error())
			}
			logOutputFile, err := os.OpenFile(fmt.Sprintf("%s/%s.txt",
				kLogOutputDir, logPrefix), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				panic(err.Error())
			}
			nr.Logger = log.New(logOutputFile,
				logPrefix+" -> ", log.Lmicroseconds|log.Lshortfile)
		}
		nr.Logger.Println("logger initialized")
	} else {
		nr.Logger = log.New(ioutil.Discard, "", 0)
	}

	// Añadir codigo de inicialización
	go raftHandler(nr)

	return nr
}

// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
func (nr *NodoRaft) para() {
	go func() { time.Sleep(5 * time.Millisecond); os.Exit(0) }()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este  nodo Raft el el conjunto de nodos
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nr.Yo
	var mandato int = nr.CurrentTerm
	var esLider bool = nr.Yo == nr.IdLider
	var idLider int = nr.IdLider

	return yo, mandato, esLider, idLider
}

// Devuelve el índice del último commit y el mandato de ese mismo
func (nr *NodoRaft) obtenerEstadoLogger() (int, int) {
	index := -1
	term := 0
	if len(nr.Logs) != 0 {
		index = nr.CommitIndex
		term = nr.Logs[index].Term
	}
	return index, term
}

// El servicio que utilice Raft (base de datos clave/valor, por ejemplo)
// Quiere buscar un acuerdo de posicion en registro para siguiente operacion
// solicitada por cliente.

// Si el nodo no es el lider, devolver falso
// Sino, comenzar la operacion de consenso sobre la operacion y devolver en
// cuanto se consiga
//
// No hay garantia que esta operacion consiga comprometerse en una entrada de
// de registro, dado que el lider puede fallar y la entrada ser reemplazada
// en el futuro.
// Primer valor devuelto es el indice del registro donde se va a colocar
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) someterOperacion(operacion TipoOperacion) (int, int,
	bool, int, string) {

	nr.Mux.Lock()
	indice := -1
	mandato := nr.CurrentTerm
	EsLider := nr.Yo == nr.IdLider
	idLider := nr.IdLider
	valorADevolver := ""

	// si es líder se intenta someter la operación
	if EsLider {
		indice = len(nr.Logs)
		mandato = nr.CurrentTerm
		entry := Entry{indice, mandato, operacion}
		nr.Logs = append(nr.Logs, entry)
		nr.Mux.Unlock()
		// Esperar a recibir confirmación de commit
		op := <-nr.OpCommited
		indice = nr.LastAplied
		valorADevolver = op
		nr.Logger.Printf("someterOperacion : Operación comprometida")

	} else {
		nr.Mux.Unlock()
	}

	return indice, mandato, EsLider, idLider, valorADevolver
}

// -----------------------------------------------------------------------
// LLAMADAS RPC al API
//
// Si no tenemos argumentos o respuesta estructura vacia (tamaño cero)
type Vacio struct{}

func (nr *NodoRaft) ParaNodo(args Vacio, reply *Vacio) error {
	defer nr.para()
	return nil
}

type EstadoParcial struct {
	Mandato int
	EsLider bool
	IdLider int
}

type EstadoRemoto struct {
	IdNodo int
	EstadoParcial
}

func (nr *NodoRaft) ObtenerEstadoNodo(args Vacio, reply *EstadoRemoto) error {
	reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider = nr.obtenerEstado()
	return nil
}

type EstadoLogger struct {
	Indice  int
	Mandato int
}

func (nr *NodoRaft) ObtenerEstadoLogger(args Vacio, reply *EstadoLogger) error {
	reply.Indice, reply.Mandato = nr.obtenerEstadoLogger()
	return nil
}

type ResultadoRemoto struct {
	ValorADevolver string
	IndiceRegistro int
	EstadoParcial
}

func (nr *NodoRaft) SometerOperacionNodo(operacion TipoOperacion, reply *ResultadoRemoto) error {
	reply.IndiceRegistro, reply.Mandato, reply.EsLider,
		reply.IdLider, reply.ValorADevolver = nr.someterOperacion(operacion)

	return nil
}

// -----------------------------------------------------------------------
// LLAMADAS RPC protocolo RAFT
//
// Structura de ejemplo de argumentos de RPC PedirVoto.
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
type ArgsPeticionVoto struct {
	CandidateTerm int
	CandidateID   int
	LastLogIndex  int
	LastLogTerm   int
}

// Structura de ejemplo de respuesta de RPC PedirVoto,
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
type RespuestaPeticionVoto struct {
	Term        int
	VoteGranted bool
}

// Función que devuelve true si el nodo solicitante tiene logs más avanzados que el nodo nr
func isBetterLeader(nr *NodoRaft, request *ArgsPeticionVoto) bool {
	isBetter := false

	if len(nr.Logs) > 0 {
		// Si el mantado de la última entrada del log del nodo soliciante es mayor que el nuestro
		// se concede el voto
		if nr.Logs[len(nr.Logs)-1].Term < request.LastLogTerm {
			isBetter = true

			// Si los mandatos son iguales se compara según la cantidad de logs almacenados
		} else if nr.Logs[len(nr.Logs)-1].Term == request.LastLogTerm &&
			request.LastLogIndex >= len(nr.Logs)-1 {

			isBetter = true
		}
	}
	return isBetter
}

// Metodo para RPC PedirVoto
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) error {
	nr.Logger.Printf("PedirVoto: Recibido solicitud de %d para mandato %d (mi mandato: %d)",
		peticion.CandidateID, peticion.CandidateTerm, nr.CurrentTerm)

	if peticion.CandidateTerm > nr.CurrentTerm { // Si me llega un mandato mayor al mío puedo darle el voto el voto
		nr.Logger.Printf("PedirVoto: Actualizando mandato a %d y votando por %d",
			peticion.CandidateTerm, peticion.CandidateID)
		if len(nr.Logs) == 0 || isBetterLeader(nr, peticion) { // Si es mejor líder le doy mi voto
			reply.Term = peticion.CandidateTerm
			reply.VoteGranted = true
			nr.CurrentTerm = peticion.CandidateTerm
			nr.VotedFor = peticion.CandidateID
			if nr.Rol == "LEADER" || nr.Rol == "CANDIDATE" {
				// Vuelvo a ser follower
				nr.Follower <- true
			}
			// No le doy el voto
		} else {
			reply.Term = nr.CurrentTerm
			reply.VoteGranted = false
		}

	} else { // Si llega un mandato menor al mío no le doy el voto
		nr.Logger.Printf("PedirVoto: Rechazando solicitud de %d (mandato solicitado: %d, mi mandato: %d)",
			peticion.CandidateID, peticion.CandidateTerm, nr.CurrentTerm)
		reply.Term = nr.CurrentTerm
		reply.VoteGranted = false
	}

	return nil
}

type ArgAppendEntries struct {
	Term     int
	LeaderID int
	// Más en P4
	PrevLogIndex int
	PrevLogTerm  int

	Entries Entry // Vector para la eficiencia ?

	LeaderCommit int
}

type Results struct {
	Term    int
	Success bool
}

// isLoggerUpdated te devuelve true si el índice y mandato del logger coincide con el solicitado
func isLoggerUpdated(nr *NodoRaft, request *ArgAppendEntries) bool {
	// Compruebo si mi logger coincide con el dado
	if len(nr.Logs)-1 < request.PrevLogIndex { // Log desactualizado
		return false
	} else if request.PrevLogIndex < len(nr.Logs) && nr.Logs[request.PrevLogIndex].Term < request.PrevLogTerm { // Log desactualizado
		return false
	} else {
		return true
	}
}

// min devuelve a si a < b, en caso contrario b
func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// Metodo de tratamiento de llamadas RPC AppendEntries
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries,
	results *Results) error {
	nr.Mux.Lock()

	nr.Logger.Printf("AppendEntries: Recibido de líder %d en mandato %d (mi mandato: %d)",
		args.LeaderID, args.Term, nr.CurrentTerm)

	if args.Term < nr.CurrentTerm { // Heartbeat atrasado
		results.Term = nr.CurrentTerm
		results.Success = false
	} else if args.Term >= nr.CurrentTerm { // Recibo heartbeat en mi mismo periodo o mayor
		nr.IdLider = args.LeaderID
		results.Term = args.Term

		// Heartbeat simple
		if args.Entries == (Entry{}) {
			results.Success = false

			// Recibido entry
		} else {
			nr.Logger.Printf("AppendEntries: Recibida entrada: %+v", args.Entries)
			// Logger vacío
			if len(nr.Logs) == 0 {
				nr.Logger.Printf("AppendEntries: Logger previo vacío y recibo el primer log")
				nr.Logs = append(nr.Logs, args.Entries)
				results.Success = true

				// Si no tengo mi logger vacío
			} else {
				results.Success = isLoggerUpdated(nr, args)
				// Mi logger coincide para PrevLogIndex en mandato
				if results.Success {
					nr.Logger.Printf("AppendEntries: Logger correcto, añado nueva entrada")
					// Me quedo con los logs desde el inicio hasta el indice actualizado
					nr.Logs = nr.Logs[0 : args.PrevLogIndex+1]
					// Añado la entrada
					nr.Logs = append(nr.Logs, args.Entries)

					nr.Logger.Printf("AppendEntries: Éxito al añadir entrada de líder %d", args.LeaderID)
				} else {
					nr.Logger.Printf("AppendEntries: Fallo logger desactualizado, error al añadir nueva entrada")
				}
			}
		}
		// Actualizo el último commit index si es necesario
		if args.LeaderCommit > nr.CommitIndex {
			nr.Logger.Printf("AppendEntries: Actualizo mi CommitIndex a %d, previamente era %d", args.LeaderCommit, nr.CommitIndex)
			nr.CommitIndex = min(args.LeaderCommit, len(nr.Logs)-1)
		}

		// Recibo heartbeat de un periodo futuro
		if args.Term > nr.CurrentTerm {
			nr.CurrentTerm = args.Term
			if nr.Rol == "LEADER" || nr.Rol == "CANDIDATE" { // Si es líder o candidato vuelve a ser follower
				nr.Follower <- true
			} else {
				nr.Heartbeat <- true
			}
		} else {
			nr.Heartbeat <- true
		}
	}
	nr.Mux.Unlock()
	return nil
}

// ----- Metodos/Funciones a utilizar como clientes
//
//

// Ejemplo de código enviarPeticionVoto
//
// nodo int -- indice del servidor destino en nr.nodos[]
//
// args *RequestVoteArgs -- argumentos para la llamada RPC
//
// reply *RequestVoteReply -- respuesta RPC
//
// Los tipos de argumentos y respuesta pasados a CallTimeout deben ser
// los mismos que los argumentos declarados en el metodo de tratamiento
// de la llamada (incluido si son punteros
//
// Si en la llamada RPC, la respuesta llega en un intervalo de tiempo,
// la funcion devuelve true, sino devuelve false
//
// la llamada RPC deberia tener un timout adecuado.
//
// Un resultado falso podria ser causado por una replica caida,
// un servidor vivo que no es alcanzable (por problemas de red ?),
// una petición perdida, o una respuesta perdida
//
// Para problemas con funcionamiento de RPC, comprobar que la primera letra
// del nombre  todo los campos de la estructura (y sus subestructuras)
// pasadas como parametros en las llamadas RPC es una mayuscula,
// Y que la estructura de recuperacion de resultado sea un puntero a estructura
// y no la estructura misma.
func (nr *NodoRaft) enviarPeticionVoto(nodo int, args *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) bool {

	err := nr.Nodos[nodo].CallTimeout("NodoRaft.PedirVoto", args, reply, 20*time.Millisecond)

	if err != nil {
		return false
	} else {
		if reply.Term > nr.CurrentTerm {
			//Si pido el voto a un nodo con mayor mandato, dejo de ser
			//candidato y vuelvo a ser follower
			nr.CurrentTerm = reply.Term
			nr.Follower <- true

		} else if reply.VoteGranted {
			//Si me dan el voto compruebo si tengo mayoría simple, en cuyo caso
			//me convierto en líder
			nr.Mux.Lock()
			nr.VotesReceived++
			nr.Mux.Unlock()
			if nr.VotesReceived > len(nr.Nodos)/2 {
				nr.Logger.Printf("Nodo %d ha ganado la elección y es ahora el líder para mandato %d", nr.Yo, nr.CurrentTerm)

				nr.Leader <- true
			}
		}
		return true
	}
}

// sendAppendEntry realiza una llamada RPC a AppendEntries a un nodo concreto y actualiza el estado del nodo emisor en función de la respuesta
func (nr *NodoRaft) sendAppendEntry(nodo int, args *ArgAppendEntries,
	results *Results) bool {

	err := nr.Nodos[nodo].CallTimeout("NodoRaft.AppendEntries", args, results, 20*time.Millisecond)
	if err != nil {
		return false
	} else {
		// PARTE DEL HEARTBEAT
		if results.Term > nr.CurrentTerm {
			//Si he enviado heartbeat a un nodo con mayor mandato dejo de ser
			//líder, actualizo mi mandato y vuelvo a ser follower
			nr.Mux.Lock()
			nr.CurrentTerm = results.Term
			nr.IdLider = -1
			nr.Follower <- true
			nr.Mux.Unlock()
			nr.Logger.Printf("sendAppendEntry: Heartbeat enviado a un nodo con mandato superior")
		}
		// PARTE DE LOS LOGS
		if results.Success { // Si se recibe confirmación del logger
			nr.Logger.Printf("sendAppendEntry: Coincidimos en logs hasta el indice %d con el nodo %d", nr.NextIndex[nodo], nodo)

			nr.MatchIndex[nodo] = nr.NextIndex[nodo]
			nr.NextIndex[nodo]++
			nr.Mux.Lock()

			if nr.MatchIndex[nodo] > nr.CommitIndex {
				nr.ACKsCommit++
				if nr.ACKsCommit >= len(nr.Nodos)/2 {
					nr.CommitIndex++
					nr.ACKsCommit = 0
					nr.Logger.Printf("sendAppendEntry: Recibidas ACK's necesario para confirmar el índice %d", nr.CommitIndex)
				}
			}
			nr.Mux.Unlock()
			// El logger del nodo solicitado es inconsistente se intenta con la Entry previa
		} else {
			if args.Entries != (Entry{}) {
				nr.NextIndex[nodo]--
			}
		}
		return true
	}
}

// requestVotes manda la solicitud de voto al resto de nodos distintos al propio
func requestVotes(nr *NodoRaft) {
	var reply RespuestaPeticionVoto
	nr.Logger.Printf("requestVotes: Nodo %d comienza solicitud de votos para mandato %d", nr.Yo, nr.CurrentTerm)

	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			nr.Logger.Printf("requestVotes: Enviando solicitud de voto a nodo %d", i)
			var lastLogIndex int
			var lastLogTerm int
			if len(nr.Logs) == 0 { // Si todavía no tiene Entries en el logger no manda ninguno
				lastLogIndex = -1
				lastLogTerm = 0
			} else { // Manda el tamaño del logger(último indice) y el mandato de ese último índice
				lastLogIndex = len(nr.Logs) - 1
				lastLogTerm = nr.Logs[lastLogIndex].Term
			}
			go nr.enviarPeticionVoto(i, &ArgsPeticionVoto{nr.CurrentTerm, nr.Yo, lastLogIndex, lastLogTerm}, &reply)
		}
	}
}

// sendAppendEntries obtiene los argumentos de entrada de envío para cada nodo (ya sea heartbeat simple o entries) y realiza el envío mediante sendAppendEntry
func sendAppendEntries(nr *NodoRaft) {
	nr.Logger.Printf("sendAppendEntries: Nodo %d (líder) enviando entradas a seguidores", nr.Yo)

	for i := 0; i < len(nr.Nodos); i++ {
		// Evita enviar a sí mismo
		if i == nr.Yo {
			continue
		}

		nr.Logger.Printf("sendAppendEntries: Preparando envío a nodo %d", i)

		// Determina el índice de la siguiente entrada a enviar
		nextIndex := nr.NextIndex[i]
		var entry Entry

		// Si nextIndex está dentro del rango del log, selecciona las entradas
		if nextIndex < len(nr.Logs) {
			entry = nr.Logs[nextIndex]
			nr.Logger.Printf("sendAppendEntries: Enviando la siguiente entry al nodo %d, %+v", i, entry)

		} else {
			nr.Logger.Printf("sendAppendEntries: Índice nextIndex %d fuera de rango, enviando heartbeat", nextIndex)
			entry = Entry{} // Entrada vacía para indicar heartbeat
		}

		// Determina el índice y el término del log previo
		prevLogIndex := nextIndex - 1
		prevLogTerm := 0
		if prevLogIndex >= 0 && prevLogIndex < len(nr.Logs) {
			prevLogTerm = nr.Logs[prevLogIndex].Term
		}

		// Construye los argumentos para AppendEntries
		args := ArgAppendEntries{
			Term:         nr.CurrentTerm,
			LeaderID:     nr.Yo,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entries:      entry,
			LeaderCommit: nr.CommitIndex,
		}

		// Log de depuración
		nr.Logger.Printf("sendAppendEntries: Nodo %d enviando AppendEntries a nodo %d (PrevLogIndex=%d, PrevLogTerm=%d, LeaderCommit=%d)",
			nr.Yo, i, prevLogIndex, prevLogTerm, nr.CommitIndex)

		// Resultado para recibir la respuesta del nodo seguidor
		var reply Results

		// Envía AppendEntries en una goroutine
		go nr.sendAppendEntry(i, &args, &reply)
	}
}

// Genera un timeout aleatorio entre 150 y 300 ms
func getRandomTimeout() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

// Funcion para gestionar el comportamiento de un nodo en el algoritmo de consenso de raft
func raftHandler(nr *NodoRaft) {
	// Descomentar para la primera prueba del test 1
	//time.Sleep(2000 * time.Millisecond)

	for {
		if nr.CommitIndex > nr.LastAplied {
			nr.LastAplied++

			op := AplicaOperacion{nr.LastAplied, nr.Logs[nr.LastAplied].Op}
			nr.ApplyOp <- op
			op = <-nr.ApplyOp
			nr.OpCommited <- op.Operacion.Valor
			nr.Logger.Printf("raftHandler: Soy %s y añado a mi máquina de estados la operación %+v, con índice %d", nr.Rol, nr.Logs[nr.LastAplied].Op, nr.LastAplied)

		}

		for nr.Rol == "FOLLOWER" {
			timerFollower := time.NewTimer(getRandomTimeout())
			select {
			case <-nr.Heartbeat: // Recibe el heartbeat
			// Sigo como follower
			case <-timerFollower.C: // Expira timeout
				nr.IdLider = -1
				nr.Rol = "CANDIDATE"
			}
		}
		for nr.Rol == "LEADER" {
			nr.IdLider = nr.Yo
			// Enviar heartbeats
			sendAppendEntries(nr)
			timerLeader := time.NewTimer(50 * time.Millisecond)
			select {
			case <-nr.Follower: // Descubre mandato mayor
				nr.Rol = "Follower"
			case <-timerLeader.C: // Expira el time out
				// Sigo como leader y compruebo si he actualizado el índice
				if nr.CommitIndex > nr.LastAplied {
					nr.LastAplied++
					op := AplicaOperacion{nr.LastAplied, nr.Logs[nr.LastAplied].Op}
					nr.ApplyOp <- op
					op = <-nr.ApplyOp
					nr.OpCommited <- op.Operacion.Valor
					nr.Logger.Printf("raftHandler: Soy LEADER y añado a mi máquina de estados la operación %+v, con índice %d", nr.Logs[nr.LastAplied].Op, nr.LastAplied)
				}
			}
		}
		for nr.Rol == "CANDIDATE" {
			nr.CurrentTerm++
			nr.VotedFor = nr.Yo
			nr.VotesReceived = 1
			requestVotes(nr)
			// Tarda 2.5 segundos hasta iniciar una nueva elección
			timerCandidate := time.NewTimer(2500 * time.Millisecond)
			select {
			case <-nr.Leader: // Se convierte en líder
				nr.Rol = "LEADER"
			case <-nr.Follower: // Se convierte en follower
				nr.Rol = "FOLLOWER"
			case <-timerCandidate.C: // Se acaba el timeout
				nr.Rol = "CANDIDATE" // Se vuelve a presentar como candidato
			}
		}
	}
}
