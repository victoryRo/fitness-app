package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"six/internal"
	"six/internal/api"
	"six/store"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var cookieStore = sessions.NewCookieStore([]byte("forDemo"))

func init() {
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func handleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Gracias a nuestro middleware, sabemos que tenemos JSON
		// lo decodificaremos en nuestro tipo de solicitud y veremos si es válido
		type loginRequest struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}

		payload := loginRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(w, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		querier := store.New(db)
		user, err := querier.GetUserByName(r.Context(), payload.Username)
		if errors.Is(err, sql.ErrNoRows) || !internal.CheckPasswordHash(payload.Password, user.PasswordHash) {
			api.JSONError(w, http.StatusForbidden, "Bad credentials")
			return
		}
		if err != nil {
			log.Println("Received error looking up user", err)
			api.JSONError(w, http.StatusInternalServerError, "Couldn't log you in due to a server error")
			return
		}

		// Somos válidos. Digámosle al usuario y configuremos una cookie.
		session, err := cookieStore.Get(r, "session-name")
		if err != nil {
			log.Println("Cookie store failed with error", err)
			api.JSONError(w, http.StatusInternalServerError, "Session error")
		}
		session.Values["userAuthenticated"] = true
		session.Values["userID"] = user.UserID
		err = session.Save(r, w)
		if err != nil {
			log.Fatalf("Erro to save session %s:", err.Error())
		}
	}
}

func checkSecret(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDetails, _ := userFromSession(r)

		querier := store.New(db)
		user, err := querier.GetUser(r.Context(), userDetails.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			api.JSONError(w, http.StatusForbidden, "User not found")
			return
		}

		api.JSONMessage(w, http.StatusOK, fmt.Sprintf("Hello there %s", user.UserName))
	}
}

func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := cookieStore.Get(r, "session-name")
		if err != nil {
			log.Println("Cookie store failed with", err)
			api.JSONError(w, http.StatusInternalServerError, "Session error")
			return
		}

		session.Options.MaxAge = -1 // deletes
		session.Values["userID"] = int64(-1)
		session.Values["userAuthenticated"] = false

		err = session.Save(r, w)
		if err != nil {
			api.JSONError(w, http.StatusInternalServerError, "Session error")
			return
		}

		api.JSONMessage(w, http.StatusOK, "logout successful")
	}
}

func handleCreateNewWorkout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDetails, ok := userFromSession(r)
		if !ok {
			api.JSONError(w, http.StatusForbidden, "Bad context")
			return
		}
		querier := store.New(db)

		res, err := querier.CreateUserWorkout(r.Context(), userDetails.UserID)
		if err != nil {
			api.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			log.Fatal("Error when encode json", err.Error())
		}
	}
}

func handleListWorkouts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDetails, ok := userFromSession(r)
		if !ok {
			api.JSONError(w, http.StatusForbidden, "Bad context")
			return
		}

		querier := store.New(db)
		workouts, err := querier.GetWorkoutsForUserID(r.Context(), userDetails.UserID)
		if err != nil {
			api.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = json.NewEncoder(w).Encode(&workouts)
		if err != nil {
			log.Fatal("Error when encode json", err.Error())
		}
	}
}

func handleAddSet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workoutID, err := strconv.Atoi(mux.Vars(r)["workout_id"])
		if err != nil {
			api.JSONError(w, http.StatusBadRequest, "Bad workout_id")
			return
		}

		type newSetRequest struct {
			ExerciseName string `json:"exercise_name,omitempty"`
			Weight       int    `json:"weight,omitempty"`
		}

		payload := newSetRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(w, http.StatusBadRequest, "Error decoding json")
			return
		}

		querier := store.New(db)

		// set variable retorna el modelo "struct" AppSet
		set, err := querier.CreateDefaultSetForExercise(r.Context(), store.CreateDefaultSetForExerciseParams{
			WorkoutID:    int64(workoutID),
			ExerciseName: payload.ExerciseName,
			Weight:       int32(payload.Weight),
		})
		if err != nil {
			api.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// NewEncoder escribe en la salida y Encode retorna un json
		err = json.NewEncoder(w).Encode(&set)
		if err != nil {
			log.Fatal("Error when encode json", err.Error())
		}
	}
}

func handleUpdateSet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}

func handleDeleteWorkout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userDetails, ok := userFromSession(r)
		if !ok {
			api.JSONError(w, http.StatusForbidden, "Bad context")
			return
		}

		workoutID, err := strconv.Atoi(mux.Vars(r)["workout_id"])
		if err != nil {
			api.JSONError(w, http.StatusBadRequest, "Bad workout_id")
			return
		}

		err = store.New(db).DeleteWorkoutByIDForUser(r.Context(), store.DeleteWorkoutByIDForUserParams{
			UserID:    userDetails.UserID,
			WorkoutID: int64(workoutID),
		})

		if err != nil {
			api.JSONError(w, http.StatusBadRequest, "Bad workout_id")
			return
		}

		api.JSONMessage(w, http.StatusOK, fmt.Sprintf("Workout %d is deleted", workoutID))
	}
}
