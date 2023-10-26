package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var (
	// Version ...
	// TrimSpace limpia los espacios al inicio y al final del string
	Version string = strings.TrimSpace(version)

	//go:embed version/version.txt
	version string

	//go:embed static/*
	staticEmbed embed.FS

	//go:embed tmpl/*.html
	tmplEmbed embed.FS

	// Un FS es una colección de archivos de solo lectura
)

type staticHandler struct {
	staticPath string
	indexPage  string
}

func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Abs devuelve una representación absoluta de la ruta
	path, err := filepath.Abs(r.URL.Path)
	log.Println("Url Path: ", r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// une cualquier número de elementos de ruta en una única ruta
	path = filepath.Join(h.staticPath, path)
	// Stat devuelve un FileInfo que describe el archivo nombrado
	_, err = os.Stat(path)

	log.Print("using embed mode")
	// Sub devuelve un FS correspondiente al subárbol con raíz en el dir de fsys
	fsys, err := fs.Sub(staticEmbed, "static")
	if err != nil {
		panic(err)
	}

	http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
}

// renderFiles renderiza el archivo y envía datos (d) a las plantillas que se van a renderizar
func renderFiles(tmpl string, w http.ResponseWriter, d interface{}) {
	// ParseFS es como ParseFiles o ParseGlob pero lee desde el sistema de archivos fsys
	// en lugar del sistema de archivos del sistema operativo host
	t, err := template.ParseFS(tmplEmbed, fmt.Sprintf("tmpl/%s.html", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	// Ejecutar aplica una plantilla analizada al objeto
	// de datos especificado y escribe la salida en wr.
	if err := t.Execute(w, d); err != nil {
		log.Fatal(err)
	}

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	result := "Login "
	r.ParseForm()

	if validateUser(r.FormValue("username"), r.FormValue("password")) {
		result = result + "successful"
	} else {
		result = result + "unsuccessful"
	}

	renderFiles("msg", w, result)
}

func validateUser(username string, password string) bool {
	return (username == "admin") && (password == "admin")
}

func main() {
	log.Println("Server Version :", Version)

	router := mux.NewRouter()

	router.HandleFunc("/login", postHandler).Methods(http.MethodPost)

	spa := staticHandler{
		staticPath: "static",
		indexPage:  "index.html",
	}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3111",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Running server on port: 3111")
	log.Fatal(srv.ListenAndServe())
}
