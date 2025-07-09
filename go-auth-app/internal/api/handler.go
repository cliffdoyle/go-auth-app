package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/cliffdoyle/go-auth-app/internal/service"
	"net/http"
)

type AuthHandler struct {
	userService service.UserService
	jwtSecret   string
	jwtExpiry   int
	validate    *validator.Validate
}

func NewAuthHandler(userService service.UserService, jwtSecret string, jwtExpiry int) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload service.SignupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.validate.Struct(payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}
	user, err := h.userService.RegisterUser(payload)
	if err != nil {
		h.respondWithError(w, http.StatusConflict, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload service.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.validate.Struct(payload); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}
	token, err := h.userService.LoginUser(payload, h.jwtSecret, h.jwtExpiry)
	if err != nil {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) GetUserDashboard(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(UserContextKey).(jwt.MapClaims)
	email := claims["email"].(string)
	data := map[string]string{"message": "Welcome to your user dashboard!", "email": email}
	h.respondWithJSON(w, http.StatusOK, data)
}

func (h *AuthHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(UserContextKey).(jwt.MapClaims)
	email := claims["email"].(string)
	data := map[string]interface{}{"message": "Welcome to the Admin Dashboard!", "admin_email": email, "secret_info": "This is top secret data for admins only."}
	h.respondWithJSON(w, http.StatusOK, data)
}