package minio

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	ErrRoleNotFound = errors.New("role not found")
)

// A role stored in the storage backend
type Role struct {
	Policy         string        `json:"policy"`
	UserNamePrefix string        `json:"user_name_prefix"`
	DefaultTTL     time.Duration `json:"default_ttl"`
	MaxTTL         time.Duration `json:"max_ttl"`
}

// List Roles
func (b *backend) ListRoles(ctx context.Context, s logical.Storage) ([]string, error) {
	roles, err := s.List(ctx, "roles/")
	if err != nil {
		return nil, errwrap.Wrapf("Unable to retrieve list of roles: {{err}}", err)
	}

	return roles, nil
}

// Get Role
func (b *backend) GetRole(ctx context.Context, s logical.Storage, role string) (*Role, error) {
	r, err := s.Get(ctx, "roles/"+role)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to retrieve role %q: {{err}}", role), err)
	}

	if r == nil {
		return nil, ErrRoleNotFound
	}

	var rv Role
	if err := r.DecodeJSON(&rv); err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("Unable to decode role %q: {{err}}", role), err)
	}

	return &rv, nil
}
