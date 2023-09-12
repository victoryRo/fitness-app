package main

import (
	"context"
	"database/sql"
	chapter01 "fitness/app/gen"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		GetAsString("DB_USER", "postgres"),
		GetAsString("DB_PASSWORD", "mysecretpassword"),
		GetAsString("DB_HOST", "localhost"),
		GetAsInt("DB_PORT", 5432),
		GetAsString("DB_NAME", "postgres"),
	)

	// abrir base de datos
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	// verificar conexion
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
	}

	// crear el almacenamiento
	st := chapter01.New(db)

	// contexto predeterminado
	ctx := context.Background()

	// pasamos los parametros para crear usuario
	_, err = st.CreateUsers(ctx, chapter01.CreateUsersParams{
		UserName:     "testuser",
		PassWordHash: "hash",
		Name:         "thanks",
	})
	if err != nil {
		log.Fatalln("Error creating user :", err)
	}

	eid, err := st.CreateExercise(ctx, "Exercise1")
	if err != nil {
		log.Fatalln("Error creating exercise :", err)
	}

	set, err := st.CreateSet(ctx, chapter01.CreateSetParams{
		ExerciseID: eid,
		Weight:     100,
	})
	if err != nil {
		log.Fatalln("Error updating exercise :", err)
	}

	set, err = st.UpdateSet(ctx, chapter01.UpdateSetParams{
		ExerciseID: eid,
		SetID:      set.SetID,
		Weight:     2000,
	})
	if err != nil {
		log.Fatalln("Error updating set :", err)
	}

	log.Println("Done!")

	u, err := st.ListUsers(ctx)
	if err != nil {
		log.Fatalln("Error listing users :", err)
	}

	for _, usr := range u {
		fmt.Printf("Name: %s, ID %d\n", usr.Name, usr.UserID)
	}
}
