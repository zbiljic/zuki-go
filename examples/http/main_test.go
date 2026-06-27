package main

import "testing"

func TestNewServerAppliesDefaultTimeouts(t *testing.T) {
	server := newServer(HTTPConfig{Addr: "127.0.0.1:0"})

	if got, want := server.Addr, "127.0.0.1:0"; got != want {
		t.Fatalf("got addr %q, want %q", got, want)
	}
	if got, want := server.ReadHeaderTimeout, defaultReadHeaderTimeout; got != want {
		t.Fatalf("got read header timeout %v, want %v", got, want)
	}
	if got, want := server.ReadTimeout, defaultReadTimeout; got != want {
		t.Fatalf("got read timeout %v, want %v", got, want)
	}
	if got, want := server.WriteTimeout, defaultWriteTimeout; got != want {
		t.Fatalf("got write timeout %v, want %v", got, want)
	}
	if got, want := server.IdleTimeout, defaultIdleTimeout; got != want {
		t.Fatalf("got idle timeout %v, want %v", got, want)
	}
}
