package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mazzael/go-api/internal/dto"
	"github.com/Mazzael/go-api/internal/entity"
	"github.com/Mazzael/go-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	GormUserRepository database.UserRepository
	Jwt                *jwtauth.JWTAuth
	JwtExpiresIn       int
}

func NewUserHandler(repo database.UserRepository, jwt *jwtauth.JWTAuth, JwtExpiresIn int) *UserHandler {
	return &UserHandler{
		GormUserRepository: repo,
		Jwt:                jwt,
		JwtExpiresIn:       JwtExpiresIn,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user dto.UserLoginInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.GormUserRepository.FindByEmail(user.Email)
	if err != nil || !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := h.Jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.GormUserRepository.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
