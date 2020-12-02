package services

import (
	"fmt"
	"os"

	"github.com/freshwebio/cloud-one/pkg/gcloud/grpc"
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
	// The file system is used for services implemented directly in Cloud::1,
	// when handing off to other services like google cloud emulators we have no
	// control over whether that uses the file system or not.
	var fs afero.Fs
	fileSystemType := os.Getenv("CLOUD_ONE_FILE_SYSTEM")
	if fileSystemType == "memory" {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}
	r.services["fs"] = fs

	// Given gRPC is a fantastic representation of a service that is usually
	// abstracted away from a REST API route handler, the default resolver will use
	// the gRPC services.
	smRootDir := fmt.Sprintf("%s/gcloud/secretmanager", os.Getenv("CLOUD_ONE_DATA_DIR"))
	r.services["gcloud.secretmanager"], err = grpc.NewSecretManager(smRootDir, fs)
	return err
}
