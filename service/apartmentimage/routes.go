package apartmentimage

import "github.com/go-refresh-practice/go-refresh-course/types"




type Handler struct {
	store types.ApartmentImagesStore
}

func NewHandler(store types.ApartmentImagesStore) *Handler {

	return &Handler{store: store}
}