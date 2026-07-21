package grpc

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nhassl3/hairdress_arz/internal/service"
	"github.com/nhassl3/hairdress_arz/internal/transport/grpc/interceptors"
	authv1 "github.com/nhassl3/hairdress_arz_52_contracts/pkg/pb/auth/v1"
	bookingv1 "github.com/nhassl3/hairdress_arz_52_contracts/pkg/pb/booking/v1"
	"github.com/nhassl3/servicehub-backend/pkg/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type Services struct {
	Auth    *service.AuthService
	Booking *service.BookingService
}

type Handlers struct {
	Auth    *AuthHandler
	Booking *BookingHandler
}

type Server struct {
	grpcServer *grpc.Server
	handlers   *Handlers
	logger     *zap.Logger
}

func NewServer(services *Services, tokenManager auth.TokenManager, log *zap.Logger) *Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RecoveryInterceptor(log),
			interceptors.LoggingInterceptor(log),
			interceptors.AuthInterceptor(tokenManager),
		),
	)

	handlers := &Handlers{
		Auth:    NewAuthHandler(services.Auth),
		Booking: NewBookingHandler(services.Booking),
	}

	registerHandlers(grpcServer, handlers)
	reflection.Register(grpcServer)

	return &Server{grpcServer, handlers, log}
}

// registerHandlers registers every service implementation with the gRPC server.
// To add a new service: implement its handler, add it to Handlers, and call
// the generated Register<Name>Server here.
func registerHandlers(srv *grpc.Server, h *Handlers) {
	authv1.RegisterAuthServiceServer(srv, h.Auth)
	bookingv1.RegisterBookingServiceServer(srv, h.Booking)
}

// Start begins accepting gRPC connections on addr (e.g. ":9090").
func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.logger.Info("gRPC server listening", zap.String("addr", addr))
	return s.grpcServer.Serve(lis)
}

// StartGateway starts the HTTP/JSON REST gateway that proxies to the local
// gRPC server. The gateway is built with grpc-gateway and maps every
// google.api.http annotation in the proto files to an HTTP endpoint.
//
// grpcAddr must be the same address the gRPC server is listening on.
func (s *Server) StartGateway(ctx context.Context, grpcAddr, httpAddr string) error {
	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// Register every service handler with the gateway mux.
	for _, fn := range []func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		authv1.RegisterAuthServiceHandler,
		bookingv1.RegisterBookingServiceHandler,
	} {
		if err := fn(ctx, mux, conn); err != nil {
			return err
		}
	}

	s.logger.Info("HTTP gateway listening", zap.String("addr", httpAddr))
	return http.ListenAndServe(httpAddr, corsMiddleware(mux))
}

// corsMiddleware adds CORS headers so that the React dev server at
// localhost:5173 (and any origin listed in allowedOrigins) can reach
// the HTTP gateway.
func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := map[string]struct{}{
		"http://localhost:5173": {},
		"http://localhost:3000": {},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		}

		// Handle preflight requests.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Shutdown gracefully drains in-flight RPCs and stops the server.
func (s *Server) Shutdown(_ context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}
