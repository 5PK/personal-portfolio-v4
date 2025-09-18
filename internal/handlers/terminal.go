package handlers

import (
	agility "kevin-portfolio/internal/agility_cms"
	"kevin-portfolio/views/partials"
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

	// Parse global sitemap
	
	sm := agility.GetSitemapFlat()
	
	for page := range sm {
		
	}

	switch command {
	case "help":
		partials.Help().Render(r.Context(), w)
		return
	case "home":
		partials.Home().Render(r.Context(), w)
		return
	case "travelpics":
		partials.TravelPics().Render(r.Context(), w)
		return
	case "gear":
		partials.Gear().Render(r.Context(), w)
		return
	case "clear":
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusOK)
		return
	default:
		partials.BadInput(command).Render(r.Context(), w)
	}
}
