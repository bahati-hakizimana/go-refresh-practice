package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/go-refresh-practice/go-refresh-course/service/apartmentimage"
	"github.com/go-refresh-practice/go-refresh-course/service/aprtment"
	"github.com/go-refresh-practice/go-refresh-course/service/booking"
	"github.com/go-refresh-practice/go-refresh-course/service/comment"
	"github.com/go-refresh-practice/go-refresh-course/service/payments"
	"github.com/go-refresh-practice/go-refresh-course/service/user"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	router := mux.NewRouter()
	
	// Serve static files from uploads directory (add this BEFORE subrouter)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Apartment routes
	apartmentStore := aprtment.NewStore(s.db)
	apartmentHandler := aprtment.NewHandler(apartmentStore)
	apartmentHandler.RegisterRoutes(subrouter)

	// Apartment images routes
	apartmentImagesStore := apartmentimage.NewStore(s.db)
	apartmentImagesHandler := apartmentimage.NewHandler(apartmentImagesStore)
	apartmentImagesHandler.RegisterImageRoutes(subrouter)

	// Booking routes
	bookingStore := booking.NewStore(s.db)
	bookingHandler := booking.NewHandler(bookingStore)
	bookingHandler.RegisterRoutes(subrouter)

	commentStore := comment.NewStore(s.db)
	commentHandler := comment.NewHandler(commentStore)
	commentHandler.RegisterRoutes(subrouter)

	// Payment routes with Pasis integration
	paymentStore := payments.NewStore(s.db)

	log.Printf("DEBUG: PasisAppKey loaded: %v (length: %d)", 
    config.Envs.PasisAppKey != "", 
    len(config.Envs.PasisAppKey))
log.Printf("DEBUG: PasisSecretKey loaded: %v (length: %d)", 
    config.Envs.PasisSecretKey != "", 
    len(config.Envs.PasisSecretKey))
	
	pasisClient := payments.NewPasisClient(
		config.Envs.PasisAppKey,
		config.Envs.PasisSecretKey,
	)
	paymentHandler := payments.NewHandler(paymentStore, pasisClient)
	paymentHandler.RegisterRoutes(subrouter)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
	})

	// Wrap router with CORS middleware
	handler := c.Handler(router)

	log.Println("Listen on", s.addr)
	log.Println("Serving static files from ./uploads")
	return http.ListenAndServe(s.addr, handler)
}
