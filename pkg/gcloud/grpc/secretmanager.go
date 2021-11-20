// Copyright (c) 2022 FRESHWEB LTD.
// Use of this software is governed by the Business Source License
// included in the file LICENSE
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/LICENSE-Apache-2.0

package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/freshwebio/cloud-uno/pkg/hosts"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/spf13/afero"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	v1Iam "google.golang.org/genproto/googleapis/iam/v1"
)

// SecretManager provides a gRPC secret manager service which can also
// be used for the HTTP API using Google's handy protojson package to translate
// between proto3 and json.
type SecretManager struct {
	dataRootDir string
	fs          afero.Fs
}

var (
	secretManagerLocalHost = "secretmanager.googleapis.local"
)

// NewSecretManager creates an instance of the Cloud::1 secret manager implementaiton.
func NewSecretManager(dataRootDir string, fs afero.Fs, ip string, hostsService hosts.Service) (secretmanagerpb.SecretManagerServiceServer, error) {
	err := fs.MkdirAll(dataRootDir, 0755)
	if err != nil {
		return nil, err
	}

	err = hostsService.Add(&hosts.HostsParams{
		IP:    &ip,
		Hosts: &secretManagerLocalHost,
	})
	if err != nil {
		return nil, err
	}
	return &SecretManager{
		dataRootDir,
		fs,
	}, nil
}

// ListSecrets deals with listing secrets for a provided project.
func (s *SecretManager) ListSecrets(ctx context.Context, req *secretmanagerpb.ListSecretsRequest) (*secretmanagerpb.ListSecretsResponse, error) {
	secrets := []*secretmanagerpb.Secret{}
	projectDir := fmt.Sprintf("%s/%s/secrets", s.dataRootDir, req.Parent)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			name := info.Name()
			secretID := strings.TrimSuffix(name, ".json")
			name = fmt.Sprintf("%s/secrets/%s", req.Parent, secretID)
			secret, err := s.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
				Name: name,
			})
			if err != nil {
				return err
			}
			secrets = append(secrets, secret)
		}
		return nil
	}
	err := afero.Walk(s.fs, projectDir, walkFn)
	if err != nil {
		return nil, err
	}
	listSecretsResponse := &secretmanagerpb.ListSecretsResponse{
		Secrets: secrets,
	}
	return listSecretsResponse, err
}

// CreateSecret deals with creating a new secret for a provided project.
func (s *SecretManager) CreateSecret(ctx context.Context, req *secretmanagerpb.CreateSecretRequest) (*secretmanagerpb.Secret, error) {
	secret := secretmanagerpb.Secret{}
	copier.Copy(&secret, req.Secret)
	secret.CreateTime = timestamppb.Now()
	secret.Name = fmt.Sprintf("%s/secrets/%s", req.Parent, req.SecretId)
	bytes, err := protojson.Marshal(&secret)
	if err != nil {
		return nil, err
	}
	dirPath := fmt.Sprintf("%s/%s/secrets/%s", s.dataRootDir, req.Parent, req.SecretId)
	err = s.fs.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s.json", dirPath, req.SecretId)
	handle, err := s.fs.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	_, err = handle.Write(bytes)
	return &secret, err
}

// Versions represents secret versions.
type Versions struct {
	Next     int             `json:"next"`
	Versions map[int]Version `json:"versions"`
}

// Version represents a secret version persisted
// to a configured file system.
type Version struct {
	File       string `json:"file"`
	Number     int    `json:"number"`
	CreateTime int    `json:"createTime"`
}

func (s *SecretManager) getVersions(secret string) (*Versions, error) {
	versionsFilePath := fmt.Sprintf("%s/%s/versions.json", s.dataRootDir, secret)
	exists, err := afero.Exists(s.fs, versionsFilePath)
	if err != nil {
		return nil, err
	}
	if exists {
		bytes, err := afero.ReadFile(s.fs, versionsFilePath)
		if err != nil {
			return nil, err
		}
		versions := &Versions{}
		err = json.Unmarshal(bytes, versions)
		if err != nil {
			return nil, err
		}
		return versions, nil
	}
	return &Versions{
		Next:     1,
		Versions: make(map[int]Version),
	}, nil
}

func (s *SecretManager) addVersion(versions *Versions, versionsDirectory string, payload []byte) (*Version, error) {
	versionsCopy := &Versions{}
	copier.Copy(versionsCopy, versions)
	fileNameUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	fileName := fileNameUUID.String()
	version := Version{
		File: fileName,
		// Ensure we take a copy instead of using a pointer
		// to make sure the correct value is serialised.
		Number:     (*versionsCopy).Next,
		CreateTime: int(time.Now().Unix()),
	}
	versionsCopy.Versions[version.Number] = version
	filePath := fmt.Sprintf("%s/%s", versionsDirectory, fileName)
	// Write the file containng the secret data.
	err = afero.WriteFile(s.fs, filePath, payload, 0755)
	if err != nil {
		return nil, err
	}
	// Before writing the update, increment the next id.
	// incrementing integers isn't great for requests being made in parallel,
	// this will most likely need to be improved depending on what software is being used locally
	// to orchestrate emulations of infrastructure.
	versionsCopy.Next = versionsCopy.Next + 1
	// Write changes to the versions object to the file system.
	versionsBytes, err := json.Marshal(versions)
	if err != nil {
		return nil, err
	}
	versionsFilePath := fmt.Sprintf("%s/versions.json", versionsDirectory)
	err = afero.WriteFile(s.fs, versionsFilePath, versionsBytes, 0755)
	return &version, nil
}

