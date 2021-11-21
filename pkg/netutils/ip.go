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
	"os/exec"
	"runtime"

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
			return "", invalidIPError(ip)
		}
		return ip, nil
	}
	if *cfg.RunOnHost {
		return DefaultHostServerIP, nil
	}
	return DefaultContainerServerIP, nil
}

// CreateLoopBackAlias deals with creating an alias for the loopback
// address so that we can assign a virtual static IP to a docker container
// running the Cloud::1 server.
// This prevents cloud uno from having conflicts with any other servers you may
// have running on port 80 on your local machine and allows us to channel
// all cloud uno host names to a separate IP.
// This is ONLY supported for linux and darwin platforms!
func CreateLoopBackAlias(ip string) error {
	// Make sure the IP is safe to inject by making sure it's a valid IP.
	if net.ParseIP(ip) == nil {
		return invalidIPError(ip)
	}

	if runtime.GOOS == "darwin" {
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo ifconfig lo0 alias %s", ip))
		return cmd.Run()
	} else if runtime.GOOS == "linux" {
		// sudo ifconfig eth0:0 {ip} netmask 255.255.255.0 up
		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo ifconfig eth0:0 %s netmask 255.255.255.0 up", ip))
		return cmd.Run()
	}
	return nil
}

func invalidIPError(ip string) error {
	return fmt.Errorf(
		"invalid IP address %s provided for the ip the server is running on",
		ip,
	)
}
