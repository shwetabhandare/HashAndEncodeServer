package main

import (
	"net/http"
	"os"
	"fmt"
	"log"
	"time"
	"strings"
	"os/signal"
	"syscall"
	"context"
	"strconv"
	"crypto/sha512"
	"encoding/base64"
)

type server struct {
	requestNum int
	router    *http.ServeMux
	hashMap map[int] string
}

func HashServer(options ...func(*server)) *server {
	s := &server{requestNum: 0, router: http.NewServeMux(), hashMap:make(map[int] string)}

	for _, f := range options {
		f(s)
	}

	s.router.HandleFunc("/hash", s.hash)
	s.router.HandleFunc("/hash/", s.gethash)
	s.router.HandleFunc("/stats", s.stats)
	s.router.HandleFunc("/shutdown", s.shutdown)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Hash and Encode Server", "Shweta Bhandare")
	s.router.ServeHTTP(w, r)
}

func (s *server) gethash(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/hash/")
	fmt.Println("Get Hash for ID:", id)

	idAsInt, _ := strconv.Atoi(id)

	if x, ok := s.hashMap[idAsInt]; ok {
		fmt.Println(x) 
	}
}

func (s *server) stats(w http.ResponseWriter, r *http.Request) {
	message := "Hello " + r.URL.Path + "\n"
	w.Write([]byte(message))
}

// 2 parameters, number, and password
// wait for 5 seconds.
// add to hash.

func (s *server) hash(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()                     // Parses the request body
	x := r.Form.Get("password") // x will be "" if parameter is not set
	if (x != "") {
		s.requestNum++;
		fmt.Println(x)
		
		// start a 5 second timer.
		// once timer expires, save password to map.
		go s.saveHashToMap(s.requestNum, x)

		w.Write([]byte(strconv.Itoa(s.requestNum)))

	} else {
		fmt.Printf("Didn't find password\n")
	}

}

func computeHash(password string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(password))

	return base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
}

func (s *server) saveHashToMap(num int, password string) {
	time.Sleep(5*time.Second)

	s.hashMap[num] = computeHash(password)

	fmt.Printf("Added %s to map at %d\n", s.hashMap[num], num)
}

func (s *server) shutdown(w http.ResponseWriter, r *http.Request) {
	message := "Hello " + r.URL.Path + "\n"
	w.Write([]byte(message))
}

func setup() (*http.Server) {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}
	fmt.Printf("Addr: " + addr + "\n")

	s := HashServer(func(s *server) {
		s.requestNum = 0
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

	hs := setup()

	fmt.Printf("done with setup\n")

	go func() {
		fmt.Printf("Running go func() \n")
		fmt.Printf("Listening on http://0.0.0.0%s\n", hs.Addr)

		if err := hs.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	graceful(hs,  5*time.Second)
}