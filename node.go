package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"bufio"
	"strconv"
	"strings"
)

const (
	serverport = 8000
)

type Perceptron struct {
	Data []Data			`json:"data"`
	Weights []float64 	`json:"weights"`
	Umbral float64		`json:"umbral"`
}

type Data struct {
	Inputs []int	`json:"inputs"`
	Outputs int		`json:"outputs"`
}

var remotehost string

func main(){
	bIn := bufio.NewReader(os.Stdin)
	fmt.Print("Current Port: ")
	port, _ := bIn.ReadString('\n')
	hostname := fmt.Sprintf("localhost:%s", strings.TrimSpace(port))

	fmt.Print("Remote Port: ")
	remoteport, _ := bIn.ReadString('\n')
	remotehost = fmt.Sprintf("localhost:%s", strings.TrimSpace(remoteport))

	fmt.Println("Listening")

	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handlerListen(conn)
	}
}

func handlerListen(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var perceptron Perceptron
	decoder.Decode(&perceptron)
	fmt.Println(perceptron)

	r := bufio.NewReader(conn)

	str, _ := r.ReadString('\n')
	id, _ := strconv.Atoi(strings.TrimSpace(str))
	fmt.Printf("Current ID %d\n", id)

	str2, _ := r.ReadString('\n')
	id2, _ := strconv.Atoi(strings.TrimSpace(str2))
	fmt.Printf("Total Iterations %d\n", id2)

	//send(perceptron)
}

func send(p Perceptron,n int) {
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	encoder.Encode(p)

	fmt.Fprintf(conn,"%d\n",n)
}