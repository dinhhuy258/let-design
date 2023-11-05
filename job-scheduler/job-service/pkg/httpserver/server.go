package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Interface interface {
	Start(ctx context.Context)
	Stop(ctx context.Context) error
	HandleFunc(method, path string, f http.HandlerFunc)
}

// server contains the server object and router object
type server struct {
	server *http.Server
	router *httprouter.Router
}

// New create a new server
func New(port int, readTimeout, writeTimeout time.Duration) Interface {
	router := httprouter.New()
	server := server{
		server: &http.Server{
			Addr:         ":" + strconv.Itoa(port),
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		router: router,
	}

	return &server
}

// Start starts the HTTP service
func (s *server) Start(ctx context.Context) {
	go func() {
		log.Printf("Serving on port %s", s.server.Addr)

		err := s.server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error from router %+v", errors.WithStack(err))
		}
	}()

	defer func() {
		_ = s.Stop(ctx)
	}()

	// Handle interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// Stop stop the server
func (s *server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// HandleFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle
func (s *server) HandleFunc(method, path string, f http.HandlerFunc) {
	s.router.HandlerFunc(method, path, f)
}
