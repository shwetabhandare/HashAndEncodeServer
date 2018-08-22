package main

import (
	"net/http"
)

func HashServer(options ...func(*server)) *server {
	s := &server{requestNum: 0, router: http.NewServeMux(), hashMap: make(map[int]string)}

	for _, f := range options {
		f(s)
	}

	s.router.HandleFunc("/hash", s.hash)
	s.router.HandleFunc("/hash/", s.gethash)
	s.router.HandleFunc("/stats", s.stats)
	s.router.HandleFunc("/shutdown", s.shutdown)

	return s
}
