# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Command line tool to interact with [Barong](https://github.com/walletera/barong) APIs. Barong is an open-source authentication/identity platform.

The tool is written in Go and uses [Cobra](https://cobra.dev/) as the CLI framework. The tool will have one command per Barong API. The available APIs and their documentation can be found in the `barong-docs/` directory.

Under the `pkg/` folder there will be one package for each API.

## Commands

```bash
# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./pkg/<package> -run TestName

# Lint (assumes golangci-lint is installed)
golangci-lint run

# Run the CLI
go run main.go <command> [flags]
```

## Architecture

### Directory Structure

```
barong-cli/
├── main.go              # Entry point, initializes root Cobra command
├── cmd/
│   ├── root.go          # Root Cobra command and --url flag
│   ├── user/            # User API commands (one file per command group)
│   │   ├── user.go      # Registers subcommands under "user"
│   │   ├── auth.go      # newAuthenticatedClient helper
│   │   ├── login.go
│   │   ├── logout.go
│   │   ├── create.go
│   │   ├── me.go
│   │   └── otp.go
│   └── management/      # Management API commands
│       ├── management.go        # Registers subcommands under "management"; --key-id and --private-key-file persistent flags
│       ├── auth.go              # newManagementClient helper (loads RSA key, creates pkg/management.Client)
│       ├── users.go
│       ├── labels.go
│       ├── profiles.go
│       ├── phones.go
│       ├── documents.go
│       ├── service_accounts.go
│       ├── otp.go
│       └── timestamp.go
├── pkg/
│   ├── user/            # User API client
│   │   ├── client.go    # HTTP client, Login/Logout/GetMe/OTP methods
│   │   └── models.go    # Response structs (UserWithFullInfo, OTPQRCode, …)
│   └── management/      # Management API client
│       ├── client.go    # HTTP client with JWT multisig auth; one method per endpoint
│       └── models.go    # Response structs (UserWithProfile, UserWithKYC, Label, Phone, …)
├── internal/
│   └── session/
│       └── session.go   # Save/Load/Delete session cookies (~/.barong-cli/session.json)
└── barong-docs/         # API reference docs and Swagger specs
    ├── barong_admin_api_v2.md
    ├── barong_user_api_v2.md
    ├── barong_management_api_v2.md
    ├── management-api-docs/     # Management API auth docs and JWT multisig reference
    └── swagger/
        ├── admin_api.json
        ├── user_api.json
        └── management_api.json
```

### Command Structure

Each Barong API maps to a top-level Cobra subcommand (e.g., `barong-cli user`, `barong-cli management`). The `cmd/` package wires Cobra commands to the API clients in `pkg/`.

### APIs

There are three Barong APIs (v2.7.0), each with Markdown docs and a Swagger/OpenAPI JSON spec in `barong-docs/swagger/`:

- **Admin API** (`admin_api.json`) — user management, document verification, KYC, user attributes (not yet implemented)
- **User API** (`user_api.json`) — session management, identity operations, OTP
- **Management API** (`management_api.json`) — label management, user/profile/phone/document/service-account operations

When implementing a new command, consult the corresponding Markdown doc and Swagger spec in `barong-docs/` for request/response shapes and authentication requirements.

## Authentication & Authorization Flow

### User API — session cookies

Endpoints under `/api/v1/auth/identity/` are public (no auth required). Endpoints under `/api/v1/auth/resource/` require a session cookie obtained through login.

**Step 1 — Login** (`POST /api/v1/auth/identity/sessions`)
Returns a session cookie that is persisted to `~/.barong-cli/session.json`.

**Step 2 — Authenticated requests**
Pass the session cookie directly on every `/api/v1/auth/resource/<path>` request. No token exchange step is needed.

Commands that need authentication use `newAuthenticatedClient` in `cmd/user/auth.go`, which loads the session cookies and passes them to `pkg/user.NewAuthenticatedClient`. The `post` and `get` helpers in `pkg/user/client.go` attach the cookies automatically when `authenticated` is true.

`internal/session` persists cookies as JSON, including their `Expires` field. `session.Load()` returns an error if any cookie is expired, prompting the user to log in again.

### Management API — JWT multisig

Every Management API request is sent as an HTTP POST (or PUT for label updates) with `Content-Type: application/json`. The body is a JWT in JWS JSON Serialization format (RFC 7515):

```json
{
  "payload": "<base64url(JSON(claims))>",
  "signatures": [
    {
      "protected": "<base64url({\"alg\":\"RS256\"})>",
      "header": { "kid": "<key-id>" },
      "signature": "<base64url(RS256 signature of 'protected.payload')>"
    }
  ]
}
```

The JWT claims include standard fields (`iat`, `exp` 30 s, `jti` random) plus a `data` object that holds the API-specific parameters.

`pkg/management.Client.buildJWT` assembles and signs the JWT using `crypto/rsa` from the standard library — no external JWT library is needed.

`cmd/management/auth.go` resolves the key ID and private key file path (from flags or env vars `BARONG_MANAGEMENT_KEY_ID` / `BARONG_MANAGEMENT_PRIVATE_KEY_FILE`), loads the PEM file (PKCS1 or PKCS8), and returns a ready-to-use `pkg/management.Client`.

## API Response Quirks

- **`OTPQRCode.Barcode`** is a base64-encoded PNG image, not an OTP URI. The `otpauth://` URI is in `OTPQRCode.URL`. Do not pass `Barcode` to a QR encoder — decode it with `base64.StdEncoding.DecodeString` and write the raw bytes to a file.
- **Management API** returns HTTP 201 for most operations (including list and get endpoints), not 200.

## Conventions

### Output

- **User-facing messages** (status, warnings, hints) go to `stderr` via `fmt.Fprintf(os.Stderr, ...)`.
- **Data/results** go to `stdout` so they can be piped or redirected.

### Sensitive files

When a command needs to write sensitive data to disk (secrets, tokens, QR codes):

- Use a fixed filename in `os.TempDir()` rather than a random temp file, so repeated runs overwrite instead of accumulate.
- Create the file with `os.OpenFile(..., 0600)` (owner read/write only).
- Print the file path to `stderr` after writing.
- Warn the user to delete the file once they are done with it.
