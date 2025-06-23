# LinkUp — Social Web Messenger

LinkUp is a social network and web messenger application built with a Golang backend and PostgreSQL database, supporting features like user authentication, threads, comments, and likes.

## Tech Stack
- **Backend:** Golang (Gin framework)
- **Database:** PostgreSQL
- **Auth:** JWT, email verification via Gmail SMTP
- **Frontend:** Vue.js *(optional frontend integration)*
- **Infrastructure:** Docker

## Features
- User registration and login with JWT
- Email verification system
- Profile view and update
- Create, update, delete, and view threads
- Comment system under threads
- Like system on threads
- User search functionality

## Project Structure
```bash
cmd/
  main.go              # Entry point
internal/
  adapter/
    handler/           # HTTP handlers
    postgres/          # DB repositories
  app/service/         # Business logic
  model/               # Entities
  usecase/             # Usecases
pkg/
  middleware/          # JWT middleware
.env                    # Environment variables
Dockerfile              # Docker configuration
```

## How to Run
### Locally
```bash
git clone https://github.com/tamaqazaq/LinkUp.git
cd LinkUp
go mod download
go run cmd/main.go
```

### With Docker
```bash
docker build -t linkup .
docker run -p 8080:8080 linkup
```

## Environment Variables (.env)
```env
DATABASE_URL=
JWT_SECRET=
SMTP_USER=
SMTP_PASSWORD=
SMTP_HOST=
SERVER_URL=
CLIENT_URL=
```

## CORS Setup
Allowing origins:
- `http://localhost:5173`
- `https://linkup-9w5.pages.dev`
- Subdomains: `*.linkup-9w5.pages.dev`

## Contact
Developed by Daulet Yermukhanov — feel free to reach out via GitHub or Telegram!
