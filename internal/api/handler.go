package api

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/sementrof/prod1/docs"
	v1 "github.com/sementrof/prod1/internal/api/v1"
	auth "github.com/sementrof/prod1/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

type TaskServer struct {
	api v1.ApiInterface
}

func NewTaskServer(api v1.ApiInterface) *TaskServer {
	return &TaskServer{api: api}
}

func (s *TaskServer) createUserPost(w http.ResponseWriter, r *http.Request) {
	s.api.CreateUserPost(w, r)
}
func (s *TaskServer) LoginUserPost(w http.ResponseWriter, r *http.Request) {
	s.api.LoginUserPost(w, r)
}
func (s *TaskServer) CreateOrganizationPost(w http.ResponseWriter, r *http.Request) {
	s.api.CreateOrganizationPost(w, r)
}

// func (s *TaskServer) getUser(w http.ResponseWriter, r *http.Request) {
// 	s.api.GetUser(w, r)
// }

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
// SetupRouter настраивает маршрутизацию API
func SetupRouter(api v1.ApiInterface) *mux.Router {
	logger := auth.LoggerFactory()

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := NewTaskServer(api)
	router.HandleFunc("/registration/user", server.createUserPost).Methods("POST")
	router.HandleFunc("/registration/organization", server.CreateOrganizationPost).Methods("POST")
	router.HandleFunc("/login", server.LoginUserPost).Methods("POST")
	// router.HandleFunc("/get/", server.getUser).Methods("GET")
	http.HandleFunc("/protected", auth.Middleware(logger, auth.ProtectedHandler))

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	return router

}
