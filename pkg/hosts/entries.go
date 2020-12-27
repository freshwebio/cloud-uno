package hosts

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

const commentChar string = "#"

// Entry provides a host entry
// in the machine's hosts file.
type Entry struct {
	IP      string
	Hosts   []string
	Raw     string
	Err     error
	Comment string
}

// NewEntry creates a new host entry.
func NewEntry(raw string) Entry {
	output := Entry{
		Raw: raw,
	}
	if output.IsComment() {
		return output
	}

	if output.HasComment() {
		// trailing comment
		commentSplit := strings.Split(output.Raw, commentChar)
		raw = commentSplit[0]
		output.Comment = commentSplit[1]
	}

	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return output
	}

	rawIP := fields[0]
	if net.ParseIP(rawIP) == nil {
		output.Err = fmt.Errorf("bad hosts entry: %q", raw)
	}

	output.IP = rawIP
	output.Hosts = fields[1:]

	return output
}

// Export converts a host entry to a line for an entry in
// a hosts file.
func (e *Entry) Export() string {
	var comment string
	if e.IsComment() { //Whole line is comment
		return e.Raw
	}

	if e.Comment != "" {
		comment = fmt.Sprintf(" %s%s", commentChar, e.Comment)
	}

	return fmt.Sprintf("%s %s%s", e.IP, strings.Join(e.Hosts, " "), comment)
}

// RemoveDuplicateHosts deals with cleaning up
// duplicate hosts.
func (e *Entry) RemoveDuplicateHosts() {
	unique := make(map[string]struct{})
	for _, h := range e.Hosts {
		unique[h] = struct{}{}
	}

	e.Hosts = []string{}
	for k := range unique {
		e.Hosts = append(e.Hosts, k)
	}
	e.Raw = e.Export()
}

// Combine deals with merging host entries.
func (e *Entry) Combine(other Entry) {
	e.Hosts = append(e.Hosts, other.Hosts...)
	if e.Comment == "" {
		e.Comment = other.Comment
	} else {
		e.Comment = fmt.Sprintf("%s %s", e.Comment, other.Comment)
	}
	e.Raw = other.Export()
}

// SortHosts deals with sorting the hosts list.
func (e *Entry) SortHosts() {
	sort.Strings(e.Hosts)
	e.Raw = e.Export()
}

// IsComment determines whether or not the entry is a comment.
func (e *Entry) IsComment() bool {
	return strings.HasPrefix(strings.TrimSpace(e.Raw), commentChar)
}

// HasComment determines whether or not the entry
// contains a comment.
func (e *Entry) HasComment() bool {
	return strings.Contains(e.Raw, commentChar)
}

// IsValid determines whether or not the entry
// is a valid.
func (e *Entry) IsValid() bool {
	return e.IP != ""
}

// IsMalformed determines whether or not the entry
// is malformed.
func (e *Entry) IsMalformed() bool {
	return e.Err != nil
}

// RegenerateExport deals with regenerating the host entry in it's raw format.
func (e *Entry) RegenerateExport() {
	e.Raw = fmt.Sprintf("%s %s", e.IP, strings.Join(e.Hosts, " "))
}
