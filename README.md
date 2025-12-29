# Sigmo (Formerly Telmo)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Sigmo is a self-hosted web UI and API for managing ModemManager-based cellular modems.
It focuses on eSIM profile operations, SMS, USSD, and network control, and ships as a
single Go binary with an embedded Vue 3 frontend.

## Advertisement

If you do not have an eUICC yet, you can purchase one from [eSTK.me](https://store.estk.me?code=esimcyou)
and use the coupon code `esimcyou` to get 10% off. We recommend [eSTK.me](https://store.estk.me?code=esimcyou)
if you need to perform profile downloads on iOS devices.

If you require more than 1MB of storage to install multiple eSIM profiles, we recommend
[9eSIM](https://www.9esim.com/?coupon=DAMON). Use the coupon code `DAMON` to also receive 10% off.

## Features

- eSIM profile list, download (SM-DP+), enable, rename, and delete.
- SIM slot switching and modem settings (alias, MSS, compatibility mode).
- SMS conversations (list, send, delete) and USSD sessions.
- Network scan and manual registration.
- OTP login via notification providers (Telegram or HTTP).
- Optional SMS forwarding to the same notification channels.

## Architecture

- Go backend serving `/api/v1` and static UI from embedded `web/dist`.
- Vue 3 + Vite frontend under `web/`.

## Requirements

- Linux with ModemManager running on the system D-Bus.
- Access to modem device nodes (often root or proper udev rules).
- Go 1.25+ to build the backend.
- Bun to build the web UI.

## Quick Start

1. Copy the example config and update credentials:
   `cp configs/config.example.toml config.toml`
2. Build the web UI:
   `cd web && bun install && bun run build`
3. Build the backend:
   `go build -o sigmo ./`
4. Run:
   `./sigmo -config config.toml`
5. Open `http://localhost:9527`.

## Configuration

Sigmo reads a TOML config file (default `config.toml`, override with `-config`).
The file is also written back when you update modem settings in the UI, so it must
be writable by the Sigmo process.

```toml
[app]
  environment = "production"
  listen_address = "0.0.0.0:9527"
  auth_providers = ["telegram"]
  otp_required = true

[channels]
  [channels.telegram]
    bot_token = "Your Telegram Bot Token"
    recipients = [123456789]

  [channels.http]
    endpoint = "https://httpbin.org/post"
    headers = { "Content-Type" = "application/json", "Authorization" = "Bearer 1234567890" }
```

Notes:

- `app.auth_providers` selects which channels are allowed for OTP login.
- `channels.*` are also used for SMS forwarding. If no channels are configured, OTP
  login and SMS forwarding are disabled.
- `modems` is keyed by ModemManager EquipmentIdentifier.
- `compatible` enables legacy modem restarts after profile changes.
- `mss` controls the APDU payload size per transfer (64-254).

## Development

- Backend: `go run ./ -config config.toml`
- Frontend dev server:
  `cd web && bun install && bun run dev`
- More frontend details in `web/README.md`.

## Service

Sample service definitions are in `init/systemd/sigmo.service` and
`init/supervisor/supervisord.conf`.

## License

MIT. See `LICENSE`.
