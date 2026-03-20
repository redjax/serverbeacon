package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/redjax/serverbeacon/api/v1"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(port string) *HttpServer {
	// Add v1 API handler
	v1_handler := v1.New()

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      v1_handler,
	}

	serv := &HttpServer{
		server: server,
	}

	return serv
}

func (s *HttpServer) ListenMulti(host string) error {
	// Listen on all IPv4 interfaces
	listener, err := net.Listen("tcp4", host+":"+s.server.Addr[1:])
	if err != nil {
		return fmt.Errorf("HTTP server listen failed: %w", err)
	}

	fmt.Printf("Server listening on %s\n", listener.Addr().String())

	// Graceful shutdown
	go func() {
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
