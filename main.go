package main

// import (
// 	"bufio"
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"math/rand"
// 	"net"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// type Perceptron struct {
// 	Data    []Data    `json:"data"`
// 	Weights []float64 `json:"weights"`
// 	Umbral  float64   `json:"umbral"`
// }

// type Data struct {
// 	Inputs  []int `json:"inputs"`
// 	Outputs int   `json:"outputs"`
// }

// func (p *Perceptron) initData() {
// 	file, _ := os.Open("mushrooms-cleaned.csv")
// 	reader := csv.NewReader(file)
// 	objects, _ := reader.ReadAll()

// 	for _, arr := range objects {
// 		var data Data
// 		for j, val := range arr {
// 			if j > 0 {
// 				data.Inputs = append(data.Inputs, int(val[0]))
// 			} else {
// 				if val == "e" {
// 					data.Outputs = 1
// 				} else {
// 					data.Outputs = -1
// 				}
// 			}
// 		}
// 		p.Data = append(p.Data, data)
// 	}

// 	for i := 0; i < len(p.Data[0].Inputs); i++ {
// 		p.Weights = append(p.Weights, rand.Float64()+0.5)
// 	}

// 	p.Umbral = 0.3
// }

// func (p *Perceptron) get(w http.ResponseWriter, r *http.Request) {
// 	jsonbytes, err := json.MarshalIndent(p, "", " ")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 	}
// 	w.Header().Add("content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonbytes)
// }

// func (p *Perceptron) post(w http.ResponseWriter, r *http.Request) {
// 	bodyBytes, err := ioutil.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 	}
// 	var data Data
// 	err = json.Unmarshal(bodyBytes, &data)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(err.Error()))
// 	}
// 	p.Data = append(p.Data, data)
// 	p.Weights = append(p.Weights, rand.Float64()+0.5)
// }

// var serverport string = "localhost:8000"
// var nodes = []string{
// 	"localhost:8001",
// 	"localhost:8002",
// 	"localhost:8003",
// 	"localhost:8004",
// 	"localhost:8005",
// }

// func send(p Perceptron) {
// 	for id := 0; id < 1; id++ {
// 		var object Perceptron
// 		object.Weights = p.Weights
// 		object.Data = append(object.Data, p.Data[id])
// 		object.Umbral = p.Umbral
// 		cont := 0
// 		for _, node := range nodes {
// 			conn, _ := net.Dial("tcp", node)
// 			defer conn.Close()
// 			encoder := json.NewEncoder(conn)
// 			encoder.Encode(object)
// 			fmt.Fprintf(conn, "%d\n", cont)
// 			cont += 5
// 		}
// 		ch <- object
// 	}
// }

// func WaitForResult() {
// 	ln, _ := net.Listen("tcp", serverport)
// 	defer ln.Close()
// 	for {
// 		conn, _ := ln.Accept()
// 		go handlerListen(conn)
// 	}
// }

// var f float64 = 0.0
// var countnodes int = 0

// func handlerListen(conn net.Conn) {
// 	defer conn.Close()
// 	r := bufio.NewReader(conn)
// 	str, _ := r.ReadString('\n')
// 	num, _ := strconv.ParseFloat(strings.TrimSpace(str), 64)
// 	f += num

// 	countnodes++

// 	if countnodes == 5 {
// 		p := <-ch
// 		close(ch)
// 		fmt.Println(p.Weights)
// 		p = PerceptronAlgorithm(f, p)
// 		fmt.Println("---------------------------------------------")
// 		fmt.Println(p.Weights)
// 		perceptron.Weights = p.Weights
// 		countnodes = 0
// 		f = 0
// 	}
// }

// func PerceptronAlgorithm(f float64, p Perceptron) Perceptron {
// 	f += p.Umbral
// 	var result int = 0
// 	if f <= 0 {
// 		result = -1
// 	} else {
// 		result = 1
// 	}

// 	if result != p.Data[0].Outputs {
// 		for i := 0; i < len(p.Weights); i++ {
// 			p.Weights[i] += float64(p.Data[0].Outputs * p.Data[0].Inputs[i])
// 		}
// 		p.Umbral += float64(p.Data[0].Outputs)
// 	}

// 	return p
// }

// var perceptron Perceptron
// var ch chan Perceptron

// func main() {
// 	perceptron.initData()
// 	http.HandleFunc("/dataset/get", perceptron.get)
// 	http.HandleFunc("/dataset/post", perceptron.post)

// 	ch = make(chan Perceptron, 1)

// 	go send(perceptron)
// 	go WaitForResult()

// 	http.ListenAndServe(":9000", nil)
// }
