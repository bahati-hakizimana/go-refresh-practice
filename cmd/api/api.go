package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/go-refresh-practice/go-refresh-course/service/apartmentimage"
	"github.com/go-refresh-practice/go-refresh-course/service/aprtment"
	"github.com/go-refresh-practice/go-refresh-course/service/booking"
	"github.com/go-refresh-practice/go-refresh-course/service/payments"
	"github.com/go-refresh-practice/go-refresh-course/service/user"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
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

	// Payment routes with Pasis integration
	paymentStore := payments.NewStore(s.db)
	pasisClient := payments.NewPasisClient(
		config.Envs.PasisAppKey,
		config.Envs.PasisSecretKey,
	)
	paymentHandler := payments.NewHandler(paymentStore, pasisClient)
	paymentHandler.RegisterRoutes(subrouter)

	log.Println("Listen on", s.addr)
	return http.ListenAndServe(s.addr, router)
}