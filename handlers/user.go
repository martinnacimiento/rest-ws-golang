package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"tincho.dev/rest-ws/middlewares"
	"tincho.dev/rest-ws/models"
	"tincho.dev/rest-ws/repositories"
	"tincho.dev/rest-ws/server"
)

const (
	EXPIRATION_TIME = time.Hour * 24
)

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=8"`
}

type SignUpResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=8"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var payload SignUpRequest
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request",
			})

			return
		}

		validate := validator.New()
		err = validate.Struct(payload)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

			return
		}

		_, err = repositories.GetUserByEmail(r.Context(), payload.Email)

		if err == nil {
			w.WriteHeader(http.StatusConflict)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email already exists",
			})

			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user := &models.User{
			Email:    payload.Email,
			Password: string(hashedPassword),
		}

		err = repositories.CreateUser(r.Context(), user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		}
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func FindAllUsersHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		users, err := repositories.FindAllUsers(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

func SignInHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var payload SignInRequest
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request",
			})

			return
		}

		validate := validator.New()
		err = validate.Struct(payload)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

			return
		}

		user, err := repositories.GetUserByEmail(r.Context(), payload.Email)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(EXPIRATION_TIME).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(s.Config().JWTSecret))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &SignInResponse{
			Token: signedToken,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		claims := r.Context().Value(middlewares.ClaimsKey).(*models.AppClaims)

		user, err := repositories.FindUserById(r.Context(), claims.UserId)

		fmt.Println(user, err)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
