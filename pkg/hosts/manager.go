// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package hosts

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/sirupsen/logrus"
)

const cloudUnoOpenComment = "Added by Cloud::1"
const cloudUnoCloseComment = "End of Cloud::1 section"
const cloudUnoEntryMark = "cloud::1"

// Manager provides a service that deals with managing
// hosts on the host machine as a part of cloud DNS emulation.
type Manager struct {
	Path               string
	Entries            []Entry
	logger             *logrus.Entry
	hasCloudUnoSection bool
}

// NewManager creates a new instance of a service that deals with managing
// os-level hosts on a machine.
func NewManager(cfg *config.Config, logger *logrus.Entry) (Service, error) {
	osHostsFilePath := os.ExpandEnv(filepath.FromSlash(HostsFilePath))
	if *cfg.HostsPath != "" {
		osHostsFilePath = os.ExpandEnv(filepath.FromSlash(*cfg.HostsPath))
	}

	mgr := &Manager{
		Path:    osHostsFilePath,
		Entries: []Entry{},
		logger:  logger,
	}
	if err := mgr.load(); err != nil {
		return mgr, err
	}
	return mgr, nil
}

func (m *Manager) load() error {
	file, err := os.Open(m.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(utfbom.SkipOnly(file))
	inSection := false
	hasOpenComment := false
	hasCloseComment := false
	for scanner.Scan() {
		entry := NewEntry(scanner.Text())

		if isOpenCommentEntry(entry) {
			inSection = true
			hasOpenComment = true
		} else if isCloseCommentEntry(entry) {
			inSection = false
			hasCloseComment = true
		}

		if inSection {
			entry.Mark(cloudUnoEntryMark)
		}

		m.Entries = append(m.Entries, entry)
	}
	m.hasCloudUnoSection = hasOpenComment && hasCloseComment
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) getIPPosition(ip string) int {
	position := -1
	i := 0
	for position == -1 && i < len(m.Entries) {
		entry := m.Entries[i]
		if !entry.IsComment() && entry.IsMarkedWith(cloudUnoEntryMark) && entry.Raw != "" && entry.IP == ip {
			position = i
		}
		i++
	}
	return position
}

// Add one or more host entries.
func (m *Manager) Add(params *Params) error {
	if net.ParseIP(*params.IP) == nil {
		return fmt.Errorf("%q is an invalid IP address", *params.IP)
	}
	hostsList := strings.Split(*params.Hosts, ",")

	position := m.getIPPosition(*params.IP)
	if position == -1 {
		// ip not already in hostsfile inside the cloud uno secton.
		entry := Entry{
			Raw:   buildRawLine(*params.IP, hostsList),
			IP:    *params.IP,
			Hosts: hostsList,
		}
		entry.Mark(cloudUnoEntryMark)
		if !m.hasCloudUnoSection {
			m.addEntryInNewCloudUnoSection(entry)
			m.hasCloudUnoSection = true
		} else {
			m.addEntryToCloudUnoSection(entry)
		}
	} else {
		// add new hosts to the correct position for the ip
		hostsCopy := m.Entries[position].Hosts
		for _, addHost := range hostsList {
			if !itemInSlice(addHost, hostsCopy) {
				hostsCopy = append(hostsCopy, addHost)
			}
		}
		m.Entries[position].Hosts = hostsCopy
		m.Entries[position].Raw = m.Entries[position].Export() // reset raw
	}
	m.clean()
	// Each host can only be configured to work for a single IP at a time,
	// to ensure the provided IP is used we need to make sure
	// we clear all other references to the same hosts.
	m.removeHostsFromOtherIPs(*params.IP, hostsList)
	err := m.flush()
	return err
}

func (m *Manager) addEntryInNewCloudUnoSection(entry Entry) {
	m.Entries = append(
		m.Entries,
		Entry{
			Raw: fmt.Sprintf("# %s", cloudUnoOpenComment),
		},
		entry,
		Entry{
			Raw: fmt.Sprintf("# %s", cloudUnoCloseComment),
		},
	)
}

func (m *Manager) addEntryToCloudUnoSection(entry Entry) {
	cloudUnoSectionClosePos := m.getCloseCloudUnoSectionPosition()
	m.Entries = InsertIntoSlice(m.Entries, cloudUnoSectionClosePos, entry)
}

func (m *Manager) getCloseCloudUnoSectionPosition() int {
	position := -1
	i := 0
	for position == -1 && i < len(m.Entries) {
		if m.Entries[i].Raw == fmt.Sprintf("# %s", cloudUnoCloseComment) {
			position = i
		}
		i++
	}
	return position
}

// Remove one or more host entries.
func (m *Manager) Remove(params *Params) error {
	var outputEntries []Entry
	hostsList := strings.Split(*params.Hosts, ",")
	if net.ParseIP(*params.IP) == nil {
		return fmt.Errorf("%q is an invalid IP address", *params.IP)
	}

	for _, entry := range m.Entries {
		// Bad lines, comments and entries outside of
		// the cloud uno section just get re-added.
		if entry.Err != nil || !entry.IsMarkedWith(cloudUnoEntryMark) || entry.IsComment() || entry.IP != *params.IP {
			outputEntries = append(outputEntries, entry)
		} else {
			var newHosts []string
			for _, checkHost := range entry.Hosts {
				if !itemInSlice(checkHost, hostsList) {
					newHosts = append(newHosts, checkHost)
				}
			}

			// If hosts is empty, skip the line completely.
			if len(newHosts) > 0 {
				newLineRaw := entry.IP

				for _, host := range newHosts {
					newLineRaw = fmt.Sprintf("%s %s", newLineRaw, host)
				}
				newEntry := NewEntry(newLineRaw)
				outputEntries = append(outputEntries, newEntry)
			}
		}
	}

	m.Entries = outputEntries
	hasMarkedEntries := m.hasMarkedEntries()
	if !hasMarkedEntries {
		m.removeCloudUnoSection()
	}
	m.clean()
	err := m.flush()
	return err
}

func (m *Manager) removeCloudUnoSection() {
	i := 0
	newEntries := []Entry{}
	for _, entry := range m.Entries {
		isCloudUnoSectionComment := isOpenCommentEntry(entry) || isCloseCommentEntry(entry)
		if !isCloudUnoSectionComment {
			newEntries = append(newEntries, entry)
		}
		i++
	}
	m.Entries = newEntries
}

func (m *Manager) hasMarkedEntries() bool {
	hasMarkedEntry := false
	i := 0
	for !hasMarkedEntry && i < len(m.Entries) {
		entry := m.Entries[i]
		// Ensure the open and close comments are excluded
		// from this check, as we are searching for contents within the Cloud
		// uno section.
		isCloudUnoSectionComment := isOpenCommentEntry(entry) || isCloseCommentEntry(entry)
		hasMarkedEntry = entry.IsMarkedWith(cloudUnoEntryMark) && !isCloudUnoSectionComment
		i++
	}
	return hasMarkedEntry
}

func (m *Manager) removeHostsFromOtherIPs(keepForIP string, hosts []string) {
	for _, host := range hosts {
		for pos, entry := range m.Entries {
			if itemInSlice(host, entry.Hosts) && entry.IP != keepForIP {
				entry.Hosts = removeFromSlice(host, entry.Hosts)
			}
			m.Entries[pos] = entry
		}
	}
}

func (m *Manager) clean() {
	for pos, entry := range m.Entries {
		entry.RemoveDuplicateHosts()
		entry.SortHosts()
		m.Entries[pos] = entry
	}
	m.hostsPerLine(HostsPerLine)
}

func (m *Manager) hostsPerLine(count int) {
	if count <= 0 {
		return
	}
	var newEntries []Entry
	for _, entry := range m.Entries {
		if len(entry.Hosts) <= count || !entry.IsMarkedWith(cloudUnoEntryMark) {
			newEntries = append(newEntries, entry)
		} else {
			for i := 0; i < len(entry.Hosts); i += count {
				entryCopy := entry
				end := len(entry.Hosts)
				if end > i+count {
					end = i + count
				}
				entryCopy.Hosts = entry.Hosts[i:end]
				entryCopy.Raw = entryCopy.Export()
				newEntries = append(newEntries, entryCopy)
			}
		}
	}
	m.Entries = newEntries
}

// Flush any changes made to the hosts file.
func (m *Manager) flush() error {
	file, err := os.Create(m.Path)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	for _, entry := range m.Entries {
		m.logger.Info(entry.Export())
		if _, err := fmt.Fprintf(w, "%s%s", entry.Export(), eol); err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return m.load()
}

func normaliseComment(rawComment string) string {
	return strings.TrimSpace(strings.Replace(rawComment, "#", "", 1))
}

func isOpenCommentEntry(entry Entry) bool {
	return entry.IsComment() && normaliseComment(entry.Raw) == cloudUnoOpenComment
}

func isCloseCommentEntry(entry Entry) bool {
	return entry.IsComment() && normaliseComment(entry.Raw) == cloudUnoCloseComment
}
