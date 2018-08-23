package main

import (
	"net/http"
	"testing"
)

func TestComputeHashValidPassword(t *testing.T) {
	password := "angryMonkey"
	expectedHash := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	passwordHash := computehash(password)

	if passwordHash != expectedHash {
		t.Errorf("Password hash was incorrect, got: %s, want: %s", passwordHash, expectedHash)
	}
}

func TestComputeHashEmptyPassword(t *testing.T) {
	password := ""
	expectedHash := ""
	passwordHash := computehash(password)

	if passwordHash != expectedHash {
		t.Errorf("Password hash was incorrect, got: %s, want: %s", passwordHash, expectedHash)
	}
}

func TestGetHashFromMapValueExists(t *testing.T) {

	s := &server{totalRequests: 0, totalTimeInNSec: 0.0, router: http.NewServeMux(), hashMap: make(map[int]string)}
	s.hashMap[0] = "hash1"
	s.hashMap[1] = "hash1"

	hashValue := s.gethashfrommap("0")
	expectedHashValue := "hash1"

	if hashValue != expectedHashValue {
		t.Errorf("Password hash value incorrect, got: %s, want: %s", hashValue, expectedHashValue)
	}
}

func TestGetHashFromMapValueDoesNotExist(t *testing.T) {

	s := &server{totalRequests: 0, totalTimeInNSec: 0.0, router: http.NewServeMux(), hashMap: make(map[int]string)}
	s.hashMap[0] = "hash1"
	s.hashMap[1] = "hash2"

	hashValue := s.gethashfrommap("2")
	expectedHashValue := ""

	if hashValue != expectedHashValue {
		t.Errorf("Password hash value incorrect, got: %s, want: %s", hashValue, expectedHashValue)
	}
}

func TestSaveHashToMapValidPasswordValue(t *testing.T) {
	s := &server{totalRequests: 0, totalTimeInNSec: 0.0, router: http.NewServeMux(), hashMap: make(map[int]string)}

	s.saveHashToMap(1, "hash1")

	hashValue := s.gethashfrommap("1")
	expectedHashValue := "mNvnDYMOdqw+tUjPFe1oyGw3soU9+Rm5evpEdJHyqe1la+Uw6uB3ylEkrRVWElrdCwnJ1ejPIXCd2i6LuGeCYA=="

	if hashValue != expectedHashValue {
		t.Errorf("Password hash value incorrect, got: %s, want: %s", hashValue, expectedHashValue)
	}
}

func TestConstructJSONNonZeroValues(t *testing.T) {

	s := &server{totalRequests: 5, totalTimeInNSec: 100.0, router: http.NewServeMux(), hashMap: make(map[int]string)}
	s.hashMap[0] = "hash1"
	s.hashMap[1] = "hash2"

	jsonReply := s.constructjson()

	expectedJSONReply := "{\"Total\":5,\"Average\":20}"

	if string(jsonReply) != expectedJSONReply {
		t.Errorf("Password hash value incorrect, got: %s, want: %s", jsonReply, expectedJSONReply)
	}
}
func TestConstructJSONZeroValues(t *testing.T) {

	s := &server{totalRequests: 0, totalTimeInNSec: 0.0, router: http.NewServeMux(), hashMap: make(map[int]string)}
	s.hashMap[0] = "hash1"
	s.hashMap[1] = "hash2"

	jsonReply := s.constructjson()

	expectedJSONReply := "{\"Total\":0,\"Average\":0}"

	if string(jsonReply) != expectedJSONReply {
		t.Errorf("Password hash value incorrect, got: %s, want: %s", jsonReply, expectedJSONReply)
	}
}
