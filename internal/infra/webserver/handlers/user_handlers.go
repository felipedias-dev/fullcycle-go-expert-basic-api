package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/dto"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/database"
	pkgEntity "github.com/felipedias-dev/fullcycle-go-expert-basic-api/pkg/entity"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB        database.UserInterface
	Jwt           *jwtauth.JWTAuth
	JwtExpiration int
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiration int) *UserHandler {
	return &UserHandler{
		UserDB:        db,
		Jwt:           jwt,
		JwtExpiration: jwtExpiration,
	}
}

func (h *UserHandler) GetJWTHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.UserDB.FindByEmail(input.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !user.ComparePassword(input.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := pkgEntity.GenerateJWT(user.ID.String(), h.Jwt, h.JwtExpiration)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
