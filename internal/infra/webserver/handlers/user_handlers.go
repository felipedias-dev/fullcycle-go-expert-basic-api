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

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GetJWT godoc
// @Summary 		Get user JWT
// @Description Get user JWT
// @Tags 				users
// @Accept  		json
// @Produce  		json
// @Param 			request		body			dto.GetJWTInput true "user credentials"
// @Success 		200 			{object}	dto.GetJWTOutput
// @Failure 		400 			{object}	Error
// @Failure 		401 			{object}	Error
// @Failure 		404 			{object}	Error
// @Failure 		500 			{object}	Error
// @Router 			/users/auth [post]
func (h *UserHandler) GetJWTHandler(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiration := r.Context().Value("jwtExpiration").(int)
	var input dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	user, err := h.UserDB.FindByEmail(input.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	if !user.ComparePassword(input.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	token, err := pkgEntity.GenerateJWT(user.ID.String(), jwt, jwtExpiration)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	accessToken := dto.GetJWTOutput{
		AccessToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary 		Create user
// @Description Create user
// @Tags 				users
// @Accept  		json
// @Produce  		json
// @Param 			request 		body			dto.CreateUserInput true "user request"
// @Success 		201
// @Failure 		400 				{object}	Error
// @Failure 		500 				{object}	Error
// @Router 			/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
