package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// handlerSlug es un Controller
func handlerSlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	if slug == "" {
		log.Println("Slug not provided")
		return
	}
	log.Println("Got slug", slug)
}

// handlerGetHelloWorld es un Controller
func handlerGetHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World against\n")
	fmt.Println("Request via", r.Method)
	fmt.Println("Current url", r.URL)
	fmt.Println("Request header", r.Header)
	fmt.Println("Request body", r.Body)
}

// handlerPostEcho ... es otro Controller
func handlerPostEcho(w http.ResponseWriter, r *http.Request) {
	log.Println("Request via", r.Method)
	log.Println(r.URL)
	log.Println(r.Header)

	// Vamos a leerlo en un buffer
	// ya que el cuerpo de la solicitud es io.ReadCloser
	// y por eso solo deberíamos leerlo una vez.
	body, err := io.ReadAll(r.Body)

	log.Println("read >", string(body), "<")

	n, err := io.Copy(w, bytes.NewReader(body))
	if err != nil {
		log.Println("Error echoing response", err)
	}
	log.Println("Wrote back", n, "bytes")
}

func main() {
	// banderas para acompañar el registro
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	port := "9002"
	if value, exists := os.LookupEnv("SERVER_PORT"); exists {
		port = value
	}

	// Desde el principio, podemos aplicar StrictSlash
	// Esta es una buena función auxiliar que significa
	// Cuando es verdadero, si la ruta es "/foo/", acceder a "/foo" realizará una redirección 301 a la primera y viceversa.
	// En otras palabras, su aplicación siempre verá la ruta especificada en la ruta.
	// Cuando es falso, si la ruta es "/foo", acceder a "/foo/" no coincidirá con esta ruta y viceversa.
	router := mux.NewRouter().StrictSlash(true)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	router.HandleFunc("/", handlerGetHelloWorld).Methods(http.MethodGet)
	router.HandleFunc("/", handlerPostEcho).Methods(http.MethodPost)
	router.HandleFunc("/{slug}", handlerSlug).Methods(http.MethodGet)

	log.Println("Starting on", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Couldn't start server ListenAndServe()", err)
	}
}
