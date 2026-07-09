# Scouti Devkit

One-click **Scouti** integration for AI coding agents. Scouti gives your product a
small virtual team that talks to your users and turns those conversations into
structured product insight. This repo builds the two release artifacts an agent
uses to wire Scouti into *your* project:

- a **skill** tarball (`skill.tar.gz`) — the agent playbook plus the product guide
  and API reference, fully self-contained, and
- the **`scouti` CLI** — a single static binary that holds your access key and
  forwards authenticated calls to the Scouti API.

Website: **https://scouti.chat**

## Install (paste this to your coding agent)

You don't clone this repo. In your own project, paste the prompt below into Claude
Code / Codex / Cursor:

> Add Scouti (an AI user-feedback system) to this project. Download its skill bundle
> — `https://github.com/scouti-chat/scouti/releases/latest/download/skill.tar.gz` —
> unpack it into a folder here (e.g. `scouti-skill/`), then open `SKILL.md` inside and
> follow it.

That's the whole bootstrap. `SKILL.md` drives the rest — CLI download, sign-in,
project setup, widget install, topic design, and reading feedback back.

## What's in this repo

- [`skill/`](skill/) — source of the published skill: `SKILL.md` (authored) plus
  `guide.md` and `api.md` (synced from the main product docs). Shipped as
  `skill.tar.gz`.
- [`cli/`](cli/) — the `scouti` CLI (Go). Shipped as per-platform binaries. See
  [`cli/README.md`](cli/README.md).

## Releases

`make dist` (from this directory) builds the full, upload-ready set into `./dist`:

- `scouti-<os>-<arch>[.exe]` — raw CLI binary per platform (direct download).
- `scouti-<os>-<arch>.tar.gz` — CLI archive per platform.
- `skill.tar.gz` — the installable skill bundle (see the prompt above).

Run `sc sync_assets` in the main repo first so `skill/guide.md` and `skill/api.md`
are current, then `make dist` and upload everything in `./dist` to a GitHub
release. CI does this automatically on a `v*` tag (`.github/workflows/release.yml`).
