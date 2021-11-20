// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package grpc

import (
	"net"

	"github.com/freshwebio/cloud-uno/pkg/types"
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
