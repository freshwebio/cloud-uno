// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package types

// Resolver provides an interface for components that want
// to provide service resolution for Cloud::1.
type Resolver interface {
	// Retrieve a service.
	Get(service string) interface{}
	// Set a service in the resolver "container",
	// should be used by packages to register their services
	// for an application.
	Set(name string, service interface{})
	// Registers all services, should be called during initialisation
	// of an application/web server.
	Register(customRegisterFunc func(r Resolver) error) error
}
