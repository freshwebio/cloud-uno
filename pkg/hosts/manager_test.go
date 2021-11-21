// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

//go:build unit

package hosts

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/freshwebio/cloud-uno/pkg/config"
	"github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

type ManagerSuite struct {
	fixtures map[string]hostsManagerTestFixture
	dir      string
}

var _ = Suite(&ManagerSuite{})

type hostsManagerTestFixture struct {
	expected string
	input    string
}

func (s *ManagerSuite) SetUpSuite(c *C) {
	// Create a temp directory that gocheck will deal with tearing
	// down for us.
	s.dir = c.MkDir()
	s.fixtures = make(map[string]hostsManagerTestFixture)
	fixtureFilePrefixes := []string{
		"add1",
		"add2",
		"add3",
		"remove1",
		"remove2",
	}
	for _, prefix := range fixtureFilePrefixes {
		expectedFilePath := fmt.Sprintf("testdata/manager/%s-expected.txt", prefix)
		inputFilePath := fmt.Sprintf("testdata/manager/%s-input.txt", prefix)
		expectedBytes, err := ioutil.ReadFile(expectedFilePath)
		if err != nil {
			c.Error(err)
			c.FailNow()
		}
		inputBytes, err := ioutil.ReadFile(inputFilePath)
		if err != nil {
			c.Error(err)
			c.FailNow()
		}
		s.fixtures[prefix] = hostsManagerTestFixture{
			expected: string(expectedBytes),
			input:    string(inputBytes),
		}
	}
}

func (s *ManagerSuite) Test_add_hosts_to_existing_cloud_uno_entry_ip(c *C) {
	s.addHostsTest(c, "add1", "172.18.0.22", "storage.googleapis.local")
}

func (s *ManagerSuite) Test_add_entry_for_hosts_file_without_a_cloud_uno_section_and_removes_duplicate_hosts(c *C) {
	s.addHostsTest(c, "add2", "172.18.0.22", "s3.aws.local")
}

func (s *ManagerSuite) Test_add_entry_for_new_ip(c *C) {
	s.addHostsTest(c, "add3", "172.18.0.24", "somethingnew.googleapis.local")
}

func (s *ManagerSuite) addHostsTest(c *C, fixtureName string, ip string, hosts string) {
	hostsPath := fmt.Sprintf("%s/%s-hosts", s.dir, fixtureName)
	manager, err := s.setUpManagerForTest(hostsPath, fixtureName)
	if err != nil {
		c.Error(err)
		c.FailNow()
	}

	err = manager.Add(&Params{
		IP:    &ip,
		Hosts: &hosts,
	})
	if err != nil {
		c.Error(err)
		c.FailNow()
	}

	s.assertHostsFileEqual(c, hostsPath, fixtureName)
}

func (s *ManagerSuite) Test_remove_host_entry(c *C) {
	s.removeHostsTest(c, "remove1", "172.18.0.22", "secretmanager.googleapis.local")
}

func (s *ManagerSuite) Test_remove_last_cloud_uno_host_entry(c *C) {
	s.removeHostsTest(c, "remove2", "172.18.0.23", "storage.googleapis.local")
}

func (s *ManagerSuite) removeHostsTest(c *C, fixtureName string, ip string, hosts string) {
	hostsPath := fmt.Sprintf("%s/%s-hosts", s.dir, fixtureName)
	manager, err := s.setUpManagerForTest(hostsPath, fixtureName)
	if err != nil {
		c.Error(err)
		c.FailNow()
	}

	err = manager.Remove(&Params{
		IP:    &ip,
		Hosts: &hosts,
	})
	if err != nil {
		c.Error(err)
		c.FailNow()
	}

	s.assertHostsFileEqual(c, hostsPath, fixtureName)
}

func (s *ManagerSuite) setUpManagerForTest(hostsPath string, fixtureName string) (Service, error) {
	err := ioutil.WriteFile(hostsPath, []byte(s.fixtures[fixtureName].input), 0644)
	if err != nil {
		return nil, err
	}

	return NewManager(
		&config.Config{
			HostsPath: &hostsPath,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
}

func (s *ManagerSuite) assertHostsFileEqual(c *C, hostsPath string, fixtureName string) {
	persistedHostsBytes, err := ioutil.ReadFile(hostsPath)
	if err != nil {
		c.Error(err)
		c.FailNow()
	}
	actual := normaliseHostsText(string(persistedHostsBytes))
	expected := normaliseHostsText(s.fixtures[fixtureName].expected)
	c.Assert(
		actual,
		Equals,
		expected,
	)
}

func normaliseHostsText(hostsText string) string {
	// We don't care about empty lines that contain only new line
	// and white space characters when determining quality between the obtained
	// and expected hosts text.
	reader := bytes.NewReader([]byte(hostsText))
	scanner := bufio.NewScanner(utfbom.SkipOnly(reader))
	normalised := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			normalised += eol
		} else if strings.HasPrefix(line, "#") {
			normalised += line + eol
		} else {
			// Make sure the hosts are ordered alphabetically
			// as per the entries sorting.
			hostPieces := strings.Split(line, " ")
			sort.Strings(hostPieces)
			normalised += strings.Join(hostPieces, " ") + eol
		}
	}
	return normalised
}
