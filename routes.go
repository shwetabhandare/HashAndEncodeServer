package main

import (
	"net/http"
	"sync"
)

func HashServer(options ...func(*server)) *server {
	s := &server{totalRequests: 0, 
		totalTimeInNSec: 0.0, 
		router: http.NewServeMux(), 
		hashMap: make(map[int]string),
		shutdownReq: make(chan bool),
		lock: new(sync.Mutex),
	}

	s.router.HandleFunc("/hash", s.hash)
	s.router.HandleFunc("/hash/", s.gethash)
	s.router.HandleFunc("/stats", s.stats)
	s.router.HandleFunc("/shutdown", s.shutdown)

	return s
}
