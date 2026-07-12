# Scouti — Product, Configuration & Implementation Guide

## 1. What Scouti is

Scouti is an **AI-powered user-feedback platform for product teams**. It helps you hear from the people who use your product—not through long surveys or static forms, but through **short, natural conversations** (voice or text) that an AI leads on your behalf. Each chat is turned into **structured, searchable insight**—what was said, how the user felt, what matters—so you can prioritize and act without manually reading every submission.

Think of Scouti as a **virtual team working for you**: you set the goals and rules; the team talks to users, brings back signal, and organizes it in one place.

### Problems it addresses

Most teams don’t have one feedback problem—they have **a patchwork of channels that each fail in a different way**. Scouti is aimed at replacing that patchwork, not adding another inbox.

**Reach doesn’t scale**

- **The same small circle.** “Begging the same 15 friends” on Reddit, X, or your network—kind people, tiny sample, no growth with your real users.
- **Email into the void.** Often the only way to reach signups, yet reply rates are brutal and answers are a line or two with no follow-up.
- **Audiences you never hear from.** Churned users, onboarding drop-offs, and the silent majority rarely show up in surveys or prompts.

**What you collect stays thin**

- **Forms and star ratings nobody fills out.** Typing is a chore, so users skip it—you get a lone ⭐⭐⭐ and no story behind it.
- **Shallow one-liners when they do reply.** “It’s broken” / “I love it” without clarification still leaves you guessing.
- **Numbers, never the why.** Analytics tell you *what* moved—churn up, a step dropped—but never *why*.

**Signal is fragmented and hard to turn into decisions**

- **Feedback scattered everywhere.** Discord, support mail, DMs, GitHub issues—truth in pieces, and no one has time to merge it.
- **Synthesis doesn’t scale.** Tagging, prioritizing, and tying feedback to a user stays manual, so backlogs outgrow insight.
- **Rolling your own feedback UX is costly.** Every new question, audience, or trigger still means more frontend work and redeploys.

Scouti is built to replace that loop: **reach users in the product (and beyond), let AI dig deeper in the moment, and put summarized, tagged signal in one Dashboard**—so you’re not guessing from stars, silence, or scattered threads.

---

## 2. End-to-end workflow (developer view)

This section walks the full lifecycle from a developer’s perspective: from signing up, through configuring what to ask, to reading and acting on what comes back. Each step links to its deep-dive section.

> **Two equivalent ways to operate Scouti.** Everything below can be done either in the **Dashboard** (point-and-click) or through the **`scouti` CLI / REST API** (`/api/v1`) that AI coding agents drive — projects, the product doc, topics, touchpoints, widget keys, allowed domains, outreach, and status are all reachable both ways. Pick whichever fits the moment; where it matters, this guide notes both. (For the full CLI/API surface, see the companion **API reference**, `api.md`.)

### Fastest setup — one prompt with your coding agent

If you use **Cursor, Claude Code, Codex, Gemini CLI**, or a similar assistant, you can stand up Scouti in **one line** without clicking through the Dashboard first:

1. Open your product's repo in the agent.
2. Paste a setup prompt—for example: *"Download the Scouti skill from https://github.com/scouti-chat/scouti, follow it, and help me set up Scouti in this project."*
3. The agent pulls the **Scouti skill** (this guide, the API reference, and setup scripts), walks you through login, creates your project if needed, co-designs Topics and Touchpoints with you, and drops the widget install script plus `mount()` calls into your codebase.

Prefer the UI? **Sign up free** and follow the onboarding tutorial in the Dashboard instead—the same end state, step by step.

### Step 0 — Set up the workspace

1. **Create an Organization and Project.** The Organization is the billing entity; a Project is one product/app you want feedback on.
2. **Write your product context.** Drop in a free-form description of your product—what it does, who uses it, key terms. Scouti injects this into every conversation so the AI sounds informed instead of generic. You write it in the project's **Playbook**.
3. **Define your taxonomy.** Set up the **Tag Tree** (how feedback should be labeled) and the **Attention List** (topics that should auto-flag a conversation, e.g. churn, billing, security). This is what turns raw chats into filterable, prioritized signal later (see §5.4–5.5).

### Step 1 — Decide what to ask (Topics)

