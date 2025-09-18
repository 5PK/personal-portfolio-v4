package main

import (
	"fmt"
	agility "kevin-portfolio/internal/agility_cms"
	"kevin-portfolio/internal/handlers"
	"kevin-portfolio/views"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)


func init() {
    // Load .env file into environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on system env")
    }

	// Set the sitemap 
	agility.RefreshSitemap()
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get the home page content from CMS
		homeContent := agility.GetHomeCommandContent()
		views.Index(homeContent).Render(r.Context(), w)
	})

	http.HandleFunc("POST /terminal/command", handlers.ExecuteCommandHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
