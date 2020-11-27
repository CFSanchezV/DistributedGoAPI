package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Perceptron - algoritmo de IA
type Perceptron struct {
	Data    []Data    `json:"data"`
	Weights []float64 `json:"weights"`
	Umbral  float64   `json:"umbral"`
	Epochs  int       `json:"epochs"`
}

// Data - data a procesar
type Data struct {
	Inputs []float64 `json:"inputs"`
	Output float64   `json:"output"`
}

// InitDownload - inicializacion con descarga
func (p *Perceptron) InitDownload() {
	resp, err := http.Get("https://github.com/stedy/Machine-Learning-with-R-datasets/raw/master/mushrooms.csv")
	if err != nil {
		panic("Error cargando datos de URL")
	}
	reader := csv.NewReader(resp.Body)
	objects, _ := reader.ReadAll()

	for i, arr := range objects {
		if i == 0 {
			continue
		}
		var data Data
		for j, val := range arr {
			if j > 0 {
				data.Inputs = append(data.Inputs, float64(val[0]))
			} else {
				if val == "e" {
					data.Output = 1
				} else {
					data.Output = -1
				}
			}
		}
		p.Data = append(p.Data, data)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(p.Data[0].Inputs); i++ {
		p.Weights = append(p.Weights, rand.Float64()+float64(rand.Intn(5)+1))
	}
}

// Init - inicializacion sin descarga
func (p *Perceptron) Init() {
	file, _ := os.Open("mushrooms-cleaned.csv")
	reader := csv.NewReader(file)
	objects, _ := reader.ReadAll()

	for _, arr := range objects {
		var data Data
		for j, val := range arr {
			if j > 0 {
				data.Inputs = append(data.Inputs, float64(val[0]))
			} else {
				if val == "e" {
					data.Output = 1
				} else {
					data.Output = -1
				}
			}
		}
		p.Data = append(p.Data, data)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(p.Data[0].Inputs); i++ {
		p.Weights = append(p.Weights, rand.Float64()+float64(rand.Intn(5)+1))
	}
}

// Train - entrenamiento de neurona
func (p *Perceptron) Train() {
	for i := range p.Data {
		f := 0.0
		for j := range p.Data[i].Inputs {
			f += p.Data[i].Inputs[j] * p.Weights[j]
		}
		f += p.Umbral
		if f <= 0 {
			f = -1
		} else {
			f = 1
		}
		if f != p.Data[i].Output {
			p.UpdateWeights(i)
		}
	}
}

// UpdateWeights - actualizar pesos
func (p *Perceptron) UpdateWeights(id int) {
	for i := range p.Weights {
		p.Weights[i] += p.Data[id].Output * p.Data[id].Inputs[i]
	}
	p.Umbral += p.Data[id].Output
}

// Predict - predecir clase
func (p *Perceptron) Predict(inputs []float64) int {
	f := 0.0
	for i := range inputs {
		f += inputs[i] * p.Weights[i]
	}
	f += p.Umbral
	if f <= 0 {
		return -1
	} else {
		return 1
	}
}

// Handler - handle HTTP request
func Handler(conn net.Conn) {
	defer conn.Close()
	<-ch
	r := bufio.NewReader(conn)
	str, _ := r.ReadString('\n')
	var input []float64
	json.Unmarshal([]byte(str), &input)
	fmt.Println("Entrada recibida y procesandola.")
	result := perceptron.Predict(input)
	Send(result)
}

// Send - enviar
func Send(result int) {
	conn, _ := net.Dial("tcp", serverport)
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", port)
	fmt.Fprintf(conn, "%d\n", result)
}

var perceptron Perceptron
var serverport string
var port string
var ch chan bool

func main() {
	fmt.Println("Descargando dataset de https://github.com/stedy/Machine-Learning-with-R-datasets/raw/master/mushrooms.csv...")
	go perceptron.InitDownload()

	fmt.Println("__________")
	fmt.Println("Ingrese los puertos para el nodo.")

	bIn := bufio.NewReader(os.Stdin)
	fmt.Print("Actual: ")
	port, _ = bIn.ReadString('\n')
	port = strings.TrimSpace(port)
	hostname := fmt.Sprintf("localhost:%s", port)

	fmt.Print("Servidor: ")
	remoteport, _ := bIn.ReadString('\n')
	remoteport = strings.TrimSpace(remoteport)
	serverport = fmt.Sprintf("localhost:%s", remoteport)

	fmt.Println("__________")
	fmt.Println("Ingrese valores para el algoritmo.")

	fmt.Print("Umbral: ")
	str, _ := bIn.ReadString('\n')
	umbral, _ := strconv.ParseFloat(strings.TrimSpace(str), 64)
	perceptron.Umbral = umbral

	fmt.Print("Epocas: ")
	str, _ = bIn.ReadString('\n')
	epochs, _ := strconv.Atoi(strings.TrimSpace(str))
	perceptron.Epochs = epochs

	ch = make(chan bool, 1)
	it := 0

	go func() {
		for it < perceptron.Epochs {
			perceptron.Train()
			time.Sleep(300 * time.Nanosecond)
			it++
		}
		ch <- true
		close(ch)
	}()

	fmt.Println("__________")
	fmt.Println("Escuchando...")
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go Handler(conn)
	}
}
