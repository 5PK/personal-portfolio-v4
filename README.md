# Kevin's Portfolio - Go Edition

A modern terminal-style personal portfolio website built with Go, templ, and HTMX.

## Features

- **Terminal Interface**: Interactive command-line style navigation
- **Dynamic Content**: HTMX-powered dynamic content loading without page refreshes
- **Travel Photos**: Random travel photo display
- **Gear Showcase**: Personal equipment and tools
- **Responsive Design**: Works on desktop and mobile devices

## Available Commands

- `help` - Show available commands
- `home` - Display home page content
- `travelpics` - Show random travel photos
- `gear` - Display personal gear and equipment
- `clear` - Clear the terminal output

## Tech Stack

- **Backend**: Go 1.23.4
- **Templating**: [templ](https://github.com/a-h/templ) - Type-safe Go templating
- **Frontend**: HTMX for dynamic interactions
- **Deployment**: Docker + Fly.io

## Project Structure

```
.
├── main.go                 # Main application entry point
├── internal/
│   ├── handlers/          # HTTP request handlers
│   │   └── terminal.go    # Terminal command handler
│   └── helpers/           # Utility functions
│       └── random_image.go # Random image selection
├── views/                 # Templ templates
│   ├── index_templ.go     # Main page template
│   └── partials/          # Partial templates
├── assets/                # Static files (CSS, JS, images)
├── Dockerfile             # Docker configuration
└── fly.toml              # Fly.io deployment config
```

## Getting Started

### Prerequisites

- Go 1.23.4 or later
- [templ](https://github.com/a-h/templ) CLI tool

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd personal-port-v4-go
```

2. Install dependencies:
```bash
go mod download
```

3. Generate templ files:
```bash
templ generate
```

4. Run the application:
```bash
go run main.go
```

The application will be available at `http://localhost:8080`

### Development

Use the provided Makefile for common tasks:
```bash
make build    # Build the application
make run      # Run the application
make clean    # Clean build artifacts
```

## Deployment

### Docker

Build and run with Docker:
```bash
docker build -t kevin-portfolio .
docker run -p 8080:8080 kevin-portfolio
```

### Fly.io

Deploy to Fly.io:
```bash
fly deploy
```

## Contributing

This is a personal portfolio project, but feel free to use it as inspiration for your own terminal-style portfolio!

## License

MIT License - feel free to use this code for your own projects.