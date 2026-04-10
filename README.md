# barong-cli

A command line tool to interact with [Barong](https://github.com/walletera/barong) APIs. Barong is an open-source authentication and identity platform.

## Requirements

- Go 1.21+
- A running Barong instance

## Installation

```bash
go install barong-cli@latest
```

Or build from source:

```bash
git clone https://github.com/walletera/barong-cli
cd barong-cli
go build -o barong-cli .
```

## Configuration

The Barong server URL can be set in two ways (in order of precedence):

1. `--url` flag on any command
2. `BARONG_URL` environment variable

If neither is set, the tool defaults to `http://localhost:9090`.

```bash
export BARONG_URL=https://barong.example.com
```

## Usage

```
barong-cli [--url <server-url>] <command> [subcommand] [flags]
```

---

## User API (`barong-cli user`)

The User API uses session-cookie authentication. Log in first with `user login` and the session is stored automatically.

### Create a user

```bash
barong-cli user create --email user@example.com --password secret
```

| Flag | Required | Description |
|------|----------|-------------|
| `--email` | yes | User email |
| `--password` | yes | User password |
| `--username` | no | Username |
| `--refid` | no | Referral UID |

### Log in

```bash
barong-cli user login --email user@example.com --password secret
```

Authenticates against Barong and saves the session to `~/.barong-cli/session.json`. The session is reused automatically by commands that require authentication.

| Flag | Required | Description |
|------|----------|-------------|
| `--email` | yes | Account email |
| `--password` | yes | Account password |
| `--otp-code` | no | Code from your authenticator app (if 2FA is enabled) |

### Log out

```bash
barong-cli user logout
```

Destroys the session on the server and deletes the local `~/.barong-cli/session.json` file.

### Show current user

```bash
barong-cli user me
```

Prints the authenticated user's profile as JSON.

### Two-factor authentication (2FA)

**Generate a QR code**

```bash
barong-cli user otp generate-qrcode
```

Writes the QR code image to a temporary file. Scan it with your authenticator app (e.g. Google Authenticator, Authy).

To get the secret key as text instead:

```bash
barong-cli user otp generate-qrcode --show-secret
```

> **Security note:** the QR code image and secret file contain your 2FA secret. Delete them once set up.

| Flag | Required | Description |
|------|----------|-------------|
| `--show-secret` | no | Write the secret key to a file instead of the QR code image |

**Enable 2FA**

```bash
barong-cli user otp enable --code 123456
```

| Flag | Required | Description |
|------|----------|-------------|
| `--code` | yes | Code from your authenticator app |

### Service accounts

```bash
# List all service accounts for the current user
barong-cli user service-account list
```

### API keys

```bash
# List API keys for the current account
barong-cli user api-key list

# List API keys for a service account
barong-cli user api-key list --service-account-uid SA456

# Create an API key
barong-cli user api-key create --algorithm RS256 --totp-code 123456

# Create an API key for a service account
barong-cli user api-key create --algorithm RS256 --totp-code 123456 --service-account-uid SA456

# Update an API key
barong-cli user api-key update <kid> --totp-code 123456 --state active

# Update an API key for a service account
barong-cli user api-key update <kid> --totp-code 123456 --service-account-uid SA456 --scope read,write

# Delete an API key
barong-cli user api-key delete <kid> --totp-code 123456

# Delete an API key for a service account
barong-cli user api-key delete <kid> --totp-code 123456 --service-account-uid SA456
```

**`list` flags**

| Flag | Required | Description |
|------|----------|-------------|
| `--service-account-uid` | no | List keys for this service account instead of the current user |
| `--page` | no | Page number |
| `--limit` | no | Results per page (max 100) |
| `--order-by` | no | Field to sort by |
| `--ordering` | no | Sort order (`asc` or `desc`) |

**`create` flags**

| Flag | Required | Description |
|------|----------|-------------|
| `--algorithm` | yes | Key algorithm (e.g. `RS256`) |
| `--totp-code` | yes | Code from your authenticator app |
| `--scope` | no | Comma-separated scopes |
| `--service-account-uid` | no | Create key for this service account instead of the current user |

**`update <kid>` flags**

| Flag | Required | Description |
|------|----------|-------------|
| `--totp-code` | yes | Code from your authenticator app |
| `--scope` | no | Comma-separated scopes |
| `--state` | no | Key state (`active` or `disabled`) |
| `--service-account-uid` | no | Update key for this service account instead of the current user |

**`delete <kid>` flags**

| Flag | Required | Description |
|------|----------|-------------|
| `--totp-code` | yes | Code from your authenticator app |
| `--service-account-uid` | no | Delete key for this service account instead of the current user |

---

## Management API (`barong-cli management`)

The Management API uses JWT multisig authentication. Every request is signed with an RSA private key; the public key must be registered in Barong's keychain configuration.

### Authentication flags

These flags are persistent across all `management` subcommands and can also be set via environment variables:

| Flag | Env var | Description |
|------|---------|-------------|
| `--key-id` | `BARONG_MANAGEMENT_KEY_ID` | Key ID registered in Barong's keychain |
| `--private-key-file` | `BARONG_MANAGEMENT_PRIVATE_KEY_FILE` | Path to the RSA private key PEM file |

```bash
export BARONG_MANAGEMENT_KEY_ID=my-backend
export BARONG_MANAGEMENT_PRIVATE_KEY_FILE=~/.barong/management.pem
```

### Generating a keypair

```bash
openssl genrsa -out management.pem 2048
openssl rsa -in management.pem -pubout -out management_pub.pem
```

Register the public key in Barong's `config/management_api_v1.yml` under `keychain` with the same key ID you pass to `--key-id`.

### Users

```bash
# Create a user
barong-cli management users create --email user@example.com --password secret

# Get user information (by uid, email, or phone)
barong-cli management users get --uid ID123

# List users
barong-cli management users list [--extended] [--from <unix>] [--to <unix>] [--page 1] [--limit 100]

# Update user role or data
barong-cli management users update --uid ID123 --role admin

# Import an existing user (with hashed password)
barong-cli management users import --email user@example.com --password-digest '$2a$...'
```

### Labels

```bash
# Create a private label for a user
barong-cli management labels create --user-uid ID123 --key kyc --value verified

# Update a label
barong-cli management labels update --user-uid ID123 --key kyc --value pending [--replace]

# Delete a label
barong-cli management labels delete --user-uid ID123 --key kyc

# List labels for a user
barong-cli management labels list --user-uid ID123

# Get users filtered by label
barong-cli management labels filter-users --key kyc --value verified
```

### Profiles

```bash
# Import a profile for a user
barong-cli management profiles import --uid ID123 --first-name Alice --last-name Smith --country US
```

### Phones

```bash
# Create a phone number
barong-cli management phones create --uid ID123 --number +1234567890

# Get phone numbers for a user
barong-cli management phones get --uid ID123

# Delete a phone number
barong-cli management phones delete --uid ID123 --number +1234567890
```

### Documents

```bash
# Push a document (base64-encoded file content)
barong-cli management documents push \
  --uid ID123 \
  --doc-type passport \
  --doc-number AB123456 \
  --filename passport \
  --file-ext jpg \
  --upload "$(base64 -w0 passport.jpg)"
```

### Service accounts

```bash
# Create a service account
barong-cli management service-accounts create --owner-uid ID123 --role service

# Get a service account
barong-cli management service-accounts get --uid SA456

# List service accounts
barong-cli management service-accounts list [--owner-uid ID123]

# Delete a service account
barong-cli management service-accounts delete --uid SA456
```

### OTP

```bash
# Sign a request with Barong OTP signature
barong-cli management otp sign --user-uid ID123 --otp-code 123456
```

### Timestamp

```bash
# Get server time (Unix epoch seconds)
barong-cli management timestamp
```

---

## Auth Debug (`barong-cli auth-debug`)

Sends a GET request to Barong's `/api/v2/auth/{path}` endpoint and prints the response headers. Useful for manually testing how Barong resolves authentication and what JWT claims it would inject into a downstream request.

```bash
barong-cli auth-debug <test-path> [flags]
```

`<test-path>` is the downstream path being tested, e.g. `api/v1/payments`.

### Authentication

Two auth methods are supported. If neither API key flag is set, the command falls back to the saved session cookie (`~/.barong-cli/session.json`).

| Flag | Description |
|------|-------------|
| `--api-key-kid` | API key ID — switches to API key auth |
| `--api-key-secret` | API key HMAC secret (required when `--api-key-kid` is set) |

### Examples

```bash
# Using the active session cookie
barong-cli auth-debug api/v1/payments

# Using an API key
barong-cli auth-debug api/v1/payments \
  --api-key-kid 71052995e443c247 \
  --api-key-secret 7ce5bf89cd53f350c797f6a31a80c435
```

### Output

Response headers are printed to stdout. The status line goes to stderr. If the response includes an `Authorization: Bearer <JWT>` header, its payload is base64-decoded and pretty-printed as JSON below the headers:

```
Authorization: Bearer eyJ...
...

Authorization JWT payload:
{
  "uid": "SI6940BA5C62",
  "email": "johndoe@example.com",
  "role": "customer",
  "level": 3,
  "state": "active",
  ...
}
```

---

## Session storage

The User API session cookie is stored at `~/.barong-cli/session.json` with `0600` permissions. It is read automatically by any command that requires authentication and deleted on logout.

## Building and testing

```bash
# Build
go build ./...

# Run tests
go test ./...

# Lint (requires golangci-lint)
golangci-lint run
```
