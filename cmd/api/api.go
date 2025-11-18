package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-refresh-practice/go-refresh-course/service/aprtment"
	"github.com/go-refresh-practice/go-refresh-course/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	 return &APIServer{
		addr: addr,
		db: db,
	 }
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1",).Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	apartmentStore := aprtment.NewStore(s.db)
	apartmentHandler := aprtment.NewHandler(apartmentStore) 
	apartmentHandler.RegisterRoutes(subrouter)
	log.Println("Listen on", s.addr)
	return http.ListenAndServe(s.addr, router)
}