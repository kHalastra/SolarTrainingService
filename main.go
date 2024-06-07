package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/solarGet/", homeHandler).Methods("GET")

	http.Handle("/solarGet/", r)
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
