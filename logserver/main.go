package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
)

var router *mux.Router

func main() {
	log.Println("Initializing logging server at port 8010.")
	runServer(":8010")
}

// runServer para ejecutar el servidor de registro
func runServer(addr string) {
	router = mux.NewRouter()
	initializeRoutes()

	scheme := pp.ColorScheme{
		String: pp.Green | pp.Bold,
		Float:  pp.Black | pp.BackgroundWhite | pp.Bold,
	}
	pp.SetColorScheme(scheme)

	log.Fatal(http.ListenAndServe(addr, router))
}

// respondWithError manejar la respuesta de error
func respondWithError(w http.ResponseWriter, code int, message string) {
	// retorna la codificacion json
	response, _ := json.Marshal(message)

	// establecemos el tipo de contenido a retornar
	w.Header().Set("Content-Type", "application/json")
	// envia el codigo de estado de la respuesta
	w.WriteHeader(code)
	w.Write(response)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	// lee desde el cuerpo de la respuesta http
	body, err := io.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	pp.Println(string(body))
	w.WriteHeader(http.StatusCreated)
}

// initializeRoutes para inicializar diferentes rutas
func initializeRoutes() {
	router.HandleFunc("/log", logHandler).Methods(http.MethodPost)
}
