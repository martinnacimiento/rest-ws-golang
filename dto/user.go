package dto

import "tincho.dev/rest-ws/models"

// SignUp DTOs

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=8"`
}

type SignUpResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

// SignIn DTOs

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=8"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

// FindAllUsers DTOs

type FindAllUsersResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

// FindOneUser DTOs

type FindOneUserResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

// UpdateUser DTOs

type UpdateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

// GetMeData DTOs

type GetMeDataResponse struct {
	Id    int64         `json:"id"`
	Email string        `json:"email"`
	Posts []models.Post `json:"posts"`
}
