package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"tincho.dev/rest-ws/dto"
	"tincho.dev/rest-ws/middlewares"
	"tincho.dev/rest-ws/models"
	"tincho.dev/rest-ws/repositories"
	"tincho.dev/rest-ws/server"
	"tincho.dev/rest-ws/utils"
)

func CreatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		payload, err := utils.Validate[dto.CreatePostRequest](r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request",
			})
			return
		}

		claims := r.Context().Value(middlewares.ClaimsKey).(*models.AppClaims)
		user, err := repositories.FindUserById(r.Context(), claims.UserId)
		fmt.Println(user)
		fmt.Println(claims.UserId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error finding user",
			})
			return
		}

		post := &models.Post{
			Title:   payload.Title,
			Content: payload.Content,
			UserID:  user.Id,
		}

		err = repositories.CreatePost(r.Context(), post)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error creating post",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

func FindAllPostsHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		posts, err := repositories.FindAllPosts(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

func FindOnePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		postIdParam := mux.Vars(r)["id"]
		postId, err := strconv.ParseInt(postIdParam, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request",
			})
			return
		}

		post, err := repositories.FindPostById(r.Context(), postId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error finding post",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func UpdateOnePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		postIdParam := mux.Vars(r)["id"]
		postId, err := strconv.ParseInt(postIdParam, 10, 64)

		if err != nil {
			fmt.Println(postIdParam, err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request",
			})
			return
		}

		payload, err := utils.Validate[dto.UpdateOnePostRequest](r)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		claims := r.Context().Value(middlewares.ClaimsKey).(*models.AppClaims)

		post, err := repositories.FindPostById(r.Context(), postId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error finding post",
			})
			return
		}

		if post.UserID != claims.UserId {
			fmt.Println(post.UserID, claims.UserId)
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Forbidden",
			})
			return
		}

		post.Title = payload.Title
		post.Content = payload.Content

		err = repositories.UpdatePost(r.Context(), post)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error updating post",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}

func DeleteOnePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		postIdParam := mux.Vars(r)["id"]
		postId, err := strconv.ParseInt(postIdParam, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request",
			})
			return
		}

		claims := r.Context().Value(middlewares.ClaimsKey).(*models.AppClaims)

		post, err := repositories.FindPostById(r.Context(), postId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error finding post",
			})
			return
		}

		if post.UserID != claims.UserId {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Forbidden",
			})
			return
		}

		err = repositories.DeletePost(r.Context(), post.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error deleting post",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Post deleted",
		})
	}
}
