package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/SSam0419/ptlbuilder/protocol"
)

type Server struct {
	listener net.Listener
	wg       sync.WaitGroup
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.wg.Done()
	}()

	for {
		msg, err := protocol.DecodeMessageFromConn(conn)
		if err != nil {
			fmt.Printf("Connection read error: %v\n", err)
			return
		}

		if m, err := msg.AsRegisterClient(); err == nil {
			fmt.Printf("Received Register Client CMD [ADDR  %s] [CONTENT %s] \n", m.Address, m.Content)
		}
		if m, err := msg.AsSendMessage(); err == nil {
			fmt.Printf("Received Send Message CMD [ ADDR  %s] [CONTENT %s] \n", m.Address, m.Content)
		}
	}
}

func (s *Server) Run(ctx context.Context) error {
	defer s.listener.Close()

	// Channel for accept errors
	errChan := make(chan error, 1)

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				errChan <- err
				return
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}()

	// Wait for context cancellation or accept error
	select {
	case <-ctx.Done():
		fmt.Println("Server shutting down gracefully...")
		s.listener.Close()
		s.wg.Wait()
		return nil
	case err := <-errChan:
		return fmt.Errorf("accept error: %w", err)
	}
}

func main() {
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create server
	server, err := NewServer(":1000")
	if err != nil {
		panic(err)
	}

	// Start server in goroutine
	go func() {
		if err := server.Run(ctx); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	client, err := net.Dial("tcp", ":1000")
	if err != nil {
		fmt.Println("Failed to dial to server : ", err)

	}

	msg, err := protocol.EncodeRegisterClientRequest("test", "another test")
	if err != nil {
		fmt.Println("failed to encode msg, ", err)
	} else {
		client.Write(msg)
	}

	msg, err = protocol.EncodeSendMessageRequest("test", "another test")
	if err != nil {
		fmt.Println("failed to encode msg, ", err)
	} else {
		client.Write(msg)
	}

	// Simulate running for 5 seconds
	time.Sleep(5 * time.Second)

	// Trigger graceful shutdown
	cancel()

	// Give some time for cleanup
	time.Sleep(time.Second)
	fmt.Println("Main process exiting")
}
