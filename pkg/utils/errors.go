// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package utils

import "errors"

var (
	// ErrMissingOrInvalidLogger provides an error that should be used
	// upon failure to get the logger from a service resolver.
	ErrMissingOrInvalidLogger = errors.New("the logger service is either missing or is of an unexpected type")
)
