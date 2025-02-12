# ShortLinks

## Description

This project is a URL shortener service that provides an API for creating shortened URLs in the following format:
- Each shortened URL must be unique, and for a given original URL, only one shortened URL should exist.
- The shortened URL must be 10 characters long.
- The shortened URL must consist of uppercase and lowercase Latin letters, digits, and the underscore (_) character.

## Technology Stack

- **Language:** Go
- **Database:** PostgreSQL (or an in-memory storage alternative)
- **Containerization:** Docker

## API Endpoints

The service provides the following HTTP endpoints:

1. **POST /shorten**
   - Saves the original URL in the database and returns a shortened URL.
   
2. **GET /{shortened_url}**
   - Accepts a shortened URL and returns the original URL.

## Deployment

The service is designed to be deployed as a **Docker image**.

## Storage Options

The service supports two storage implementations, specified as a parameter at runtime:
1. **PostgreSQL**
2. **In-memory storage** (custom implementation for storing links within the application)

## Testing

The project includes **unit tests** to ensure the correctness of the implemented functionality.

---

