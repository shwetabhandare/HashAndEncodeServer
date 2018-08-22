package main

import (
	"fmt"
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
	expectedHash := "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg=="
	passwordHash := computehash(password)
	fmt.Println(passwordHash)

	if passwordHash != expectedHash {
		t.Errorf("Password hash was incorrect, got: %s, want: %s", passwordHash, expectedHash)
	}
}

func TestGetHashFromMapValidID(t *testing.T) {

	s := &server{requestNum: 0, router: http.NewServeMux(), hashMap: make(map[int]string)}}

	s.hashMap[requestNum] = "hash1"

}
