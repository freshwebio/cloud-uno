package httpapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/freshwebio/cloud-one/pkg/httputils"
	"github.com/freshwebio/cloud-one/pkg/types"
	"github.com/gorilla/mux"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	// SecretManagerHost specifies the host on which Cloud::1 will accept
	// API requests for Google Cloud Secret Manager.
	SecretManagerHost = "secretmanager.googleapis.local"
)

// RegisterSecretManager deals with registering the routes for the secret manager api.
func RegisterSecretManager(router *mux.Router, resolver types.Resolver) {
	secretManager := resolver.Get("gcloud.secretmanager").(secretmanagerpb.SecretManagerServiceServer)
	c := &secretManagerController{
		secretManager,
	}
	router.HandleFunc("/v1/projects/{project}/secrets/{secret:.*:addVersion}", c.AddVersion).
		Methods("POST").Host("secretmanager.googleapis.local")

	router.HandleFunc("/v1/projects/{project}/secrets", c.Create).
		Methods("POST").Host("secretmanager.googleapis.local").
		Queries("secretId", "{secretId:.+}")

	router.HandleFunc("/v1/projects/{project}/secrets", c.ListSecrets).
		Methods("GET").Host("secretmanager.googleapis.local")
}

type secretManagerController struct {
	secretManager secretmanagerpb.SecretManagerServiceServer
}

func (c *secretManagerController) AddVersion(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	secret := mux.Vars(r)["secret"]
	parent := fmt.Sprintf("projects/%s/secrets/%s", project, secret)
	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			fmt.Sprintf("Invalid request: %s", err.Error()),
		)
		return
	}
	secretVersionRequest := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
	}
	err = protojson.Unmarshal(requestBytes, secretVersionRequest)
	secretVersion, err := c.secretManager.AddSecretVersion(
		r.Context(),
		secretVersionRequest,
	)
	if err != nil {
		log.Println(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(secretVersion)
	if err != nil {
		log.Println(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected error occurred: failed when preparing response"),
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseBytes))
}

func (c *secretManagerController) Create(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	secretID := mux.Vars(r)["secretId"]
	parent := fmt.Sprintf("projects/%s", project)
	requestBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputils.HTTPError(
			w, http.StatusBadRequest,
			fmt.Sprintf("Invalid request: %s", err.Error()),
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
		log.Println(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(secretResponse)
	if err != nil {
		log.Println(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected error occurred: failed when preparing response"),
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(responseBytes))
}

func (c *secretManagerController) ListSecrets(w http.ResponseWriter, r *http.Request) {
	project := mux.Vars(r)["project"]
	parent := fmt.Sprintf("projects/%s", project)
	listSecretsRequest := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}
	listSecretsResponse, err := c.secretManager.ListSecrets(
		r.Context(),
		listSecretsRequest,
	)
	if err != nil {
		log.Println(err)
		httputils.HTTPErrorFromGRPC(w, err)
		return
	}
	responseBytes, err := protojson.Marshal(listSecretsResponse)
	if err != nil {
		log.Println(err)
		httputils.HTTPError(
			w, http.StatusBadRequest,
			fmt.Sprintf("Unexpected error occurred: failed when preparing response"),
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBytes))
}
