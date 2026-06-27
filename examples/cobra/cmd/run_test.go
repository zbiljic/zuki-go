package cmd

import (
	"testing"
	"time"

	"github.com/go-toho/contrib/config/vipero"

	"github.com/zbiljic/zuki-go/examples/cobra/internal/app"
)

func TestRunCommandBindsTimeoutFlags(t *testing.T) {
	opts := &app.Options{Viper: vipero.New(AppName)}
	cmd := newRunCommand(AppName, opts)

	if err := cmd.Flags().Set("read-header-timeout", "3s"); err != nil {
		t.Fatalf("set read header timeout: %v", err)
	}
	if err := cmd.Flags().Set("read-timeout", "4s"); err != nil {
		t.Fatalf("set read timeout: %v", err)
	}
	if err := cmd.Flags().Set("write-timeout", "5s"); err != nil {
		t.Fatalf("set write timeout: %v", err)
	}
	if err := cmd.Flags().Set("idle-timeout", "6s"); err != nil {
		t.Fatalf("set idle timeout: %v", err)
	}

	if got, want := opts.Viper.GetDuration("http.read_header_timeout"), 3*time.Second; got != want {
		t.Fatalf("got read header timeout %v, want %v", got, want)
	}
	if got, want := opts.Viper.GetDuration("http.read_timeout"), 4*time.Second; got != want {
		t.Fatalf("got read timeout %v, want %v", got, want)
	}
	if got, want := opts.Viper.GetDuration("http.write_timeout"), 5*time.Second; got != want {
		t.Fatalf("got write timeout %v, want %v", got, want)
	}
	if got, want := opts.Viper.GetDuration("http.idle_timeout"), 6*time.Second; got != want {
		t.Fatalf("got idle timeout %v, want %v", got, want)
	}
}
