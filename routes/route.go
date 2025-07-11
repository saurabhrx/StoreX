package routes

import (
	"github.com/gorilla/mux"
	"net/http"
	"storeX/handler"
	"storeX/middleware"
	"storeX/models"
	"storeX/utils"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()

	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseError(w, http.StatusOK, "server is running")
	})

	public := srv.PathPrefix("/api/v1").Subrouter()
	public.HandleFunc("/user/login", handler.LoginUser).Methods("POST")

	protected := public.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	private := protected.NewRoute().Subrouter()
	private.Use(middleware.AuthRole(models.RoleAdmin))
	private.HandleFunc("/user", handler.CreateUser).Methods("POST")
	private.HandleFunc("/asset", handler.CreateAsset).Methods("POST")
	private.HandleFunc("/asset/assign", handler.AssignAsset).Methods("POST")
	public.HandleFunc("/employees", handler.GetEmployees).Methods("GET")

	return srv
}
