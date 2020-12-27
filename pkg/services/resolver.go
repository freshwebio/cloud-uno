package services

import (
	"fmt"

	"github.com/freshwebio/cloud-one/pkg/config"
	"github.com/freshwebio/cloud-one/pkg/gcloud/grpc"
	"github.com/freshwebio/cloud-one/pkg/hosts"
	"github.com/freshwebio/cloud-one/pkg/netutils"
	"github.com/freshwebio/cloud-one/pkg/types"
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
		hostsService, err = hosts.NewManager()
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

	// Given gRPC is a fantastic representation of a service that is usually
	// abstracted away from a REST API route handler, the default resolver will use
	// the gRPC services.
	smRootDir := fmt.Sprintf("%s/gcloud/secretmanager", *cfg.DataDirectory)
	r.services["gcloud.secretmanager"], err = grpc.NewSecretManager(smRootDir, fs, serverIP, hostsService)
	return
}
