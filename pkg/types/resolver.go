package types

// Resolver provides an interface for components that want
// to provide service resolution for Cloud::1.
type Resolver interface {
	// Retrieve a service.
	Get(service string) interface{}
	// Registers all services, should be called during initialisation
	// of an application/web server.
	Register() error
}
