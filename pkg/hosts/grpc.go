// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package hosts

import (
	context "context"
	"errors"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

var (
	// GRPCMinConnectTimeout provides the minimum timeout
	// for a connection to the hosts agent to complete.
	GRPCMinConnectTimeout = 10 * time.Second
)

// NewGRPCClient creates a client to connect to the hosts
// agent over gRPC.
func NewGRPCClient(hostAgentAddr string) (Service, error) {
	// Add some retry behaviour as hosts agent is not guaranteed
	// to be running as the server starts up.
	conn, err := grpc.Dial(
		hostAgentAddr,
		// This is designed for a developer's local machine, so insecure gRPC
		// is acceptable as long as this software is not used on a public network!
		grpc.WithInsecure(),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				MinConnectTimeout: 10 * time.Second,
				Backoff:           backoff.DefaultConfig,
			},
		),
	)
	if err != nil {
		return nil, err
	}
	mgrClient := NewManagerClient(conn)
	return &GRPCClient{
		client: mgrClient,
	}, nil
}

// GRPCClient is an implementation of a client that the server uses
// to communicate with the hosts agent.
type GRPCClient struct {
	client ManagerClient
}

var (
	// ErrFailedToAddHosts is returned when the hosts agent failed to
	// to add hosts due to an error that has been handled on the hosts agent
	// side.
	ErrFailedToAddHosts = errors.New("failed to add the provided hosts to the given IP")
	// ErrFailedToRemoveHosts is returned when the hosts agent failed to
	// to remove hosts from an IP due to an error that has been handled on the hosts agent
	// side.
	ErrFailedToRemoveHosts = errors.New("failed to remove the provided hosts from the given IP")
)

// Add deals with making a request to a gRPC server
// to add a list of hosts to a given IP for local
// DNS emulation.
func (m *GRPCClient) Add(params *Params) (err error) {
	response, err := m.client.Add(context.Background(), &HostsRequest{
		Ip:    *params.IP,
		Hosts: *params.Hosts,
	})
	if err != nil {
		return
	}
	if !response.GetApplied() {
		return ErrFailedToAddHosts
	}
	return nil
}

// Remove deals with making a request to a gRPC server
// to remove a list of hosts from a given IP for
// local DNS emulation.
func (m *GRPCClient) Remove(params *Params) (err error) {
	response, err := m.client.Remove(context.Background(), &HostsRequest{
		Ip:    *params.IP,
		Hosts: *params.Hosts,
	})
	if err != nil {
		return
	}
	if !response.GetApplied() {
		return ErrFailedToRemoveHosts
	}
	return nil
}

// GRPCServer provides the gRPC service running in the hosts agent.
type GRPCServer struct {
	// We must embed the unimplemented interface
	// for a hosts manager server.
	UnimplementedManagerServer
	// This is the real implementation.
	Impl Service
}

// Add deals with adding a list of hosts to a given IP for local
// DNS emulation.
func (m *GRPCServer) Add(ctx context.Context, req *HostsRequest) (*HostsResponse, error) {
	err := m.Impl.Add(&Params{
		IP:    &req.Ip,
		Hosts: &req.Hosts,
	})
	if err != nil {
		return nil, err
	}
	return &HostsResponse{Applied: true}, nil
}

// Remove deals with removing a list of hosts from a given IP for
// local DNS emulation.
func (m *GRPCServer) Remove(ctx context.Context, req *HostsRequest) (*HostsResponse, error) {
	err := m.Impl.Remove(&Params{
		IP:    &req.Ip,
		Hosts: &req.Hosts,
	})
	if err != nil {
		return nil, err
	}
	return &HostsResponse{Applied: true}, nil
}
