package kosample

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type UserResource struct {
	DB *sqlx.DB
}
type User struct {
	ID   string `json:"id" DB:"id"`
	Name string `json:"name" DB:"name"`
}

func NewUserResource(DB *sqlx.DB) *UserResource {
	return &UserResource{
		DB: DB,
	}
}

func (ur *UserResource) Router() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Get("/user", ur.List)
	r.Get("/liveness", ur.Health)
	r.Get("/readiness", ur.DeepHealth)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	return r
}

func (ur *UserResource) List(w http.ResponseWriter, r *http.Request) {
	user := []User{}
	ur.DB.Select(&user, "SELECT id, name FROM users")

	render.JSON(w, r, user)
}

func (ur *UserResource) Health(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}
	render.JSON(w, r, data)
}

func (ur *UserResource) DeepHealth(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := ur.DB.Ping()
	if err != nil {
		data = map[string]string{
			"status": "error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, data)
		return
	}
	data = map[string]string{
		"status": "ok",
	}
	render.JSON(w, r, data)
}