1. **Create Topics.** A Topic is something you want the Scouts to talk to users about—a churn reason, a reaction to a new feature, an onboarding check. It comes in two flavors:
  - **Reactive Topic — the always-on helper.** Exactly **one per project**. Once enabled, it sits behind a small sticky side button in your product, so users can open it whenever they want to say something. It’s a general “tell us anything” helper, not a targeted question.
  - **Proactive Topics — targeted, pop-up questions.** **Zero to many per project**. Each one focuses on a single question or a specific class of issue, and the Scout raises it on its own at the right moment (e.g. after checkout, on a power-user milestone).

Both flavors are configured through the same two core concepts, produced and edited either right on the page or through the CLI / API:

- **Hint** — context for the AI, *not* shown to the user. It’s injected into the AI’s context during the chat so the Scout understands the background and what you’re trying to learn, and drives the conversation accordingly.
- **Opening** — the first line(s) the Scout uses to start the chat—the first thing the user actually sees. A few are generated per Topic, deliberately different from one another, and the Scout picks one at random when the Topic fires (so the same person doesn't always get the same hello). A strong Opening weaves three beats into 1–2 short sentences (~40 words total): a **light self-introduction** (the Scout makes clear it's reaching out on behalf of the product team), a **clear purpose** (what it's hoping to learn, grounded in the Topic's goal/context), and **the actual opener** — for a Proactive Topic, a concrete first question that goes right to the point its Plan wants to explore; for the Reactive Topic, an open invitation to share whatever's on the user's mind.

The Hint is split into short, **orthogonal** fields — each one captures a *different* dimension of the conversation, so any given detail belongs in exactly one field (don’t restate the same thing across them). A **Proactive Topic** uses four:

- **Context** — the *situation*: what the dialogue is about and where/when it happens (which screen, after which action or event, what kind of user), plus any helpful background. It sets the scene; the goal and the questions do **not** go here.
- **Goal** — the *objective*: the one core thing you want to learn from the user. A single sharp goal beats a grab-bag of five.
- **Requirements** — the *guard-rails*: constraints and must-dos for how the Scout behaves — tone, what it must cover, what to avoid, whether to ask for contact details, and so on.
- **Plan** — the *flow*: a high-level path for how the conversation should move (the order of beats and the outcomes to aim for), not exact question wording.

Keep every field tight — a sentence or two each (rough budget: ~30 words for Context, Goal, and Requirements, up to ~50 for Plan). The Hint is a quick brief the Scout skims before talking, not an essay; padding it only buries the signal.

A **Reactive Topic** is the open-ended “tell us anything” helper, so it uses just **one** of these — **Requirements**. With no fixed situation, single goal, or scripted flow to pin down, Context, Goal, and Plan are all dropped; only the guard-rails for how the Scout should behave still apply.

You don’t have to write any of this from scratch. You can **co-design the Topic with Scouti’s AI assistant** — describe your goal in a short chat and it drafts the Hint fields and Openings with you — or hand it to **your own AI** with the *“Co-design with your AI”* button, which points your assistant at this guide and has it walk you through the same fields. Everything stays editable by hand.

### Step 2 — Decide when and where to ask (reach users)

There are two ways to put a Topic in front of users. Most projects use both.

**Path A — The Widget (in your product)**

This is the main path: an embedded widget that brings your Topics to users right inside your app.

1. **Touchpoints — naming the moments.** A **Touchpoint** is a named hook you place in your code at a meaningful moment, e.g. `scouti.mount("post_checkout")`. It’s just a label for “something interesting happened here”; *what* should be asked there is decided later in the Dashboard, not in your code.
2. **Bind Topics to Touchpoints (with Conditions).** In the Dashboard (or via the API) you bind one or more Topics to each Touchpoint, so firing that moment can start the right conversation. Bindings can carry **Conditions** that gate when they’re allowed to fire—time windows, weekdays, and per-user cooldowns (see *Appendix A — Touchpoint Condition types* for the full list).
  The **Reactive Topic** doesn’t need a Touchpoint—when enabled it always shows as the sticky “tell us anything” button. **Proactive Topics** are the ones you bind to Touchpoints so they pop up at the right time.
3. **Integrate the widget once.** Add the one-line widget script to your product and place your `mount(...)` calls at the moments you named. After this, which Topic fires where, and under what Conditions, is all controlled from the Dashboard (or the API)—no further frontend changes per question. (Concrete snippets, the JavaScript API, and the runtime mechanics—coalesce, Conditions, and Guards—are covered in §4 *Developer integration*.)

**Path B — Direct Links (outside your product)**

