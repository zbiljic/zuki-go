package server

import (
	"testing"
	"time"
)

func TestNewAppliesConfiguredTimeouts(t *testing.T) {
	config := Config{
		Addr:              "127.0.0.1:0",
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       4 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       6 * time.Second,
	}

	server := New(config)

	if got, want := server.Addr, config.Addr; got != want {
		t.Fatalf("got addr %q, want %q", got, want)
	}
	if got, want := server.ReadHeaderTimeout, config.ReadHeaderTimeout; got != want {
		t.Fatalf("got read header timeout %v, want %v", got, want)
	}
	if got, want := server.ReadTimeout, config.ReadTimeout; got != want {
		t.Fatalf("got read timeout %v, want %v", got, want)
	}
	if got, want := server.WriteTimeout, config.WriteTimeout; got != want {
		t.Fatalf("got write timeout %v, want %v", got, want)
	}
	if got, want := server.IdleTimeout, config.IdleTimeout; got != want {
		t.Fatalf("got idle timeout %v, want %v", got, want)
	}
}
