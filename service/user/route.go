package user

import (
	"fmt"
	"net/http"
	"strings"

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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
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

	// Trim email to avoid accidental space
	payload.Email = strings.TrimSpace(payload.Email)
	// check if user exists

	_,err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User with this email %s already exist", payload.Email))
		return
	}



	// validate
	 if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Envalid payload %v", error))
		return
	 }

	//  Hash password

	hashedPassword, err  := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Failed to hash password %v", err))
		return
	}

	// Create user

	user:= types.User{

		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: hashedPassword,
		Role: "user",

	}
	
	 err = h.store.CreateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return response

	response := map[string]interface{}{
		"message": "User registered successfully",
		"email":   user.Email,
		"firstName": user.FirstName,
		"lastName": user.LastName,

	}

	utils.WriteJson(w, http.StatusCreated, response)


}