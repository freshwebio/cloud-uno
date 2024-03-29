// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

//go:build windows

package hosts

const (
	// HostsPerLine is the amount of hosts that can appear on a single line.
	HostsPerLine = 9
	// HostsFilePath provides the location to the hosts file on darwin and linux
	// systems.
	HostsFilePath = "${SystemRoot}/System32/drivers/etc/hosts"
	eol           = "\r\n"
)
