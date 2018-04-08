package client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/suzuki-shunsuke/go-graylog"
)

// CreateRole creates a new role.
func (client *Client) CreateRole(role *graylog.Role) (*ErrorInfo, error) {
	return client.CreateRoleContext(context.Background(), role)
}

// CreateRoleContext creates a new role with a context.
func (client *Client) CreateRoleContext(
	ctx context.Context, role *graylog.Role,
) (*ErrorInfo, error) {
	if role == nil {
		return nil, fmt.Errorf("role is nil")
	}
	return client.callPost(ctx, client.Endpoints.Roles, role, role)
}

// GetRoles returns all roles.
func (client *Client) GetRoles() ([]graylog.Role, int, *ErrorInfo, error) {
	return client.GetRolesContext(context.Background())
}

// GetRolesContext returns all roles with a context.
func (client *Client) GetRolesContext(ctx context.Context) (
	[]graylog.Role, int, *ErrorInfo, error,
) {
	roles := &graylog.RolesBody{}
	ei, err := client.callGet(
		ctx, client.Endpoints.Roles, nil, roles)
	return roles.Roles, roles.Total, ei, err
}

// GetRole returns a given role.
func (client *Client) GetRole(name string) (*graylog.Role, *ErrorInfo, error) {
	return client.GetRoleContext(context.Background(), name)
}

// GetRoleContext returns a given role with a context.
func (client *Client) GetRoleContext(
	ctx context.Context, name string,
) (*graylog.Role, *ErrorInfo, error) {
	if name == "" {
		return nil, nil, errors.New("name is empty")
	}
	role := &graylog.Role{}
	ei, err := client.callGet(
		ctx, client.Endpoints.Role(name), nil, role)
	return role, ei, err
}

// UpdateRole updates a given role.
func (client *Client) UpdateRole(name string, role *graylog.Role) (
	*ErrorInfo, error,
) {
	return client.UpdateRoleContext(context.Background(), name, role)
}

// UpdateRoleContext updates a given role with a context.
func (client *Client) UpdateRoleContext(
	ctx context.Context, name string, role *graylog.Role,
) (*ErrorInfo, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}
	if role == nil {
		return nil, fmt.Errorf("role is nil")
	}
	return client.callPut(ctx, client.Endpoints.Role(name), role, role)
}

// DeleteRole deletes a given role.
func (client *Client) DeleteRole(name string) (*ErrorInfo, error) {
	return client.DeleteRoleContext(context.Background(), name)
}

// DeleteRoleContext deletes a given role with a context.
func (client *Client) DeleteRoleContext(
	ctx context.Context, name string,
) (*ErrorInfo, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}
	return client.callDelete(ctx, client.Endpoints.Role(name), nil, nil)
}
