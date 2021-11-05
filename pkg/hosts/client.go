package hosts

import (
	"net/rpc"

	"github.com/freshwebio/cloud-uno/pkg/connect"
)

// NewClient creates a new client for managing os hosts.
func NewClient() (Service, error) {
	rpcClient, err := rpc.DialHTTP("unix", connect.SockAddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		rpcClient,
	}, nil
}

// Client represents a service to manage os hosts by interacting
// with another process with access to os hosts over a unix socket.
type Client struct {
	rpcClient *rpc.Client
}

// Add deals with adding a new set of host names for the specified ip address.
func (c *Client) Add(ip *string, hosts *string) error {
	return c.rpcClient.Call("Manager.Add", ip, hosts)
}

// Remove deals with removing a provided set of host names for the
// specified ip address.
func (c *Client) Remove(ip *string, hosts *string) error {
	return c.rpcClient.Call("Manager.Remove", ip, hosts)
}
