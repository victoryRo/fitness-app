package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"six/internal"
	"six/internal/api"
	"six/store"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		internal.GetAsString("DB_USER", "victory"),
		internal.GetAsString("DB_PASS", "programmingExpert"),
		internal.GetAsString("DB_HOST", "localhost"),
		internal.GetAsInt("DB_PORT", 5432),
		internal.GetAsString("DB_NAME", "fullstackdb"),
	)

	// open DB
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	} else {
		log.Println("Starting database correctly ->")
	}

	// check DB connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database -> ", err)
	} else {
		log.Println("Check database connection successful ->")
	}

	// create our demo user
	createUserInDB(db)

	// Start our server
	server := api.NewServer(internal.GetAsInt("SERVER_PORT", 9002))
	server.MustStart()
	defer func(server *api.Server) {
		err := server.Stop()
		if err != nil {
			log.Fatalf("Error stopping server %s:", err.Error())
		}
	}(server)

	defaultMiddleware := []mux.MiddlewareFunc{
		api.JSONMiddleware,
		api.CORSMiddleware(internal.GetAsSlice("CORS_WHITELIST",
			[]string{
				"http://localhost:9000",
				"http://0.0.0.0:9000",
			}, ","),
		),
	}

	// handlers
	server.AddRoute("/login", handleLogin(db), http.MethodPost, defaultMiddleware...)
	server.AddRoute("/logout", handleLogout(), http.MethodGet, defaultMiddleware...)

	// Our session protected middleware
	protectedMiddleware := append(defaultMiddleware, validCookieMiddleware(db))
	server.AddRoute("/checkSecret", checkSecret(db), http.MethodGet, protectedMiddleware...)

	// Workouts
	server.AddRoute("/workout", handleCreateNewWorkout(db), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/workout", handleListWorkouts(db), http.MethodGet, protectedMiddleware...)
	server.AddRoute("/workout/{workout_id}", handleDeleteWorkout(db), http.MethodDelete, protectedMiddleware...)
	server.AddRoute("/workout/{workout_id}", handleAddSet(db), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/workout/{workout_id}/{set_id}", handleUpdateSet(db), http.MethodPut, protectedMiddleware...)

	// Wait for CTRL-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// Bloqueamos aquí hasta que se reciba un CTRL-C/SigInt
	// Una vez recibido, salimos y se limpia el servidor
	<-sigChan
}

func createUserInDB(db *sql.DB) {
	ctx := context.Background()
	querier := store.New(db)

	fmt.Println("Creating user@user...")
	hasPass := internal.HashPassword("password")

	_, err := querier.CreateUsers(ctx, store.CreateUsersParams{
		UserName:     "user@user",
		PasswordHash: hasPass,
		Name:         "Dummy user",
	})

	// Esto es interesante de ver, la biblioteca sql/pq recomienda que usemos
	// este patrón para comprender los errores. Podríamos usar el ErrorCode directamente
	// o busque el tipo específico. Sabemos que estaremos violando la violación_única
	// si nuestro usuario ya existe en la base de datos
	if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		log.Println("Dummy user already present")
		return
	}

	if err != nil {
		log.Println("Failed to create user:", err)
	}
}
