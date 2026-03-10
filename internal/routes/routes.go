package routes

import (
	"lab2-terrylneal/internal/db"
	"lab2-terrylneal/internal/handlers"
	"lab2-terrylneal/internal/middleware"
	"log"
	"net/http"
	"strings"
)

func SetupRoutes(mux *http.ServeMux, database *db.DB) *http.ServeMux {
	if mux == nil {
		mux = http.NewServeMux()
	}

	// HTML
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/about", handlers.About)
	mux.HandleFunc("/contact", handlers.Contact)
	mux.HandleFunc("/hobby", handlers.Hobby)

	// JSON
	mux.HandleFunc("/api/info", handlers.APIInfo)

	if database != nil {
		userHandlers := handlers.NewUserHandlers(database)

		mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimSuffix(r.URL.Path, "/")

			if path == "/api/users" {
				switch r.Method {
				case http.MethodGet:
					userHandlers.GetAllUsers(w, r)
				case http.MethodPost:
					userHandlers.CreateUser(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			} else if strings.HasPrefix(path, "/api/users/") {
				switch r.Method {
				case http.MethodGet:
					userHandlers.GetUser(w, r)
				case http.MethodPut:
					userHandlers.UpdateUser(w, r)
				case http.MethodPatch:
					userHandlers.PartialUpdateUser(w, r)
				case http.MethodDelete:
					userHandlers.DeleteUser(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			} else {
				http.NotFound(w, r)
			}
		})
	} else {
		log.Println("Warning: No database connection - API endpoints disabled")
	}

	return mux
}

func ApplyMiddleware(handler http.Handler) http.Handler {
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.TimingMiddleware(handler)
	return handler
}
