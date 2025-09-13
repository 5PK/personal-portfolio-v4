# Binaries
TAILWIND = ./tailwindcss
GO       = go
AIR      = air
TEMPL    = templ

# Paths
INPUT_CSS  = assets/input.css
OUTPUT_CSS = assets/output.css

# Default target
all: build

# Build Tailwind CSS once
css:
	$(TAILWIND) -i $(INPUT_CSS) -o $(OUTPUT_CSS)

# Watch Tailwind (auto rebuild on changes)
css-watch:
	$(TAILWIND) -i $(INPUT_CSS) -o $(OUTPUT_CSS) --watch

# Generate templ files once
templ-generate:
	$(TEMPL) generate

# Watch templ files and regenerate on changes
templ-watch:
	$(TEMPL) generate --watch

# Run Go server normally
run:
	$(GO) run .

# Run with hot reload (if air is installed)
dev:
	# Run tailwind, templ, and air in background
	$(MAKE) css-watch & \
	$(MAKE) templ-watch & \
	$(AIR)

# Clean build artifacts
clean:
	rm -f $(OUTPUT_CSS)

