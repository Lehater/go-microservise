// handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-microservice/models"
	"go-microservice/services"
	"go-microservice/storage"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUserRoutes регистрирует все CRUD-роуты.
func RegisterUserRoutes(r *mux.Router, service *services.UserService) {
	h := NewUserHandler(service)

	r.HandleFunc("/users", h.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.GetUserByID).Methods(http.MethodGet)
	r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", h.DeleteUser).Methods(http.MethodDelete)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}
	writeJSON(w, users, http.StatusOK)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	user, err := h.service.GetByID(id)
	if err != nil {
		if err == storage.ErrUserNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}
	writeJSON(w, user, http.StatusOK)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, created, http.StatusCreated)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(id, user)
	if err != nil {
		if err == storage.ErrUserNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, updated, http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		if err == storage.ErrUserNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Вспомогательные функции

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id must be integer", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
