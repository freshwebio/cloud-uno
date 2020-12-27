// +build windows

package hosts

const (
	// HostsPerLine is the amount of hosts that can appear on a single line.
	HostsPerLine = 9
	// HostsFilePath provides the location to the hosts file on darwin and linux
	// systems.
	HostsFilePath = "${SystemRoot}/System32/drivers/etc/hosts"
	eol           = "\r\n"
)
