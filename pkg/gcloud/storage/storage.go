// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package storage

import "github.com/docker/docker/client"

// Storage provides a common interface for an implementation of a service
// that implements a backend for a Google cloud storage API emulation.
type Storage interface {
	BucketAccessControls() BucketAccessControls
	Buckets() Buckets
	Channels() Channels
	DefaultObjectAccessControls() DefaultObjectAccessControls
	Notifications() Notifications
	ObjectAccessControls() ObjectAccessControls
	Objects() Objects
	ProjectsHMACKeys() ProjectsHMACKeys
	ProjectsServiceAccounts() ProjectsServiceAccounts
}

// New creates an instance of a backend service
// for a google cloud storage emulator.
func New(dockerClient *client.Client) (Storage, error) {
	return nil, nil
}
