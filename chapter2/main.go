package main

import (
	"context"
	"database/sql"
	chapter2 "fitness/app/gen"
	"fitness/app/logger"
	"flag"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// retorna un booleano
	l := flag.Bool("local", false, "true - send to stdout, false - send to logging server")
	// analiza la linea de comandos en busca de la flag
	flag.Parse()

	// enviamos el valor bool para establecer Registro local o remoto
	logger.SetLoggingOutput(*l)

	logger.Logger.Debugf("Application logging to stdout = %v", *l)
	logger.Logger.Info("Starting the application...")

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		GetAsString("DB_USER", "postgres"),
		GetAsString("DB_PASSWORD", "SecretPassword"),
		GetAsString("DB_HOST", "localhost"),
		GetAsInt("DB_PORT", 5432),
		GetAsString("DB_NAME", "postgres"),
	)

	// abrimos la base de datos
	db, err := sql.Open("postgres", dbURI)
	// envie el error al logger pkg
	if err != nil {
		logger.Logger.Errorf("Error opening database : %s", err.Error())
	}

	// verificamos la conexion
	if err := db.Ping(); err != nil {
		logger.Logger.Errorf("Error from database ping : %s password %s", err.Error(), GetAsString("DB_PASSWORD", "SecretPassword"))
	}

	logger.Logger.Info("Database connection fine")

	// crear almacenamiento
	st := chapter2.New(db)

	ctx := context.Background()

	chuser, err := st.CreateUsers(ctx, chapter2.CreateUsersParams{
		UserName:     "testuser",
		PassWordHash: "hash",
		Name:         "test",
	})

	if err != nil {
		// logger.Logger.Fatal("Error creating user")
		logger.Logger.Fatalf("Error creating user: %v", err)
	}
	logger.Logger.Info("Success - user creation")

	eid, err := st.CreateExercise(ctx, "Exercise 2")

	if err != nil {
		logger.Logger.Errorf("Error creating exercise: %s", err)
	}
	logger.Logger.Info("Success - exercise creation")

	sid, err := st.UpsertSet(ctx, chapter2.UpsertSetParams{
		ExerciseID: eid,
		Weight:     100,
	})

	if err != nil {
		logger.Logger.Errorf("Error updating sets %s", err)
	}

	_, err = st.UpsertWorkout(ctx, chapter2.UpsertWorkoutParams{
		UserID:    chuser.UserID,
		SetID:     sid,
		StartDate: time.Time{},
	})

	if err != nil {
		logger.Logger.Error("Error updating workouts")
	}
	logger.Logger.Info("Success - updating workout")

	logger.Logger.Info("Application complete")

	// sentry implement something similar
	// https://github.com/getsentry/sentry-go/blob/master/example/basic/main.go#L50
	defer time.Sleep(1 * time.Second)
}