// AddSecretVersion deals with adding a new version for a specified secret.
func (s *SecretManager) AddSecretVersion(ctx context.Context, req *secretmanagerpb.AddSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	secret := req.Parent
	versions, err := s.getVersions(secret)
	if err != nil {
		return nil, err
	}
	versionsDirectory := fmt.Sprintf("%s/%s", s.dataRootDir, req.Parent)
	fmt.Println("versionsDirectory:", versionsDirectory)
	version, err := s.addVersion(versions, versionsDirectory, req.Payload.Data)
	if err != nil {
		return nil, err
	}
	versionName := fmt.Sprintf("%s/versions/%d", req.Parent, versions.Next)
	return &secretmanagerpb.SecretVersion{
		Name:       versionName,
		CreateTime: timestamppb.New(time.Unix(int64(version.CreateTime), 0)),
	}, nil
}

func (s *SecretManager) createSecretFilePath(name string) string {
	pathPieces := strings.Split(name, "/")
	fileName := fmt.Sprintf("%s.json", pathPieces[len(pathPieces)-1])
	filePath := fmt.Sprintf("%s/%s/%s", s.dataRootDir, name, fileName)
	return filePath
}

// GetSecret deals with retrieving a specified secret.
func (s *SecretManager) GetSecret(ctx context.Context, req *secretmanagerpb.GetSecretRequest) (*secretmanagerpb.Secret, error) {
	filePath := s.createSecretFilePath(req.Name)
	bytes, err := afero.ReadFile(s.fs, filePath)
	if err != nil {
		return nil, err
	}
	secret := &secretmanagerpb.Secret{}
	err = protojson.Unmarshal(bytes, secret)
	if err != nil {
		return nil, err
	}
	secret.Name = req.Name
	return secret, nil
}

// UpdateSecret deals with updating a subset of fields for the specified secret.
func (s *SecretManager) UpdateSecret(ctx context.Context, req *secretmanagerpb.UpdateSecretRequest) (*secretmanagerpb.Secret, error) {
	err := validateUpdateMask(req.UpdateMask)
	if err != nil {
		return nil, err
	}
	storedSecret, err := s.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: req.Secret.Name,
	})
	if err != nil {
		return nil, err
	}
	storedSecret.Labels = req.Secret.Labels
	bytes, err := protojson.Marshal(storedSecret)
	if err != nil {
		return nil, err
	}
	filePath := s.createSecretFilePath(req.Secret.Name)
	err = afero.WriteFile(s.fs, filePath, bytes, 0755)
	return storedSecret, err
}

// DeleteSecret deals with deleting a secret for the provided project.
func (*SecretManager) DeleteSecret(context.Context, *secretmanagerpb.DeleteSecretRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}

// ListSecretVersions deals with listing all versions of a given secret.
func (*SecretManager) ListSecretVersions(context.Context, *secretmanagerpb.ListSecretVersionsRequest) (*secretmanagerpb.ListSecretVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecretVersions not implemented")
}

// GetSecretVersion deals with retrieving metadata about a secret version.
func (*SecretManager) GetSecretVersion(context.Context, *secretmanagerpb.GetSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecretVersion not implemented")
}

// AccessSecretVersion deals with retrieving the raw data for a specified secret version.
func (*SecretManager) AccessSecretVersion(context.Context, *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccessSecretVersion not implemented")
}

// DisableSecretVersion deals with disabling the specified secret version.
func (*SecretManager) DisableSecretVersion(context.Context, *secretmanagerpb.DisableSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableSecretVersion not implemented")
}

// EnableSecretVersion deals with enabling the specified secret version.
func (*SecretManager) EnableSecretVersion(context.Context, *secretmanagerpb.EnableSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableSecretVersion not implemented")
}

// DestroySecretVersion deals with permanently destroying the specified secret version.
func (*SecretManager) DestroySecretVersion(context.Context, *secretmanagerpb.DestroySecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroySecretVersion not implemented")
}

// SetIamPolicy deals with setting an IAM policy for the specified secret.
func (*SecretManager) SetIamPolicy(context.Context, *v1Iam.SetIamPolicyRequest) (*v1Iam.Policy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetIamPolicy not implemented")
}

// GetIamPolicy retrieves the IAM policy for the specified secret.
func (*SecretManager) GetIamPolicy(context.Context, *v1Iam.GetIamPolicyRequest) (*v1Iam.Policy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIamPolicy not implemented")
}

// TestIamPermissions checks the permissions the caller has for the specified secret.
func (*SecretManager) TestIamPermissions(context.Context, *v1Iam.TestIamPermissionsRequest) (*v1Iam.TestIamPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestIamPermissions not implemented")
}

func validateUpdateMask(updateMask *fieldmaskpb.FieldMask) error {
	foundInvalidField := false
	i := 0
	for !foundInvalidField && i < len(updateMask.Paths) {
		if updateMask.Paths[i] != "labels" {
			foundInvalidField = true
		}
		i = i + 1
	}
	if foundInvalidField {
		return status.Errorf(codes.InvalidArgument, "Update mask must only contain mutable fields")
	}
	return nil
}
