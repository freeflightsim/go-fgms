
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
	Flights []*Player `json:"flights"`
}

// Handle /flights.json
func JsonFlightsHandler(resp http.ResponseWriter, req *http.Request){

	payload := FlightsPayload{Success: true}
	payload.Flights = make([]*Player, 0)
	for _, p := range Server.Players {
		payload.Flights = append(payload.Flights, p)
	}
	bits, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return
	}
	resp.Write( bits )
}



