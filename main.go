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
	"sync"
)

type server struct {
	totalRequests   int
	totalTimeInNSec int64
	router          *http.ServeMux
	hashMap         map[int]string
	shutdownReq     chan bool
	lock            *sync.Mutex
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

func setup() (*http.Server, *server) {

	addr := getaddr()

	s := HashServer(func(s *server) {
		s.totalRequests = 0
		s.totalTimeInNSec = 0
	})

	hs := &http.Server{Addr: addr, Handler: s}

	return hs, s
}

func (s *server) waituntilshutdown(hs *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-stop:
		fmt.Printf("Shutdown request signal: %v\n", sig)
	case sig := <-s.shutdownReq:
		fmt.Printf("Shutdown request received: %v\n", sig)
	}


	for true {

		s.lock.Lock()

		numHashed := len(s.hashMap)
		numReq  := s.totalRequests

		s.lock.Unlock()

    	if numHashed < numReq {
			time.Sleep(time.Second)
    	} else {
			break
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := hs.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf(string(s.constructjson()))
		fmt.Println("Server stopped")
	}
}

func main() {

	httpServer, s := setup()

	go func() {
		fmt.Printf("Listening on http://0.0.0.0%s\n", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	s.waituntilshutdown(httpServer, 5*time.Second)
}
