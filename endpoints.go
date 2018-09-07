package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"fmt"
)

type statsresponse struct {
	Total   int `json:"total"`
	Average float64 `json:"average"`
}

func (s *server) retrievehash() {
	for {

	}
}

func getidfromurl(urlpath string) int {

	id := strings.TrimPrefix(urlpath, "/hash/")
	if idAsInt, err := strconv.Atoi(id); err == nil {
		return idAsInt
	}
	return -1
}
func (s *server) handlehashrequest() {
	for {
		payloadReceived := <- s.payloadChan

		id := payloadReceived.reqid;

		fmt.Printf("Getting hash for id: %d, %s\n", id, s.hashMap[id])

		payloadReceived.hash = s.hashMap[id]

		s.getHashChan  <- payloadReceived.hash
	}
}

func (s *server) gethash(w http.ResponseWriter, r *http.Request) {

	id := getidfromurl(r.URL.Path)
	if id == -1 {
		message := "ERROR: Invalid request number provided in the URL: " + r.URL.Path + "\n"
		w.Write([]byte(message))
		return
	}

	payload := payload{reqid: id, password: "", hash: ""}

	s.payloadChan <- payload

	passwordHash := <- s.getHashChan

	fmt.Printf("Received hash: %s\n", passwordHash)


	if passwordHash != "" {
		w.Write([]byte(passwordHash))
	} else {
		message := "ERROR: Password Hash for request id: " + strconv.Itoa(id) + " does not exist.\n"
		w.Write([]byte(message))
	}
}

func (s *server) handlepasswordhashing() {
	fmt.Printf("Starting handle password hashing")

	for {

		payloadReceived := <- s.hashChan

		fmt.Printf("Received id: %d, password: %s to hash\n", payloadReceived.reqid, payloadReceived.password)

		time.Sleep(5 * time.Second)

		// Calculating time to process requests.
		start := time.Now()
		t := time.Now()

		password := payloadReceived.password
		id := payloadReceived.reqid

		s.hashMap[id] = computehash(password)

		fmt.Printf("Hashed id: %d, password: %s to hash\n", id, s.hashMap[id])

		elapsed := t.Sub(start)

	//fmt.Printf("Added %s to map at: %d\n", password, num)
		s.totalTimeInNSec += elapsed.Nanoseconds()

	}	
}

func (s *server) getid() {
	fmt.Printf("Starting handle payload")

	for {

		s.totalRequests = s.totalRequests + 1;
		fmt.Printf("sending id: %d \n", s.totalRequests)
		s.getId  <- s.totalRequests;

	}	
}

func (s *server) hash(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	passwordFromForm := r.Form.Get("password")

	reqid := <- s.getId

	w.Write([]byte(strconv.Itoa(reqid)))

	payload := payload{reqid: reqid, password: passwordFromForm, hash: ""}

	s.hashChan <- payload
}

func (s *server) getnumberhashed() int {
	return len(s.hashMap)
}

func (s *server) constructjson() []byte {

	averageTime := 0.0


	if s.totalRequests > 0 {
		timeInMicroSeconds := s.totalTimeInNSec / 1000;
		averageTime = float64(timeInMicroSeconds) / float64(s.getnumberhashed())
	} 

	stats := &statsresponse{Total: s.totalRequests, Average: averageTime}

	b, _ := json.Marshal(stats)
	return b
}

func (s *server) stats(w http.ResponseWriter, r *http.Request) {
	b := s.constructjson()
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func computehash(password string) string {
	if password != "" {
		sha512 := sha512.New()
		sha512.Write([]byte(password))

		return base64.StdEncoding.EncodeToString(sha512.Sum(nil))
	}
	return ""
}

func (s *server) savetohashmap(num int, password string) {
	time.Sleep(5 * time.Second)

	// Calculating time to process requests.
	start := time.Now()
	t := time.Now()
	s.hashMap[num] = computehash(password)
	elapsed := t.Sub(start)

	//fmt.Printf("Added %s to map at: %d\n", password, num)
	s.totalTimeInNSec += elapsed.Nanoseconds()
}

func (s *server) shutdown(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(string(s.constructjson()))
	s.shutdownReq <- true
}
