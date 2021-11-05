package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/connect"
	"github.com/freshwebio/cloud-uno/pkg/hosts"
)

func main() {
	if err := os.RemoveAll(connect.SockAddr); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	hostsManager, err := hosts.NewManager(cfg)
	if err != nil {
		log.Fatal(err)
	}

	rpc.Register(hostsManager)
	rpc.HandleHTTP()
	listener, err := net.Listen("unix", connect.SockAddr)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	fmt.Println("Serving Cloud::1 Host Agent ...")
	http.Serve(listener, nil)
}
