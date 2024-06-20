# Project Setup and Launch Instructions

## Prerequisites

1. Create a `.env` file based on the example provided in `.env.example`.
2. Download and install [ngrok](https://ngrok.com/).
3. Run `ngrok http TELEGRAM_PORT` and add the resulting URL to your `.env` file under the key `TELEGRAM_DOMAIN`.

## Launch with Docker Compose

1. Build and start the containers:
   ```bash
   docker compose up --detach --build
   ```

2. Run the telegram bot:
   ```bash
   docker compose exec ten_tarot_bot go run ./cmd/telegram_bot
   ```
   
3. Run the http server:
   ```bash
   docker compose exec ten_tarot_server go run ./cmd/http_server
   ```
   
   The HTTP server will be running on the port specified by `HTTP_PORT`, and the Telegram bot will be running on the
   port specified by `TELEGRAM_PORT`.