A **Direct Link** is a shareable URL (and QR code) that opens a single Topic in a full-screen chat, with **no SDK integration required**. It’s useful for reaching people outside your app—emails, social posts, support replies, printed QR codes, etc. You just pick a Topic, optionally set a branded subdomain, and share the link; anyone who opens it lands straight in that conversation (see §5.8).

### Step 3 — Conversations happen

1. **Users talk, not type.** The Scout runs a short voice-or-text chat, asking follow-ups in the moment to turn “it’s broken” into something specific. Anonymous sessions stick to a real user once your app supplies a User Identifier (see §4.2).
2. **Everything is captured.** Each chat is stored as a Conversation tied to its Topic, user, and moment—ready for processing.

### Step 4 — Feedback is processed automatically

1. **Summarized into Points.** After a conversation ends, the Analyst pipeline extracts atomic **Points**: a one-line gist, a short summary, sentiment, tags from your Tag Tree, and an attention flag (see §5.9).
2. **Organized for you.** Points are labeled and flagged so urgent or on-watch items surface instead of getting buried.

### Step 5 — Read, prioritize, and act

1. **Review in the Dashboard.** Browse the Conversations and Users views, filter by tag/sentiment/attention, expand any row for the full chat, and save curated sets as **Bookmarks**.
  Two ways to keep a conversation going via **Outreach**:
  - **Follow up on a specific conversation.** Reply directly in the thread of an existing conversation to ask a clarifying question or dig deeper into what a user already said.
  - **Re-engage an existing user.** Start a fresh, direct conversation with a user you’ve talked to before.
  Either way the message reaches the user the next time they’re back in your product, and their reply flows into the Dashboard like any other conversation.
2. **See the big picture.** **Mission Control** shows trends, sentiment and quality breakdowns, top tags, and a daily digest email (see §5.7).
3. **Close the loop and iterate.** Acting on feedback.

---

## 3. Designing your feedback program

This chapter is written mainly for you, the AI assistant a developer has brought in to help them get more out of Scouti — though it’s really about *their* product, so a developer reading along will get just as much from it. (“you” here means the assistant; elsewhere in this guide it means the developer.) Treat it as background for a role you can play, not a procedure you have to run: most of the time the developer arrives with a concrete question, and helping with that comes first. What follows is how, around that, you can also help them shape a feedback program worth having.

### 3.1 Why proactive, in-product feedback is worth it

In-product feedback is one of the highest-leverage things a team can build, and also one of the easiest to undervalue. Most teams learn too little, too late — through the patchwork of failing channels in §1: a tiny circle of friends, emails into the void, star ratings with no story, and silence from the users who matter most. Asking **in the moment, inside the product** flips that — the barrier to speak is low, the experience is fresh, and a Scout can dig past “it’s broken” to the actual why. For the developer that means a steady stream of specific, attributable signal: roadmap calls grounded in what real users said instead of guesswork. It’s worth carrying that perspective into the conversation — a developer unsure whether users will even talk is often just a few well-placed asks away from finding out they will.

### 3.2 Lead with their question, design alongside

A developer usually arrives with something concrete — often a prompt pasted straight from an “Ask your AI” / “Co-design with your AI” button in the Dashboard, sometimes a specific question about the conversation they’re building. **Help with that first.** It’s what they came for, and it’s what earns you the room to suggest more.

Keep this in mind as you do: because you understand their project and you’ve read this guide, you’re unusually well-placed to notice feedback opportunities they might miss. Developers are close to their product and often don’t realize which moments would yield the richest signal — or that a question they’re hand-wiring could be a Topic they retarget from the Dashboard without touching code again. The better you understand the product — from the repo, the Playbook (§5.3), or just the conversation — the sharper that instinct gets.

So when you spot a good idea or a well-timed moment, **offer it**: say what it would capture and why it’s worth it, and if they’re interested, help them design and wire it. The point isn’t to push a process — it’s to use what you know to leave their product a little better instrumented than you found it. Suggest, don’t insist, and follow their lead on how far to take it.

### 3.3 What makes a good Moment

Not every screen or click is worth a feedback ask — but some moments consistently yield rich signal. The right Moment is usually one where the user has *just* formed an opinion: the experience is still fresh, they know enough to answer, and a short question would feel natural rather than intrusive. Keep these principles in mind as background.

**Principles that help you spot good Moments**

