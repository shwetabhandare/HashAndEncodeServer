package main

import (
	"net/http"
)

func HashServer(options ...func(*server)) *server {
	s := &server{totalRequests: 0, 
		totalTimeInNSec: 0.0, 
		router: http.NewServeMux(), 
		hashMap:        make(map[int]string),
		shutdownReq:    make(chan bool),
		passwordToHash: make(chan string),
		nextId: 	    make(chan int),
	}

	s.router.HandleFunc("/hash", s.hash)
	s.router.HandleFunc("/hash/", s.gethash)
	s.router.HandleFunc("/stats", s.stats)
	s.router.HandleFunc("/shutdown", s.shutdown)

	go s.getnextid()
	return s
}
