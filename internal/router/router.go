package router

import (
	"net/http"

	"github.com/avran02/kode/internal/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	swagger "github.com/swaggo/http-swagger"
)

type Router struct {
	c controller.Controller
	*chi.Mux
}

func (r *Router) getAuthRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/register", r.c.Register)
	router.Post("/login", r.c.Login)

	return router
}

func (r *Router) getNoteRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(r.c.AuthenticationMiddleware)

	router.Post("/", r.c.CreateNote)
	router.Get("/", r.c.GetNotes)

	return router
}

func New(c controller.Controller) Router {
	router := Router{
		c: c,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/docs/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/openapi.yml")
	})
	r.Get("/swagger/*", swagger.Handler(
		swagger.URL("/docs/openapi.yml"),
	))
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})

	r.Mount("/", router.getAuthRoutes())
	r.Mount("/notes", router.getNoteRoutes())

	return Router{
		c:   c,
		Mux: r,
	}
}
