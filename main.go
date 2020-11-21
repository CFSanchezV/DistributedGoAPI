package main

import (
	"encoding/json"
	"encoding/csv"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"bufio"
	"fmt"
	"net"
	"os"
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

func (p *Perceptron) initData() {
	file,_ := os.Open("mushrooms-cleaned.csv")
	reader := csv.NewReader(file)
	objects,_ := reader.ReadAll()

	for _,arr := range objects {
		var data Data
		for j,val := range arr {
			if j > 0{
				data.Inputs = append(data.Inputs,int(val[0]))
			} else {
				if val == "e"{
					data.Outputs = 1
				} else {
					data.Outputs = -1
				} 
			}
		}
		p.Data = append(p.Data,data)
	}

	for i:=0; i<len(p.Data[0].Inputs); i++{
		p.Weights = append(p.Weights,rand.Float64()+0.5)
	}
}

func (p *Perceptron) get(w http.ResponseWriter, r *http.Request) {
	jsonbytes,err := json.MarshalIndent(p,""," ")
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonbytes)
}

//func (p *Perceptron) post(w http.ResponseWriter, r *http.Request) {
//	bodybytes,err := ioutil.ReadAll(r.Body)
//	defer r.Body.Close()
//
//	if err != nil { 
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte(err.Error()))
//	}
//	// complete
//}

var hostname string = "localhost:8000"
var remotehost string = "localhost:8001"

func waitForResult(){
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handlerResult(conn)
	}
}

func handlerResult(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	str, _ := r.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(str))

	fmt.Println("Result: %d",num)
}

func send(p Perceptron,id int) {
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	encoder.Encode(p)

	fmt.Fprintf(conn,"%d\n%d\n",id,len(p.Weights))
}

var perceptron Perceptron

func main(){
	perceptron.initData()
	http.HandleFunc("/dataset",perceptron.get)

	send(perceptron,0)

	http.ListenAndServe(":9000",nil)
}