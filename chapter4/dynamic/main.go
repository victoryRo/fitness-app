package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type staticHandler struct {
	staticPath string
	indexPage  string
}

// ServeHTTP ...
func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Abs devuelve una representación absoluta de la ruta. Si la ruta no es absoluta,
	// se unirá al directorio de trabajo actual para convertirla en una ruta absoluta
	path, err := filepath.Abs(r.URL.Path)
	log.Println(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Join une cualquier número de elementos de ruta en una única ruta,
	// separándolos con un separador específico del sistema operativo. Los elementos vacíos se ignoran.
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)

	// FileServer devuelve un controlador que atiende solicitudes HTTP con el contenido del sistema de archivos enraizado en la raíz.
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	res := "Login "
	// Para todas las solicitudes,
	// ParseForm analiza la consulta sin formato de la URL y actualiza r.Form.
	r.ParseForm()

	// FormValue devuelve el primer valor del componente nombrado en la consulta.
	if validateUser(r.FormValue("username"), r.FormValue("password")) {
		res = res + "Successfull"
	} else {
		res = res + "Unsuccessful"
	}

	// ParseFiles crea una nueva plantilla y analiza las definiciones de plantilla de los archivos nombrados.
	t, err := template.ParseFiles("static/tmpl/msg.html")
	if err != nil {
		fmt.Fprintf(w, "error processing")
		return
	}

	tpl := template.Must(t, err)
	// Execute aplica una plantilla analizada al objeto de datos especificado y escribe la salida en wr.
	tpl.Execute(w, res)
}

// valida si el nombre de usuario y el password son los esperados
func validateUser(username string, password string) bool {
	return (username == "admin") && (password == "admin")
}

func main() {
	// NewRouter devuelve una nueva instancia de enrutador.
	router := mux.NewRouter()

	router.HandleFunc("/login", postHandler).Methods(http.MethodPost)

	spa := staticHandler{
		staticPath: "static",
		indexPage:  "index.html",
	}

	// PathPrefix registra una nueva ruta con un comparador para el prefijo de ruta URL
	router.PathPrefix("/").Handler(spa)

	// Un servidor define parámetros para ejecutar un servidor HTTP.
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3333",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
