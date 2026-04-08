package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jaxron/axonet/pkg/client"
	"github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/simpleWrapper/userWrapper/core/types"
)

type Users struct {
	client   *client.Client
	validate *validator.Validate
	baseURL  string
}

func NewUsers(client *client.Client, baseURL string) *Users {
	return &Users{
		client:   client,
		validate: validator.New(),
		baseURL:  baseURL,
	}
}

func (u *Users) GetAll(ctx context.Context) ([]types.User, error) {
	var result []types.User
	resp, err := u.client.NewRequest().
		Method(http.MethodGet).
		URL(u.baseURL + "/users").
		Result(&result).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, nil
}

func (u *Users) GetByID(ctx context.Context, id string) (*types.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	var result types.User
	resp, err := u.client.NewRequest().
		Method(http.MethodGet).
		URL(u.baseURL + "/users/" + id).
		Result(&result).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &result, nil
}

func (u *Users) Create(ctx context.Context, req types.CreateUserRequest) (*types.User, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	var result types.User
	resp, err := u.client.NewRequest().
		Method(http.MethodPost).
		URL(u.baseURL + "/users").
		Header("Content-Type", "application/json").
		MarshalBody(req).
		Result(&result).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &result, nil
}

func (u *Users) Update(ctx context.Context, req types.UpdateUserRequest) (*types.User, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	body := make(map[string]string)
	if req.Email != "" {
		body["email"] = req.Email
	}
	if req.Name != "" {
		body["name"] = req.Name
	}

	var result types.User
	resp, err := u.client.NewRequest().
		Method(http.MethodPatch).
		URL(u.baseURL + "/users/" + req.ID).
		Header("Content-Type", "application/json").
		MarshalBody(body).
		Result(&result).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &result, nil
}

func (u *Users) Delete(ctx context.Context, req types.DeleteUserRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	resp, err := u.client.NewRequest().
		Method(http.MethodDelete).
		URL(u.baseURL + "/users/" + req.ID).
		Do(ctx)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
