package minio

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/minio/minio/pkg/madmin"
)

const minioKeyType = "minio_access_key"
const minioSecretKeyLength = 24

func (b *backend) minioAccessKeys() *framework.Secret {
	return &framework.Secret{
		Type: minioKeyType,
		Fields: map[string]*framework.FieldSchema{
			"accessKeyId": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Minio Access Key ID.",
			},
			"secretAccessKey": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Minio Secret Access Key.",
			},
			"accountStatus": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Minio account status.",
			},
			"policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Minio policy attached to access key.",
			},
		},

		Revoke: b.minioAccessKeyRevoke,
	}
}

func (b *backend) minioAccessKeyCreate(ctx context.Context, s logical.Storage, accessKeyId, policy string) (*madmin.UserInfo, error) {

	client, err := b.getMadminClient(ctx, s)
	if err != nil {
		return nil, err
	}

	b.Logger().Info("Adding minio user", "accessKeyId", accessKeyId)

	secretAccessKey, err := base62.Random(minioSecretKeyLength)
	if err != nil {
		b.Logger().Error("Generating random secret key", "accessKeyId", accessKeyId, "error", err)
		return nil, err
	}

	err = client.AddUser(ctx, accessKeyId, secretAccessKey)
	if err != nil {
		b.Logger().Error("Adding minio user failed", "accessKeyId", accessKeyId,
			"error", err)
		return nil, err
	}

	b.Logger().Info("Adding policy to minio user", "accessKeyId", accessKeyId,
		"policy", policy)
	err = client.SetPolicy(ctx, policy, accessKeyId, false) //using updated method
	if err != nil {
		b.Logger().Error("Setting minio user policy failed", "accessKeyId", accessKeyId,
			"policy", policy, "error", err)
		return nil, err
	}

	newUser := &madmin.UserInfo{
		SecretKey:  secretAccessKey,
		PolicyName: policy,
		Status:     madmin.AccountEnabled,
	}

	return newUser, nil
}

func (b *backend) minioAccessKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	client, err := b.getMadminClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	accessKeyIdRaw, ok := req.Secret.InternalData["accessKeyId"]
	if !ok {
		return nil, fmt.Errorf("secret is missing internal accessKeyId")
	}

	accessKeyId, ok := accessKeyIdRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing internal accessKeyId")
	}

	b.Logger().Info("Revoking access key", "accessKeyId", accessKeyId)

	if err = client.RemoveUser(ctx, accessKeyId); err != nil {
		return nil, err
	}

	return nil, nil
}
