# Scouti API (`/api/v1`)

The gate every agent and the `scouti` CLI talk to. One base URL, one auth header,
one predictable envelope. All calls below are shown as `scouti request`, but the
CLI is just a thin forwarder — the HTTP contract is identical if you call it
directly.

- **Base URL:** `https://scouti.chat/api/v1` (override with `SCOUTI_API_URL`).
- **Auth:** `Authorization: Bearer <access-key>` on every call. The CLI attaches
  it for you; you never pass it by hand.
- **Content type:** requests and responses are JSON.

## Response envelope

- **Success** → the resource itself, no wrapper. `200` for reads/updates, `201`
  for creates, `{ "deleted": true }` for deletes.
- **Error** → `{ "error": { "code": "...", "message": "..." } }` with a non-2xx
  status. `scouti request` prints this and exits non-zero.

| Status | `code`            | Meaning                                        |
| ------ | ----------------- | ---------------------------------------------- |
| 400    | `invalid_request` | Bad/missing fields in the body or query        |
| 401    | `unauthorized`    | Missing or invalid access key                  |
| 403    | `forbidden`       | Authenticated, but not a member of that org    |
| 404    | `not_found`       | Project / topic / … doesn't exist or isn't yours |
| 402/*  | `billing_*`       | Org can't run (see `verify` → `billing_reason`)  |
| 500    | `db_error` / `internal_error` | Something broke server-side        |

## Conventions

- Everything except `/me` and the auth endpoints is **project-scoped** under
  `/projects/{projectId}/…`. Get `projectId` from `GET /me`.
- A wrong or inaccessible `projectId` (or `topicId` / `touchpointId`) always
  returns `404 not_found` — never a `200` with empty or zeroed data. So zeros in
  a `status` / list response mean "no activity yet", not "bad id". If you get a
  `not_found`, re-check the id against `GET /me`.
- Basic config (product doc, tag tree, attention list) is **read + modify only** —
  every project is created with these, so there is no create/delete for them.
- CSV query params take comma-separated values, e.g. `?status=summarized,archived`.

## Endpoint index

| Method & path | Purpose |
| --- | --- |
| `GET /me` | List orgs/projects you can act on |
| `GET /projects/{id}` | Read a project (name, domains, metadata, doc) |
| `PATCH /projects/{id}` | Update name / domains / metadata |
| `GET·PATCH /projects/{id}/doc` | Read/replace the product doc (plain text) |
| `GET·PATCH /projects/{id}/tag-tree` | Read/replace the tag tree (≤3 levels deep) |
| `GET·PATCH /projects/{id}/attention-list` | Read/replace the attention list (≤5) |
| `GET·POST·DELETE /projects/{id}/keys` | List / mint / revoke widget keys (`pk_`) |
| `GET·POST /projects/{id}/topics` | List / create topics |
| `PATCH·DELETE /projects/{id}/topics/{topicId}` | Update / delete a topic |
| `GET·POST /projects/{id}/touchpoints` | List / create touchpoints |
| `PATCH·DELETE /projects/{id}/touchpoints/{tpId}` | Update / delete a touchpoint |
| `GET /projects/{id}/conversations` | Search summarized conversations |
| `GET /projects/{id}/users` | List end-users (outreach targets) |
| `POST /projects/{id}/outreach` | Start an outreach conversation with a user |
| `POST /projects/{id}/conversations/{cid}/outreach` | Append to an outreach thread |
| `GET /projects/{id}/status` | Mission-Control metrics for a time window |
| `GET /projects/{id}/verify` | Config + billing health check |

---

## Identity

### `GET /me`

Your first call after login. Returns the orgs and projects you can act on.

```bash
scouti request GET /me
```

```json
{
  "user": { "id": "uuid" },
  "orgs": [
    {
      "id": "org-uuid",
      "name": "alice's Organization",
      "role": "owner",
      "projects": [{ "id": "proj-uuid", "name": "alice's Project" }]
    }
  ]
}
```

---

## Project & basic config

### `GET /projects/{id}`

```bash
scouti request GET /projects/PROJECT_ID
```

Returns `{ id, name, organization_id, metadata, client_auth, doc, created_at }`.
`client_auth` holds allowed domains and widget keys; `metadata` holds the tag tree
and attention list (edit those via their dedicated endpoints below); `doc` is the
product doc as **plain text** (the server stores it base64-encoded, but you never
deal with base64).

### `PATCH /projects/{id}`

Update any of `name`, `domains` (allow-list for the widget, ≤ 50 entries), or
`metadata` (free-form extras, shallow-merged). At least one field is required.
`tag_tree` and `attention_list` are **not** accepted here — they're validated, so
set them through their dedicated endpoints below.

```bash
scouti request PATCH /projects/PROJECT_ID '{"domains":["example.com","app.example.com"]}'
```

### `GET·PATCH /projects/{id}/doc`

The **product doc** that grounds every Scout conversation. Send and receive
**plain text** — the server handles base64 for storage. Max 20,000 characters.

```bash
scouti request GET   /projects/PROJECT_ID/doc
scouti request PATCH /projects/PROJECT_ID/doc '{"doc":"# My Product\n\nWhat it does, who it is for, current status…"}'
```

Both return `{ "doc": "..." }` (plain text).

### `GET·PATCH /projects/{id}/tag-tree`

The taxonomy conversations get tagged against. `tag_tree` is an array of nodes:

```jsonc
{ "label": "usability", "description": "…", "color": "#14b8a6",
  "children": [ { "label": "onboarding", "description": "…" } ] }
```

Rules (rejected with `invalid_request` otherwise): nesting **≤ 3 levels** deep,
every node needs a non-empty `label` (≤ 60 chars); `description` ≤ 500 chars;
`color` optional. `PATCH` replaces the whole tree.

```bash
scouti request PATCH /projects/PROJECT_ID/tag-tree @tag-tree.json
```

Returns `{ "tag_tree": ... }`.

### `GET·PATCH /projects/{id}/attention-list`

`attention_list` is up to **5** short strings (≤ 300 chars each) the Scout should
keep an ear out for. `PATCH` replaces the whole list; blank entries are dropped.

```bash
scouti request PATCH /projects/PROJECT_ID/attention-list '{"attention_list":["pricing confusion","mobile bugs"]}'
```

Returns `{ "attention_list": [...] }`.

---

## Widget keys (`pk_`)

Publishable keys embedded in front-end code — safe to expose, unlike your access
key.

```bash
scouti request GET    /projects/PROJECT_ID/keys
scouti request POST   /projects/PROJECT_ID/keys '{"name":"web"}'       # → 201, { key, name, ... }
scouti request DELETE "/projects/PROJECT_ID/keys?key=pk_live_xxx"      # → { deleted: true }
```

---

## Topics

A **topic** is one thing you want to learn from users.

**Fields you can send** (create/patch). Anything else is rejected with
`invalid_request`:

| Field | Type | Notes |
| --- | --- | --- |
| `name` | string | Required on create. Display name (≤ 80 chars). |
| `type` | `reactive` \| `proactive` | Required on create. `reactive` = waits for the user; `proactive` = the Scout opens. |
| `hint` | object | The design brief — **exactly these four keys**, all optional strings (≤ 4,000 chars each): `context`, `goal`, `requirements`, `plan`. Unknown keys are rejected. See below. |
| `openings` | string[] | Opener lines (≤ 500 chars each); the Scout uses one per conversation. Blank entries are dropped; **max 10**. |
| `status` | `active` \| `archived` | Patch only. |
| `enabled` | boolean | Whether the topic can fire. |
| `direct_link_enabled` | boolean | Turns on the shareable direct link. |
| `metadata` | object | Free-form extras, shallow-merged into the row's metadata. |

**The `hint` brief** is not a free sentence — it's a structured object with four
fixed sections:

| Key | What goes here |
| --- | --- |
| `context` | The product moment / situation the topic stays grounded in. |
| `goal` | What you want to learn from the user. |
| `requirements` | Constraints or must-dos for the conversation. |
| `plan` | The interview path / follow-up sequence to prioritize. |

Reactive topics really only use `requirements` (the rest stays open-ended);
proactive topics use all four.

### `GET·POST /projects/{id}/topics`

`GET` returns rows shaped `{ id, created_at, hint, openings, metadata, status }`
(`hint` is the JSON string above; `metadata` carries `name`, `type`, `enabled`,
`direct_link_enabled`).

```bash
scouti request POST /projects/PROJECT_ID/topics @topic.json
```

```json
{
  "name": "Onboarding friction",
  "type": "proactive",
  "hint": {
    "context": "First session, right after signup.",
    "goal": "Find where new users get stuck.",
    "requirements": "Keep it short and specific; don't lead the user.",
    "plan": "1. Ask how setup went. 2. Dig into any friction. 3. Thank them."
  },
  "openings": ["Hey! Mind if I ask how setup went for you?"],
  "metadata": { "enabled": true }
}
```

### `PATCH·DELETE /projects/{id}/topics/{topicId}`

Patch any subset of the fields in the table above (at least one required).

```bash
scouti request PATCH  /projects/PROJECT_ID/topics/TOPIC_ID '{"status":"archived"}'
scouti request PATCH  /projects/PROJECT_ID/topics/TOPIC_ID '{"hint":{"requirements":"Only ask on mobile."}}'
scouti request DELETE /projects/PROJECT_ID/topics/TOPIC_ID
```

**Direct link:** enable it, and the topic opens at `https://scouti.chat/t/TOPIC_ID`:

```bash
scouti request PATCH /projects/PROJECT_ID/topics/TOPIC_ID '{"direct_link_enabled":true}'
```

---

## Touchpoints

A **touchpoint** is where/when topics are surfaced in your product. Rows are
`{ id, name, config, status }`.

**Top-level fields:**

| Field | Type | Notes |
| --- | --- | --- |
| `name` | string | Required on create. A slug: `^[_a-z0-9][a-z0-9_-]*$`, ≤48 chars, unique per project. |
| `config` | object | The behaviour below. Unknown keys are rejected. |
| `status` | `enabled` \| `disabled` | Patch only. |

**`config` fields** (all optional; anything else → `invalid_request`):

| Key | Type | Notes |
| --- | --- | --- |
| `mode` | `reactive` \| `proactive` | Must match the `type` of the topics bound here (enforced server-side). |
| `force_show` | boolean | Bypass the normal show/skip heuristics. |
| `probability` | number `0`–`1` | Mount-time sampling rate — the touchpoint only fires if a per-mount roll passes. Omitted = always fire. |
| `y_position` | number `0`–`1` | Vertical placement of the widget button, as a fraction of viewport height (0 = top). |
| `topics` | array | Topic **bindings** — which topics this touchpoint can surface (see below). |

**Each `topics[]` binding:**

| Key | Type | Notes |
| --- | --- | --- |
| `topic_id` | string | Required. The topic to surface. |
| `weight` | number ≥ 0 | Relative pick weight when several topics compete. |
| `conditions` | array | Gates on when the topic may fire. Each item needs a `type` of `time_range`, `weekdays`, or `cooldown`; its other params are passed through as-is. |

```bash
scouti request GET    /projects/PROJECT_ID/touchpoints
scouti request POST   /projects/PROJECT_ID/touchpoints @touchpoint.json
scouti request PATCH  /projects/PROJECT_ID/touchpoints/TP_ID '{"status":"disabled"}'
scouti request DELETE /projects/PROJECT_ID/touchpoints/TP_ID
```

```json
{
  "name": "post_checkout",
  "config": {
    "mode": "proactive",
    "probability": 0.3,
    "topics": [
      { "topic_id": "TOPIC_ID", "weight": 1, "conditions": [{ "type": "cooldown", "minutes": 1440 }] }
    ]
  }
}
```

---

## Conversations

### `GET /projects/{id}/conversations`

Search summarized conversations (same data as the dashboard list). All params
optional.

| Param | Meaning |
| --- | --- |
| `from`, `to` | ISO date bounds |
| `q` | full-text search |
| `sentiment`, `quality`, `status`, `topic`, `tag` | CSV filters (`topic` = topic ids) |
| `limit` | default `100` |
| `offset` | default `0` |

```bash
scouti request GET "/projects/PROJECT_ID/conversations?status=summarized&limit=20"
```

---

## Users & outreach

### `GET /projects/{id}/users`

The pool of end-users you can reach out to. Params: `page` (1), `page_size` (12),
`query`/`q`, `sort` (`time_desc`).

```bash
scouti request GET "/projects/PROJECT_ID/users?page=1&page_size=20"
```

### `POST /projects/{id}/outreach`

Start a new outreach conversation on a user's most recent deliverable session.
Pick `session_ids` from the users list.

```bash
scouti request POST /projects/PROJECT_ID/outreach '{"message":"Hi! Got 30s for a quick question?","session_ids":["sess-uuid"]}'
```

→ `201 { conversation_id, session_id, pending_count, queued_at }`.

### `POST /projects/{id}/conversations/{cid}/outreach`

Append another message to an existing outreach thread.

```bash
scouti request POST /projects/PROJECT_ID/conversations/CONV_ID/outreach '{"message":"Following up — still keen to hear your thoughts!"}'
```

→ `{ pending_count, queued_at }`.

---

## Status & verify

### `GET /projects/{id}/status`

A Mission-Control snapshot for a time `window` (`7d`, `2w`, `1m`, or a bare day
count; default `7d`).

```bash
scouti request GET "/projects/PROJECT_ID/status?window=7d"
```

→ `{ window_days, depth, users, top_tags, notifications }`.

### `GET /projects/{id}/verify`

"Is the widget wired up and can it run?" Fix anything `false`, then re-run.

```bash
scouti request GET /projects/PROJECT_ID/verify
```

```json
{
  "domains_ok": true,
  "has_active_topic": true,
  "has_enabled_touchpoint": true,
  "has_project_key": true,
  "doc_sufficient": true,
  "billing_ok": true,
  "billing_reason": null
}
```

---

## Auth (device flow) — handled by the CLI

You normally never call these; `scouti login` and the browser approval page do.
Documented for completeness.

- `POST /auth/device/start` → `{ device_code, user_code, verification_uri, verification_uri_complete, interval, expires_in }`. Unauthenticated.
- `POST /auth/device/token` body `{ device_code }` → `{ status: "pending", interval }` while waiting, or `{ access_key, token_type: "uak", user_id }` once approved. Unauthenticated; the key is returned exactly once.
- `POST /auth/device/approve` body `{ user_code }` → `{ approved: true }`. Requires a logged-in browser session; also provisions a default workspace.
