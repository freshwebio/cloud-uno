// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package main

import (
	"log"
	"net"
	"net/http"

	"github.com/freshwebio/cloud-uno/internal/coresvc"
	"github.com/freshwebio/cloud-uno/internal/gcloud/grpc"
	"github.com/freshwebio/cloud-uno/internal/gcloud/httpapi"
	"github.com/freshwebio/cloud-uno/internal/webserver"
	"github.com/freshwebio/cloud-uno/pkg/services"
	"github.com/freshwebio/cloud-uno/pkg/types"
	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

func httpServe(l net.Listener, resolver types.Resolver) error {
	mux := mux.NewRouter()
	httpapi.RegisterSecretManager(mux, resolver)
	err := webserver.RegisterStatic(mux, resolver)
	if err != nil {
		return err
	}

	s := &http.Server{Handler: mux}
	return s.Serve(l)
}

func main() {
	resolver := services.NewDefaultResolver()
	err := resolver.Register(coresvc.RegisterServices)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", ":5988")
	if err != nil {
		log.Fatal(err)
	}
	m := cmux.New(listener)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g := new(errgroup.Group)
	g.Go(func() error { return grpc.Serve(grpcListener, resolver) })
	g.Go(func() error { return httpServe(httpListener, resolver) })
	g.Go(func() error { return m.Serve() })

	log.Println("Running Cloud::1 Server on port 5988 ...")
	g.Wait()
}
