package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-refresh-practice/go-refresh-course/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {

	userStore := &mockUserStore{}
	handler := NewHandler(userStore)


	t.Run("Should fail if the user payload is invalid", func(t *testing.T) {

	payload := types.RegisterUserPayload{
		FirstName: "user",
		LastName: "mam",
		Email:"invalid",
		Password: "scscs",
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()

   router.HandleFunc("/register", handler.handlerRegister)	
   router.ServeHTTP(rr, req)

   if rr.Code != http.StatusBadRequest{
	t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
   }
	
})
	t.Run("Should register user ", func(t *testing.T) {

	payload := types.RegisterUserPayload{
		FirstName: "user",
		LastName: "mam",
		Email:"valid@gmail.com",
		Password: "scscs",
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()

   router.HandleFunc("/register", handler.handlerRegister)	
   router.ServeHTTP(rr, req)

   if rr.Code != http.StatusCreated{
	t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
   }
	
})
}





type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail (email string)(*types.User, error){

	return nil, fmt.Errorf("User not found")
}
func (m *mockUserStore) GetUserById (id int)(*types.User, error){

	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
} 