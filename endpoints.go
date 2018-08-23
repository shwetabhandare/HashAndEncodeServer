package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type statsresponse struct {
	Total   int
	Average float64
}

func (s *server) gethashfrommap(id string) string {

	if idAsInt, err := strconv.Atoi(id); err == nil {
		if passwordHash, ok := s.hashMap[idAsInt]; ok {
			return passwordHash
		}
	}
	return ""
}

func (s *server) gethash(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/hash/")
	fmt.Println("Get Hash for ID:", id)

	passwordHash := s.gethashfrommap(id)
	if passwordHash != "" {
		fmt.Println(passwordHash)
		w.Write([]byte(passwordHash))
	} else {
		message := "ERROR: Password Hash for request id: : " + id + " does not exist.\n"
		w.Write([]byte(message))
	}
}

func (s *server) hash(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	passwordFromForm := r.Form.Get("password")
	if passwordFromForm != "" {
		s.totalRequests++
		fmt.Println(passwordFromForm)

		go s.saveHashToMap(s.totalRequests, passwordFromForm)

		w.Write([]byte(strconv.Itoa(s.totalRequests)))

	} else {
		message := "ERROR: Unable to find password in the POST request. Found: " + r.Form.Encode() + " instead" + "\n"
		w.Write([]byte(message))
	}

}

func (s *server) constructjson() []byte {

	averageTime := 0.0
	if s.totalRequests > 0 {
		averageTime = float64(s.totalTimeInNSec) / float64(s.totalRequests)
	} else {
		averageTime = 0.0
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

func (s *server) saveHashToMap(num int, password string) {
	time.Sleep(5 * time.Second)

	start := time.Now()
	t := time.Now()
	s.hashMap[num] = computehash(password)
	elapsed := t.Sub(start)

	s.totalTimeInNSec += elapsed.Nanoseconds()
}

func (s *server) shutdown(w http.ResponseWriter, r *http.Request) {
	message := "Hello " + r.URL.Path + "\n"
	w.Write([]byte(message))
}
