---
name: scouti-integration
description: Integrate the Scouti user-feedback / thought-capture system into a product — sign in, provision a project, install the widget, design topics & touchpoints, and verify.
---

# Scouti one-click integration

You are helping a developer wire **Scouti** into their product. Scouti pairs the
product with a small virtual team that talks to end-users in real time and turns
those conversations into product insight. Your job is to take the developer from
"nothing" to "widget live + feedback topics configured + verified", with as few
manual steps as possible.

Website: https://scouti.chat

## The two rules that keep this safe

1. **Never handle the access key (UAK).** The `scouti` CLI stores and sends it for
   you. Never print it, never write it into code or config, never paste it into
   the chat. If you ever see a `uak_…` value, stop and tell the developer.
2. **Never change their code silently.** Every edit to the developer's app
   (installing the widget, adding `mount()`) is proposed as a diff and applied
   only after they've seen it. Explain any write to Scouti before you make it.

## Division of labor

- **Deterministic work → the API, via `scouti request`.** Auth, project lookup,
  creating topics/touchpoints, minting keys, verifying — all live server-side.
  You just call them.
- **Judgment work → you.** Detect the tech stack, decide where the widget and
  `mount()` belong, **interview the developer about what they want to learn from
  their users**, and turn that into well-written topics.

## Getting the `scouti` CLI

`scouti` is a single self-contained binary — no runtime to install. Fetch the one
matching this machine and use it for every step below.

```bash
# Detect os (linux|darwin|windows) and arch (amd64|arm64), then:
curl -fsSL "https://github.com/scouti-chat/scouti/releases/latest/download/scouti-<os>-<arch>" -o scouti
chmod +x scouti     # Windows: download scouti-windows-<arch>.exe, no chmod
```

If the download isn't reachable and Go is installed, you can build from source in
this devkit instead: `make -C ../cli build`. Don't install anything that would
make the developer approve a system change without telling them first.

## Workflow

Run these in order. Each `scouti request` prints JSON; read it before moving on.

1. **Log in.** Run `scouti login`. It opens a browser to sign in / sign up and
   provisions a default workspace. On success the CLI holds the key locally.
2. **Locate the project.** `scouti request GET /me` → pick the `orgs[].projects[]`
   to work in (usually the only one). Call its id `PROJECT_ID` below.
3. **Mint an embeddable key.** `scouti request POST /projects/PROJECT_ID/keys`
   → a `pk_…` client key for the widget. (This is *not* the UAK; it's safe in
   front-end code.)
4. **Register the dev domain.** `scouti request PATCH /projects/PROJECT_ID`
   with the app's domain(s) so the widget is allowed to load.
5. **Install the widget.** Add the loader + `mount()` to the app. Detect the
   framework first; propose the change as a diff. (See "Installing the widget".)
6. **Interview, then design topics.** Ask the developer what they most want to
   learn from users. Turn each answer into a topic:
   `scouti request POST /projects/PROJECT_ID/topics @topic.json`, and a
   touchpoint that surfaces it: `POST /projects/PROJECT_ID/touchpoints`.
   (See "Writing good topics".)
7. **Verify.** `scouti request GET /projects/PROJECT_ID/verify` → confirms
   domains, an active topic, an enabled touchpoint, a key, the product doc, and
   billing are all OK. Fix anything that comes back false, then re-run.

Later / ongoing: `GET /projects/PROJECT_ID/status?window=7d` for a health
snapshot, `GET …/conversations?status=summarized` to read summarized sessions,
`GET …/users` + `POST …/outreach` to reach out to a specific user.

## Calling the API

The CLI is a thin authenticated forwarder (think `gh api`):

```bash
scouti request <METHOD> <PATH> [body]
# body: @file.json  |  -  (stdin)  |  an inline JSON string
```

- Base path is always `/api/v1`; write just the `PATH` (e.g. `/me`).
- Large payloads go through a file (`@topic.json`) or stdin, not the command line.
- A non-zero exit code means the call failed — read the printed
  `{ "error": { "code", "message" } }` and react; don't blindly retry.

Endpoint reference: [`../api.md`](../api.md) — every endpoint, its parameters, and
example payloads.

## Installing the widget

TODO: fill in per-framework loader snippets (plain script tag, React/Next,
Vue, …) and where `mount()` should go. Principle: the loader is inert until
`mount()` runs, so place `mount()` where the feedback prompt makes sense
(e.g. after a key action), not blindly in the root layout.

## Writing good topics

A **topic** is one thing you want to learn from users; a **touchpoint** is where
and when it gets surfaced. Good topics are specific, open-ended, and tied to a
real product decision — not "Any feedback?".

- `metadata.type` is `reactive` (waits for the user) or `proactive` (Scouti
  starts the conversation). A touchpoint's `config.mode` must match its topic's
  type.
- Write from the "developer's scout" voice: curious, concrete, one focus per
  topic.

TODO: embed 2–3 worked topic examples + the field-by-field schema.

## When something blocks you

Stop and ask the developer if: login can't complete, a call returns 402/billing,
`verify` keeps failing after fixes, or you're unsure where to place `mount()`.
Don't guess on their codebase or spend the key on retries.
