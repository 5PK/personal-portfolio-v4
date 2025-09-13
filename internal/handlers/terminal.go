package handlers

import (
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

	switch command {
	case "help":
		partials.Help().Render(r.Context(), w)
		return
	case "home":
		partials.Home().Render(r.Context(), w)
		return
	case "clear":
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusOK)
		return
	default:
		partials.BadInput(command).Render(r.Context(), w)
	}

	// For now, just return a simple response
	// w.Header().Set("Content-Type", "text/plain")
	// w.Write([]byte("Command received: " + command))
}
