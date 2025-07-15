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
	public.HandleFunc("/refresh", handler.Refresh).Methods("POST")

	protected := public.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/users", handler.GetUsers).Methods("GET")
	protected.HandleFunc("/user/{user-id}/timeline", handler.UserTimeline).Methods("GET")
	protected.HandleFunc("/assets", handler.GetAssets).Methods("GET")
	protected.HandleFunc("/asset/{asset-id}/timeline", handler.AssetTimeline).Methods("GET")
	protected.HandleFunc("/asset/stats", handler.AssetStats).Methods("GET")
	protected.HandleFunc("/dashboard", handler.Dashboard).Methods("GET")

	employeeRoutes := protected.NewRoute().Subrouter()
	employeeRoutes.Use(middleware.AuthRole(models.RoleAdmin, models.RoleEmployeeManager))
	employeeRoutes.HandleFunc("/user", handler.CreateUser).Methods("POST")
	employeeRoutes.HandleFunc("/user/{user-id}", handler.UpdateUserDetails).Methods("PUT")
	employeeRoutes.HandleFunc("/user/{user-id}/delete", handler.DeleteUser).Methods("DELETE")

	assetRoutes := protected.NewRoute().Subrouter()
	assetRoutes.Use(middleware.AuthRole(models.RoleAdmin, models.RoleAssetManager))
	assetRoutes.HandleFunc("/asset", handler.CreateAsset).Methods("POST")
	assetRoutes.HandleFunc("/asset/assign", handler.AssignAsset).Methods("POST")
	assetRoutes.HandleFunc("/asset/{asset-id}/unassign", handler.UnassignAsset).Methods("POST")
	assetRoutes.HandleFunc("/vendor", handler.CreateVendor).Methods("POST")
	assetRoutes.HandleFunc("/asset/{asset-id}/service", handler.Service).Methods("POST")
	assetRoutes.HandleFunc("/asset/{asset-id}/delete", handler.DeleteAsset).Methods("DELETE")

	adminOnly := protected.NewRoute().Subrouter()
	adminOnly.Use(middleware.AuthRole(models.RoleAdmin))
	adminOnly.HandleFunc("/user/role/change", handler.RoleChange).Methods("PUT")
	adminOnly.HandleFunc("/user/type/change", handler.TypeChange).Methods("PUT")
	return srv
}
