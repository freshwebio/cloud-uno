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
