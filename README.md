# Project Setup and Launch Instructions

## Prerequisites

1. Create a `.env` file based on the example provided in `.env.example`.
2. Download and install [ngrok](https://ngrok.com/).
3. Run `ngrok http TELEGRAM_PORT` and add the resulting URL to your `.env` file under the key `TELEGRAM_DOMAIN`.

## Launch with Docker Compose

1. Build the bot:
   ```bash
   make build
   ```

2. Start the bot:
   ```bash
   make start
   ```

   The HTTP server will be running on the port specified by `HTTP_PORT`, and the Telegram bot will be running on the
   port specified by `TELEGRAM_PORT`.
