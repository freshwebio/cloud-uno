// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package logging

import "github.com/sirupsen/logrus"

// CreateLogger deals with creating a logger to be used
// throughout the application.
func CreateLogger() *logrus.Entry {
	return logrus.New().WithFields(logrus.Fields{})
}
