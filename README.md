# FounderPing Devkit

One-click **FounderPing** integration for AI coding agents. FounderPing gives your product a
small virtual team that talks to your users and turns those conversations into
structured product insight. This repo builds the two release artifacts an agent
uses to wire FounderPing into *your* project:

- a **skill** tarball (`skill.tar.gz`) — the agent playbook plus the product guide
  and API reference, fully self-contained, and
- the **`founderping` CLI** — a single static binary that holds your access key and
  forwards authenticated calls to the FounderPing API.

Website: **https://founderping.app**

## Install (paste this to your coding agent)

You don't clone this repo. In your own project, paste the prompt below into Claude
Code / Codex / Cursor:

> Add FounderPing (an AI user-feedback system) to this project. Download its skill bundle
> — `https://github.com/founderping/founderping/releases/latest/download/skill.tar.gz` —
> unpack it into a folder here (e.g. `founderping-skill/`), then open `SKILL.md` inside and
> follow it.

That's the whole bootstrap. `SKILL.md` drives the rest — CLI download, sign-in,
project setup, widget install, topic design, and reading feedback back.

## What's in this repo

- [`skill/`](skill/) — source of the published skill: `SKILL.md` (authored) plus
  `guide.md` and `api.md` (synced from the main product docs). Shipped as
  `skill.tar.gz`.
- [`cli/`](cli/) — the `founderping` CLI (Go). Shipped as per-platform binaries. See
  [`cli/README.md`](cli/README.md).

## Releases

`make dist` (from this directory) builds the full, upload-ready set into `./dist`:

- `founderping-<os>-<arch>[.exe]` — raw CLI binary per platform (direct download).
- `founderping-<os>-<arch>.tar.gz` — CLI archive per platform.
- `skill.tar.gz` — the installable skill bundle (see the prompt above).

Run `sc sync_assets` in the main repo first so `skill/guide.md` and `skill/api.md`
are current, then `make dist` and upload everything in `./dist` to a GitHub
release. CI does this automatically on a `v*` tag (`.github/workflows/release.yml`).
