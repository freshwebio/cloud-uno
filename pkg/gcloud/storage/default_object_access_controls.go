// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package storage

// DefaultObjectAccessControls represents a service
// that deals with managing default object access controls
// in a Google Cloud Storage API emulation.
type DefaultObjectAccessControls interface {
	DeleteDefaultObjectAccessControl()
	GetDefaultObjectAccessControl()
	CreateDefaultObjectAccessControl()
	ListDefaultObjectAccessControls()
}
