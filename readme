# Notifications Server

This Go web server showcases Server-Sent Events (SSE) for real-time notifications. It utilizes the [Chi](https://github.com/go-chi/chi) router and provides HTML templates for a generator and a display. Clients can connect to the `/events` endpoint to receive real-time notifications.

## Features

- **Server-Sent Events (SSE):** The server implements SSE to push real-time notifications to connected clients.

- **Generator and Display Templates:** The server serves HTML templates for a generator (Button to generate new event) and a display page. Clients can generate events through the generator, and the display page will receive and display notifications.

## Getting Started

1. **Install Dependencies:**
   ```bash
   go get -u github.com/go-chi/chi/v5
## Endpoints

- **`/events`:** SSE endpoint for real-time notifications.

- **`/generate`:** POST endpoint to generate events and push notifications.

- **`/generate`:** GET endpoint to render the generator template.

- **`/display`:** GET endpoint to render the display template.

## Usage

1. Open the generator page by visiting [http://localhost:3000/generate](http://localhost:3000/generate).

2. Open the display page by visiting [http://localhost:3000/display](http://localhost:3000/display).

3. Generate events using the generator page, and observe real-time notifications displayed on the display page.

## Dependencies

- [Chi Router](https://github.com/go-chi/chi): Lightweight, idiomatic router for building Go HTTP services.

