package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"bufio"
	"net"
	"fmt"
	"os"
)

var nodes []string			// __ nodos enlazados al servidor
var serverhost string		// __ host del servidor actual
var data []Data 			// __ almacena toda la data recibida de post
var input []float64			// __ último input ingresado

var countnodes int = 0		// __ auxiliar que valida la llegada de las salidas de todos los nodos
var countedible int = 0		// __ auxiliar para contar salidas comestibles de los nodos
var countpoisonous int = 0	// __ auxiliar para contar salidas venenosas de los nodos
var resultado string        // __ Mensjae de salida del algoritmo 

////////////////////////////////////////////////////////////////////////////////////////////////////

type Data struct {
	Inputs []float64	`json:"inputs"`
	Output float64		`json:"output"`
}

type Json struct {
	Inputs []string	`json:"inputs"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func get(w http.ResponseWriter, r *http.Request) {
	jsonbytes,err := json.Marshal(data)
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonbytes)
}

func post(w http.ResponseWriter, r *http.Request) {
	bodyBytes,err:=ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	var j Json
	err = json.Unmarshal(bodyBytes,&j)
	if err != nil { 
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	var in []float64
	for i:= range(j.Inputs){
		in = append(in,float64(j.Inputs[i][0]))
	}
	Send(in)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Send(in []float64){
	input = in
	for _,node := range(nodes){
		conn, _ := net.Dial("tcp",node)
		defer conn.Close()	
		encoder := json.NewEncoder(conn)
		encoder.Encode(in)
	}
}

func Listen(){
	ln, _ := net.Listen("tcp", serverhost)
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
	port = strings.TrimSpace(port)

	str,_ := r.ReadString('\n')
	output,_ := strconv.Atoi(strings.TrimSpace(str))

	if output == -1{
		resultado="Venenoso"
		countpoisonous++
	} else {
		resultado="Comestible"
		countedible++
	}

	fmt.Println("__________")
	fmt.Println("Información del puerto ",port,": resultado:",resultado)

	countnodes++
	if countnodes == len(nodes){
		var pe float64 = float64(countedible)/float64(countnodes)*100
		var pp float64 = float64(countpoisonous)/float64(countnodes)*100
		fmt.Println("__________")
		fmt.Println("Poisibilidad de acierto (comestible):",pe,"%")
		fmt.Println("Poisibilidad de acierto (venenoso):",pp,"%")
		fmt.Println("__________")

		var d Data
		d.Inputs = input
		if pe > pp {d.Output = 1} else {d.Output = -1}
		data = append(data,d)

		countnodes = 0
		countedible = 0
		countpoisonous = 0
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func main(){
	fmt.Println("__________")
	fmt.Println("Ingrese el puerto para el servidor.")
	bIn := bufio.NewReader(os.Stdin)
	fmt.Print("Actual: ")
	serverport,_ := bIn.ReadString('\n')
	serverhost = fmt.Sprintf("localhost:%s", strings.TrimSpace(serverport))

	fmt.Println("__________")
	fmt.Println("Ingrese los puertos asociados al servidor.")
	fmt.Print("Cantidad de nodos: ")
	str,_ := bIn.ReadString('\n')
	n,_ := strconv.Atoi(strings.TrimSpace(str))

	for i:=0;i<n;i++{
		fmt.Println("Información del nodo",i+1)
		fmt.Print("Puerto: ")
		port,_ := bIn.ReadString('\n')
		port = strings.TrimSpace(port)
		remotehost := fmt.Sprintf("localhost:%s", port)
		nodes = append(nodes,remotehost)
	}

	// input := []float64{'x','s','n','t','p','f','c','n','k','e','e','s','s','w','w','p','w','o','p','k','s','u'}
	// in := []float64{'s','y','h','w','j','l','c','n','p','e','e','k','w','g','a','l','w','s','f','j','q','a'}
	go Listen()
	// Send(in)

	http.HandleFunc("/get",get)
	http.HandleFunc("/post",post)
	http.ListenAndServe(":9000",nil)
}