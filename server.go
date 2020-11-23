package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"bufio"
	"net"
	"fmt"
	"os"
)

func Send(input []float64){
	for _,node := range(nodes){
		conn, _ := net.Dial("tcp",node)
		defer conn.Close()	
		encoder := json.NewEncoder(conn)
		encoder.Encode(input)
	}
}

func Listen(){
	ln, _ := net.Listen("tcp", "localhost:8000")
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	port, _ := r.ReadString('\n')
	result,_ := r.ReadString('\n')
	fmt.Println("__________\n")
	fmt.Print("El puerto ",port)
	fmt.Print("dice:",result)
}

var nodes []string
var remotehost string

func main(){
	bIn := bufio.NewReader(os.Stdin)
	fmt.Print("Cantidad de Nodos: ")
	str,_ := bIn.ReadString('\n')
	n,_ := strconv.Atoi(strings.TrimSpace(str))

	for i:=0;i<n;i++{
		fmt.Print("Puerto del Nodo ",i+1,":")
		port,_ := bIn.ReadString('\n')
		port = strings.TrimSpace(port)
		remotehost = fmt.Sprintf("localhost:%s", port)
		nodes = append(nodes,remotehost)
	}

	input := []float64{'x','s','n','t','p','f','c','n','k','e','e','s','s','w','w','p','w','o','p','k','s','u'}
	Send(input)
	Listen()
}