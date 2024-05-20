package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"server-template/internal/database"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	certFile string = ".cert/server.crt"
	keyFile  string = ".cert/server.key"
)

type FixedSizeKey [4]byte

type Server struct {
	port       int
	db         database.Service
	lock       sync.Mutex
	lockedKeys map[FixedSizeKey]struct{}
	uploadids  map[string]bool
	httpServer *http.Server
	grpcServer *grpc.Server
	useTLS 	   bool
}

type Option func(*Server)

func NewServer(options ...Option) *Server {
	s := &Server{
		port:       loadPortFromEnv(),
		lock:       sync.Mutex{},
		lockedKeys: map[FixedSizeKey]struct{}{},
		useTLS: 	false,
	}
	for _, option := range options {
		option(s)
	}

	s.initHTTPServer()
	s.initGRPCServer() // Initialize gRPC server

	return s
}

func (s *Server) initHTTPServer() {
	var cred *tls.Config = nil
	
	if s.useTLS {
		cred = loadTLSCertificate(certFile, keyFile)
	}

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(), 
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		TLSConfig:    cred,
	}

}


func loadTLSCertificate(certFile, keyFile string) (*tls.Config) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        log.Fatalf("failed to load server certificate: %v", err)
    }

    certPool := x509.NewCertPool()
    certBytes, err := os.ReadFile(certFile)
    if err != nil {
        log.Fatalf("failed to read server certificate: %v", err)
    }
    certPool.AppendCertsFromPEM(certBytes)

    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientAuth:   tls.NoClientCert,
        ClientCAs:    certPool,
    }
}

func (s *Server) initGRPCServer() {
	var opts []grpc.ServerOption

    if s.useTLS {
		creds := loadTLSCertificate(certFile, keyFile)
        opts = append(opts, grpc.Creds(credentials.NewTLS(creds)))
    } 

	opts = append(opts, grpc.ConnectionTimeout(time.Second))
	opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: time.Second * 10,
		Timeout:           time.Second * 20,
	}))
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(
		keepalive.EnforcementPolicy{
				MinTime:             time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	opts = append(opts, grpc.MaxConcurrentStreams(5))

	s.grpcServer = grpc.NewServer(opts...)
}

func (s *Server) startGRPCServer(wg *sync.WaitGroup) {
	defer wg.Done()
	grpcAddr := fmt.Sprintf(":%d", s.port+997)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	defer grpcListener.Close()

	log.Printf("gRPC server starting on port %d\n", s.port+997)
	if err := s.grpcServer.Serve(grpcListener); err != nil {
		log.Fatalf("gRPC server failed to start: %v", err)
	}
}

func (s *Server) startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("HTTP server starting on port %d\n", s.port)

	var err error
	if s.useTLS {
		err = s.httpServer.ListenAndServeTLS("", "")
	} else {
		err = s.httpServer.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed to start on %s: %v", addr, err)
	}
}

func (s *Server) Start() {
	fmt.Println("start servers...")
	var wg sync.WaitGroup
	wg.Add(2)
	{
		go s.startHTTPServer(&wg)
		go s.startGRPCServer(&wg)
	}
	wg.Wait()
}

func loadPortFromEnv() int {
	portStr, exists := os.LookupEnv("PORT")
	if !exists {
		return 8080
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Warning: Invalid PORT environment variable '%s', falling back to default port 8080.\n", portStr)
		return 8080
	}
	return port
}

func GetPortOrDefault(defaultPort int) int {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return defaultPort
	}
	return port
}

func (a *Server) unlockKey(key []byte) {
	a.lock.Lock()
	defer a.lock.Unlock()

	var keyToCheck FixedSizeKey
	copy(keyToCheck[:], key) 
	delete(a.lockedKeys, keyToCheck)
}

func (a *Server) lockKey(key []byte) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	var keyToCheck FixedSizeKey
	copy(keyToCheck[:], key) 

	_, ok := a.lockedKeys[keyToCheck]
	if ok {
		return false 
	}
	a.lockedKeys[keyToCheck] = struct{}{}
	return true
}

// Options
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithTSL(useTSL bool) Option {
	return func(s *Server) {
		s.useTLS = useTSL
	}
}