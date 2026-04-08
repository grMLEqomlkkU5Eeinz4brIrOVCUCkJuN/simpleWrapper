package users_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jaxron/axonet/pkg/client"
	users "github.com/simpleWrapper/core/resources"
	"github.com/simpleWrapper/core/types"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *users.Users) {
	t.Helper()
	server := httptest.NewServer(handler)
	c := client.NewClient()
	u := users.NewUsers(c, server.URL)
	return server, u
}

func sampleUser() types.User {
	return types.User{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Email:     "john@example.com",
		Name:      "John Doe",
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestGetAllUsers(t *testing.T) {
	expected := []types.User{sampleUser()}
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/users" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	result, err := u.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 user, got %d", len(result))
	}
	if result[0].ID != expected[0].ID {
		t.Errorf("expected ID %s, got %s", expected[0].ID, result[0].ID)
	}
	if result[0].Email != expected[0].Email {
		t.Errorf("expected email %s, got %s", expected[0].Email, result[0].Email)
	}
}

func TestGetUserByID(t *testing.T) {
	expected := sampleUser()
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/users/"+expected.ID {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	result, err := u.GetByID(context.Background(), expected.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != expected.ID {
		t.Errorf("expected ID %s, got %s", expected.ID, result.ID)
	}
	if result.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, result.Name)
	}
}

func TestGetUserByID_EmptyID(t *testing.T) {
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach server")
	})
	defer server.Close()

	_, err := u.GetByID(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestCreateUser(t *testing.T) {
	expected := sampleUser()
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/users" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req types.CreateUserRequest
		json.Unmarshal(body, &req)
		if req.Email != "john@example.com" {
			t.Errorf("expected email john@example.com, got %s", req.Email)
		}
		if req.Name != "John Doe" {
			t.Errorf("expected name John Doe, got %s", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	result, err := u.Create(context.Background(), types.CreateUserRequest{
		Email: "john@example.com",
		Name:  "John Doe",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != expected.ID {
		t.Errorf("expected ID %s, got %s", expected.ID, result.ID)
	}
}

func TestCreateUser_ValidationError(t *testing.T) {
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach server")
	})
	defer server.Close()

	_, err := u.Create(context.Background(), types.CreateUserRequest{
		Email: "",
		Name:  "",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUpdateUser(t *testing.T) {
	expected := sampleUser()
	expected.Name = "Jane Doe"
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/users/"+expected.ID {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]string
		json.Unmarshal(body, &req)
		if req["name"] != "Jane Doe" {
			t.Errorf("expected name Jane Doe, got %s", req["name"])
		}
		if _, ok := req["email"]; ok {
			t.Error("email should not be in body when empty")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	})
	defer server.Close()

	result, err := u.Update(context.Background(), types.UpdateUserRequest{
		ID:   expected.ID,
		Name: "Jane Doe",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Jane Doe" {
		t.Errorf("expected name Jane Doe, got %s", result.Name)
	}
}

func TestUpdateUser_ValidationError(t *testing.T) {
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach server")
	})
	defer server.Close()

	_, err := u.Update(context.Background(), types.UpdateUserRequest{
		ID: "",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestDeleteUser(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/users/"+userID {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	err := u.Delete(context.Background(), types.DeleteUserRequest{
		ID: userID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUser_ValidationError(t *testing.T) {
	server, u := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not reach server")
	})
	defer server.Close()

	err := u.Delete(context.Background(), types.DeleteUserRequest{
		ID: "",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
