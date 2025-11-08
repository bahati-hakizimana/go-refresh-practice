package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-refresh-practice/go-refresh-course/service/auth"
	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"
)

type Handler struct {

	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/login", h.handlerLogin).Methods("POST")
	router.HandleFunc("/register", h.handlerRegister).Methods("POST")

}

func (h *Handler) handlerLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handlerRegister(w http.ResponseWriter, r *http.Request){

	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.PulseJson(r, &payload); err != nil {

		utils.WriteError(w, http.StatusBadRequest, err)
		return

	}
	// check if user exists

	_,err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User with this email %s already exist", payload.Email))
		return
	}


	// if not exist we craete new user

	// validate
	 if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Envalid payload %v", error))
		return
	 }

	hashedPassword, err  := auth.HashPassword(payload.Password)

	 err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, nil)


}