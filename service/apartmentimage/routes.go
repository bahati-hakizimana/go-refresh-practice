package apartmentimage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-refresh-practice/go-refresh-course/middleware"
	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/go-refresh-practice/go-refresh-course/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ApartmentImagesStore
}

func NewHandler(store types.ApartmentImagesStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterImageRoutes(router *mux.Router) {

    // GET all images (auth required)
    router.Handle("/apartment-images",
        middleware.AuthMiddleware(http.HandlerFunc(h.handleGetApartmentImages)),
    ).Methods(http.MethodGet)

    // POST image (admin only)
    router.Handle("/apartment-images",
        middleware.AuthMiddleware(middleware.AdminOnly(
            http.HandlerFunc(h.handleAddApartmentImage),
        )),
    ).Methods(http.MethodPost)
}


// ----------------------------------------------------
// GET images for apartment
// ----------------------------------------------------

func (h *Handler) handleGetApartmentImages(w http.ResponseWriter, r *http.Request) {
    // Optional query param for future filtering
    apartmentIDStr := r.URL.Query().Get("apartmentId")

    var images []types.ApartmentImage
    var err error

    if apartmentIDStr == "" {
        // Get all images
        images, err = h.store.GetAllImages()
    } else {
        // Future: filter by apartmentId
        apartmentID, convErr := strconv.Atoi(apartmentIDStr)
        if convErr != nil {
            utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("apartmentId must be an integer"))
            return
        }
        images, err = h.store.GetImagesByApartmentID(apartmentID)
    }

    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

    utils.WriteJson(w, http.StatusOK, images)
}



// ----------------------------------------------------
// POST add image to apartment
// ----------------------------------------------------

func (h *Handler) handleAddApartmentImage(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form
    if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid form data"))
        return
    }

    // Get apartment ID
    apartmentIDStr := r.FormValue("apartmentId")
    apartmentID, err := strconv.Atoi(apartmentIDStr)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("apartmentId must be an integer"))
        return
    }

    // Get file from form-data
    file, header, err := r.FormFile("imageFile") // <-- FormData key
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image file is required"))
        return
    }
    defer file.Close()

    // Optional: Get caption
    caption := r.FormValue("caption")
    if caption == "" {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("caption is required"))
        return
    }

    // Save file locally (or upload to S3)
    dst := "./uploads/" + header.Filename
    out, err := os.Create(dst)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unable to save file"))
        return
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unable to save file"))
        return
    }

    // Save in database
    img := types.ApartmentImage{
        ApartmentID: apartmentID,
        ImageURL:    "/uploads/" + header.Filename, // path or public URL
        Caption:     caption,
    }

    newImage, err := h.store.CreateApartmentImage(img)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, err)
        return
    }

    utils.WriteJson(w, http.StatusCreated, newImage)
}



