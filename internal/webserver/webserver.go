// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package webserver

import (
	"net/http"

	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/hosts"
	"github.com/freshwebio/cloud-uno/pkg/netutils"
	"github.com/freshwebio/cloud-uno/pkg/types"
	"github.com/gorilla/mux"
)

var (
	// WebServerHost specifies the host on which Cloud::1 UI will be served.
	WebServerHost = "console.clouduno.local"
)

// RegisterStatic deals with registering the web server to serve the Cloud Uno UI.
func RegisterStatic(router *mux.Router, resolver types.Resolver) (err error) {
	fileServer := http.FileServer(http.Dir("./client/build/"))
	router.PathPrefix("/").Handler(fileServer).Host(WebServerHost)
	// As the web server doesn't have a service in the same way cloud provider API emulators
	// do, we'll register the host for the web server here.
	hostsManager := resolver.Get("hosts").(hosts.Service)
	cfg := resolver.Get("config").(*config.Config)

	serverIP, err := netutils.SelectServerIP(cfg)
	if err != nil {
		return
	}

	err = hostsManager.Add(&hosts.Params{
		IP:    &serverIP,
		Hosts: &WebServerHost,
	})
	return
}
