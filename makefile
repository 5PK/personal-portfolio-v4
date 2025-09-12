# Binaries
TAILWIND = ./tailwindcss
GO       = go
AIR      = air

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

# Run Go server normally
run:
	$(GO) run .

# Run with hot reload (if air is installed)
dev:
	# Run tailwind in background, air for Go reload
	$(MAKE) css-watch & \
	$(AIR)

# Clean build artifacts
clean:
	rm -f $(OUTPUT_CSS)

