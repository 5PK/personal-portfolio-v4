package main

import (
	"fmt"
	"kevin-portfolio/internal/handlers"
	"kevin-portfolio/views"
	"net/http"
)

func main() {
	// Serve static files (CSS, JS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))
	http.Handle("/htmx/", http.StripPrefix("/htmx/", http.FileServer(http.Dir("htmx/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})

	http.HandleFunc("POST /terminal/command", handlers.ExecuteCommandHandler)
	// r.HandleFunc("/terminal/output", handlers.GetOutputHandler).Methods("GET")

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
