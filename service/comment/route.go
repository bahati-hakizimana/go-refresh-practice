package comment

import (
	"encoding/json"
	"net/http"

	"github.com/go-refresh-practice/go-refresh-course/middleware"
	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.CommentStore
}

func NewHandler(store types.CommentStore) *Handler {
	return &Handler{store: store}
}

func(h *Handler) RegisterRoutes(router *mux.Router) {
	router.Handle("/comments",
	middleware.AuthMiddleware(middleware.AdminOnly(http.HandlerFunc(h.handleGetComments))),
).Methods(http.MethodGet)

// create comment

router.Handle("/comments",

http.HandlerFunc(h.handleCreateComments),
).Methods(http.MethodPost)
}


// Get all comments

func(h *Handler)handleGetComments(w http.ResponseWriter, r *http.Request) {
	comment, err := h.store.GetComments()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, comment)
}


func (h *Handler) handleCreateComments(w http.ResponseWriter, r *http.Request) {

	var payload types.CommentPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	comment := types.Comment{
		Name:    payload.Name,
		Comment: payload.Comment,
	}

	created, err := h.store.CreateComment(comment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, created)
}

