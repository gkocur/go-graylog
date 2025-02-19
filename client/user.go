package client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/suzuki-shunsuke/go-graylog"
)

// CreateUser creates a new user account.
func (client *Client) CreateUser(user *graylog.User) (*ErrorInfo, error) {
	return client.CreateUserContext(context.Background(), user)
}

// CreateUserContext creates a new user account with a context.
func (client *Client) CreateUserContext(
	ctx context.Context, user *graylog.User,
) (*ErrorInfo, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}
	return client.callPost(ctx, client.Endpoints().Users(), user, nil)
}

// GetUsers returns all users.
func (client *Client) GetUsers() ([]graylog.User, *ErrorInfo, error) {
	return client.GetUsersContext(context.Background())
}

// GetUsersContext returns all users with a context.
func (client *Client) GetUsersContext(ctx context.Context) ([]graylog.User, *ErrorInfo, error) {
	users := &graylog.UsersBody{}
	ei, err := client.callGet(ctx, client.Endpoints().Users(), nil, users)
	return users.Users, ei, err
}

// GetUser returns a given user.
func (client *Client) GetUser(name string) (*graylog.User, *ErrorInfo, error) {
	return client.GetUserContext(context.Background(), name)
}

// GetUserContext returns a given user with a context.
func (client *Client) GetUserContext(
	ctx context.Context, name string,
) (*graylog.User, *ErrorInfo, error) {
	if name == "" {
		return nil, nil, errors.New("name is empty")
	}
	u, err := client.Endpoints().User(name)
	if err != nil {
		return nil, nil, err
	}
	user := &graylog.User{}
	ei, err := client.callGet(ctx, u.String(), nil, user)
	return user, ei, err
}

// UpdateUser updates a given user.
func (client *Client) UpdateUser(prms *graylog.UserUpdateParams) (*ErrorInfo, error) {
	return client.UpdateUserContext(context.Background(), prms)
}

// UpdateUserContext updates a given user with a context.
func (client *Client) UpdateUserContext(
	ctx context.Context, prms *graylog.UserUpdateParams,
) (*ErrorInfo, error) {
	if prms == nil {
		return nil, fmt.Errorf("user is nil")
	}
	if prms.Username == "" {
		return nil, errors.New("name is empty")
	}
	u, err := client.Endpoints().User(prms.Username)
	if err != nil {
		return nil, err
	}
	return client.callPut(ctx, u.String(), prms, nil)
}

// DeleteUser deletes a given user.
func (client *Client) DeleteUser(name string) (*ErrorInfo, error) {
	return client.DeleteUserContext(context.Background(), name)
}

// DeleteUserContext deletes a given user with a context.
func (client *Client) DeleteUserContext(
	ctx context.Context, name string,
) (*ErrorInfo, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}
	u, err := client.Endpoints().User(name)
	if err != nil {
		return nil, err
	}
	return client.callDelete(ctx, u.String(), nil, nil)
}
