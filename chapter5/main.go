package main

import (
	"context"
	"database/sql"
	"embed"
	dbpostgres "five/database"
	"five/pkg"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var (
	// Version app
	Version string = strings.TrimSpace(version)
	//go:embed version/version.txt
	version string

	//go:embed src/static
	staticEmbed embed.FS

	//go:embed src/css/*
	cssEmbed embed.FS

	//go:embed src/tmpl/*.html
	tmplEmbed embed.FS

	dbQuery *dbpostgres.Queries

	store = sessions.NewCookieStore([]byte("forDemo"))
)

// renderFiles renderiza el archivo y envía datos (d) a las plantillas que se van a renderizar
func renderFiles(tmpl string, w http.ResponseWriter, d interface{}) {
	t, err := template.ParseFS(tmplEmbed, fmt.Sprintf("src/tmpl/%s.html", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Function renderFiles running ...")
	if err := t.Execute(w, d); err != nil {
		log.Fatal(err)
	}
}

// securityMiddleware es un middleware para garantizar que
// todas las solicitudes tengan un sesión válida y autenticada
func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// primero, toda solicitud DEBE tener una sesión válida
		if sessionValid(w, r) {
			if r.URL.Path == "/login" {
				next.ServeHTTP(w, r)
				return
			}
		}

		// si tiene una sesión válida, asegúrese de que haya sido autenticado
		if hasBeenAuthenticated(w, r) {
			next.ServeHTTP(w, r)
			return
		}

		// de lo contrario será necesario redirigirlo a /login
		storeAuthenticated(w, r, false)
		http.Redirect(w, r, "/login", 307)
	})
}

// sessionValid comprueba si la sesión es válida
func sessionValid(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, "session_token")
	if err != nil {
		log.Fatalf("Problem getting session: %s", err)
	}

	return !session.IsNew
}

// hasBeenAuthenticated comprueba si la sesión contiene la flag para indicar
// que la sesión ha pasado por el proceso de autenticación
func hasBeenAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session_token")
	a, _ := session.Values["authenticated"]

	if a == nil {
		return false
	}

	return a.(bool)
}

// storeAuthenticated para almacenar el valor autenticado
func storeAuthenticated(w http.ResponseWriter, r *http.Request, v bool) {
	session, _ := store.Get(r, "session_token")

	session.Values["authenticated"] = v
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// authenticationHandler maneja la autenticación
func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	result := "Login "
	r.ParseForm()

	if validateUser(r.FormValue("username"), r.FormValue("password")) {
		storeAuthenticated(w, r, true)
		result = result + " Successful OK"
	} else {
		result = result + "Unsuccessful what happened"
	}

	renderFiles("msg", w, result)
}

// validateUser comprueba si el nombre de usuario/contraseña existe en la base de datos
func validateUser(username, password string) bool {
	// consulta los datos de la base de datos
	ctx := context.Background()
	u, _ := dbQuery.GetUserByName(ctx, username)

	// el usuario no existe
	if u.UserName != username {
		return false
	}

	return pkg.CheckPasswordHash(password, u.PassWordHash)
}

func basicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware called on", r.URL.Path)
		// hacer algo
		h.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Server Version :", Version)

	initDatabase()

	router := mux.NewRouter()

	// POST handler for /login
	router.HandleFunc("/login", authenticationHandler).Methods(http.MethodPost)

	// embed handler for /css path
	csscontentStatic, _ := fs.Sub(cssEmbed, "src/css")
	css := http.FileServer(http.FS(csscontentStatic))
	router.PathPrefix("/css").Handler(http.StripPrefix("/css", css))

	// embed handler for /app path
	contentStatic, _ := fs.Sub(staticEmbed, "src/static")
	static := http.FileServer(http.FS(contentStatic))
	router.PathPrefix("/app").Handler(securityMiddleware(http.StripPrefix("/app", static)))

	// add /login path
	router.PathPrefix("/login").Handler(securityMiddleware(http.StripPrefix("/login", static)))

	// root will redirect to /app
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app", http.StatusPermanentRedirect) //308
	})

	// Use our basicMiddleware
	router.Use(basicMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3222",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening at port :3222")
	log.Fatal(srv.ListenAndServe())

}

func initDatabase() {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		GetAsString("DB_USER", "victory"),
		GetAsString("DB_PASSWORD", "yosoyelcamino"),
		GetAsString("DB_HOST", "localhost"),
		GetAsInt("DB_PORT", 5432),
		GetAsString("DB_NAME", "fitness"),
	)

	// Abrir la base de datos
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	// Vericar conexion
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	} else {
		log.Println("Database running well...")
	}

	// Crear almacenamiento
	dbQuery = dbpostgres.New(db)

	ctx := context.Background()

	createUserDb(ctx)

	if err != nil {
		os.Exit(1)
	}
}

func createUserDb(ctx context.Context) {
	// se ha creado el usuario
	u, _ := dbQuery.GetUserByName(ctx, "user@user")

	if u.UserName == "user@user" {
		log.Println("user@user exist...")
		return
	}

	log.Println("Creating user@user...")
	hashPwd, _ := pkg.HashPassword("password")

	_, err := dbQuery.CreateUsers(ctx, dbpostgres.CreateUsersParams{
		UserName:     "user@user",
		PassWordHash: hashPwd,
		Name:         "Dummy user",
	})
	if err != nil {
		log.Println("error getting user@dummyuser.domain", err)
		os.Exit(1)
	}
}
