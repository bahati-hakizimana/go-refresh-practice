package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-refresh-practice/go-refresh-course/config"
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

	var payload types.LoginUserPayload
	if err := utils.PulseJson(r, &payload); err != nil {

		utils.WriteError(w, http.StatusBadRequest, err)

		return

	}

	if err := utils.Validate.Struct(payload); err != nil{
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not invalid email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User not found, invalid Email or Password"))
		return
	}
    secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJwt(secret, u.ID, u.Email, u.Role)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token, "email": u.Email, "firstName": u.FirstName, "lastName": u.LastName, "role": u.Role})





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