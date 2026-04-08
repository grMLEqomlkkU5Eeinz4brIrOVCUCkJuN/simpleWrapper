package types

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"     validate:"required,email"`
	Name      string    `json:"name"      validate:"required,min=1"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetUserByIdResponse struct {
	User
}

type GetAllUsersResponse struct {
	Users []User
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name"  validate:"required,min=1"`
}

type UpdateUserRequest struct {
	ID    string `json:"id"              validate:"required"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
	Name  string `json:"name,omitempty"  validate:"omitempty,min=1"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}
