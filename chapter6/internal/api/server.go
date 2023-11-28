package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	port   string
	server http.Server
	router *mux.Router
	wg     sync.WaitGroup
}

func NewServer(port int) *Server {
	router := mux.NewRouter().StrictSlash(true)
	return &Server{
		router: router,
		port:   fmt.Sprintf(":%d", port),
	}
}

func (s *Server) AddRoute(path string, handler http.HandlerFunc, method string, mwf ...mux.MiddlewareFunc) {
	subRouter := s.router.PathPrefix(path).Subrouter()
	subRouter.Use(mwf...)
	subRouter.HandleFunc("", handler).Methods(method)
	log.Printf("Added route: [%v] [%v]\n", path, method)
}

// MustStart iniciará el servidor y si no puede conectarse al puerto
// saldrá con un mensaje de registro fatal
func (s *Server) MustStart() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.server = http.Server{
		Addr:           fmt.Sprintf("0.0.0.0%s", s.port),
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   0 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	// Add to the WaitGroup for the listener goroutine
	s.wg.Add(1)

	// start the listener
	go func() {
		fmt.Printf("Api server stated at %v on http://%s\n", time.Now().Format(time.Stamp), s.server.Addr)

		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Api server failed to start with error: %v\n", err)
		}
		log.Println("Api server stopped")
		s.wg.Done()
	}()
}

// Stop stops the API Server
func (s *Server) Stop() error {
	// Crea un contexto para intentar un apagado elegante de 5 segundos
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("Api server stopping")

	// Intenta el cierre elegante cerrando el oyente
	// y completar todas las solicitudes a bordo
	if err := s.server.Shutdown(ctx); err != nil {
		if e := s.server.Close(); e != nil {
			// Parece que se agotó el tiempo de cierre elegante. Forzar cierre
			log.Printf("API server stop with error: %v", e)
			return e
		}
		return err
	}

	s.wg.Wait()
	return nil
}
