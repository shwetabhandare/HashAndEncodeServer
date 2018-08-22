package main

import (
	"os"
	"testing"
)

func TestGetAddrValueNotSet(t *testing.T) {
	addr := getaddr()
	if addr != ":8080" {
		t.Errorf("Addr was incorrect, got: %s, want: %s.", addr, ":8080")
	}
}

func TestGetAddrValueSet(t *testing.T) {
	os.Setenv("PORT", "8081")
	addr := getaddr()
	if addr != ":8081" {
		t.Errorf("Sum was incorrect, got: %s, want: %s.", addr, ":8081")
	}
}
