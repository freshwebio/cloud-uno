// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HTTPError writes http error responses with a message represented in a JSON object.
func HTTPError(w http.ResponseWriter, statusCode int, message string) {
	HTTPErrorWithFields(w, statusCode, message, map[string]interface{}{})
}

// HTTPErrorWithFields writes http error responses with a message represented in a JSON object
// along with extra fields.
func HTTPErrorWithFields(w http.ResponseWriter, statusCode int, message string, fields map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fields["message"] = message
	errorResponse, _ := json.Marshal(fields)
	w.Write(errorResponse)
}

// HTTPErrorFromGRPC deals with taking an error from a gRPC service call
// and converting it to a HTTP status code.
func HTTPErrorFromGRPC(w http.ResponseWriter, err error) {
	message := "Unexpected error occurred"
	httpStatusCode := http.StatusInternalServerError
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.PermissionDenied:
			message = "Permission denied"
			httpStatusCode = http.StatusForbidden
		case codes.Unauthenticated:
			message = "Unauthenticated"
			httpStatusCode = http.StatusUnauthorized
		case codes.Unimplemented:
			message = "Not implemented"
			httpStatusCode = http.StatusNotImplemented
		case codes.InvalidArgument:
			message = fmt.Sprintf("Invalid argument: %s", err.Error())
			httpStatusCode = http.StatusBadRequest
		}
	}
	HTTPError(
		w, httpStatusCode,
		message,
	)
}
