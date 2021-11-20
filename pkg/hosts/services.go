// Copyright (c) 2022 FRESHWEB LTD. and Contributors
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package hosts

import (
	"errors"

	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/connect"
	"github.com/freshwebio/cloud-uno/pkg/types"
	"github.com/freshwebio/cloud-uno/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	// ErrMissingOrInvalidHostsService provides the error to be used
	// when the hosts service is not where it is expected to be in the service resolver.
	ErrMissingOrInvalidHostsService = errors.New("hosts service missing in resolver container or is of an unexpected type")
)

// RegisterServices deals with registering host-specific
// services to be used throughout the application.
func RegisterServices(resolver types.Resolver) (err error) {
	cfg, ok := resolver.Get("config").(*config.Config)
	if !ok {
		return config.ErrMissingOrInvalidConfigService
	}

	logger, ok := resolver.Get("logger").(*logrus.Entry)
	if !ok {
		return utils.ErrMissingOrInvalidLogger
	}

	var hostsService Service
	// If the cloud::1 server is running directly on the host machine then there
	// is no need for communicating with a separate process with access to the os hosts file,
	// with the right permissions it can interact with the os hosts file directly.
	if *cfg.RunOnHost {
		hostsService, err = NewManager(cfg, logger)
		if err != nil {
			return
		}
	} else {
		hostsService, err = NewGRPCClient(connect.DeriveHostAgentAddr())
		if err != nil {
			return
		}
	}

	resolver.Set("hosts", hostsService)
	return
}
