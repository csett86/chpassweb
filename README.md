# Conversational Password Change Web App


A single Go binary web application that provides a conversational password change interface, mapping PAM's conversational API 1:1 to HTTP POST requests. Uses in-memory sessions and pure HTML/CSS (no JavaScript).


## Features

- Conversational PAM flow: Each PAM prompt becomes an HTML form
- In-memory sessions: Tracks PAM conversations per user
- No JavaScript: Pure HTML/CSS forms
- Apache integration: Uses Apache for TLS and PAM authentication
- Single binary: Easy deployment
- Direct PAM messages: PAM prompts and errors are passed through verbatim


## Architecture

User -> Apache (PAM Auth + TLS) -> Go Web App (PAM chauthtok) -> PAM


1. Apache handles TLS termination and user authentication via PAM (mod_auth_pam)
2. Go App handles PAM chauthtok conversation with in-memory sessions


## Quick Start

### Build Dependencies
- libpam0g-dev (Debian/Ubuntu)

### Build
make build

### Run (as root)
sudo ./chpass-web

The server listens on 127.0.0.1:8080.

### Test
make test


## Deployment

### Apache Configuration
Use the provided apache.conf.example with mod_auth_pam.

### Docker
docker build -t chpass-web .
docker run -p 8080:8080 --privileged chpass-web


## Security
- Run as root (required for pam_chauthtok)
- Apache handles authentication
- HttpOnly, SameSite=Strict cookies
- 5-minute session timeout