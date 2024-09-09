package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"tincho.dev/rest-ws/dto"
	"tincho.dev/rest-ws/middlewares"
	"tincho.dev/rest-ws/models"
	"tincho.dev/rest-ws/repositories"
	"tincho.dev/rest-ws/server"
	"tincho.dev/rest-ws/utils"
)

const (
	EXPIRATION_TIME = time.Hour * 24
)

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload, err := utils.Validate[dto.SignUpRequest](r)

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

		response := dto.SignUpResponse{
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

		response := make([]dto.FindAllUsersResponse, 0)

		for _, user := range users {
			response = append(response, dto.FindAllUsersResponse{
				Id:    user.Id,
				Email: user.Email,
			})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func ListUsersHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		offset, limit, err := utils.GetPagination(r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		users, err := repositories.ListUsers(r.Context(), offset, limit)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]dto.FindAllUsersResponse, 0)

		for _, user := range users {
			response = append(response, dto.FindAllUsersResponse{
				Id:    user.Id,
				Email: user.Email,
			})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func FindOneUserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userIdParam := mux.Vars(r)["id"]
		userId, err := strconv.ParseInt(userIdParam, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := repositories.FindUserById(r.Context(), userId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &dto.FindOneUserResponse{
			Id:    user.Id,
			Email: user.Email,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func SignInHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload, err := utils.Validate[dto.SignInRequest](r)

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

		response := &dto.SignInResponse{
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

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &dto.GetMeDataResponse{
			Id:    user.Id,
			Email: user.Email,
			Posts: user.Posts,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateUserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		claims := r.Context().Value(middlewares.ClaimsKey).(*models.AppClaims)

		payload, err := utils.Validate[dto.UpdateUserRequest](r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error:": "Invalid request",
			})
			return
		}

		user, err := repositories.FindUserById(r.Context(), claims.UserId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Email = payload.Email

		err = repositories.UpdateOneUser(r.Context(), user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &dto.UpdateUserResponse{
			Id:    user.Id,
			Email: user.Email,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteUserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		claims := r.Context().Value(middlewares.ClaimsKey).(*models.AppClaims)

		err := repositories.DeleteOneUser(r.Context(), claims.UserId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
