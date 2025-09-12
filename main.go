package main

import (
	"fmt"
	"net/http"
	"time"

	"kevin-portfolio/htmx-templ/views"
)

func main() {
	// Serve static files (CSS, JS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index("Kevin").Render(r.Context(), w)
	})

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<div id="time">%s</div>`, time.Now().Format(time.RFC1123))
	})

	fmt.Println("Listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
