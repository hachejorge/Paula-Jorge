// Escribir vuestro código de funcionalidad Raft en este fichero
//

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
	kLogOutputDir = "./logs_raft/"
)

type TipoOperacion struct {
	Operacion string // La operaciones posibles son "leer" y "escribir"
	Clave     string
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
	IdLider int
	// Utilización opcional de este logger para depuración
	// Cada nodo Raft tiene su propio registro de trazas (logs)
	Logger *log.Logger

	// Vuestros datos aqui.
	// Canales para comunicarse
	Leader    chan bool
	Follower  chan bool
	Heartbeat chan bool

	CurrentTerm int // Periodo actual
	VotedFor    int // A quién he votado

	VotesReceived int    // Número de votos recibidos
	Rol           string // Rol del nodo, ("LEADER", "FOLLOWER", "CANDIDATE")

	CommitIndex int // Índice de la última entrada comprometida (commited)
	LastAplied  int // Índice de la última entrada aplicada ??

	NextIndex  []int // Índice para cada servidor del log que tenemos que enviarle
	MatchIndex []int // Índice del log máximo que este replicado en esa máquina
	// Nodo 0     Lider    [1,1,2]  -> MatchIndex [x,2,1]
	// Nodo 1	  Follower [1,1]
	// Nodo 2     Follower [1]

	Logs []Entry

	// mirar figura 2 para descripción del estado que debe mantenre un nodo Raft
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

	nr.Follower = make(chan bool)
	nr.Leader = make(chan bool)
	nr.Heartbeat = make(chan bool)

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
	mandato := -1
	EsLider := nr.Yo == nr.IdLider
	idLider := nr.IdLider
	valorADevolver := ""

	if EsLider {
		indice = len(nr.Logs)
		mandato = nr.CurrentTerm
		entry := Entry{indice, mandato, operacion}
		nr.Logs = append(nr.Logs, entry)
		nr.Mux.Unlock()
		// Esperar a recibir confirmación de commit
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

type ResultadoRemoto struct {
	ValorADevolver string
	IndiceRegistro int
	EstadoParcial
}

func (nr *NodoRaft) SometerOperacionRaft(operacion TipoOperacion, reply *ResultadoRemoto) error {
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
	// Si el mantado de la última entrada del log del nodo soliciante es mayor que el nuestro
	// se concede el voto
	if nr.Logs[len(nr.Logs)-1].Term < request.LastLogTerm {
		isBetter = true

		// Si los mandatos son iguales se compara según la cantidad de logs almacenados
	} else if nr.Logs[len(nr.Logs)-1].Term == request.LastLogTerm &&
		request.LastLogIndex >= len(nr.Logs)-1 {

		isBetter = true
	}

	return isBetter
}

// Metodo para RPC PedirVoto
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) error {

	if peticion.CandidateTerm > nr.CurrentTerm { // Si me llega un mandato mayor al mío le doy el voto

		if len(nr.Logs) == 0 || isBetterLeader(nr, peticion) {
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
			// nr.CurrentTerm = peticion.CandidateTerm no es necesario ??
			reply.Term = nr.CurrentTerm
			reply.VoteGranted = false
		}

	} else { // Si llega un mandato menor al mío no le doy el voto
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

	entries []Entry // Vector para la eficiencia ?

	LeaderCommit int
}

type Results struct {
	Term    int
	Success bool
}

// Metodo de tratamiento de llamadas RPC AppendEntries
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries,
	results *Results) error {
	if args.Term < nr.CurrentTerm { // Heartbeat atrasado
		results.Term = nr.CurrentTerm
		results.Success = false
	} else if args.Term == nr.CurrentTerm { // Recibo heartbeat en mi mismo periodo
		nr.IdLider = args.LeaderID
		results.Term = args.Term
		nr.Heartbeat <- true
	} else { // Recibo heartbeat de un periodo futuro
		nr.IdLider = args.LeaderID
		nr.CurrentTerm = args.Term
		results.Term = args.Term
		if nr.Rol == "FOLLOWER" {
			nr.Heartbeat <- true
		} else { // Si es líder o candidato vuelve a ser follower
			nr.Follower <- true
		}
	}

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
				nr.Leader <- true
			}
		}
		return true
	}
}

func (nr *NodoRaft) enviarLatido(nodo int, args *ArgAppendEntries,
	results *Results) bool {

	err := nr.Nodos[nodo].CallTimeout("NodoRaft.AppendEntries", args, results, 20*time.Millisecond)
	if err != nil {
		return false
	} else {
		if results.Term > nr.CurrentTerm {
			//Si he enviado heartbeat a un nodo con mayor mandato dejo de ser
			//líder, actualizo mi mandato y vuelvo a ser follower
			nr.CurrentTerm = results.Term
			nr.IdLider = -1
			nr.Follower <- true
		}
		return true
	}
}

func requestVotes(nr *NodoRaft) {
	var reply RespuestaPeticionVoto
	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			go nr.enviarPeticionVoto(i, &ArgsPeticionVoto{nr.CurrentTerm, nr.Yo}, &reply)
		}
	}
}

func sendHeartbeats(nr *NodoRaft) {
	var reply Results
	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			go nr.enviarLatido(i, &ArgAppendEntries{nr.CurrentTerm, nr.Yo}, &reply)
		}
	}
}

func getRandomTimeout() time.Duration {
	// Genera un timeout aleatorio entre 150 y 300 ms
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

// Funcion para gestionar el comportamiento de un nodo en el algoritmo de consenso de raft
func raftHandler(nr *NodoRaft) {
	//var timerFollower time.Timer
	//var timerLeader time.Timer
	//timerLeader := time.NewTimer(50 * time.Millisecond)

	// Descomentar para la primera prueba del test 1
	// time.Sleep(6000 * time.Millisecond)

	for {
		if nr.Rol == "FOLLOWER" {
			timerFollower := time.NewTimer(getRandomTimeout())
			select {
			case <-nr.Heartbeat: // Recibe el heartbeat
			// Sigo como follower
			case <-timerFollower.C: // Expira timeout
				nr.IdLider = -1
				nr.Rol = "CANDIDATE"
			}
		} else if nr.Rol == "LEADER" {
			nr.IdLider = nr.Yo
			// Enviar heartbeats
			sendHeartbeats(nr)
			timerLeader := time.NewTimer(50 * time.Millisecond)
			select {
			case <-nr.Follower: // Descubre mandato mayor
				nr.Rol = "Follower"
			case <-timerLeader.C: // Expira el time out
				// Sigo como leader
				// Vuelvo a mandar HeartBeats
			}
		} else { // nr.Rol == "CANDIDATE"
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
