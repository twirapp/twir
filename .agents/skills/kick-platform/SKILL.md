---
name: kick-platform
description: Complete reference for integrating with the Kick streaming platform API. Covers authentication, webhooks/events, REST API endpoints, and security best practices.
metadata:
  author: satont
  version: "1.0.0"
---

# Kick Platform Integration Skill

Complete reference for integrating with the Kick streaming platform API. Covers authentication, webhooks/events, REST API endpoints, and security best practices.

## Table of Contents

- [Overview](#overview)
- [Authentication](#authentication)
  - [OAuth 2.1 Flow](#oauth-21-flow)
  - [Token Types](#token-types)
  - [Scopes](#scopes)
  - [Token Management](#token-management)
- [Webhook Events](#webhook-events)
  - [Event Types](#event-types)
  - [Webhook Security](#webhook-security)
  - [Event Payloads](#event-payloads)
- [REST API](#rest-api)
  - [Users](#users)
  - [Channels](#channels)
  - [Chat](#chat)
  - [Moderation](#moderation)
  - [Livestreams](#livestreams)
  - [Categories](#categories)
  - [Channel Rewards](#channel-rewards)
  - [KICKs](#kicks)
- [Event Subscriptions](#event-subscriptions)
- [Rate Limits & Best Practices](#rate-limits--best-practices)
- [Common Patterns](#common-patterns)

---

## Overview

Kick's Public API provides programmatic access to the Kick streaming platform. The API consists of:

- **OAuth Server**: `https://id.kick.com` — authentication and token management
- **API Server**: `https://api.kick.com` — REST API endpoints
- **Webhooks**: Event-driven notifications for real-time updates

### Base URLs

```
OAuth:  https://id.kick.com
API:    https://api.kick.com
```

### Response Format

All API responses follow a consistent envelope:

```json
{
  "data": { ... },
  "message": "OK"
}
```

---

## Authentication

### OAuth 2.1 Flow

Kick uses OAuth 2.1 with PKCE (Proof Key for Code Exchange) for user authentication.

#### Step 1: Generate PKCE Code Verifier & Challenge

```typescript
import crypto from "crypto";

// Generate code verifier (43-128 chars)
const codeVerifier = crypto.randomBytes(32).toString("base64url");

// Generate code challenge (S256)
const codeChallenge = crypto
  .createHash("sha256")
  .update(codeVerifier)
  .digest("base64url");
```

#### Step 2: Authorization Request

```
GET https://id.kick.com/oauth/authorize
```

**Query Parameters:**

| Parameter               | Required | Description                           |
| ----------------------- | -------- | ------------------------------------- |
| `client_id`             | Yes      | Your application's client ID          |
| `response_type`         | Yes      | Must be `code`                        |
| `redirect_uri`          | Yes      | Registered redirect URI               |
| `state`                 | Yes      | Random state string (CSRF protection) |
| `scope`                 | Yes      | Space-separated scopes                |
| `code_challenge`        | Yes      | PKCE code challenge                   |
| `code_challenge_method` | Yes      | Must be `S256`                        |

**Example URL:**

```
https://id.kick.com/oauth/authorize?
  response_type=code&
  client_id=YOUR_CLIENT_ID&
  redirect_uri=https://yourapp.com/callback&
  scope=user:read%20channel:read%20chat:write&
  code_challenge=CHALLENGE&
  code_challenge_method=S256&
  state=RANDOM_STATE
```

**Note on 127.0.0.1 Redirect URIs:**

If using `127.0.0.1` as redirect host, add a sacrificial query parameter before `redirect_uri` due to a Next.js bug:

```
https://id.kick.com/oauth/authorize?
  response_type=code&
  client_id=CLIENT_ID&
  redirect=127.0.0.1&          // sacrificial parameter
  redirect_uri=http://127.0.0.1/callback&
  scope=SCOPES&
  code_challenge=CHALLENGE&
  code_challenge_method=S256&
  state=STATE
```

#### Step 3: Exchange Code for Token

```
POST https://id.kick.com/oauth/token
Content-Type: application/x-www-form-urlencoded
```

**Body:**

```
grant_type=authorization_code
client_id=YOUR_CLIENT_ID
client_secret=YOUR_CLIENT_SECRET
redirect_uri=https://yourapp.com/callback
code_verifier=CODE_VERIFIER
code=CODE_FROM_CALLBACK
```

**Response:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "refresh_token": "eyJ...",
  "expires_in": 3600,
  "scope": "user:read channel:read"
}
```

### Token Types

#### App Access Token (Client Credentials)

For server-to-server API calls without user context.

```
POST https://id.kick.com/oauth/token
Content-Type: application/x-www-form-urlencoded
```

**Body:**

```
grant_type=client_credentials
client_id=YOUR_CLIENT_ID
client_secret=YOUR_CLIENT_SECRET
```

**Response:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

**Use cases:**

- Public data access (categories, livestreams)
- Webhook event subscriptions

#### User Access Token (Authorization Code)

For acting on behalf of a user.

**Use cases:**

- Sending chat messages
- Updating channel metadata
- Moderation actions
- Subscribing to channel events

### Scopes

| Scope                            | Description                                            |
| -------------------------------- | ------------------------------------------------------ |
| `user:read`                      | View user information (username, streamer ID, etc.)    |
| `channel:read`                   | View channel information (description, category, etc.) |
| `channel:write`                  | Update livestream metadata                             |
| `channel:rewards:read`           | Read channel point rewards                             |
| `channel:rewards:write`          | Create, edit, delete channel rewards                   |
| `chat:write`                     | Send chat messages                                     |
| `streamkey:read`                 | Read stream URL and stream key                         |
| `events:subscribe`               | Subscribe to channel events (webhooks)                 |
| `moderation:ban`                 | Execute ban/unban actions                              |
| `moderation:chat_message:manage` | Delete chat messages                                   |
| `kicks:read`                     | View KICKs leaderboard information                     |

### Token Management

#### Refresh Token

```
POST https://id.kick.com/oauth/token
Content-Type: application/x-www-form-urlencoded
```

**Body:**

```
grant_type=refresh_token
client_id=YOUR_CLIENT_ID
client_secret=YOUR_CLIENT_SECRET
refresh_token=REFRESH_TOKEN
```

**Response:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "refresh_token": "eyJ...",
  "expires_in": 3600,
  "scope": "user:read channel:read"
}
```

**Note:** Refresh tokens are reusable/flexible (updated Nov 2025).

#### Revoke Token

```
POST https://id.kick.com/oauth/revoke?token=TOKEN&token_hint_type=access_token
Content-Type: application/x-www-form-urlencoded
```

#### Token Introspection

```
POST https://id.kick.com/oauth/token/introspect
Authorization: Bearer <access_token>
```

**Response:**

```json
{
  "data": {
    "active": true,
    "client_id": "your_client_id",
    "token_type": "user",
    "scope": "user:read channel:read",
    "exp": 1771046347
  },
  "message": "OK"
}
```

**Note:** The old `/public/v1/token/introspect` endpoint is deprecated. Use `/oauth/token/introspect` instead.

---

## Webhook Events

### Event Types

| Event Name           | Event ID                            | Version | Description                   |
| -------------------- | ----------------------------------- | ------- | ----------------------------- |
| Chat Message         | `chat.message.sent`                 | 1       | Message sent in chat          |
| Channel Follow       | `channel.followed`                  | 1       | User follows a channel        |
| Subscription Renewal | `channel.subscription.renewal`      | 1       | Subscription renewed          |
| Subscription Gifts   | `channel.subscription.gifts`        | 1       | Subscriptions gifted          |
| Subscription Created | `channel.subscription.new`          | 1       | New subscription              |
| Reward Redemption    | `channel.reward.redemption.updated` | 1       | Reward redeemed               |
| Livestream Status    | `livestream.status.updated`         | 1       | Stream started/ended          |
| Livestream Metadata  | `livestream.metadata.updated`       | 1       | Stream title/category changed |
| Moderation Banned    | `moderation.banned`                 | 1       | User banned/timed out         |
| Kicks Gifted         | `kicks.gifted`                      | 1       | Kicks gifted to channel       |

### Webhook Security

#### Headers

| Header                         | Type    | Description                            |
| ------------------------------ | ------- | -------------------------------------- |
| `Kick-Event-Message-Id`        | ULID    | Unique message ID (idempotent key)     |
| `Kick-Event-Subscription-Id`   | ULID    | Subscription ID                        |
| `Kick-Event-Signature`         | Base64  | RSA signature for verification         |
| `Kick-Event-Message-Timestamp` | RFC3339 | Message timestamp                      |
| `Kick-Event-Type`              | string  | Event type (e.g., `chat.message.sent`) |
| `Kick-Event-Version`           | string  | Event version                          |

#### Signature Verification

The signature is created by concatenating three values separated by `.`:

```
signature_payload = message_id + "." + timestamp + "." + raw_body
```

Then signed with Kick's private key using RSA-SHA256 (PKCS1v15).

#### Public Key

```
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq/+l1WnlRrGSolDMA+A8
6rAhMbQGmQ2SapVcGM3zq8ANXjnhDWocMqfWcTd95btDydITa10kDvHzw9WQOqp2
MZI7ZyrfzJuz5nhTPCiJwTwnEtWft7nV14BYRDHvlfqPUaZ+1KR4OCaO/wWIk/rQ
L/TjY0M70gse8rlBkbo2a8rKhu69RQTRsoaf4DVhDPEeSeI5jVrRDGAMGL3cGuyY
6CLKGdjVEM78g3JfYOvDU/RvfqD7L89TZ3iN94jrmWdGz34JNlEI5hqK8dd7C5EF
BEbZ5jgB8s8ReQV8H+MkuffjdAj3ajDDX3DOJMIut1lBrUVD1AaSrGCKHooWoL2e
twIDAQAB
-----END PUBLIC KEY-----
```

Also available at: `GET https://api.kick.com/public/v1/public-key`

#### Verification Example (Go)

```go
package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func ParsePublicKey(bs []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bs)
	if block == nil {
		return nil, fmt.Errorf("not decodable key")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("not public key")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not expected public key interface")
	}
	return publicKey, nil
}

func VerifyWebhookSignature(
	publicKey *rsa.PublicKey,
	messageID string,
	timestamp string,
	body []byte,
	signatureB64 string,
) error {
	// Decode base64 signature
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(signatureB64)))
	n, err := base64.StdEncoding.Decode(decoded, []byte(signatureB64))
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	signature := decoded[:n]

	// Build payload: message_id.timestamp.body
	payload := []byte(fmt.Sprintf("%s.%s.%s", messageID, timestamp, body))

	// Hash and verify
	hashed := sha256.Sum256(payload)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
}
```

#### Verification Example (TypeScript)

```typescript
import crypto from "crypto";

const KICK_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq/+l1WnlRrGSolDMA+A8
6rAhMbQGmQ2SapVcGM3zq8ANXjnhDWocMqfWcTd95btDydITa10kDvHzw9WQOqp2
MZI7ZyrfzJuz5nhTPCiJwTwnEtWft7nV14BYRDHvlfqPUaZ+1KR4OCaO/wWIk/rQ
L/TjY0M70gse8rlBkbo2a8rKhu69RQTRsoaf4DVhDPEeSeI5jVrRDGAMGL3cGuyY
6CLKGdjVEM78g3JfYOvDU/RvfqD7L89TZ3iN94jrmWdGz34JNlEI5hqK8dd7C5EF
BEbZ5jgB8s8ReQV8H+MkuffjdAj3ajDDX3DOJMIut1lBrUVD1AaSrGCKHooWoL2e
twIDAQAB
-----END PUBLIC KEY-----`;

function verifyWebhookSignature(
  messageId: string,
  timestamp: string,
  body: string,
  signatureB64: string,
): boolean {
  const payload = `${messageId}.${timestamp}.${body}`;

  const verifier = crypto.createVerify("RSA-SHA256");
  verifier.update(payload);

  return verifier.verify(KICK_PUBLIC_KEY, signatureB64, "base64");
}
```

### Event Payloads

#### Chat Message (`chat.message.sent`)

```json
{
  "message_id": "unique_message_id_123",
  "replies_to": {
    "message_id": "unique_message_id_456",
    "content": "This is the parent message!",
    "sender": {
      "is_anonymous": false,
      "user_id": 12345,
      "username": "parent_sender_name",
      "is_verified": false,
      "profile_picture": "https://example.com/parent_sender_avatar.jpg",
      "channel_slug": "parent_sender_channel",
      "identity": null
    }
  },
  "broadcaster": {
    "is_anonymous": false,
    "user_id": 123456789,
    "username": "broadcaster_name",
    "is_verified": true,
    "profile_picture": "https://example.com/broadcaster_avatar.jpg",
    "channel_slug": "broadcaster_channel",
    "identity": null
  },
  "sender": {
    "is_anonymous": false,
    "user_id": 987654321,
    "username": "sender_name",
    "is_verified": false,
    "profile_picture": "https://example.com/sender_avatar.jpg",
    "channel_slug": "sender_channel",
    "identity": {
      "username_color": "#FF5733",
      "badges": [
        {
          "text": "Moderator",
          "type": "moderator"
        },
        {
          "text": "Sub Gifter",
          "type": "sub_gifter",
          "count": 5
        },
        {
          "text": "Subscriber",
          "type": "subscriber",
          "count": 3
        }
      ]
    }
  },
  "content": "Hello [emote:4148074:HYPERCLAP] [emote:4148074:HYPERCLAP] [emote:37226:KEKW]",
  "emotes": [
    {
      "emote_id": "4148074",
      "positions": [
        { "s": 6, "e": 30 },
        { "s": 32, "e": 56 }
      ]
    },
    {
      "emote_id": "37226",
      "positions": [{ "s": 58, "e": 75 }]
    }
  ],
  "created_at": "2025-01-14T16:08:06Z"
}
```

**Chat Message Notes:**

- `replies_to` is present only for reply messages
- Emotes are referenced as `[emote:{emote_id}:{emote_name}]` in content
- `identity` contains user-specific chat styling (badges, color)
- `created_at` added July 2025

#### Channel Follow (`channel.followed`)

```json
{
  "broadcaster": {
    "is_anonymous": false,
    "user_id": 123456789,
    "username": "broadcaster_name",
    "is_verified": true,
    "profile_picture": "https://example.com/broadcaster_avatar.jpg",
    "channel_slug": "broadcaster_channel",
    "identity": null
  },
  "follower": {
    "is_anonymous": false,
    "user_id": 987654321,
    "username": "follower_name",
    "is_verified": false,
    "profile_picture": "https://example.com/sender_avatar.jpg",
    "channel_slug": "follower_channel",
    "identity": null
  }
}
```

#### Subscription Created (`channel.subscription.new`)

```json
{
  "broadcaster": { ... },
  "subscriber": { ... },
  "duration": 1,
  "created_at": "2025-01-14T16:08:06Z",
  "expires_at": "2025-02-14T16:08:06Z"
}
```

#### Subscription Renewal (`channel.subscription.renewal`)

```json
{
  "broadcaster": { ... },
  "subscriber": { ... },
  "duration": 3,
  "created_at": "2025-01-14T16:08:06Z",
  "expires_at": "2025-02-14T16:08:06Z"
}
```

#### Subscription Gifts (`channel.subscription.gifts`)

```json
{
  "broadcaster": { ... },
  "gifter": {
    "is_anonymous": false,
    "user_id": 987654321,
    "username": "gifter_name",
    "is_verified": false,
    "profile_picture": "https://example.com/sender_avatar.jpg",
    "channel_slug": "gifter_channel",
    "identity": null
  },
  "giftees": [
    {
      "is_anonymous": false,
      "user_id": 561654654,
      "username": "giftee_name",
      "is_verified": true,
      "profile_picture": "https://example.com/broadcaster_avatar.jpg",
      "channel_slug": "giftee_channel",
      "identity": null
    }
  ],
  "created_at": "2025-01-14T16:08:06Z",
  "expires_at": "2025-02-14T16:08:06Z"
}
```

**Note:** If `gifter.is_anonymous` is true, all gifter fields except `is_anonymous` are `null`.

#### Reward Redemption (`channel.reward.redemption.updated`)

```json
{
  "id": "01KBHE78QE4HZY1617DK5FC7YD",
  "user_input": "unban me",
  "status": "rejected",
  "redeemed_at": "2025-12-02T22:54:19.323Z",
  "reward": {
    "id": "01KBHE7RZNHB0SKDV1H86CD4F3",
    "title": "Uban Request",
    "cost": 1000,
    "description": "Only good reasons pls"
  },
  "redeemer": {
    "user_id": 123,
    "username": "naughty-user",
    "is_verified": false,
    "profile_picture": "",
    "channel_slug": "naughty_user"
  },
  "broadcaster": {
    "user_id": 333,
    "username": "gigachad",
    "is_verified": true,
    "profile_picture": "",
    "channel_slug": "gigachad"
  }
}
```

**Status values:** `pending`, `accepted`, `rejected`

#### Livestream Status (`livestream.status.updated`)

**Stream Started:**

```json
{
  "broadcaster": { ... },
  "is_live": true,
  "title": "Stream Title",
  "started_at": "2025-01-01T11:00:00+11:00",
  "ended_at": null
}
```

**Stream Ended:**

```json
{
  "broadcaster": { ... },
  "is_live": false,
  "title": "Stream Title",
  "started_at": "2025-01-01T11:00:00+11:00",
  "ended_at": "2025-01-01T15:00:00+11:00"
}
```

#### Livestream Metadata (`livestream.metadata.updated`)

```json
{
  "broadcaster": { ... },
  "metadata": {
    "title": "Stream Title",
    "language": "en",
    "has_mature_content": true,
    "category": {
      "id": 123,
      "name": "Category name",
      "thumbnail": "http://example.com/image123"
    }
  }
}
```

#### Moderation Banned (`moderation.banned`)

```json
{
  "broadcaster": { ... },
  "moderator": { ... },
  "banned_user": { ... },
  "metadata": {
    "reason": "banned reason",
    "created_at": "2025-01-14T16:08:05Z",
    "expires_at": "2025-01-14T16:10:06Z"
  }
}
```

**Note:** `expires_at` is `null` for permanent bans.

#### Kicks Gifted (`kicks.gifted`)

```json
{
  "broadcaster": {
    "user_id": 123456789,
    "username": "broadcaster_name",
    "is_verified": true,
    "profile_picture": "https://example.com/broadcaster_avatar.jpg",
    "channel_slug": "broadcaster_channel"
  },
  "sender": {
    "user_id": 987654321,
    "username": "gift_sender",
    "is_verified": false,
    "profile_picture": "https://example.com/sender_avatar.jpg",
    "channel_slug": "gift_sender_channel"
  },
  "gift": {
    "amount": 500,
    "name": "Rage Quit",
    "type": "LEVEL_UP",
    "tier": "MID",
    "message": "w",
    "pinned_time_seconds": 600
  },
  "created_at": "2025-10-20T04:00:08.634Z"
}
```

---

## REST API

### Authentication

All API requests require an `Authorization: Bearer <token>` header.

### Users

#### Get Users

```
GET /public/v1/users?id[]=123&id[]=456
```

**Parameters:**

- `id` (array, optional): User IDs. If omitted, returns current authenticated user.

**Response:**

```json
{
  "data": [
    {
      "email": "user@example.com",
      "name": "username",
      "profile_picture": "https://...",
      "user_id": 123
    }
  ],
  "message": "OK"
}
```

**Scopes:** `user:read` (for sensitive data), or App Access Token

### Channels

#### Get Channels

```
GET /public/v1/channels?broadcaster_user_id[]=123&broadcaster_user_id[]=456
GET /public/v1/channels?slug=channel1&slug=channel2
```

**Parameters (mutually exclusive):**

- `broadcaster_user_id` (array): Up to 50 IDs
- `slug` (array): Up to 50 slugs, max 25 chars each

**Response:**

```json
{
  "data": [
    {
      "broadcaster_user_id": 123,
      "slug": "channel_slug",
      "channel_description": "Description",
      "stream_title": "Stream Title",
      "banner_picture": "https://...",
      "category": {
        "id": 1,
        "name": "Category",
        "thumbnail": "https://..."
      },
      "stream": {
        "is_live": true,
        "is_mature": false,
        "language": "en",
        "viewer_count": 100,
        "thumbnail": "https://...",
        "url": "https://...",
        "key": "stream_key",
        "start_time": "2025-01-01T11:00:00+11:00",
        "custom_tags": ["tag1", "tag2"]
      },
      "active_subscribers_count": 1000,
      "canceled_subscribers_count": 50
    }
  ],
  "message": "OK"
}
```

**Scopes:** `channel:read` or App Access Token

#### Update Channel

```
PATCH /public/v1/channels
```

**Body:**

```json
{
  "stream_title": "New Title",
  "category_id": 123,
  "custom_tags": ["tag1", "tag2"]
}
```

**Scopes:** `channel:write`

### Chat

#### Send Chat Message

```
POST /public/v1/chat
```

**Body:**

```json
{
  "content": "Hello chat!",
  "type": "user",
  "broadcaster_user_id": 123,
  "reply_to_message_id": "uuid-of-message"
}
```

**Parameters:**

- `content` (string, required): Message content, max 500 chars
- `type` (enum, required): `"user"` or `"bot"`
- `broadcaster_user_id` (integer): Required for `type: "user"`, ignored for `type: "bot"`
- `reply_to_message_id` (string, optional): UUID of message to reply to

**Response:**

```json
{
  "data": {
    "is_sent": true,
    "message_id": "uuid"
  },
  "message": "OK"
}
```

**Scopes:** `chat:write`

#### Delete Chat Message

```
DELETE /public/v1/chat/{message_id}
```

**Scopes:** `moderation:chat_message:manage`

### Moderation

#### Ban/Timeout User

```
POST /public/v1/moderation/bans
```

**Body:**

```json
{
  "broadcaster_user_id": 123,
  "user_id": 456,
  "duration": 60,
  "reason": "Spam"
}
```

**Parameters:**

- `broadcaster_user_id` (integer, required)
- `user_id` (integer, required)
- `duration` (integer, optional): Timeout in minutes (1-10080). Omit for permanent ban.
- `reason` (string, optional): Max 100 chars

**Scopes:** `moderation:ban`

#### Unban User

```
DELETE /public/v1/moderation/bans
```

**Body:**

```json
{
  "broadcaster_user_id": 123,
  "user_id": 456
}
```

**Scopes:** `moderation:ban`

### Livestreams

#### Get Livestreams

```
GET /public/v1/livestreams
```

**Parameters:**

- `broadcaster_user_id` (array): Up to 50 IDs
- `category_id` (integer)
- `language` (string)
- `limit` (integer): Default 25, max 100
- `sort` (enum): `viewer_count` or `started_at`

**Response:**

```json
{
  "data": [
    {
      "broadcaster_user_id": 123,
      "channel_id": 456,
      "slug": "channel_slug",
      "stream_title": "Title",
      "language": "en",
      "has_mature_content": false,
      "viewer_count": 100,
      "thumbnail": "https://...",
      "profile_picture": "https://...",
      "started_at": "2025-01-01T11:00:00+11:00",
      "category": {
        "id": 1,
        "name": "Category",
        "thumbnail": "https://..."
      },
      "custom_tags": ["tag1"]
    }
  ],
  "message": "OK"
}
```

#### Get Livestreams Stats

```
GET /public/v1/livestreams/stats
```

**Response:**

```json
{
  "data": {
    "total_count": 1500
  },
  "message": "OK"
}
```

### Categories

#### Get Categories (V2) — Recommended

```
GET /public/v2/categories?cursor=&limit=25&name[]=Gaming&tag[]=fps&id[]=1
```

**Parameters:**

- `cursor` (string): Pagination cursor (4-28 chars)
- `limit` (integer): Default 25, max 1000
- `name` (array): Category names (min 3 chars each, max 100)
- `tag` (array): Category tags (min 3 chars each, max 100)
- `id` (array): Category IDs

**Response:**

```json
{
  "data": [
    {
      "id": 1,
      "name": "Category Name",
      "tags": ["tag1", "tag2"],
      "thumbnail": "https://..."
    }
  ],
  "message": "OK",
  "pagination": {
    "next_cursor": "cursor_string"
  }
}
```

**Note:** V1 categories endpoints are deprecated.

### Channel Rewards

#### Get Channel Rewards

```
GET /public/v1/channels/rewards
```

**Response:**

```json
{
  "data": [
    {
      "id": "01KBHE7RZNHB0SKDV1H86CD4F3",
      "title": "Reward Title",
      "cost": 100,
      "description": "Description",
      "background_color": "#00e701",
      "is_enabled": true,
      "is_paused": false,
      "is_user_input_required": false,
      "should_redemptions_skip_request_queue": false
    }
  ],
  "message": "OK"
}
```

**Scopes:** `channel:rewards:read` or `channel:rewards:write`

#### Create Channel Reward

```
POST /public/v1/channels/rewards
```

**Body:**

```json
{
  "title": "New Reward",
  "cost": 100,
  "description": "Description",
  "background_color": "#00e701",
  "is_enabled": true,
  "is_user_input_required": false,
  "should_redemptions_skip_request_queue": false
}
```

**Limits:** Max 15 rewards per channel (including disabled).

**Scopes:** `channel:rewards:write`

#### Update Channel Reward

```
PATCH /public/v1/channels/rewards/{id}
```

**Body:** (at least one field required)

```json
{
  "title": "Updated Title",
  "cost": 200,
  "is_enabled": false
}
```

**Note:** Only the app that created the reward can update it.

**Scopes:** `channel:rewards:write`

#### Delete Channel Reward

```
DELETE /public/v1/channels/rewards/{id}
```

**Note:** Only the app that created the reward can delete it.

**Scopes:** `channel:rewards:write`

#### Get Channel Reward Redemptions

```
GET /public/v1/channels/rewards/redemptions?status=pending&reward_id=ID&cursor=
```

**Parameters:**

- `reward_id` (string, optional): Filter by specific reward
- `status` (enum, optional): `pending`, `accepted`, `rejected` (default: `pending`)
- `id` (array, optional): Filter by redemption IDs (mutually exclusive with other filters)
- `cursor` (string, optional): Pagination cursor

**Response:**

```json
{
  "data": [
    {
      "reward": {
        "id": "01KBHE7RZNHB0SKDV1H86CD4F3",
        "title": "Reward Title",
        "cost": 100,
        "description": "Description",
        "can_manage": true,
        "is_deleted": false
      },
      "redemptions": [
        {
          "id": "01KBHE78QE4HZY1617DK5FC7YD",
          "redeemed_at": "2025-12-02T22:54:19.323Z",
          "redeemer": {
            "user_id": 123
          },
          "status": "pending",
          "user_input": "user message"
        }
      ]
    }
  ],
  "message": "OK",
  "pagination": {
    "next_cursor": "cursor"
  }
}
```

**Scopes:** `channel:rewards:read` or `channel:rewards:write`

#### Accept Redemptions

```
POST /public/v1/channels/rewards/redemptions/accept
```

**Body:**

```json
{
  "ids": ["01KBHE78QE4HZY1617DK5FC7YD"]
}
```

**Limits:** Max 25 redemptions per request, IDs must be unique.

**Response:** Only returns data for failed redemptions.

**Scopes:** `channel:rewards:write`

#### Reject Redemptions

```
POST /public/v1/channels/rewards/redemptions/reject
```

**Body:**

```json
{
  "ids": ["01KBHE78QE4HZY1617DK5FC7YD"]
}
```

**Limits:** Max 25 redemptions per request.

**Scopes:** `channel:rewards:write`

### KICKs

#### Get KICKs Leaderboard

```
GET /public/v1/kicks/leaderboard?top=10
```

**Parameters:**

- `top` (integer): Number of entries to return (1-100, default: 10)

**Response:**

```json
{
  "data": {
    "lifetime": [
      {
        "user_id": 123,
        "username": "user1",
        "gifted_amount": 5000,
        "rank": 1
      }
    ],
    "month": [...],
    "week": [...]
  },
  "message": "OK"
}
```

**Scopes:** `kicks:read`

---

## Event Subscriptions

### Get Subscriptions

```
GET /public/v1/events/subscriptions?broadcaster_user_id=123
```

**Response:**

```json
{
  "data": [
    {
      "id": "subscription_id",
      "app_id": "app_id",
      "broadcaster_user_id": 123,
      "event": "chat.message.sent",
      "version": 1,
      "method": "webhook",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ],
  "message": "OK"
}
```

### Subscribe to Events

```
POST /public/v1/events/subscriptions
```

**Body:**

```json
{
  "broadcaster_user_id": 123,
  "events": [
    {
      "name": "chat.message.sent",
      "version": 1
    },
    {
      "name": "channel.followed",
      "version": 1
    }
  ],
  "method": "webhook"
}
```

**Parameters:**

- `broadcaster_user_id` (integer): Required for App Access Token. Ignored for User Access Token (inferred from token).
- `events` (array, required): List of events to subscribe to
- `method` (enum): Must be `"webhook"`

**Limits:**

- Max 10,000 subscriptions per event type per app
- Max 1,000 `chat.message.sent` subscriptions for unverified apps

**Response:**

```json
{
  "data": [
    {
      "name": "chat.message.sent",
      "version": 1,
      "subscription_id": "sub_id",
      "error": null
    }
  ],
  "message": "OK"
}
```

### Delete Subscriptions

```
DELETE /public/v1/events/subscriptions?id=sub1&id=sub2
```

**Parameters:**

- `id` (array, required): Subscription IDs to delete

**Note:** Disabling webhooks for an app automatically unsubscribes from all events.

---

## Rate Limits & Best Practices

### Webhook Reliability

- **Auto-unsubscribe:** Apps that fail to process webhooks for over a day are automatically unsubscribed from that event
- **Idempotency:** Use `Kick-Event-Message-Id` to deduplicate events
- **Verification:** Always verify webhook signatures before processing
- **Response:** Return 2xx status quickly. Process events asynchronously.

### App Verification

- Unverified apps: 1,000 `chat.message.sent` subscriptions
- Verified apps: 10,000 subscriptions per event type
- Verification requires emailing `developers@kick.com` with:
  - Client ID
  - App Name
  - Bot verification need
  - Reason for verification
  - Supporting evidence

### General Best Practices

1. **Store tokens securely** — never expose client secrets in client-side code
2. **Handle token expiration** — refresh before expiry
3. **Verify webhooks** — always validate signatures
4. **Use idempotency keys** — handle duplicate events gracefully
5. **Subscribe minimally** — only subscribe to events you need
6. **Handle failures** — implement retry logic for API calls
7. **Respect rate limits** — implement exponential backoff

---

## Common Patterns

### Complete OAuth Flow (TypeScript)

```typescript
import crypto from "crypto";

class KickOAuth {
  private clientId: string;
  private clientSecret: string;
  private redirectUri: string;

  constructor(clientId: string, clientSecret: string, redirectUri: string) {
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.redirectUri = redirectUri;
  }

  generatePKCE() {
    const verifier = crypto.randomBytes(32).toString("base64url");
    const challenge = crypto
      .createHash("sha256")
      .update(verifier)
      .digest("base64url");
    return { verifier, challenge };
  }

  getAuthorizationUrl(
    scopes: string[],
    state: string,
    pkce: { challenge: string },
  ) {
    const params = new URLSearchParams({
      response_type: "code",
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      scope: scopes.join(" "),
      code_challenge: pkce.challenge,
      code_challenge_method: "S256",
      state,
    });
    return `https://id.kick.com/oauth/authorize?${params.toString()}`;
  }

  async exchangeCode(code: string, codeVerifier: string) {
    const body = new URLSearchParams({
      grant_type: "authorization_code",
      client_id: this.clientId,
      client_secret: this.clientSecret,
      redirect_uri: this.redirectUri,
      code_verifier: codeVerifier,
      code,
    });

    const res = await fetch("https://id.kick.com/oauth/token", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: body.toString(),
    });

    return res.json();
  }

  async refreshToken(refreshToken: string) {
    const body = new URLSearchParams({
      grant_type: "refresh_token",
      client_id: this.clientId,
      client_secret: this.clientSecret,
      refresh_token: refreshToken,
    });

    const res = await fetch("https://id.kick.com/oauth/token", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: body.toString(),
    });

    return res.json();
  }
}
```

### Webhook Handler (Express)

```typescript
import express from "express";
import crypto from "crypto";

const app = express();

// Raw body middleware for signature verification
app.use("/webhooks/kick", express.raw({ type: "application/json" }));

app.post("/webhooks/kick", (req, res) => {
  // Verify signature first
  const signature = req.headers["kick-event-signature"] as string;
  const messageId = req.headers["kick-event-message-id"] as string;
  const timestamp = req.headers["kick-event-message-timestamp"] as string;
  const eventType = req.headers["kick-event-type"] as string;

  const isValid = verifyWebhookSignature(
    messageId,
    timestamp,
    req.body,
    signature,
  );
  if (!isValid) {
    return res.status(401).send("Invalid signature");
  }

  // Parse body
  const payload = JSON.parse(req.body);

  // Handle event
  switch (eventType) {
    case "chat.message.sent":
      handleChatMessage(payload);
      break;
    case "channel.followed":
      handleFollow(payload);
      break;
    case "livestream.status.updated":
      handleStreamStatus(payload);
      break;
    // ... etc
  }

  // Return 2xx quickly
  res.status(200).send("OK");
});
```

### API Client (Go)

```go
package kick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL   = "https://api.kick.com"
	OAuthURL  = "https://id.kick.com"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    BaseURL,
		token:      token,
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader *bytes.Reader
	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bs)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *Client) SendChatMessage(ctx context.Context, content string, broadcasterID int, msgType string) error {
	body := map[string]interface{}{
		"content": content,
		"type":    msgType,
	}
	if broadcasterID > 0 {
		body["broadcaster_user_id"] = broadcasterID
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/public/v1/chat", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("kick API error: %d", resp.StatusCode)
	}

	return nil
}
```

---

## Changelog

Recent notable changes:

- **Jan 2026**: Categories V2 endpoint added, V1 deprecated. Token introspect moved to `/oauth`.
- **Dec 2025**: Auto-unsubscribe for failing webhooks. App deletion added.
- **Nov 2025**: Refresh tokens made reusable. `pinned_time_seconds` added to kicks.gifted.
- **Oct 2025**: KICKs leaderboard endpoint added. `kicks.gifted` webhook event added.
- **Aug 2025**: `replies_to` field added to chat messages.
- **Jul 2025**: Multiple `broadcaster_user_id` params on livestreams. `created_at` added to chat messages.
- **May 2025**: Moderation banned webhook. Moderation endpoints.
- **Apr 2025**: Reply chat messages. Channel subscriptions `expires_at`. Livestreams endpoint.

For full changelog: https://docs.kick.com/changelog.md
