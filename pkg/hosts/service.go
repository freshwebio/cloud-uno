package hosts

// Service provides a common interface for a service that manages
// os hosts, can be a server or a client.
type Service interface {
	Add(ip *string, hosts *string) error
	Remove(ip *string, hosts *string) error
}
