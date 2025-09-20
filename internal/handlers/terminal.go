package handlers

import (
	agility "kevin-portfolio/internal/agility_cms"
	"kevin-portfolio/views/partials"
	"log"
	"net/http"
	"strings"
)

func ExecuteCommandHandler(w http.ResponseWriter, r *http.Request) {
	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Command is required", http.StatusBadRequest)
		return
	}

	// Parse command and arguments
	parts := strings.Fields(command)
	if len(parts) == 0 {
		http.Error(w, "Invalid command", http.StatusBadRequest)
		return
	}

	if command == "clear" {
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusOK)
		return
	}

	commandRoute := "/terminal/commands/" + command

	// get global sitemap
	sm := agility.GetCurrentSitemap()

	for route, sitemapPage := range sm {
		// if the route is a part of terminal commands,
		log.Println("route" + route)
		log.Println("command" + commandRoute)
		if route == commandRoute {
			// Get the page from agility
			page := agility.GetPage(sitemapPage.PageID)
			agility.RenderPage(r.Context(), w, page)
			return
		}
	}

	// no route found
	partials.BadInput(command).Render(r.Context(), w)

}
