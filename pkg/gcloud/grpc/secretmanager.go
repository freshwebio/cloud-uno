package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/freshwebio/cloud-one/pkg/hosts"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/spf13/afero"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
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
	err = hostsService.Add(&ip, &secretManagerLocalHost)
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
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			filePath := fmt.Sprintf("%s/%s", projectDir, info.Name())
			bytes, err := afero.ReadFile(s.fs, filePath)
			if err != nil {
				return err
			}
			secret := &secretmanagerpb.Secret{}
			err = protojson.Unmarshal(bytes, secret)
			if err != nil {
				return err
			}
			secretID := strings.TrimSuffix(info.Name(), ".json")
			secret.Name = fmt.Sprintf("%s/secrets/%s", req.Parent, secretID)
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
	dirPath := fmt.Sprintf("%s/%s/secrets", s.dataRootDir, req.Parent)
	err = s.fs.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s/%s.json", dirPath, req.SecretId, req.SecretId)
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
	versionsCopy.Next = versionsCopy.Next + 1
	fileNameUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	fileName := fileNameUUID.String()
	version := &Version{
		File:       fileName,
		Number:     versionsCopy.Next,
		CreateTime: int(time.Now().Unix()),
	}
	filePath := fmt.Sprintf("%s/%s", versionsDirectory, fileName)
	// Write the file containng the secret data.
	err = afero.WriteFile(s.fs, filePath, payload, 0755)
	if err != nil {
		return nil, err
	}
	return version, nil
}

// AddSecretVersion deals with adding a new version for a specified secret.
func (s *SecretManager) AddSecretVersion(ctx context.Context, req *secretmanagerpb.AddSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	secret := req.Parent
	versions, err := s.getVersions(secret)
	if err != nil {
		return nil, err
	}
	versionsDirectory := fmt.Sprintf("%s/%s", s.dataRootDir, req.Parent)
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

func (*SecretManager) GetSecret(context.Context, *secretmanagerpb.GetSecretRequest) (*secretmanagerpb.Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecret not implemented")
}
func (*SecretManager) UpdateSecret(context.Context, *secretmanagerpb.UpdateSecretRequest) (*secretmanagerpb.Secret, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSecret not implemented")
}
func (*SecretManager) DeleteSecret(context.Context, *secretmanagerpb.DeleteSecretRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}
func (*SecretManager) ListSecretVersions(context.Context, *secretmanagerpb.ListSecretVersionsRequest) (*secretmanagerpb.ListSecretVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecretVersions not implemented")
}
func (*SecretManager) GetSecretVersion(context.Context, *secretmanagerpb.GetSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSecretVersion not implemented")
}
func (*SecretManager) AccessSecretVersion(context.Context, *secretmanagerpb.AccessSecretVersionRequest) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccessSecretVersion not implemented")
}
func (*SecretManager) DisableSecretVersion(context.Context, *secretmanagerpb.DisableSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableSecretVersion not implemented")
}
func (*SecretManager) EnableSecretVersion(context.Context, *secretmanagerpb.EnableSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableSecretVersion not implemented")
}
func (*SecretManager) DestroySecretVersion(context.Context, *secretmanagerpb.DestroySecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroySecretVersion not implemented")
}
func (*SecretManager) SetIamPolicy(context.Context, *v1Iam.SetIamPolicyRequest) (*v1Iam.Policy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetIamPolicy not implemented")
}
func (*SecretManager) GetIamPolicy(context.Context, *v1Iam.GetIamPolicyRequest) (*v1Iam.Policy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIamPolicy not implemented")
}
func (*SecretManager) TestIamPermissions(context.Context, *v1Iam.TestIamPermissionsRequest) (*v1Iam.TestIamPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestIamPermissions not implemented")
}
