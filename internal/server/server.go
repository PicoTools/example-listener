package server

import (
	"context"
	"crypto/tls"
	"example_listener/internal/config"

	listener "github.com/PicoTools/pico-shared/proto/gen/listener/v1"
	"github.com/PicoTools/pico-shared/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	svc  listener.ListenerServiceClient
	conn *grpc.ClientConn
	ctx  context.Context
}

func New() (*Server, error) {
	ctx := context.Background()

	conn, err := grpc.NewClient(
		config.ServerAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
		grpc.WithUnaryInterceptor(unaryClientInterceptor(config.Token)),
		grpc.WithStreamInterceptor(streamClientInterceptor(config.Token)),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(shared.MaxProtobufMessageSize),
			grpc.MaxCallSendMsgSize(shared.MaxProtobufMessageSize),
		),
	)
	if err != nil {
		return nil, err
	}

	svc := listener.NewListenerServiceClient(conn)

	return &Server{
		svc:  svc,
		conn: conn,
		ctx:  ctx,
	}, nil
}
