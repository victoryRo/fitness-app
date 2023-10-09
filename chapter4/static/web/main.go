package main

import (
	"log"
	"net/http"
)

func main() {
	// devuelve un controlador que atiende solicitudes HTTP
	// con el contenido del la raiz de sistema de archivos
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Starting up server on port 3333 ...")
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatalf("error occurred starting up server %s : ", err)
	}
}
