
package fgms


import (

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)


func InitHttp() {
	r := mux.NewRouter()
	r.HandleFunc("/flights.json", JsonFlightsHandler)
	//r.HandleFunc("/products", ProductsHandler)
	//r.HandleFunc("/articles", ArticlesHandler)

	http.Handle("/", r)
	go http.ListenAndServe(":8888", nil)
	fmt.Println("Started http------------------------")
}


func JsonFlightsHandler(resp http.ResponseWriter, req *http.Request){

	resp.Write( []byte("Yes") )
}


