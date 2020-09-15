package app

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRouter(a App) http.Handler {
	router := mux.NewRouter()

	router.Handle("/subscribers", manageSubscriptionHandler{
		service: a.subscriberService,
	}).Methods("POST")

	// wrap the handlers in CORS
	origins := []string{
		"*",
	}

	return handlers.CORS(handlers.AllowedOrigins(origins))(router)
}
