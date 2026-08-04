package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"ticket-triage-api/interfaces"
	"ticket-triage-api/middleware"
	"ticket-triage-api/service"
	"ticket-triage-api/util"
)

type AuthController struct {
	authService interfaces.AuthService
}

func NewAuthController(authService interfaces.AuthService) *AuthController{
	return &AuthController{authService: authService}
}

type registerRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string){
	writeJSON(w, status, map[string]string{"error" : message})
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		writeError(w, http.StatusBadRequest, "email, password, and full name are required")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	customer, err := c.authService.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered){
			writeError(w, http.StatusConflict, "email already registered")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to register")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id" : customer.ID,
		"email" : customer.Email,
		"full_name" : customer.FullName,
	})
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	token, customer, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token" : token,
		"email" : customer.Email,
		"full_name" : customer.FullName,
	})
}

func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*util.Claims)

	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"email" : claims.Email,
	})
}

