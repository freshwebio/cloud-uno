// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package netutils

import (
	"fmt"
	"net"

	"github.com/freshwebio/cloud-uno/pkg/config"
)

// DefaultContainerServerIP provides the default ip address
// that the Cloud::1 server is behind when running in a container.
const DefaultContainerServerIP = "172.18.0.22"

// DefaultHostServerIP provides the default ip address
// the the Cloud::1 server is behind when running as a process directly
// on a host machine.
const DefaultHostServerIP = "127.0.0.1"

// SelectServerIP deals with selecting the correct IP the Cloud::1
// server is running on for the current environment.
func SelectServerIP(cfg *config.Config) (string, error) {
	if *cfg.ServerIP != DefaultContainerServerIP {
		ip := *cfg.ServerIP
		// Only validate IP when a custom IP has been provided.
		if net.ParseIP(ip) == nil {
			return "", fmt.Errorf(
				"Invalid IP address %s provided for the ip the server is running on",
				ip,
			)
		}
		return ip, nil
	}
	if *cfg.RunOnHost {
		return DefaultHostServerIP, nil
	}
	return DefaultContainerServerIP, nil
}
