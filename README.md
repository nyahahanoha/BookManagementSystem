# BookManagementSystem

A book management system that uses barcodes for registration and management. It consists of a backend server, a web frontend, and a local scanner application.

## Overview

- **Web Interface**: Manage and view books via a web browser.
- **Barcode Scanning**: Register books using a local scanner application connected to a Bluetooth barcode scanner.
- **Backend API**: Handles business logic, data storage (MySQL), and integrations (Google Books API).

## Architecture

The system is composed of the following services aimed to be deployed with Docker Compose:

- **Proxy (`booksystem_proxy`)**: [Pomerium](https://www.pomerium.com/) handles authentication and reverse proxying.
- **Frontend (`booksystem_ui`)**: Built with [Deno Fresh](https://fresh.deno.dev/), providing the user interface.
- **Backend (`booksystem_api`)**: A Go-based API server.
- **Database (`booksystem_db`)**: MySQL database for storing book data.
- **Certbot (`certbot`)**: Automatically handles SSL certificates via Cloudflare DNS.

## Project Structure

- `backend/`: Go API server source code.
- `frontend/`: Deno Fresh web application source code.
- `scanner/`: Local scanner application source code (currently for macOS).
- `proxy/`: Pomerium configuration and certificate storage.
- `api/`: API definitions (Protobuf/gRPC).

## Configuration

### Backend Configuration (`backend/config.yaml`)
Configures the API server, database connection, and external services.

- **`books`**: Search settings (NDL, Google Books API).
- **`store`**: Data storage settings (MySQL, FileSystem).
- **`address`**: Server listening port (default `:8080`).
- **`admin_email`**: Administrator email list.
- **`pomerium_jwks_url`**: URL for Pomerium JWKS (for authentication verification).

### Scanner Configuration (`scanner/mac/config.yaml`)
Configures the local Bluetooth scanner application.

- **`bluetooth`**: Bluetooth device settings (Device Name, Service UUID, Characteristic UUID).
- **`api`**: The URL of the BookManagementSystem backend.
- **`callback_port`**: Local port for receiving temporary callbacks.

## Getting Started

### Prerequisites

- Docker & Docker Compose
- direnv (recommended for environment variable management)

### Installation & Setup

1.  **Environment Variables**:
    Copy `.envrc.example` to `.envrc` and fill in the required values (Database credentials, Google Books API key, domain settings, etc.).
    ```bash
    cp .envrc.example .envrc
    allow .envrc
    ```

2.  **Start Services**:
    Run the system using Docker Compose.
    ```bash
    docker-compose up -d
    ```

3.  **Access the Application**:
    Open your browser and navigate to `https://books.nyahahanoha.net` (or your configured domain).

## Local Scanner Usage

The scanner application (`scanner/mac`) runs locally to bridge a Bluetooth barcode scanner with the web API.

1.  **Configuration**:
    Edit `scanner/mac/config.yaml` to match your Bluetooth scanner's name and service UUIDs.
2.  **Run**:
    ```bash
    cd scanner/mac
    go run main.go
    ```
