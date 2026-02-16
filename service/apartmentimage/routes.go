package apartmentimage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	// GET all images PUBLIC (no auth required)
	router.HandleFunc("/apartment-images/public", h.handleGetApartmentImagesPublic).Methods(http.MethodGet)

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

	// DELETE image (admin only)
	router.Handle("/apartment-images/{id}",
		middleware.AuthMiddleware(middleware.AdminOnly(
			http.HandlerFunc(h.handleDeleteApartmentImage),
		)),
	).Methods(http.MethodDelete)

	// Serve uploads from Fly.io volume
	router.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/", http.FileServer(http.Dir("/root/uploads"))),
	)
}

// ----------------------------------------------------
// GET images PUBLIC (no authentication)
// ----------------------------------------------------
func (h *Handler) handleGetApartmentImagesPublic(w http.ResponseWriter, r *http.Request) {
	apartmentIDStr := r.URL.Query().Get("apartmentId")

	var images []types.ApartmentImage
	var err error

	if apartmentIDStr == "" {
		images, err = h.store.GetAllImages()
	} else {
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
// GET images for apartment (protected)
// ----------------------------------------------------
func (h *Handler) handleGetApartmentImages(w http.ResponseWriter, r *http.Request) {
	apartmentIDStr := r.URL.Query().Get("apartmentId")

	var images []types.ApartmentImage
	var err error

	if apartmentIDStr == "" {
		images, err = h.store.GetAllImages()
	} else {
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
	apartmentIDStr := strings.TrimSpace(r.FormValue("apartmentId"))
	if apartmentIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("apartmentId is required"))
		return
	}
	apartmentID, err := strconv.Atoi(apartmentIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("apartmentId must be an integer"))
		return
	}

	// Get file from form-data
	file, header, err := r.FormFile("imageFile")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image file is required"))
		return
	}
	defer file.Close()

	// Get caption
	caption := r.FormValue("caption")
	if caption == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("caption is required"))
		return
	}

	// Sanitize filename (remove unsafe characters)
	re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	sanitized := re.ReplaceAllString(header.Filename, "_")
	sanitized = strings.ToLower(sanitized)

	// Add timestamp to make filename unique
	timestamp := time.Now().Unix()
	sanitizedFilename := fmt.Sprintf("%d_%s", timestamp, sanitized)

	// Save to Fly.io volume
	dst := filepath.Join("/root/uploads", sanitizedFilename)
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

	// Construct frontend-accessible URL
	fullURL := fmt.Sprintf("https://%s/uploads/%s", os.Getenv("FLY_APP_NAME")+".fly.dev", sanitizedFilename)

	// Save in database
	img := types.ApartmentImage{
		ApartmentID: apartmentID,
		ImageURL:    fullURL,
		Caption:     caption,
	}

	newImage, err := h.store.CreateApartmentImage(img)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, newImage)
}

func (h *Handler) handleDeleteApartmentImage(w http.ResponseWriter, r *http.Request) {
	// Get image ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	imageID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid image ID"))
		return
	}

	// Delete from DB
	img, err := h.store.DeleteApartmentImage(imageID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Remove file from Fly.io volume
	if img.ImageURL != "" {
		// Extract filename from URL
		filename := filepath.Base(img.ImageURL)
		filepath := filepath.Join("/root/uploads", filename)
		_ = os.Remove(filepath) // ignore error if file already missing
	}

	utils.WriteJson(w, http.StatusOK, map[string]interface{}{
		"message": "Image deleted successfully",
		"image":   img,
	})
}