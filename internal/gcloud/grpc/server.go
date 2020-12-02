package grpc

import (
	"net"

	"github.com/freshwebio/cloud-one/pkg/types"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc"
)

// Serve deals with serving all the gRPC servers for the subset of google cloud services
// implemented with gRPC.
func Serve(l net.Listener, resolver types.Resolver) error {
	s := grpc.NewServer()
	secretmanagerpb.RegisterSecretManagerServiceServer(
		s,
		resolver.Get("gcloud.secretmanager").(secretmanagerpb.SecretManagerServiceServer),
	)
	return s.Serve(l)
}
