// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package gcloud

import (
	"fmt"
	"log"

	"github.com/docker/docker/client"
	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/freshwebio/cloud-uno/pkg/gcloud/grpc"
	"github.com/freshwebio/cloud-uno/pkg/gcloud/storage"
	"github.com/freshwebio/cloud-uno/pkg/hosts"
	"github.com/freshwebio/cloud-uno/pkg/netutils"
	"github.com/freshwebio/cloud-uno/pkg/types"
	"github.com/freshwebio/cloud-uno/pkg/utils"
	"github.com/spf13/afero"
	"google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	// GCloudSecretManagerName provides the name used to identify
	// the google cloud secret manager service.
	GCloudSecretManagerName = "secretmanager"
	// GCloudStorageName provides the name used to identify
	// the google cloud storage service.
	GCloudStorageName = "storage"
)

// RegisterServices deals with registering google cloud
// services to be used for handling gRPC and HTTP requests.
func RegisterServices(resolver types.Resolver) (err error) {
	cfg, ok := resolver.Get("config").(*config.Config)
	if !ok {
		return config.ErrMissingOrInvalidConfigService
	}

	hostsService, ok := resolver.Get("hosts").(hosts.Service)
	if !ok {
		return hosts.ErrMissingOrInvalidHostsService
	}

	fs := resolver.Get("fs").(afero.Fs)
	// if !ok {
	// 	return coresvc.ErrMissingOrInvalidFS
	// }

	// Configuration will always contain the user-provided server IP,
	// we need to do a little more work to make sure the correct server IP
	// is selected based on the way in which the application is being run
	// as well as providing some validation.
	serverIP, err := netutils.SelectServerIP(cfg)
	if err != nil {
		return
	}

	// Docker is used to orchestrate and manage both vendor-managed
	// emulators and open source software used as the backend for some cloud services.
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Given gRPC is a fantastic representation of a service that is usually
	// abstracted away from a REST API route handler, the default resolver will use
	// the gRPC services for Google Cloud APIs that support gRPC.
	if utils.CommaSeparatedListContains(*cfg.GCloudServices, GCloudSecretManagerName) {
		smRootDir := fmt.Sprintf("%s/gcloud/secretmanager", *cfg.DataDirectory)
		var secretmgr secretmanager.SecretManagerServiceServer
		secretmgr, err = grpc.NewSecretManager(smRootDir, fs, serverIP, hostsService)
		if err != nil {
			return
		}
		resolver.Set("gcloud.secretmanager", secretmgr)
	}

	if utils.CommaSeparatedListContains(*cfg.GCloudServices, GCloudStorageName) {
		var storageService storage.Storage
		storageService, err = storage.New(dockerClient)
		if err != nil {
			return
		}
		resolver.Set("gcloud.storage", storageService)
	}

	return
}
