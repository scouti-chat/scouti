# Scouti Devkit

One-click **Scouti** integration for AI coding agents. Scouti gives your product a
small virtual team that talks to your users and turns those conversations into
product insight. This kit lets an agent (Claude Code / Codex / Cursor) sign you
in, provision a project, install the widget, and set up your feedback topics.

Website: **https://scouti.chat**

## Quickstart

Tell your coding agent:

> Set up Scouti for my product — read `skill/SKILL.md` and follow it.

## What's here

- [`skill/`](skill/) — the agent playbook (`SKILL.md`). Start here.
- [`cli/`](cli/) — the `scouti` CLI (Go): holds your access key and forwards
  authenticated calls to the Scouti API. See [`cli/README.md`](cli/README.md).

`CLAUDE.md`, `AGENTS.md`, and `.cursor/rules/` just point agents at the skill.
