package grpc

import (
	"context"
	"fmt"

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

// NewSecretManager creates an instance of the Cloud::1 secret manager implementaiton.
func NewSecretManager(dataRootDir string, fs afero.Fs) (secretmanagerpb.SecretManagerServiceServer, error) {
	err := fs.MkdirAll(dataRootDir, 0755)
	if err != nil {
		return nil, err
	}
	return &SecretManager{
		dataRootDir,
		fs,
	}, nil
}

func (s *SecretManager) ListSecrets(context.Context, *secretmanagerpb.ListSecretsRequest) (*secretmanagerpb.ListSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecrets not implemented")
}

func (s *SecretManager) CreateSecret(ctx context.Context, req *secretmanagerpb.CreateSecretRequest) (*secretmanagerpb.Secret, error) {
	secret := secretmanagerpb.Secret{}
	copier.Copy(req.Secret, &secret)
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
	filePath := fmt.Sprintf("%s/%s.json", dirPath, req.SecretId)
	handle, err := s.fs.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	_, err = handle.Write(bytes)
	return &secret, err
}

func (*SecretManager) AddSecretVersion(context.Context, *secretmanagerpb.AddSecretVersionRequest) (*secretmanagerpb.SecretVersion, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSecretVersion not implemented")
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
