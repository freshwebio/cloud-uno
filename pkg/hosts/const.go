// +build !windows

package hosts

const (
	// HostsPerLine is the amount of hosts that can appear on a single line.
	HostsPerLine = -1 // unlimited
	// HostsFilePath provides the location to the hosts file on darwin and linux
	// systems.
	HostsFilePath = "/etc/hosts"
	eol           = "\n"
)
