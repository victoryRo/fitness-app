package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func setCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			return
		}

		if r.PostFormValue("username") == "user@user" && r.PostForm.Get("password") == "password" {
			log.Println("Setting cookie")
			// Note: set the cookie before writing the response
			http.SetCookie(w, &http.Cookie{
				Name:  "user-session",
				Value: "user@user:password",
			})
			_, err := fmt.Fprintf(w, "Successful login ...")
			if err != nil {
				log.Println("error with login password and username")
			}
			return
		}

		_, err = fmt.Fprintf(w, "Bad login")
		if err != nil {
			return
		}
	}
}

func checkCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking cookies:")
		for _, c := range r.Cookies() {
			log.Println(c)
		}
	}
}

// unsetCookie disabled the cookie
func unsetCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Deleting cookie")
		http.SetCookie(w, &http.Cookie{
			Name:   "user-session",
			Value:  "",
			MaxAge: -1,
			Expires: time.Date(
				1983, 7, 26, 20, 34, 58, 651387237, time.UTC),
		})
		_, err := fmt.Fprintf(w, "Successful logout")
		if err != nil {
			return
		}
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", setCookie()).Methods(http.MethodPost)
	router.HandleFunc("/", checkCookie()).Methods(http.MethodGet)
	router.HandleFunc("/logout", unsetCookie()).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:           router,
		Addr:              "127.0.0.1:3221",
		WriteTimeout:      time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
	}

	fmt.Println("Starting server at", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
