package main

import (
	"log"
	"net"
	"net/http"

	"github.com/freshwebio/cloud-one/internal/gcloud/grpc"
	"github.com/freshwebio/cloud-one/internal/gcloud/httpapi"
	"github.com/freshwebio/cloud-one/pkg/services"
	"github.com/freshwebio/cloud-one/pkg/types"
	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
)

func httpServe(l net.Listener, resolver types.Resolver) error {
	mux := mux.NewRouter()
	httpapi.RegisterSecretManager(mux, resolver)

	s := &http.Server{Handler: mux}
	return s.Serve(l)
}

func main() {
	resolver := services.NewDefaultResolver()
	err := resolver.Register()
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
