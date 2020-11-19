package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

//Mushroom entry in dataset
type Mushroom struct {
	ID                    string `json:"id"`
	class                 string `json:"class"`
	capShape              string `json:"capShape"`
	capSurface            string `json:"capSurface"`
	capColor              string `json:"capColor"`
	bruises               string `json:"bruises"`
	odor                  string `json:"odor"`
	gillAttachment        string `json:"gillAttachment"`
	gillSpacing           string `json:"gillSpacing"`
	gillSize              string `json:"gillSize"`
	gillColor             string `json:"gillColor"`
	stalkShape            string `json:"stalkShape"`
	stalkRoot             string `json:"stalkRoot"`
	stalkSurfaceAboveRing string `json:"stalkSurfaceAboveRing"`
	stalkSurfaceBelowRing string `json:"stalkSurfaceBelowRing"`
	stalkColorAboveRing   string `json:"stalkColorAboveRing"`
	stalkColorBelowRing   string `json:"stalkColorBelowRing"`
	veilType              string `json:"veilType"`
	veilColor             string `json:"veilColor"`
	ringNumber            string `json:"ringNumber"`
	ringType              string `json:"ringType"`
	sporePrintColor       string `json:"sporePrintColor"`
	population            string `json:"population"`
	habitat               string `json:"habitat"`
}

type mushroomHandlers struct {
	sync.Mutex
	bd map[string]Mushroom
}

func (h *mushroomHandlers) mushrooms(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.get(writer, req)
		return
	case "POST":
		// h.post(w, req)
		return
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("Petición no permitida"))
		return
	}
}

func (h *mushroomHandlers) get(w http.ResponseWriter, r *http.Request) {
	mushrooms := make([]Mushroom, len(h.bd))

	h.Lock()
	i := 0
	for _, mushroom := range h.bd {
		mushrooms[i] = mushroom
		i++
	}
	h.Unlock()

	jsonBytes, err := json.MarshalIndent(mushrooms, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//constructor
func newMushroomHandlers() *mushroomHandlers {
	return &mushroomHandlers{
		bd: map[string]Mushroom{
			"anyid24": Mushroom{
				ID:                    "01",
				class:                 "ss",
				capShape:              "ss",
				capSurface:            "ss",
				capColor:              "ss",
				bruises:               "ss",
				odor:                  "ss",
				gillAttachment:        "ss",
				gillSpacing:           "ss",
				gillSize:              "ss",
				gillColor:             "ss",
				stalkShape:            "ss",
				stalkRoot:             "ss",
				stalkSurfaceAboveRing: "ss",
				stalkSurfaceBelowRing: "ss",
				stalkColorAboveRing:   "ss",
				stalkColorBelowRing:   "ss",
				veilType:              "ss",
				veilColor:             "ss",
				ringNumber:            "ss",
				ringType:              "ss",
				sporePrintColor:       "ss",
				population:            "ss",
				habitat:               "ss",
			},
		},
	}
}

func main() {	
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}
