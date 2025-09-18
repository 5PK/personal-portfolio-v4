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

	// Get the Sitemap
	var sitemap = agility.GetSitemapFlat()
	// Set the sitemap 
	agility.SetSitemap(sitemap)
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})

	http.HandleFunc("POST /terminal/command", handlers.ExecuteCommandHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
