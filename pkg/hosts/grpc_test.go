// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

//go:build unit

package hosts

import (
	context "context"
	"log"
	"net"
	"strings"
	"testing"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type GRPCSuite struct {
	listener   *bufconn.Listener
	server     *grpc.Server
	ipHostsMap map[string]string
}

var _ = Suite(&GRPCSuite{})

func (s *GRPCSuite) SetUpSuite(c *C) {
	s.listener = bufconn.Listen(1024 * 1024)
	s.server = grpc.NewServer()
	s.ipHostsMap = make(map[string]string)
	RegisterManagerServer(s.server, &GRPCServer{
		Impl: &mockManager{
			ipHostsMap: s.ipHostsMap,
		},
	})
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func (s *GRPCSuite) bufDialer(context.Context, string) (net.Conn, error) {
	return s.listener.Dial()
}

func (s *GRPCSuite) Test_add_hosts_for_ip(c *C) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	if err != nil {
		c.Error(err)
		c.FailNow()
	}
	defer conn.Close()
	client := NewManagerClient(conn)
	resp, err := client.Add(ctx, &HostsRequest{
		Ip:    "172.1.0.22",
		Hosts: "api.google.local,api.aws.local,api.azure.local,api.example.local",
	})
	if err != nil {
		c.Error(err)
		c.FailNow()
	}
	c.Assert(resp.GetApplied(), Equals, true)
	c.Assert(s.ipHostsMap["172.1.0.22"], Equals, "api.google.local,api.aws.local,api.azure.local,api.example.local")
}

func (s *GRPCSuite) Test_remove_hosts_from_ip(c *C) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(s.bufDialer), grpc.WithInsecure())
	if err != nil {
		c.Error(err)
		c.FailNow()
	}
	defer conn.Close()
	client := NewManagerClient(conn)
	_, err = client.Add(ctx, &HostsRequest{
		Ip:    "172.1.0.23",
		Hosts: "api.google.local,api.aws.local,api.azure.local,api.example.local",
	})
	if err != nil {
		c.Error(err)
		c.FailNow()
	}

	resp, err := client.Remove(ctx, &HostsRequest{
		Ip:    "172.1.0.23",
		Hosts: "api.google.local,api.azure.local",
	})
	c.Assert(resp.GetApplied(), Equals, true)
	c.Assert(s.ipHostsMap["172.1.0.23"], Equals, "api.aws.local,api.example.local")
}

type mockManager struct {
	ipHostsMap map[string]string
}

func (m *mockManager) Add(params *Params) error {
	m.ipHostsMap[*params.IP] = *params.Hosts
	return nil
}

func (m *mockManager) Remove(params *Params) error {
	hostsBefore := m.ipHostsMap[*params.IP]
	hostsBeforeList := strings.Split(hostsBefore, ",")
	hostsToRemoveList := strings.Split(*params.Hosts, ",")
	finalHosts := []string{}
	for _, host := range hostsBeforeList {
		if !itemInSlice(host, hostsToRemoveList) {
			finalHosts = append(finalHosts, host)
		}
	}
	m.ipHostsMap[*params.IP] = strings.Join(finalHosts, ",")
	return nil
}
