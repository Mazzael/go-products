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

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	GormUserRepository database.UserRepository
}

func NewUserHandler(repo database.UserRepository) *UserHandler {
	return &UserHandler{
		GormUserRepository: repo}
}

// Login godoc
// @Summary      User login
// @Description  User login
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body     dto.UserLoginInput  true  "user credentials"
// @Success      200  {object}  dto.UserLoginOutput
// @Failure      401
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users/auth [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("token").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

	var user dto.UserLoginInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.GormUserRepository.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := dto.UserLoginOutput{
		AccessToken: tokenString,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}

	err = h.GormUserRepository.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)

		return
	}
	w.WriteHeader(http.StatusCreated)
}
