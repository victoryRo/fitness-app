package main

import (
	"bytes"
	"encoding/json"
	"log"
)

// ejemplo para mostrar el registro usando la biblioteca estándar
func main() {
	// returns the standard logger
	ol := log.Default()

	// establece el formato de registro en - dd/mm/aa hh:mm:ss
	ol.SetFlags(log.LstdFlags)    // retorna la fecha y la hora cuando se ejecuto el log
	ol.Println("Just a log text") // texto, error o comentario del log
	lognumber(ol)
	logjson(ol)
}

// logjson para registrar json como registro
func logjson(ol *log.Logger) {
	ol.SetFlags(log.Ltime) // retorna hora:minutos:segundos

	ex := `{"name": "Cake","batters":{"batter":[{ "id": "001", "type": "Good Food" }]},"topping":[{ "id": "002", "type": "Syrup" }]}`

	var prettyJSON bytes.Buffer
	// indenta el json dando un formato mas facil de leer
	error := json.Indent(&prettyJSON, []byte(ex), "", "\t")
	if error != nil {
		ol.Fatalf("Error parsing : %s", error.Error())
	}

	// toma el buffer de bytes y devuelve un string en formato indentado
	ol.Println(string(prettyJSON.Bytes()))
}

// lognumber para registrar numero como registro
func lognumber(ol *log.Logger) {
	ol.SetFlags(log.Lshortfile) // mostrará el formato de nombre de archivo: número de línea
	ol.Printf("This is number %d", 1)
}
