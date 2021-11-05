package services

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

const (
	// GCloudSecretManagerName provides the name used to identify
	// the google cloud secret manager service.
	GCloudSecretManagerName = "secretmanager"
	// GCloudStorageName provides the name used to identify
	// the google cloud storage service.
	GCloudStorageName = "storage"
)

func (r *defaultResolver) Register() (err error) {
	cfg, err := config.Load()
	if err != nil {
		return
	}
	r.services["config"] = cfg

	// The file system is used for services implemented directly in Cloud::1,
	// when handing off to other services like google cloud emulators we have no
	// control over whether that uses the file system or not.
	var fs afero.Fs
	if *cfg.FileSystem == "memory" {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}
	r.services["fs"] = fs

	var hostsService hosts.Service
	// If the cloud::1 server is running on the host machine then there
	// is no need for communicating with a separate process over a unix socket,
	// with the right permissions it can interact with the os hosts file directly.
	if *cfg.RunOnHost {
		hostsService, err = hosts.NewManager(cfg)
		if err != nil {
			return
		}
	} else {
		hostsService, err = hosts.NewClient()
		if err != nil {
			return
		}
	}

	r.services["hosts"] = hostsService

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
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	// Given gRPC is a fantastic representation of a service that is usually
	// abstracted away from a REST API route handler, the default resolver will use
	// the gRPC services for Google Cloud APIs that support gRPC.
	if utils.CommaSeparatedListContains(*cfg.GCloudServices, GCloudSecretManagerName) {
		smRootDir := fmt.Sprintf("%s/gcloud/secretmanager", *cfg.DataDirectory)
		r.services["gcloud.secretmanager"], err = grpc.NewSecretManager(smRootDir, fs, serverIP, hostsService)
	}

	if utils.CommaSeparatedListContains(*cfg.GCloudServices, GCloudStorageName) {
		r.services["gcloud.storage"], err = storage.New(dockerClient)
	}
	return
}
