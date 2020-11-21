package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
	"strconv"
)

const (
	serverport = "localhost:8000"
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

func main(){
	bIn := bufio.NewReader(os.Stdin)
	fmt.Print("Current Port: ")
	port, _ := bIn.ReadString('\n')
	hostname := fmt.Sprintf("localhost:%s", strings.TrimSpace(port))

	fmt.Println("Listening...")

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

	r := bufio.NewReader(conn)
	str, _ := r.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(str))

	var f float64 = 0
	n := len(perceptron.Weights)
	for i:=num; i<num+5;i++{
		if i < n{
			f += float64(perceptron.Data[0].Inputs[i])*perceptron.Weights[i]
		}
	}
	fmt.Println(f)
	send(f)
}

func send(f float64) {
	conn, _ := net.Dial("tcp", serverport)
	defer conn.Close()
	fmt.Fprintf(conn,"%f\n",f)
}