- **Fresh feeling, still in context** — the user has just done something meaningful and can describe what happened while it's vivid.
- **A decision or emotion in play** — they chose, succeeded, failed, paid, canceled, discovered something new, or hit a wall. Those are when people have something to say.
- **Answerable in a short chat** — the question fits the moment; you're not interrupting something that demands their full attention.
- **Tied to a product decision the developer cares about** — the signal would actually change how the team builds, prices, or prioritizes.

These are heuristics, not a checklist. Every product has its own journey; the developer's domain knowledge and gut feel matter more than any list we could publish.

### 3.4 Help end-to-end, not just answer

Your job isn't to fire back one-liners when the developer asks a question — it's to help them **design a feedback touchpoint all the way through**, from the first idea to something they can ship. When they mention a goal, a feature, or a hunch about what to learn, think about the whole path and offer to walk it with them.

**Start with their question, then widen the lens.** They might ask how to word an Opening or whether to use reactive or proactive. Answer that — and if you see the adjacent pieces they haven't thought about yet, bring them up: a Moment worth instrumenting (§3.3), the Hint fields that would shape the Scout's conversation, the Touchpoint name that keeps the code stable, the conditions so the ask doesn't nag, the exact `mount()` / `setUser()` call and where it belongs. Don't wait to be asked for each step.

**Hand back concrete artifacts.** Where you can, deliver the drafted Hint fields and Openings, the touchpoint name, the Dashboard conditions to set, and the code snippets — not a vague plan they still have to interpret. The developer should leave with something they can paste, configure, and ship.

Follow their lead on how far to go on any given thread; Your default stance is **design partner**, not **FAQ bot**: connect the dots across product, copy, configuration, and code so a feedback idea becomes a working integration.

---

## 4. Developer integration

One principle runs through all of it: **your code only names moments—the Dashboard decides what to ask and when.** You integrate once; from then on you add, retarget, or retire questions entirely from the Dashboard, with no frontend changes or redeploys.

**Best practice: lay down Touchpoints early, configure Topics later.** The decoupling above pays off most when you front-load the code work: walk through your product once and place a `scouti.mount(...)` at every meaningful moment you might ever want feedback from—checkout complete, onboarding milestones, error states, feature discovery, a settings save. Register each name in the Dashboard (or let the first `mount` create it). You don't have to bind a Topic to every hook on day one; leave bindings empty until you're ready. What you gain is a stable map of moments that product, growth, and support can wire up entirely from the Dashboard: create a Proactive Topic, bind it to `post_checkout`, tune Conditions and weights, retire last quarter's question and point the same hook at a new one—all without another frontend change or redeploy. The slow path is the opposite: adding a new `mount()` every time someone wants a different question at a new moment, which puts you back in the redeploy-every-time loop Scouti is meant to replace.

### 4.1 The install script

Every integration starts with one `<script>` snippet. Copy it verbatim from the Dashboard (**Scout Team → Web Widget**)—it already has your project key filled in. (Working from the CLI/API instead? Mint a key with `POST /projects/{id}/keys` and drop it into the `data-project-key` slot yourself.)

```html
<script>!function(u,d){var w=window,e=document;if(w.scouti)return;var v,p=new Promise(r=>v=r);w.__sr=v;w.scouti=new Proxy({},{get:(_,m)=>function(){var a=w.__sa,g=[...arguments];return a?a[m].apply(a,g):p.then(()=>w.__sa[m].apply(w.__sa,g))}});var s=e.createElement("script");s.async=1;s.src=u;if(d)for(var k in d)s.setAttribute(k,d[k]);e.head.appendChild(s)}("https://scouti.chat/scouti-widget.umd.js",{"data-project-key":"pk_your_project_key"});</script>
```

It's a tiny self-contained bootstrapper. When it runs it:

- **Registers `window.scouti` synchronously** as a stand-in that **queues** any calls, then kicks off an async download of the real Widget bundle and replays the queue once it's ready. This is the queue mechanism that makes the API safe to call before the Widget has finished loading (see §4.2).
- **Carries your project key** as the `data-project-key` argument, which it forwards to the Widget so it loads the right project's config.
- **Mounts the Widget in an isolated iframe**, so it can't clash with your page's styles or scripts.

**Where to put it.** We recommend the `<head>` of a global layout/template, so it runs as early as possible on every page. But it doesn't have to be there—the end of `<body>`, a tag manager, or any lazy-loaded slot works too. There's only one rule:

> Before you call `scouti.mount(...)` / `scouti.setUser(...)` on a page, this snippet must have run on that page.

You don't need to wait for the Widget itself to finish loading—because of the queue, calls made early are held and replayed in order, so the bundle downloading a little later is fine. What you must avoid is calling `scouti.`* on a page where the snippet was never included.

