package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	totalRequests   int
	totalTimeInNSec int64
	router          *http.ServeMux
	hashMap         map[int]string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Hash and Encode Server", "Shweta Bhandare")
	s.router.ServeHTTP(w, r)
}

func getaddr() string {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}
	return addr
}

func setup() *http.Server {

	addr := getaddr()

	s := HashServer(func(s *server) {
		s.totalRequests = 0
		s.totalTimeInNSec = 0
	})

	hs := &http.Server{Addr: addr, Handler: s}

	return hs
}

func graceful(hs *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Server stopped")
	}
}

func main() {

	httpServer := setup()

	go func() {
		fmt.Printf("Listening on http://0.0.0.0%s\n", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	graceful(httpServer, 5*time.Second)
}
