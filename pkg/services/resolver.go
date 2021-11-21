// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package services

import (
	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/gcloud"
	"github.com/freshwebio/cloud-uno/pkg/hosts"
	"github.com/freshwebio/cloud-uno/pkg/types"
)

// NewDefaultResolver produces an instance of the built-in default
// service resolver.
func NewDefaultResolver() types.Resolver {
	return &defaultResolver{
		services: make(map[string]interface{}),
	}
}

type defaultResolver struct {
	services map[string]interface{}
}

func (r *defaultResolver) Get(service string) interface{} {
	return r.services[service]
}

func (r *defaultResolver) Register(customRegisterFunc func(r types.Resolver) error) (err error) {
	// Config should come before all other services as pretty
	// much every other service will make use of config.
	err = config.RegisterServices(r)
	if err != nil {
		return
	}

	err = customRegisterFunc(r)
	if err != nil {
		return
	}

	err = hosts.RegisterServices(r)
	if err != nil {
		return
	}

	err = gcloud.RegisterServices(r)

	return
}

func (r *defaultResolver) Set(name string, service interface{}) {
	r.services[name] = service
}