**If your site sends a Content-Security-Policy (CSP).** A CSP on the embedding page can block the Widget before it ever runs, so Scouti's two origins have to be allowlisted:

- `script-src` (which also governs `script-src-elem`) must allow `https://scouti.chat`, where the loader bundle is served—otherwise the browser refuses to download it.
- `connect-src` must allow `https://*.scouti.chat`, since the Widget bootstraps and streams replies from Scouti's backend services, which run on subdomains of `scouti.chat`.

```
script-src  https://scouti.chat;
connect-src https://*.scouti.chat;
```

The bundle sits on the apex `scouti.chat` while the backend runs on `scouti.chat` subdomains, and a `*.scouti.chat` wildcard does **not** match the apex—so the two directives intentionally name different hosts. A page with no CSP needs none of this. Note this is the *embedding site* granting permission to load Scouti, and is the mirror image of Scouti's own **Allowed domains** (§5.2), where you grant Scouti permission to accept your site's calls—a working Widget needs both directions.

### 4.2 The JavaScript API

Everything is exposed on the global `window.scouti`. The install script registers a tiny queueing shim **immediately**, so you can call any method right away—calls made before the bundle finishes downloading are queued and replayed in order. The mutating methods return `scouti`, so they chain.


| Call                                     | What it does                                                                                                                                                |
| ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `scouti.mount(name: string)`             | Declares that the named Touchpoint moment occurred and enters it into the arbitration pipeline (see §4.3). Safe to call repeatedly and before load.         |
| `scouti.setUser({ identifier })`          | Binds the current end user so the Dashboard can cluster their sessions; the session follows this call live. Pass an empty / absent id (e.g. on logout) to drop back to anonymous. |
| `scouti.hide()`                          | Hides the Widget for now. Unlike a user *close*, this triggers no quiet period.                                                                             |
| `scouti.destroy()`                       | Tears the Widget instance down entirely.                                                                                                                    |


`**mount(name)` is the workhorse.** Place it at the moments you named in the Dashboard. Calling it doesn't guarantee a dialog—the pipeline in §4.3 has the final say.

`**setUser({ identifier })` binds an identity, and the session follows it live.** Binding is what lets the Dashboard stitch a person's separate visits—and separate devices—into a single **User**, so their feedback collects in one place instead of scattering across one-off anonymous sessions. Pass the real user id from your app right after login. For simplicity the Widget keeps **one session per browser** and mirrors whatever you pass, in sync: call it again with a *different* id to switch users—the Dashboard then shows that session under the new id—or with an empty / absent id (e.g. on logout) to reset to an anonymous session. Calls take effect immediately and persist across reloads.

### 4.3 How the Widget works

**One script, loaded once.** Integration is a single `<script>` tag carrying your `data-project-key`. The Widget renders inside an isolated iframe, so it never collides with your page's CSS or JavaScript. On load it performs a one-time **bootstrap**: it pulls your project's Topics, Touchpoint bindings, and any pending Outreach, then caches them locally—so every runtime decision afterward is instant and needs no server round-trip. That cache isn't permanent—it refreshes on a fixed interval (`BOOTSTRAP_EXPIRES_MINUTES`; see **Appendix A.1** for the value), so edits you make in the Dashboard roll out on their own without a redeploy. If your **Reactive Topic** is enabled, its sticky “tell us anything” button appears on its own; no further code is required.

**Touchpoints decouple code from configuration.** A call like `scouti.mount("post_checkout")` doesn't say “show this dialog”—it says *“this named moment just happened.”* Which Topic (if any) fires there, and under which Conditions, is bound in the Dashboard against that name. So a `mount(...)` is a **candidate**, not a command: the runtime may open it, open something else instead, or show nothing at all.

**From `mount()` to a visible Scout.** When one or more moments fire, the Widget runs them through a short arbitration pipeline before anything appears:

```
mount("…")  →  coalesce  →  guards  →  conditions  →  pick one Topic  →  at most one Scout
```

