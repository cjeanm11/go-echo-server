package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"server-template/internal/database"
	"strconv"
	"sync"
	"time"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc"
	"server-template/pkg/util"
)

var (
	portNumber = util.GetEnvOrDefault("PORT","8080")
)

type FixedSizeKey [4]byte

type Server struct {
	port       int
	domain     string
	db         *database.Service
	lock       sync.Mutex
	lockedKeys map[FixedSizeKey]struct{}
	httpServer *http.Server
	grpcServer *grpc.Server
}


type Option func(*Server)

func NewServer(options ...Option) *Server {
	s := &Server{
		port:       loadPortFromEnv(),
		lock:       sync.Mutex{},
		lockedKeys: map[FixedSizeKey]struct{}{},
		db : database.New(),
	}
	for _, option := range options {
		option(s)
	}

	s.initHTTPServer()
	s.initGRPCServer() 

	return s
}

func (s *Server) initHTTPServer() {

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) initGRPCServer() {
	var opts []grpc.ServerOption

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
	err = s.httpServer.ListenAndServe()

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
	port, err := strconv.Atoi(portNumber)
	if err != nil {
		log.Printf("Warning: Invalid PORT environment variable '%s', falling back to default port 8080.\n", portNumber)
		return 8080
	}
	return port
}


func (a *Server) UnlockKey(key []byte) {
	a.lock.Lock()
	defer a.lock.Unlock()

	var keyToCheck FixedSizeKey
	copy(keyToCheck[:], key)
	delete(a.lockedKeys, keyToCheck)
}

func (a *Server) LockKey(key []byte) bool {
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

func WithDomain(domain string) Option {
	return func(s *Server) {
		s.domain = domain
	}
}

