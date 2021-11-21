// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

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
	marks   []string
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

// Mark adds a mark to an entry so it can be filtered
// on, this is especially useful for the hosts manager
// so it only manipulates hosts that were created by
// cloud uno.
func (e *Entry) Mark(mark string) {
	e.marks = append(e.marks, mark)
}

// IsMarkedWith checks whether an entry
// has been marked with the provided string.
func (e *Entry) IsMarkedWith(mark string) bool {
	hasMark := false
	i := 0
	for !hasMark && i < len(e.marks) {
		hasMark = e.marks[i] == mark
		i++
	}
	return hasMark
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

// InsertIntoSlice inserts a given hosts entry at a specified
// index into a slice.
func InsertIntoSlice(slice []Entry, index int, value Entry) []Entry {
	if len(slice) == index {
		return append(slice, value)
	}
	newSlice := append(slice[:index+1], slice[index:]...)
	newSlice[index] = value
	return newSlice
}
