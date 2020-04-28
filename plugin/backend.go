package minio

import (
	"context"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/minio/minio/pkg/madmin"
)

type backend struct {
	*framework.Backend

	client *madmin.AdminClient

	clientMutex sync.RWMutex
}

// Factory returns a configured instance of the minio backend
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}

	b.Logger().Info("Plugin successfully initialized")
	return b, nil
}

// Backend returns a configured minio backend
func Backend() *backend {
	var b backend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        "The minio secrets backend provisions users on a Minio server",

		Paths: []*framework.Path{
			b.pathConfigCRUD(),
			b.pathRoles(),
			b.pathRolesCRUD(),
			b.pathKeysRead(),
		},

		Secrets: []*framework.Secret{
			b.minioAccessKeys(),
		},
	}

	b.client = (*madmin.AdminClient)(nil)

	return &b
}
