// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package httpapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/freshwebio/cloud-uno/pkg/httputils"
	"github.com/freshwebio/cloud-uno/pkg/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	// SecretManagerHost specifies the host on which Cloud::1 will accept
	// API requests for Google Cloud Secret Manager.
	SecretManagerHost              = "secretmanager.googleapis.local"
	failedPreparingResponseMessage = "Unexpected error occurred: failed when preparing response"
)

// RegisterSecretManager deals with registering the routes for the secret manager api.
func RegisterSecretManager(router *mux.Router, resolver types.Resolver) {
	secretManager := resolver.Get("gcloud.secretmanager").(secretmanagerpb.SecretManagerServiceServer)
	logger := resolver.Get("logger").(*logrus.Entry)
	c := &secretManagerController{
		secretManager,
		logger,
	}
	router.HandleFunc("/v1/projects/{project}/secrets/{secret:.*:addVersion}", c.AddVersion).
		Methods("POST").Host(SecretManagerHost)

	router.HandleFunc("/v1/projects/{project}/secrets", c.Create).
		Methods("POST").Host(SecretManagerHost).
		Queries("secretId", "{secretId:.+}")

	router.HandleFunc("/v1/projects/{project}/secrets", c.ListSecrets).
		Methods("GET").Host(SecretManagerHost)

	router.HandleFunc("/v1/projects/{project}/secrets/{secret}", c.GetSecret).
		Methods("GET").Host(SecretManagerHost)

	router.HandleFunc("/v1/projects/{project}/secrets/{secret}", c.UpdateSecret).
		Methods("PATCH").Host(SecretManagerHost).
		Queries("updateMask", "{updateMask:.+}")
}

type secretManagerController struct {
	secretManager secretmanagerpb.SecretManagerServiceServer
	logger        *logrus.Entry
}

func (c *secretManagerController) AddVersion(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	secret := strings.TrimSuffix(mux.Vars(r)["secret"], ":addVersion")
	parent := fullyQualifiedSecretName(project, secret)
	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			httputils.InvalidRequestMessage(err),
		)
		return
	}
	secretVersionRequest := &secretmanagerpb.AddSecretVersionRequest{}
	err = protojson.Unmarshal(requestBytes, secretVersionRequest)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	// Set parent after unmarshalling so it doesn't get overridden.
	secretVersionRequest.Parent = parent
	secretVersion, err := c.secretManager.AddSecretVersion(
		r.Context(),
		secretVersionRequest,
	)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(secretVersion)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			failedPreparingResponseMessage,
		)
		return
	}
	httputils.SetResponseAsJSON(w)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseBytes))
}

func (c *secretManagerController) Create(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	secretID := mux.Vars(r)["secretId"]
	parent := fullyQualifiedProject(project)
	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			httputils.InvalidRequestMessage(err),
		)
		return
	}
	createSecretRequest := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
		Secret:   &secretmanagerpb.Secret{},
	}
	err = protojson.Unmarshal(requestBytes, createSecretRequest.Secret)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			fmt.Sprintf("Invalid request: %s", err.Error()),
		)
		return
	}

	secretResponse, err := c.secretManager.CreateSecret(
		r.Context(),
		createSecretRequest,
	)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(secretResponse)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			failedPreparingResponseMessage,
		)
		return
	}
	httputils.SetResponseAsJSON(w)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseBytes))
}

func (c *secretManagerController) ListSecrets(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	parent := fullyQualifiedProject(project)
	listSecretsRequest := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}
	listSecretsResponse, err := c.secretManager.ListSecrets(
		r.Context(),
		listSecretsRequest,
	)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(listSecretsResponse)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			failedPreparingResponseMessage,
		)
		return
	}
	httputils.SetResponseAsJSON(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBytes))
}

func (c *secretManagerController) GetSecret(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	secret := mux.Vars(r)["secret"]
	name := fullyQualifiedSecretName(project, secret)
	getSecretRequest := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}
	getSecretResponse, err := c.secretManager.GetSecret(
		r.Context(),
		getSecretRequest,
	)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(getSecretResponse)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			failedPreparingResponseMessage,
		)
		return
	}
	httputils.SetResponseAsJSON(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBytes))
}

func (c *secretManagerController) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	secret := mux.Vars(r)["secret"]
	mask := mux.Vars(r)["updateMask"]
	name := fullyQualifiedSecretName(project, secret)
	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			httputils.InvalidRequestMessage(err),
		)
		return
	}
	updateSecretRequest := &secretmanagerpb.UpdateSecretRequest{
		Secret: &secretmanagerpb.Secret{},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: strings.Split(mask, ","),
		},
	}
	err = protojson.Unmarshal(requestBytes, updateSecretRequest.Secret)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			httputils.InvalidRequestMessage(err),
		)
		return
	}
	updateSecretRequest.Secret.Name = name
	updateSecretResponse, err := c.secretManager.UpdateSecret(
		r.Context(),
		updateSecretRequest,
	)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(updateSecretResponse)
	if err != nil {
		c.logger.Error(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			failedPreparingResponseMessage,
		)
		return
	}
	httputils.SetResponseAsJSON(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBytes))
}

func fullyQualifiedSecretName(project string, secret string) string {
	return fmt.Sprintf("projects/%s/secrets/%s", project, secret)
}

func fullyQualifiedProject(project string) string {
	return fmt.Sprintf("projects/%s", project)
}
