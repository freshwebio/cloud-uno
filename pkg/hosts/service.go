// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package hosts

// HostParams provides the parameters required
// to add or remove a list of hosts to an IP.
type HostsParams struct {
	IP    *string
	Hosts *string
}

// Service provides a common interface for a service that manages
// os hosts, can be a server or a client.
type Service interface {
	Add(params *HostsParams) error
	Remove(params *HostsParams) error
}
