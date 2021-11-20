// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package connect

import (
	"fmt"

	"github.com/freshwebio/cloud-uno/pkg/utils"
)

// HostAgentPort provides the tcp port the host
// agent is running on.
const HostAgentPort = "5989"

// DeriveHostAgentAddr determines the correct
// address to connect to for the host agent.
func DeriveHostAgentAddr() string {
	if utils.IsRunningInDockerContainer() {
		return fmt.Sprintf("host.docker.internal:%s", HostAgentPort)
	}
	return fmt.Sprintf("127.0.0.1:%s", HostAgentPort)
}
