package main

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/mcpherrinm/teyes/model"
)

func main() {
	host, ok := os.LookupEnv("LISTEN_HOST")
	if !ok {
		host = "0.0.0.0"
	}
	port, ok := os.LookupEnv("LISTEN_PORT")
	if !ok {
		port = "2222"
	}

	key, ok := os.LookupEnv("SSH_PRIVATE_KEY")
	if !ok {
		log.Fatal("SSH_PRIVATE_KEY not set")
	}

	// Wish wants the host key in PEM, so we need to convert
	rawKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Fatalf("Decoding private key: %s", err)
	}

	// Encode rawKey to PEM:
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: rawKey,
	})

	s, err := wish.NewServer(
		wish.WithHostKeyPEM(keyPem),
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware()))
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return model.Model{}, model.Options
}