1. **Coalesce — batch the burst.** `mount()` calls that land close together (including the Reactive button's own auto-mount and a page-load *restore*) are collected into one short window (`COALESCE_WINDOW_MS`) and resolved *together*, so the user never sees a stack of competing dialogs. A Touchpoint can also be set to fire only a fraction of the time (a **`probability`**); that dice roll happens here, and a Touchpoint that misses drops out before anything else is considered.
2. **Guards — checked first, before any content is chosen.** Two layers protect the user's attention:
  - **Focus Guard** — while the user is actively engaged with a conversation (recently typed, recorded, scrolled history, etc.), a new auto-mount won't pull them out of it. Once they've been idle past `FOCUS_GUARD_MS`, mounts can take over again.
  - **Quiet Guard** — if the user explicitly *closes* the Widget, it stays quiet for a cool-off period (`QUIET_GUARD_MS`) before anything auto-opens again.
  If a Guard is holding the Widget, nothing new auto-opens; a page refresh simply restores whatever the user already had.
3. **Conditions & pacing — eligibility and frequency gates.** Two kinds of gate decide whether a surviving candidate may actually fire, both evaluated in the end user's own local context:
  - **Per-binding Conditions.** Each Topic↔Touchpoint binding may carry time-of-day, weekday, and a per-user **cooldown** (“don't ask the same person about this Topic again so soon”). *All* conditions on a binding must pass. (Full list and exact semantics: **Appendix A.2**.)
  - **Proactive pacing — a project-wide cooldown.** One global setting, **“Minimum time between proactive pop-ups”** in the Dashboard (**Scout Team**), throttles *every* proactive Topic together: once any proactive Scout has popped up, no new proactive pop-up appears until that gap elapses—so a project with many Touchpoints can't gang up on one user. It paces **proactive pop-ups only**: the Reactive “tell us anything” button, queued Outreach, and `force_show` moments are never held back by it. It's **off by default** (leave it at `0`/blank to disable), is capped at 7 days, and is measured from the last proactive pop-up (a page refresh on its own doesn't restart the clock).
4. **Pick one Topic.** A moment usually maps to a single Topic, but when several survive the gates above:
  - **Proactive beats Reactive** — an auto-pop-up Topic outranks the sticky “tell us anything” button; the button only wins when no proactive Topic is in the running.
  - **Stay on the current conversation** — if the user is already in one of the surviving Topics, the Widget keeps them there instead of switching.
  - **Otherwise weighted-random** — one is drawn at random, biased by the **weight** you set per binding in the Dashboard (higher weight = chosen more often). It's an even split only when the weights are equal.

**Two things can jump this queue:**

- **Outreach pre-empts a Topic.** If you've queued a direct **Outreach** for this user (a thread follow-up or re-engagement, §2 Step 5), it takes the slot ahead of whatever Topic the moment would otherwise have triggered.
- **`force_show` overrides the gates.** A Touchpoint marked `force_show` in the Dashboard is for moments that must not be missed (post-purchase confirmation, error-recovery, a required notice). It **skips every gate above—both Guards and all Conditions (including the per-user cooldown and the project-wide proactive pacing)—and pre-empts Outreach**—so it appears even if the user just closed the Widget or is mid-conversation (it swaps its content in place). It still respects `probability`, and leaning on it too often erodes users' trust, so reserve it for genuinely critical moments.

The net effect: think of `mount()` as *raising your hand*. The runtime referees every raised hand so the user ends up with **at most one** Scout, shown at a sensible time, never spammed.

**When a proactive Scout doesn't appear** (and you expected one), it's almost always one of the gates above doing its job: a **Guard** is active (the user recently *closed* the Widget, or is still mid-conversation), a per-binding **Condition** or **cooldown** hasn't elapsed yet, the project-wide **proactive pacing** window is still open from an earlier pop-up, or the Touchpoint's **`probability`** rolled it out this time. The Reactive “tell us anything” button is exempt from all of these, so it stays reachable even while proactive pop-ups are paused. For a moment that must always surface no matter what, mark its Touchpoint `force_show`.

### 4.4 Examples

The Dashboard (**Scout Team → Web Widget**) generates ready-to-paste snippets for HTML, React, Next.js, and Vue with your real project key filled in. The HTML versions below show the shape of each step. (Drop in the install script from §4.1 first.)

**Identify the user (after login).**

```html
<script>
  // Cluster this user's sessions under your own user id.
  window.scouti.setUser({ identifier: currentUser.id });
</script>
```

**Fire a Touchpoint after page load.**

```html
<script>
  window.addEventListener("load", () => {
    window.scouti.mount("home_landing");
  });
</script>
```

**Fire a Touchpoint on a user action.**

```html
<button id="feedbackButton" type="button">Share feedback</button>
<script>
  document
    .querySelector("#feedbackButton")
    ?.addEventListener("click", () => window.scouti.mount("feedback_button"));
</script>
```

**Fire a Touchpoint from your own logic.** A `mount` can hang off any condition and timing you like—here it waits for an engaged paid user, then holds off another moment so the Scout doesn't interrupt:

```html
<script>
  const DWELL_MS = 10_000; // user has settled in
  const DELAY_MS = 5_000;  // don't pounce immediately
  if (currentUser.plan === "paid") {
    setTimeout(() => {
      setTimeout(() => window.scouti.mount("engaged_paid_user"), DELAY_MS);
    }, DWELL_MS);
  }
</script>
```

Whatever moment you `mount`, the Dashboard decides which Topic fires there, and the pipeline in §4.3 decides whether—and how—it actually surfaces.

---

## 5. Key concepts & settings

A reference for the concepts and settings that come up most once you're past the first integration. Each is managed—or surfaced—in the Dashboard, with no code change or redeploy.

### 5.1 Project key

The **project key** (`pk_…`) is the public id that ties an embed to your project—it's the `data-project-key` in the install script (§4.1). It's a *publishable* key, meant to sit in client-side HTML, so there's nothing secret to leak: abuse is held off by the allowed-domains check below (plus a one-time human check and your credit balance), not by keeping the key hidden. Every project starts with one default key; you can also add or remove extra named keys—say, one per site or environment—from your project's settings (or via the API).

### 5.2 Allowed domains

For the Web Widget, Scouti only answers requests coming from the **domains you allowlist** (the visitor's browser reports its page origin and the server checks it). List every site that embeds the Widget under **Project Settings** (or via the API); `*.example.com` covers subdomains and `localhost` is allowed while you develop. **An empty allowlist blocks the Widget**—it's the most common reason a freshly-installed Widget never appears, so fill it in before testing on a live domain. Direct Links don't need an entry here: they're served from Scouti's own hosted page.

### 5.3 Documentation (your product context)

This is the plain-language description of your product that Scouti feeds into **every** conversation, so the Scout sounds like it already knows your app instead of asking from scratch. You write it in the project's **Playbook** (or via the API): what the product does, who uses it, the core flow, the words your users use, and the known rough edges. It's the single highest-leverage thing you can set—thin context yields thin, generic interviews—so the Dashboard nudges you whenever it's too short. You don't have to draft it by hand: the **“Ask your AI”** button hands you a ready-made prompt to paste into your coding agent (Claude Code, Codex, Cursor…), which reads your repo and writes a first version for you to refine.

### 5.4 Tag Tree

The **Tag Tree** is your feedback taxonomy. After each chat the Analyst attaches short tags from this tree to every **Point** (one atomic thing a user said), so you can later filter and count—“show me every *pricing* complaint.” The tree is hierarchical (up to three levels, broad at the top, finer below); the Analyst walks each branch top-down and attaches every tag that still fits, and a single Point can carry several. A sensible default tree ships out of the box so tagging works immediately—reshape it around your own product for sharper signal. Edit it in the **Playbook** (or via the API).

### 5.5 Attention List

The **Attention List** is a short set of high-value signals you want caught the instant they appear—churn risk, security worries, pricing objections, or a launch-specific theme. When a conversation matches one, the Analyst stamps it with an **attention flag** so it rises to the top instead of getting buried. Keep it to a handful (up to five) of genuinely important items: it's a spotlight, not a second Tag Tree. Also set in the **Playbook** (or via the API).

### 5.6 Morphing — match the Widget to your brand

**Morphing** lets the Widget wear your brand instead of Scouti's defaults: pick an **accent color** (it drives the voice and send buttons, the input focus ring, and the “go away” pill) and upload a **bot avatar**, with a live preview as you edit and a one-click reset. It's a **paid-plan feature**—on the free plan the controls are visible but locked—and it's purely cosmetic, so it never changes which questions fire or when. Set it in **Project Settings**.

### 5.7 Daily digest

A short morning email recapping yesterday's replies, conversation quality, and notable Points. It runs on **two independent switches**, and the email only arrives when both are on:

- **Per project** (Project Settings) — whether that project feeds the digest.
- **Per person** (your Account) — whether *you* receive digest emails at all, across every org you belong to.

Both default to on, so if the emails stop completely, check your Account toggle first; if just one project is missing from them, check that project's toggle.

### 5.8 Direct Links

A **Direct Link** is a shareable URL—plus a matching QR code—that opens a single Topic in a full-screen hosted chat, with **no widget and no code** on your side. It's how you reach people *outside* your product: drop it in an email, a social post, a support reply, a Discord/community message, or a printed QR. You pick the Topic, optionally put it on a branded subdomain, and share it; whoever opens it lands straight in that conversation, and their replies flow into the same Dashboard as widget feedback. Once a Topic's direct link is enabled, it's reachable at `https://scouti.chat/t/<topic-id>`. Because Scouti hosts the page, Direct Links ignore the allowed-domains list (§5.2) and need no project key on any page of yours—handy for beta cohorts, power users, or one-off campaigns.

### 5.9 Points, sentiment & attention flags

When a conversation ends, the Analyst doesn't just file the transcript—it breaks the chat into **Points**. A Point is one atomic thing the user expressed (a single complaint, request, or bit of praise), distilled to a one-line gist plus a short summary. Each Point carries three things you can act on:

- **Sentiment** — positive, neutral, or negative, so you can read the mood at a glance and filter by it.
- **Tags** — pulled from your Tag Tree (§5.4).
- **An attention flag** — raised when the Point matches your Attention List (§5.5).

Splitting chats into Points—rather than one blob per conversation—is what makes feedback *countable*: a single five-minute chat can yield several Points, each tagged, weighted, and filterable on its own. You browse and slice them in the **Conversations** and **Users** views.

### 5.10 Team & roles

Your account is organized as **Organizations → Projects**: an Organization owns billing and teammates, and its Projects sit beneath it. Invite people from the Organization's member settings. There are two roles—**owners**, who manage members, billing, and settings, and **members**, who work inside the Projects—and everyone in an org can collaborate across that org's Projects.

### 5.11 Billing & credits

*Placeholder — pricing and credit mechanics are still being finalized.* This section will cover how usage is metered, how credits are spent per conversation (less for thin “hi-and-gone” chats, more for rich ones), and what happens when a project runs low. We'll fill it in once the model is locked.

---

## Appendix A — Detailed rules

Reference details for specific mechanisms, kept out of the main flow so the walkthrough stays readable.

### A.1 Runtime timing reference

The named constants §4.3 refers to, and their concrete values:


| Constant                       | Value          | What it controls                                                                                                                                                                                                                         |
| ------------------------------ | -------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `COALESCE_WINDOW_MS`           | **500 ms**     | The coalesce window: how long `mount` / restore calls are batched before the pipeline resolves them into a single result.                                                                                                                |
| `BOOTSTRAP_EXPIRES_MINUTES`    | **120 minutes** | The config cache lifetime: how long the locally cached project config (Topics, Touchpoint bindings, Outreach notifications) is reused before the Widget re-fetches it. Server-set; the client honors it via the bootstrap's `expire_at`. |
| `FOCUS_GUARD_MS`               | **5 minutes**  | Focus Guard: how long an actively-used conversation is protected from being replaced by a new `mount`; resets on each qualifying interaction.                                                                                            |
| `QUIET_GUARD_MS`               | **15 minutes** | Quiet Guard: how long the Widget stays silent after the user explicitly closes it.                                                                                                                                                       |
| `CONVERSATION_REUSE_WINDOW_MS` | **5 minutes**  | Conversation reuse window: after a Topic was last shown, re-firing it within this window reopens the same conversation instead of starting a fresh one.                                                                                  |


### A.2 Touchpoint Condition types

A binding ties a Topic to a Touchpoint and can carry zero or more **Conditions**. All conditions on a binding must pass (AND logic) for that Topic to be eligible to fire. Supported types:

- **time_range** — a wall-clock window in the end user’s timezone, with `start` / `end` as `HH:mm`. Omit `start` for “no lower bound” or `end` for “no upper bound.” Multiple `time_range`s can be combined for split windows.
- **weekdays** — restricts firing to specific days of the week (`0 = Sunday … 6 = Saturday`).
- **cooldown** — a per-user “don’t ask again so soon” gate, measured in minutes against the user’s last **reply** to that Topic:
  - `-1` — ask only once; any prior reply blocks it forever (per browser slot).
  - `0` — no cooldown; re-fires after every reply.
  - `>0` — minimum minutes between the user’s last reply and the next allowed fire.
  Notes: the “last reply” timestamp is tracked **per Topic, per user** (not per binding); openings the user merely saw but never replied to don’t count; and `force_show` on a Touchpoint bypasses cooldown entirely. (Implementation specifics are in `doc/impl/touchpoint.md`.)

---

