package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"storeX/utils"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()

	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseError(w, http.StatusOK, "server is running")
	})

	//api := srv.PathPrefix("/api/v1").Subrouter()

	return srv
}
