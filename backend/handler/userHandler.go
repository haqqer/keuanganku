package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/haqqer/keuanganku/repo"
	"github.com/haqqer/keuanganku/utils/response"
)

type UsersHandler struct{}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	googleID := r.PathValue("id")
	if googleID == "" {
		response.Error(w, http.StatusBadRequest, "ID empty")
		return
	}

	userRepo := repo.UserRepo{}
	userData, err := userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			response.Success(w, map[string]string{})
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, userData)
	// switch r.Method {
	// case "GET":
	// 	// Handle GET request
	// 	fmt.Fprintf(w, "List of users")
	// case "POST":
	// 	// Handle POST request
	// 	fmt.Fprintf(w, "Create a new user")
	// default:
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// }
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var p map[string]string
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(w, p)
	// switch r.Method {
	// case "GET":
	// 	// Handle GET request
	// 	fmt.Fprintf(w, "List of users")
	// case "POST":
	// 	// Handle POST request
	// 	fmt.Fprintf(w, "Create a new user")
	// default:
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// }
}
