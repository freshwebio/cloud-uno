// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package coresvc

import (
	"errors"

	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/types"
	"github.com/spf13/afero"
)

var (
	// ErrMissingOrInvalidFS provides the error to be used
	// when the file system service is not where it is expected to be in the service resolver.
	ErrMissingOrInvalidFS = errors.New("file system service missing in resolver container or is of an unexpected type")
)

// RegisterServices deals with registering core
// services to be used throughout the application.
func RegisterServices(resolver types.Resolver) error {
	cfg, ok := resolver.Get("config").(*config.Config)
	if !ok {
		return config.ErrMissingOrInvalidConfigService
	}

	// The file system is used for services implemented directly in Cloud::1,
	// when handing off to other services like google cloud emulators we have no
	// control over whether that uses the file system or not.
	var fs afero.Fs
	if *cfg.FileSystem == "memory" {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}
	resolver.Set("fs", fs)
	return nil
}
