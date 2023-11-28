package main

import (
	"context"
	"database/sql"
	"net/http"
	"six/internal/api"
	"six/store"
)

// Este middleware no es reutilizable fuera de nuestra aplicación ya que utiliza
// genera almacenes y es mejor extraerlo y mantenerlo acoplado al resto
// de nuestro código

type UserSession struct {
	UserID int64
}

// Definimos esto para que no pueda chocar fuera de nuestro paquete
// con cualquier otra cosa.
type ourCustomKey string

const sessionKey ourCustomKey = "unique-session-key-for-our-example"

func validCookieMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := cookieStore.Get(r, "session-name")
			if err != nil {
				api.JSONError(w, http.StatusInternalServerError, "Session error")
				return
			}

			userID, userIDOK := session.Values["userID"].(int64)
			isAuth, isAuthOK := session.Values["userAuthenticated"].(bool)

			// Podríamos seguir con lo anterior pero mantengamos nuestra lógica simple
			if !userIDOK || !isAuthOK {
				api.JSONError(w, http.StatusInternalServerError, "Session error")
				return
			}
			if !isAuth || userID < 1 {
				api.JSONError(w, http.StatusForbidden, "Bad credentials")
				return
			}

			querier := store.New(db)
			user, err := querier.GetUser(r.Context(), int64(userID))
			if err != nil || user.UserID < 1 {
				api.JSONError(w, http.StatusForbidden, "Bad credentials")
				return
			}

			ctx := context.WithValue(r.Context(), sessionKey, UserSession{UserID: user.UserID})
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func userFromSession(r *http.Request) (UserSession, bool) {
	session, ok := r.Context().Value(sessionKey).(UserSession)
	if session.UserID < 1 {
		// No debería suceder
		return UserSession{}, false
	}
	return session, ok
}
