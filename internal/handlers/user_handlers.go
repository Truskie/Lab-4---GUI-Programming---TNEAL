package handlers

import (
	"encoding/json"
	"lab2-terrylneal/internal/db"
	"lab2-terrylneal/internal/jsonhelper"
	"lab2-terrylneal/internal/models"
	"lab2-terrylneal/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type UserHandlers struct {
	repo *repository.UserRepository
}

func NewUserHandlers(database *db.DB) *UserHandlers {
	return &UserHandlers{
		repo: repository.NewUserRepository(database.DB),
	}
}

// CreateUser
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var user models.User
	if err := jsonhelper.ReadJSON(w, r, &user); err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	if err := h.repo.Create(&user); err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	}

	jsonhelper.WriteJSON(w, http.StatusCreated, response)
}

// GetUser
func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.repo.GetByID(id)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if user == nil {
		jsonhelper.WriteJSONError(w, http.StatusNotFound, "User not found")
		return
	}

	response := models.UserResponse{
		Success: true,
		Data:    user,
	}

	jsonhelper.WriteJSON(w, http.StatusOK, response)
}

// GetAllUsers
func (h *UserHandlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	onlyActive := r.URL.Query().Get("active") == "true"

	users, err := h.repo.GetAll(onlyActive)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	response := models.UserResponse{
		Success: true,
		Data:    users,
	}

	jsonhelper.WriteJSON(w, http.StatusOK, response)
}

// UpdateUser
func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := jsonhelper.ReadJSON(w, r, &user); err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	user.ID = id

	if err := h.repo.Update(&user); err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	}

	jsonhelper.WriteJSON(w, http.StatusOK, response)
}

// PartialUpdateUser
func (h *UserHandlers) PartialUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var updates map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updates); err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	user, err := h.repo.PartialUpdate(id, updates)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := models.UserResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	}

	jsonhelper.WriteJSON(w, http.StatusOK, response)
}

// DeleteUser
func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonhelper.WriteJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonhelper.WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// hard delete check
	hardDelete := r.URL.Query().Get("hard") == "true"

	if err := h.repo.Delete(id, hardDelete); err != nil {
		if err.Error() == "user not found" {
			jsonhelper.WriteJSONError(w, http.StatusNotFound, "User not found")
			return
		}
		jsonhelper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := "User deleted successfully"
	if hardDelete {
		message = "User permanently deleted successfully"
	} else {
		message = "User deactivated successfully"
	}

	response := models.UserResponse{
		Success: true,
		Message: message,
	}

	jsonhelper.WriteJSON(w, http.StatusOK, response)
}
