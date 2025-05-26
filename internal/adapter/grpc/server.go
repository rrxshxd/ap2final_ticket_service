package grpc

import (
	"fmt"
	grpccfg "github.com/sorawaslocked/ap2final_base/pkg/grpc"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/ticket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

type Server struct {
	s             *grpc.Server
	cfg           grpccfg.Config
	addr          string
	log           *slog.Logger
	ticketUseCase TicketUseCase
}

func New(
	cfg grpccfg.Config,
	log *slog.Logger,
	ticketUseCase TicketUseCase,
) *Server {
	server := &Server{
		cfg:           cfg,
		addr:          fmt.Sprintf(":%d", cfg.Port),
		log:           log,
		ticketUseCase: ticketUseCase,
	}

	server.register()

	return server
}

func (s *Server) MustRun() {
	go func() {
		if err := s.run(); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) Stop() {
	s.log.Info("stopping grpc server", slog.String("addr", s.addr))

	s.s.GracefulStop()
}

func (s *Server) register() {
	s.s = grpc.NewServer()

	svc.RegisterTicketServiceServer(s.s, NewTicketServer(s.ticketUseCase, s.log))

	reflection.Register(s.s)
}

func (s *Server) run() error {
	const op = "grpc.run"

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("starting grpc server", slog.String("addr", s.addr))

	if err := s.s.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
