# scouti CLI

A tiny Go CLI with exactly two jobs:

1. **hold your Scouti access key** securely on this machine, and
2. **forward authenticated requests** to the Scouti API (`/api/v1`).

Think of it as "`curl` that carries your key" (à la `gh api`). It ships as a
single static binary — no Node, no Python, no runtime to install. All product
logic lives server-side; the CLI stays deliberately thin so there's nothing to
learn beyond `login` and `request`.

- Full endpoint reference: [`../skill/api.md`](../skill/api.md)
- Agent playbook (how to drive an integration end-to-end): [`../skill/SKILL.md`](../skill/SKILL.md)
- Website: https://scouti.chat

---

## Install

### Download a prebuilt binary (recommended)

Grab the file for your platform from the
[latest release](https://github.com/scouti-chat/scouti/releases/latest) and put
it on your `PATH`:

```bash
# Linux / macOS — auto-detects os + arch
os=$(uname -s | tr '[:upper:]' '[:lower:]')          # linux | darwin
arch=$(uname -m); [ "$arch" = "x86_64" ] && arch=amd64; [ "$arch" = "aarch64" ] && arch=arm64
curl -fsSL "https://github.com/scouti-chat/scouti/releases/latest/download/scouti-${os}-${arch}" -o scouti
chmod +x scouti
sudo mv scouti /usr/local/bin/            # or anywhere on PATH
```

On Windows, download `scouti-windows-amd64.exe` (or `-arm64`) and put it on your
`PATH`.

### Build from source

Requires Go 1.22+. The Makefile lives at the **devkit root**, so run these from
there (`cd ..`):

```bash
make build      # build ./cli/scouti for the current platform
make dist       # cross-compile every platform + archives + skill → ../dist
make clean
```

---

## Authentication model

- Your identity is a **user access key** (`uak_…`). It's secret — treat it like a
  password.
- The CLI stores it at `~/.scouti/credentials.json` with `0600` permissions and
  **never prints it**.
- Every `scouti request` attaches it as `Authorization: Bearer …` automatically.
- In CI or headless jobs, skip the file and pass the key via the
  `SCOUTI_ACCESS_KEY` environment variable instead.

---

## `scouti login`

Authorize this machine and save the key.

```bash
scouti login                 # interactive: browser device flow
scouti login --token uak_xxx # store a pre-issued key (CI / headless)
```

**Device flow (default).** No key is ever typed or pasted:

1. The CLI starts an authorization and prints a URL plus a short `user_code`.
2. It tries to open your browser; if it can't (e.g. a remote/SSH session), copy
   the URL into any browser yourself.
3. Sign in / sign up and approve — confirm the `user_code` on screen matches the
   one the CLI printed. A default workspace is provisioned for new accounts.
4. The CLI polls in the background and, once approved, saves the key. Done.

The `user_code` is only for eyeball confirmation; the CLI polls with a separate,
high-entropy code that never leaves your machine.

---

## `scouti request`

Forward one authenticated call and print the JSON response.

```
scouti request <METHOD> <PATH> [body]
```

- `<METHOD>` — `GET`, `POST`, `PATCH`, `DELETE`, … (case-insensitive).
- `<PATH>` — relative to `/api/v1`, always starting with `/` (e.g. `/me`,
  `/projects/PID/topics`). Quote paths that contain query strings.
- `[body]` — optional request body, from one of these sources:

| Form | Meaning |
| --- | --- |
| `path/to/file.json` | read the body from that file (bare path) |
| `@path/to/file.json` | read the body from that file (explicit; missing file errors) |
| `-` | read the body from **stdin** |
| `'{"inline":"json"}'` | use the argument as a literal JSON string |

Prefer a **file** (or stdin) for anything non-trivial — it keeps large payloads
off the command line and out of your shell history. The body is validated as JSON
before the request is sent, so you fail fast on typos.

### Exit codes

| Code | Meaning |
| --- | --- |
| `0` | HTTP 2xx |
| `1` | HTTP non-2xx, network error, or not logged in |
| `2` | Usage error (missing method/path) |

On failure the server's `{ "error": { "code", "message" } }` is printed — read it
and react; don't blindly retry.

### Examples

```bash
# Read
scouti request GET /me
scouti request GET "/projects/PID/status?window=7d"

# Create from an inline string
scouti request POST /projects/PID/keys '{"name":"web"}'

# Create from a file (bare path or @path both work)
scouti request POST /projects/PID/topics topic.json
scouti request POST /projects/PID/topics @topic.json

# Create from stdin
echo '{"name":"post_checkout","config":{"mode":"proactive"}}' \
  | scouti request POST /projects/PID/touchpoints -

# Update / delete
scouti request PATCH  /projects/PID/topics/TID '{"status":"archived"}'
scouti request DELETE "/projects/PID/keys?key=pk_live_xxx"
```

See [`../skill/api.md`](../skill/api.md) for every endpoint, its parameters, and
example payloads.

---

## Environment variables

| Variable | Purpose |
| --- | --- |
| `SCOUTI_API_URL` | Override the API origin (default `https://scouti.chat`). Use for staging / local dev. |
| `SCOUTI_ACCESS_KEY` | Use this key instead of `~/.scouti/credentials.json` (CI / headless). Takes precedence over the file. |

---

## Typical end-to-end flow

```bash
scouti login
scouti request GET /me                                   # find your PROJECT_ID
scouti request POST /projects/PID/keys '{"name":"web"}'  # widget key (pk_)
scouti request PATCH /projects/PID '{"domains":["example.com"]}'
scouti request POST /projects/PID/topics topic.json      # design topics
scouti request POST /projects/PID/touchpoints tp.json    # surface them
scouti request GET /projects/PID/verify                  # confirm it's live
```

---

## Troubleshooting

- **`Not logged in`** — run `scouti login`, or export `SCOUTI_ACCESS_KEY`.
- **Browser didn't open** — expected on remote/headless boxes; copy the printed
  URL into a browser manually.
- **`request body is not valid JSON`** — the body (inline, file, or stdin) wasn't
  valid JSON. Check the source named in the error.
- **`forbidden` / `not_found`** — you're not a member of that org/project, or the
  id is wrong. Re-check `GET /me`.
- **Point at a different backend** — set `SCOUTI_API_URL`, e.g.
  `SCOUTI_API_URL=http://localhost:3000 scouti request GET /me`.
