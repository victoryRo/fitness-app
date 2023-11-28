package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"store/rd/gen"
	"store/rd/pkg"
	"strings"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"

	// el paquete contiene el controlador y la API para comunicarse con Redis
	"github.com/go-redis/redis/v8"
	// contiene una API simple para leer, escribir y eliminar datos de Redis:
	rstore "github.com/rbcervilla/redisstore/v8"
)

var (
	Version string = strings.TrimSpace(version)
	//go:embed version/version.txt
	version string

	//go:embed src/static
	staticEmbed embed.FS

	//go:embed src/css/*
	cssEmbed embed.FS

	//go:embed src/tmpl/*.html
	tmplEmbed embed.FS

	dbQuery *gen.Queries

	store *rstore.RedisStore
)

// renderFiles renderiza el archivo y pone los datos (d) dentro del template para ser rederizado
func renderFiles(tmpl string, w http.ResponseWriter, d interface{}) {
	t, err := template.ParseFS(tmplEmbed, fmt.Sprintf("src/tmpl/%s.html", tmpl))
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(w, d); err != nil {
		log.Fatal(err)
	}
}

func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// primero que todo la solicitud debe tener al sesion valida
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

		// de otra manera esto necesita ser redireccionado
		storeAuthenticated(w, r, false)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})
}

// logoutHandler desloguea al usuario
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if hasBeenAuthenticated(w, r) {
		session, _ := store.Get(r, "session_token")
		session.Options.MaxAge = -1
		err := session.Save(r, w)
		if err != nil {
			log.Println("failed to delete session", err)
		}
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

// loginHandler maneja la autenticacion
func loginHandler(w http.ResponseWriter, r *http.Request) {
	result := "Login "
	r.ParseForm()

	if validateUser(r.FormValue("username"), r.FormValue("password")) {
		storeAuthenticated(w, r, true)
		result += result + "Successful"
	} else {
		result += result + "Unsuccessful"
	}

	renderFiles("msg", w, result)
}

// sessionValid verifica si la sesion es valida
func sessionValid(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session_token")
	return !session.IsNew
}

// hasBeenAuthenticated verifica si la sesion contiene la flag para indicar
// que la sesion ha pasado por el proceso de autenticacion
func hasBeenAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session_token")
	a, _ := session.Values["authenticated"]

	if a == nil {
		return false
	}
	return a.(bool)
}

// storeAuthenticated almacena el valor de la autenticacion
func storeAuthenticated(w http.ResponseWriter, r *http.Request, v bool) {
	session, _ := store.Get(r, "session_token")

	session.Values["authenticated"] = v
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validateUser verifica si el usuario existe en la base de datos
func validateUser(username, password string) bool {
	ctx := context.Background()
	// solicitud de datos desde la base de datos
	u, _ := dbQuery.GetUserByName(ctx, username)

	// username no existe
	if u.UserName != username {
		return false
	}

	return pkg.CheckPasswordHash(password, u.PassWordHash)
}

func basicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middlware called on", r.URL.Path)
		// hacer cosas
		h.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting server... version", Version)

	initDB()
	initRedis()

	log.Println("Starting database and redis store... correctly")

	router := mux.NewRouter()

	// POST manejador del login
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)

	// NOTE: embed handler for /css path

	// Sub devuelve un FS correspondiente al subárbol enraizado en el dir de fsys.
	cssContent, _ := fs.Sub(cssEmbed, "src/css")
	// FileServer devuelve un controlador que sirve solicitudes HTTP con el contenido del sistema de archivos rooteado en la raíz.
	// FS convierte fsys en una implementación de FileSystem, y, para usar con FileServer
	css := http.FileServer(http.FS(cssContent))
	// PathPrefix agrega un matcher para el prefijo de ruta de URL
	// Haddler establece un manejador para la ruta
	// StripPrefix devuelve un controlador que sirve solicitudes HTTP
	// eliminando el prefijo dado de la ruta de la URL de solicitud
	router.PathPrefix("/css").Handler(http.StripPrefix("/css", css))

	// NOTE: embed handler for /app path
	contentStatic, _ := fs.Sub(staticEmbed, "src/static")
	static := http.FileServer(http.FS(contentStatic))
	router.PathPrefix("/app").Handler(securityMiddleware(http.StripPrefix("/app", static)))

	// NOTE: add /login path
	router.PathPrefix("/login").Handler(securityMiddleware(http.StripPrefix("/login", static)))

	// NOTE: add /logout path
	router.HandleFunc("/logout", logoutHandler).Methods(http.MethodGet)

	// NOTE: root will redirect to /app
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app", http.StatusPermanentRedirect)
	})

	log.Println("roting working...")

	// use our basic middleware
	router.Use(basicMiddleware)

	srv := http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3011",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Serving... server running on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func initRedis() {
	var err error

	// retorna el client de redis especificado por Options
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// retorna un nuevo almacenamiento de datos de redis
	store, err = rstore.NewRedisStore(context.Background(), client)
	if err != nil {
		log.Fatal("Failed to create redis store ", err)
	} else {
		log.Println("Successfully created redis store")
	}

	// establece un prefix para guardar la sesion en redis
	store.KeyPrefix("session_token")
}

func initDB() {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		GetAsString("DB_USER", "expert"),
		GetAsString("DB_PASSWORD", "asuncionAutosostenida"),
		GetAsString("DB_HOST", "localhost"),
		GetAsInt("DB_PORT", 5432),
		GetAsString("DB_NAME", "fitness"),
	)

	fmt.Printf("dbURI: %s\n", dbURI)

	// open DB
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	// Check DB connection
	if err := db.Ping(); err != nil {
		log.Fatalln("Error pinging database", err)
	} else {
		log.Println("Successfully connected to database")
	}

	// create store
	dbQuery = gen.New(db)

	ctx := context.Background()

	createUserDB(ctx)

	if err != nil {
		os.Exit(1)
	}
}

func createUserDB(ctx context.Context) {
	// el usuario ha sido creado
	u, _ := dbQuery.GetUserByName(ctx, "user@user")

	if u.UserName == "user@user" {
		log.Println("User already exists")
		return
	}

	log.Println("Creating user...")
	hasPass, _ := pkg.HashPassword("password")

	_, err := dbQuery.CreateUsers(ctx, gen.CreateUsersParams{
		UserName:     "user@user",
		PassWordHash: hasPass,
		Name:         "Dummy User",
	})

	if err != nil {
		log.Println("Failed to create user on database method createUserDB()", err)
		os.Exit(1)
	}
}
