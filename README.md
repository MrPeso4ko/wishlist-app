# Wishlist App Backend

A Go backend service with PostgreSQL for managing user wishlists, featuring authentication and public sharing.

## Features

- **User Management**
  - Registration with login/password
  - JWT authentication
  - Password hashing

- **Wishlist Functionality**
  - Create, read, update, delete wishes
  - Optional fields: comments, images, prices
  - Public view by username

- **Technical**
  - PostgreSQL database with GORM
  - Automatic migrations
  - HTTP API with Gin
  - Logging and metrics
  - Docker support

## Quick Start

### Prerequisites
- Go 1.19+
- PostgreSQL 13+
- Docker (optional)

### Local Development
1. Clone repository
2. Configure database in `config.yaml`
3. Install dependencies and run:
```bash
make deps && make run
```

### Docker
```bash
docker-compose up --build
```

## API Documentation

### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - Login and get JWT token

### Wishes
- `GET /api/wishes/:username` - Public view
- `POST /api/wishes` - Create new (authenticated)
- `PUT /api/wishes/:id` - Update (authenticated)
- `DELETE /api/wishes/:id` - Delete (authenticated)
- `GET /api/wishes` - User's wishes (authenticated)

## Testing
Run unit and integration tests:
```bash
make test
```

## Deployment
1. Build production image:
```bash
docker build -t wishlist-app .
```
2. Deploy with database:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Configuration
See `.env.example` for:
- Database connection
- JWT settings
- Log levels