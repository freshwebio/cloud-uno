// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package config

import (
	"errors"

	"github.com/freshwebio/cloud-uno/pkg/types"
)

var (
	// ErrMissingOrInvalidConfigService provides the error to be used
	// when the config struct is not where it is expected to be in the service resolver.
	ErrMissingOrInvalidConfigService = errors.New("config service missing in resolver container or is of an unexpected type")
)

// RegisterServices deals with registering config-specific
// services to be used throughout the application.
func RegisterServices(r types.Resolver) (err error) {
	cfg, err := Load()
	if err != nil {
		return
	}
	r.Set("config", cfg)
	return nil
}
