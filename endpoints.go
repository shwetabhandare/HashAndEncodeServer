package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

	r.ParseForm()               // Parses the request body
	x := r.Form.Get("password") // x will be "" if parameter is not set
	if x != "" {
		s.requestNum++
		fmt.Println(x)

		go s.saveHashToMap(s.requestNum, x)

		w.Write([]byte(strconv.Itoa(s.requestNum)))

	} else {
		message := "ERROR: Unable to find password in the POST request. Found: " + r.Form.Encode() + " instead" + "\n"
		w.Write([]byte(message))
	}

}

func (s *server) stats(w http.ResponseWriter, r *http.Request) {
	message := "Hello " + r.URL.Path + "\n"
	w.Write([]byte(message))
}

func computehash(password string) string {
	sha512 := sha512.New()
	sha512.Write([]byte(password))

	return base64.StdEncoding.EncodeToString(sha512.Sum(nil))
}

func (s *server) saveHashToMap(num int, password string) {
	time.Sleep(5 * time.Second)

	s.hashMap[num] = computehash(password)

	fmt.Printf("Added %s to map at %d\n", s.hashMap[num], num)
}

func (s *server) shutdown(w http.ResponseWriter, r *http.Request) {
	message := "Hello " + r.URL.Path + "\n"
	w.Write([]byte(message))
}
