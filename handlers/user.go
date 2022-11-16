package handlers

import (
	"encoding/json"
	"francocorrea/go/rest-ws/models"
	"francocorrea/go/rest-ws/repositories"
	"francocorrea/go/rest-ws/server"
	"net/http"

	"github.com/segmentio/ksuid"
)

type SingUpRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type SingUpResponse struct {
	Id string `json:"id"`
	Email string `json:"email"`
}

func SingUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SingUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		id, err := ksuid.NewRandom()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Email: request.Email,
			Password: request.Password,
			Id: id.String(),
		}

		err = repositories.InsertUser(r.Context(), &user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SingUpResponse{
			Id: user.Id,
			Email: user.Email,
		})
	}
}