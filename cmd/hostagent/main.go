// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/connect"
	"github.com/freshwebio/cloud-uno/pkg/hosts"
	"github.com/freshwebio/cloud-uno/pkg/logging"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.CreateLogger()

	grpcServer := grpc.NewServer()
	managerImpl, err := hosts.NewManager(cfg, logger)
	if err != nil {
		log.Fatal("Create hosts manager error: ", err)
	}
	hosts.RegisterManagerServer(
		grpcServer,
		&hosts.GRPCServer{
			Impl: managerImpl,
		},
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", connect.HostAgentPort))
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	fmt.Println("Serving Cloud::1 Host Agent on port 5989 ...")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Serve error: ", err)
	}
}
