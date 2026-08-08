![preview](./assets/preview.png)

___

## Index

1. [Chirpy](#chirpy)
2. [Migrations](#migrations)
3. [API Docs](#api-docs)
   - [Auth conventions](#auth-conventions)
   - [Error shape](#error-shape)
   - [Admin](#admin)
   - [Health](#health)
   - [Users](#users)
   - [Auth tokens](#auth-tokens)
   - [Chirps](#chirps)
   - [Polka](#polka)
   - [Static files](#static-files)

## Chirpy

Chirpy is a social network similar to Twitter.

A "chirp" is just a short message that a user can post to the API, like a "tweet".

## Migrations

Migrations are being handled by goose, and sqlc. See:

- [../docs/goose.md](../docs/goose.md)
- [../docs/sqlc.md](../docs/sqlc.md)

___

## API Docs

Base URL: `http://localhost:8080`

### Auth conventions

| Scheme | Header format | Used by |
| --- | --- | --- |
| Access JWT | `Authorization: Bearer <jwt>` | Update user, create/delete chirps |
| Refresh token | `Authorization: Bearer <refresh_token>` | Refresh access token, revoke refresh token |
| Polka API key | `Authorization: ApiKey <polka_key>` | Polka webhooks |

- Access JWTs expire in **1 hour** (`iss: chirpy-access`).
- Refresh tokens expire in **60 days** and can be revoked.
- `POLKA_KEY` is loaded from `.env`.

### Error shape

JSON error responses look like:

```json
{
  "error": "message here"
}
```

---

### Admin

#### `GET /admin/metrics/`

Returns an HTML page with the number of hits to `/app/`.

| | |
| --- | --- |
| Headers | none |
| Query | none |
| Body | none |
| Success | `200` `text/plain` HTML metrics page |

#### `POST /admin/reset/`

Resets metrics to `0` and deletes all users (cascades chirps / refresh tokens). Only allowed when `PLATFORM=dev`.

| | |
| --- | --- |
| Headers | none |
| Query | none |
| Body | none |
| Success | `200` `text/plain` e.g. `Hits: 0` |
| Errors | `403` if not in `dev`; `500` on DB failure |

---

### Health

#### `GET /api/healthz/`

| | |
| --- | --- |
| Headers | none |
| Query | none |
| Body | none |
| Success | `200` `text/plain` body: `OK` |

---

### Users

#### `POST /api/users`

Creates a user.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "plain-text-password"
}
```

| | |
| --- | --- |
| Headers | `Content-Type: application/json` |
| Query | none |
| Success | `201` |

**Response body**

```json
{
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2026-08-08T12:00:00Z",
  "updated_at": "2026-08-08T12:00:00Z",
  "is_chirpy_red": false
}
```

| Errors | |
| --- | --- |
| `400` | Invalid JSON |
| `500` | Hashing or DB failure |

Password is hashed with Argon2id before storage. The hashed password is never returned.

#### `PUT /api/users`

Updates the authenticated user's email and password.

**Headers**

```http
Authorization: Bearer <access_jwt>
Content-Type: application/json
```

**Request body**

```json
{
  "email": "new@example.com",
  "password": "new-password"
}
```

| | |
| --- | --- |
| Query | none |
| Success | `200` — same `User` JSON shape as create |

| Errors | |
| --- | --- |
| `400` | Invalid JSON |
| `401` | Missing/invalid access JWT |
| `500` | Hashing or DB failure |

#### `POST /api/login`

Logs in with email/password and returns tokens.

**Request body**

```json
{
  "email": "user@example.com",
  "password": "plain-text-password"
}
```

| | |
| --- | --- |
| Headers | `Content-Type: application/json` |
| Query | none |
| Success | `200` |

**Response body**

```json
{
  "token": "<access_jwt>",
  "refresh_token": "<refresh_token>",
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2026-08-08T12:00:00Z",
  "updated_at": "2026-08-08T12:00:00Z",
  "is_chirpy_red": false
}
```

| Errors | |
| --- | --- |
| `400` | Invalid JSON |
| `401` | Wrong email/password |
| `500` | Failed to create refresh token |

---

### Auth tokens

#### `POST /api/refresh`

Exchanges a valid (non-revoked) refresh token for a new access JWT.

**Headers**

```http
Authorization: Bearer <refresh_token>
```

| | |
| --- | --- |
| Query | none |
| Body | none |
| Success | `200` |

**Response body**

```json
{
  "token": "<new_access_jwt>"
}
```

| Errors | |
| --- | --- |
| `401` | Missing/invalid/revoked refresh token, or user missing |
| `500` | Failed to create access token |

#### `POST /api/revoke`

Revokes a refresh token (`revoked_at` + `updated_at` set).

**Headers**

```http
Authorization: Bearer <refresh_token>
```

| | |
| --- | --- |
| Query | none |
| Body | none |
| Success | `204` |

| Errors | |
| --- | --- |
| `401` | Missing/invalid refresh token |
| `500` | DB failure |

---

### Chirps

#### `POST /api/chirps`

Creates a chirp for the authenticated user. Body max length **140**. Words `kerfuffle`, `sharbert`, and `fornax` are replaced with `****`.

**Headers**

```http
Authorization: Bearer <access_jwt>
Content-Type: application/json
```

**Request body**

```json
{
  "body": "Hello world"
}
```

> `user_id` / `token` fields may appear on the request type, but authorship comes from the JWT, not the body.

| | |
| --- | --- |
| Query | none |
| Success | `201` |

**Response body**

```json
{
  "id": "uuid",
  "body": "Hello world",
  "user_id": "uuid",
  "created_at": "2026-08-08T12:00:00Z",
  "updated_at": "2026-08-08T12:00:00Z"
}
```

| Errors | |
| --- | --- |
| `400` | Invalid JSON or body too long |
| `401` | Missing/invalid access JWT |
| `500` | DB failure |

#### `GET /api/chirps`

Lists chirps.

**Query parameters**

| Param | Required | Values | Default |
| --- | --- | --- | --- |
| `author_id` | no | user UUID | all authors |
| `sort` | no | `asc`, `desc` | `asc` (by `created_at`) |

Examples:

```http
GET /api/chirps
GET /api/chirps?sort=desc
GET /api/chirps?author_id=<uuid>
GET /api/chirps?author_id=<uuid>&sort=desc
```

| | |
| --- | --- |
| Headers | none |
| Body | none |
| Success | `200` — JSON array of chirp objects (same shape as create) |

| Errors | |
| --- | --- |
| `400` | Invalid `author_id` |
| `500` | DB failure |

#### `GET /api/chirps/{chirpId}`

Fetches one chirp by UUID path param.

| | |
| --- | --- |
| Headers | none |
| Query | none |
| Body | none |
| Path | `chirpId` — chirp UUID |
| Success | `200` — single chirp object |

| Errors | |
| --- | --- |
| `400` | Invalid `chirpId` |
| `404` | Chirp not found |

#### `DELETE /api/chirps/{chirpId}`

Deletes a chirp. Only the author may delete it.

**Headers**

```http
Authorization: Bearer <access_jwt>
```

| | |
| --- | --- |
| Query | none |
| Body | none |
| Path | `chirpId` — chirp UUID |
| Success | `204` |

| Errors | |
| --- | --- |
| `400` | Invalid `chirpId` |
| `401` | Missing/invalid access JWT |
| `403` | Authenticated user is not the author |
| `404` | Chirp not found |
| `500` | DB failure |

---

### Polka

#### `POST /api/polka/webhooks`

Webhook used by Polka to upgrade a user to Chirpy Red.

**Headers**

```http
Authorization: ApiKey <POLKA_KEY>
Content-Type: application/json
```

**Request body**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
  }
}
```

| | |
| --- | --- |
| Query | none |
| Success | `204` empty body |

Behavior:

- If `event` is anything other than `user.upgraded` → immediate `204` (ignored).
- If `event` is `user.upgraded` → sets `is_chirpy_red = true` for that user, then `204`.

| Errors | |
| --- | --- |
| `400` | Invalid JSON or invalid `user_id` |
| `401` | Missing/invalid `ApiKey` |
| `404` | User not found |
| `500` | DB failure |

---

### Static files

#### `GET /app/...`

Serves files from the project root (file server), with hit metrics middleware.

| | |
| --- | --- |
| Headers | none |
| Query | none |
| Body | none |
| Success | static file response (`index.html`, assets, etc.) |
