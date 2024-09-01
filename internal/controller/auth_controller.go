package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	contextkeys "github.com/avran02/kode/internal/context_keys"
	"github.com/avran02/kode/internal/dto"
	"github.com/avran02/kode/internal/service"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)

	AuthenticationMiddleware(next http.Handler) http.Handler
}

type authController struct {
	authService service.AuthService
}

func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := c.authService.Register(req.Username, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dto.LoginResponse{
		Token: token,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *authController) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing or malformed token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "missing or malformed token", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		userID, err := c.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func newAuthController(s service.AuthService) AuthController {
	return &authController{
		authService: s,
	}
}
