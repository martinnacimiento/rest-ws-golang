package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"tincho.dev/rest-ws/dto"
)

type StructConstraint interface {
	dto.SignUpRequest | dto.SignInRequest | dto.CreatePostRequest | dto.UpdateOnePostRequest | dto.UpdateUserRequest
}

func Validate[T StructConstraint](r *http.Request) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(payload)

	if err != nil {
		return nil, err
	}

	return &payload, nil
}
