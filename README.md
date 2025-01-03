# Event Planner

A lightweight web application written in Go that allows users to manage events through a simple CRUD interface. The project serves as a practical implementation of essential web application features using Go's standard library and minimal dependencies.

## Screenshots

Here's a visual walkthrough of the application:

### Home Page
![Home Page Screenshot](ui/assets/img/home.png)
The home page provides quick access to event management and user authentication.

### Events List
![Events List Screenshot](ui/assets/img/events.png)
View and manage all your events in one place, with options to create, edit, and delete events.

### Event Details
![Event Details Screenshot](ui/assets/img/view.png)
Detailed view of individual events showing all relevant information.

## Features

- **Event Management**
    - Create, read, update, and delete events
    - View event details and listings

- **User Authentication & Security**
    - User registration and login
    - Session management
    - CSRF protection
    - Secure password handling

## Dependencies

### Backend Dependencies
- **github.com/alexedwards/scs/v2**: Session management middleware for secure user sessions
- **github.com/alexedwards/scs/sqlite3store**: SQLite3 session store for SCS
- **github.com/go-playground/form/v4**: Form decoder and validator for processing HTTP form data
- **github.com/justinas/alice**: Middleware chaining for clean HTTP handler composition
- **github.com/justinas/nosurf**: CSRF protection middleware
- **github.com/mattn/go-sqlite3**: SQLite3 database driver
- **golang.org/x/crypto/bcrypt**: Password hashing and verification

### Frontend Dependencies
- **TailwindCSS**: Utility-first CSS framework for styling
    - **@tailwindcss/forms**: Form styling plugin
    - **@tailwindcss/typography**: Typography styling plugin
- **Iconify**: Icon system with multiple icon sets
    - Used via CDN for UI icons (material-symbols, lucide, etc.)

### Development Tools
- **Air**: Live reload for Go applications during development
- **Browser-sync**: Browser synchronization and auto-reload
- **Just**: Command runner for development tasks
- **Goose**: Database migration tool

## Tech Stack

- **Backend**
    - Go (Golang) for server-side logic
    - SQLite for data persistence
    - Goose for database migrations

- **Frontend**
    - HTML templates
    - TailwindCSS for styling
    - JavaScript for enhanced interactivity

- **Development Tools**
    - Air for live reloading
    - Browser-sync for development server
    - Just for task automation
    - Docker for containerization

## Getting Started

### Running with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/madalinpopa/go-event-planner.git && cd go-event-planner
   ```

2. Build the Docker image:
   ```bash
   docker build -t go-event-planner .
   ```

3. Run the container:
   ```bash
   docker run -p 4000:4000 \
       -v $(pwd)/database:/app/database \
       go-event-planner
   ```

The application will be available at `http://localhost:4000`

### Development Setup

1. Install dependencies:
    - Go 1.x or higher
    - TailwindCSS CLI
    - Just command runner
    - Air (for live reloading)
    - Browser-sync
    - Goose (for database migrations)

2. Install project dependencies:
   ```bash
   go mod download
   ```

3. Run database migrations:
   ```bash
   just migrate up
   ```

4. Start the development server:
   ```bash
   just dev
   ```

This will start:
- The Go server with live reloading
- TailwindCSS compiler in watch mode
- Browser-sync for automatic browser refreshing

## Development Commands

- Update Go dependencies: `just update`
- Build production CSS: `just build`
- Run migrations: `just migrate [command]`
- Create new migration: `just makemigrations [name]`

## Project Structure

```
.
├── cmd/
│   └── web/                    # Web application code
│       ├── context.go          # Request context definitions
│       ├── forms.go            # Form handling and validation
│       ├── handlers.go         # HTTP request handlers
│       ├── helpers.go          # Helper functions
│       ├── main.go             # Application entry point
│       ├── middleware.go       # HTTP middleware
│       ├── routes.go           # Route definitions
│       └── templates.go        # Template handling
├── database/                   # Database related files
│   ├── events.db               # SQLite database
│   └── migrations/             # Database migrations
├── internal/                   # Private application packages
│   ├── models/                 # Data models
│   │   ├── errors.go
│   │   ├── event.go
│   │   └── user.go
│   └── validator/              # Validation logic
│       └── validator.go
├── ui/                         # User interface related code
│   ├── assets/                 # Source assets
│   │   └── input.css           # TailwindCSS input file
│   ├── embed.go                # File embedding for Go
│   ├── html/         
│   │   ├── base.tmpl           # Base template
│   │   ├── pages/              # Page templates
│   │   │   ├── auth/           # Authentication pages
│   │   │   ├── events/         # Event management pages
│   │   │   └── home.tmpl       # Homepage
│   │   └── partials/           # Reusable template parts
│   │       ├── components/     # UI components
│   │       └── forms/          # Form templates
│   └── static/                 # Static assets
│       ├── css/                # Compiled CSS
│       ├── fonts/              # Custom fonts
│       ├── img/                # Images
│       └── js/                 # JavaScript files
├── Dockerfile                  # Docker configuration
├── go.mod                      # Go module definition
├── justfile                    # Task runner commands
└── tailwind.config.js          # TailwindCSS configuration
```

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Future Development

Some ideas for future enhancements:

- Add event categories and tags
- Implement user roles and permissions
- Add event sharing functionality
- Create API endpoints for programmatic access
- Add event reminders and notifications
- Implement event search and filtering
- Add event location with map integration
- Create a REST API
- Add unit tests and integration tests