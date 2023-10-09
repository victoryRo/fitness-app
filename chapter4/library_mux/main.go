package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func handlerGetHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!\n")
	log.Println(r.Method)
	log.Println(r.URL)
	log.Println(r.Header)
	log.Println(r.Body)
}

type foo struct{}

func (f foo) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello alternate\n")
	log.Println(r.Method)
	log.Println(r.URL)
	log.Println(r.Header)
	log.Println(r.Body)
}

func main() {
	// establece algunas banderas para un facil debugging
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	// obtiene un puerto de la SERVER_PORT o un valor por default
	port := "9002"
	if value, exists := os.LookupEnv("SERVER_PORT"); exists {
		port = value
	}

	// Podríamos usar el mux predeterminado y luego usar http.HandleFunc y http.ListenAndServe(port,nil) pero
	// Creo que es mejor crear el tuyo propio. Como veremos con patrones posteriores, cubriremos elegantes
	// apagado también.
	router := http.NewServeMux()

	// instancias del servidor
	srv := http.Server{
		Addr:           ":" + port, // Addr especifica opcionalmente la dirección de escucha del servidor en el formato "host:puerto"
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20, // Esto es 1 MB, es una buena práctica limitar la cantidad de datos que aceptará de un cliente.
	}

	// Esto es sólo para mostrar una forma alternativa de declarar un controlador
	// al tener una estructura que implementa la interfaz ServeHTTP(...)
	// dummyHandler := foo{}

	router.HandleFunc("/", handlerGetHelloWorld)
	// router.Handle("/1", dummyHandler)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Couldn't ListenAndServe()", err)
	}
	fmt.Println("Running server in port 9002")

	// http://localhost:9002/
}
