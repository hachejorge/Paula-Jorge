// Implementacion de despliegue en ssh de multiples nodos
//
// Unica funcion exportada :
//		func ExecMutipleHosts(cmd string,
//							  hosts []string,
//							  results chan<- string,
//							  privKeyFile string)
//

package despliegue

import (
	"bufio"
	"bytes"

	//"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Fatalln("ERROR CONEXION SSH", err)
	}
	defer conn.Close()

	//fmt.Printf("APRES CONN %#v\n", config)

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalln("ERROR SESSION", err)
	}
	defer session.Close()

	//fmt.Println("APRES SESSION")

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stdoutBuf

	//fmt.Println("ANTES RUN", cmd)

	session.Run(cmd)

	//fmt.Println("TRAS RUN", cmd)

	return hostname + ": \n" + stdoutBuf.String()
}

// func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
// 	log.Printf("Conectando al host %s...", hostname)
// 	conn, err := ssh.Dial("tcp", hostname+":22", config)
// 	if err != nil {
// 		log.Printf("Error conectando al host %s: %v", hostname, err)
// 		return fmt.Sprintf("Error conectando al host %s: %v", hostname, err)
// 	}
// 	defer conn.Close()

// 	log.Printf("Conexión establecida con host %s. Creando sesión...", hostname)
// 	session, err := conn.NewSession()
// 	if err != nil {
// 		log.Printf("Error creando sesión SSH para host %s: %v", hostname, err)
// 		return fmt.Sprintf("Error creando sesión SSH para host %s: %v", hostname, err)
// 	}
// 	defer session.Close()

// 	var stdoutBuf bytes.Buffer
// 	session.Stdout = &stdoutBuf
// 	session.Stderr = &stdoutBuf

// 	log.Printf("Ejecutando comando en host %s: %s", hostname, cmd)
// 	err = session.Run(cmd)
// 	if err != nil {
// 		log.Printf("Error ejecutando comando en host %s: %v", hostname, err)
// 		return fmt.Sprintf("Error ejecutando comando en host %s: %v\nOutput: %s", hostname, err, stdoutBuf.String())
// 	}

// 	log.Printf("Comando ejecutado correctamente en host %s", hostname)
// 	return fmt.Sprintf("%s: \n%s", hostname, stdoutBuf.String())
// }

func buildSSHConfig(signer ssh.Signer) *ssh.ClientConfig {

	return &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		// verify host public key
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}
}

func execOneHost(hostname string, results chan<- string,
	signer ssh.Signer, cmd string) {
	// get host public key
	// ssh_config must have option "HashKnownHosts no" !!!!
	//hostKey := getHostKey(hostname)
	//config := buildSSHConfig(signer, hostKey)
	config := buildSSHConfig(signer)

	//fmt.Println(cmd)
	//fmt.Println(hostname)

	results <- executeCmd(cmd, hostname, config)
}

// func execOneHost(hostname string, results chan<- string, signer ssh.Signer, cmd string) string {
// 	config := buildSSHConfig(signer)

// 	result := executeCmd(cmd, hostname, config)
// 	if results != nil {
// 		results <- result
// 	}
// 	return result
// }

// Ejecutar un mismo comando en múltiples hosts mediante ssh
func ExecMutipleHosts(cmd string,
	hosts []string,
	results chan<- string,
	privKeyFile string) {

	//results := make(chan string, 1000)

	//Read private key file for user
	pkey, err := ioutil.ReadFile(
		filepath.Join(os.Getenv("HOME"), ".ssh", privKeyFile))

	//fmt.Println("PrivKey: ", string(pkey))

	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(pkey)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	for _, hostname := range hosts {
		go execOneHost(hostname, results, signer, cmd)
	}
}

// func ExecMutipleHosts(cmd string, hosts []string, results chan<- string, privKeyFile string) {
// 	// Leer clave privada
// 	pkeyPath := filepath.Join(os.Getenv("HOME"), ".ssh", privKeyFile)
// 	pkey, err := ioutil.ReadFile(pkeyPath)
// 	if err != nil {
// 		log.Fatalf("Error leyendo clave privada (%s): %v", pkeyPath, err)
// 	}

// 	// Crear el signer para la clave privada
// 	signer, err := ssh.ParsePrivateKey(pkey)
// 	if err != nil {
// 		log.Fatalf("Error parseando clave privada: %v", err)
// 	}

// 	var wg sync.WaitGroup

// 	// Ejecutar comandos en paralelo
// 	for _, hostname := range hosts {
// 		wg.Add(1)
// 		go func(host string) {
// 			defer wg.Done()
// 			log.Printf("Ejecutando comando en host: %s", host)
// 			result := execOneHost(host, results, signer, cmd)
// 			log.Printf("Resultado de host %s:\n%s", host, result)
// 		}(hostname)
// 	}

// 	// Esperar a que todas las goroutines terminen
// 	wg.Wait()
// 	log.Println("Todos los comandos han sido ejecutados.")
// }
