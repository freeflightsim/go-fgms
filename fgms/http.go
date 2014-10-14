
package fgms


import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)


// Setup HTTP Handlers and start HTTP server
func SetupHttp() {
	r := mux.NewRouter()
	r.HandleFunc("/flights.json", JsonFlightsHandler)
	//r.HandleFunc("/products", ProductsHandler)
	//r.HandleFunc("/articles", ArticlesHandler)

	http.Handle("/", r)
	go http.ListenAndServe(":8888", nil)
	fmt.Println("Started http------------------------")
}


type FlightsPayload struct {
	Success bool `json:"success"`
	Flights []FlightPayload `json:"flights"`
}

type FlightPayload struct {
	Callsign string `json:"callsign"`
	Model string `json:"model"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt_ft"`
	Speed float32 `json:"kt"`
}

// Handle /flights.json
func JsonFlightsHandler(resp http.ResponseWriter, req *http.Request){

	payload := FlightsPayload{Success: true}
	payload.Flights = make([]FlightPayload, 0)

	for _, p := range Server.Players {

		fl := FlightPayload{Callsign: p.Callsign, Model: p.ModelName}
		fl.Lat, fl.Lon, fl.Alt = p.GetLatLonAlt()
		payload.Flights = append(payload.Flights, fl)
	}

	bits, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return
	}
	resp.Write( bits )
}



