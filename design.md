# sitNL 2027 design system

Source of truth for every later build step. Direction: **Bouwlamp** (build lamp),
aligned with the TechEd 2026 keynote analysis in the claude.ai project
(`design/teched-2026-visuele-stijl.md`, `design/sitnl-teched-tokens.css`, `design/motifs/`).
If code and this file disagree, this file wins until it is updated on purpose.
Where this file deviates from the TechEd analysis, it says so and why.
Live on sitnl.nl since 20 November 2026 (repo `jpenninkhof/sitnl`, deployed by Netlify).

---

## 1. Brand & thesis

**Brand:** SAP Inside Track Netherlands 2027. Always written **sitNL** (small s,
capital NL), also at the start of a sentence and in headings; `#sitNL` as a hashtag.
Never "SitNL" or "SITNL". (The logo artwork spells "SIT NL" in capitals; that is the
artwork, not the word mark.) Community-run, free, in person.

**Thesis line:** Build things that matter.

The line is SAP TechEd 2026's theme, used on purpose. sitNL is where that theme lands
locally: TechEd tells you what to build, sitNL is the room where Dutch builders show
each other what they built. Visually we take TechEd's system almost whole (the
navy-to-violet field, 72 lettering, the staggered headline with an accent full stop,
the pill trio) and make one swap: TechEd's teal accent becomes the **orange of the
sitNL lions**. That orange (`#fca400`) sits right next to SAP gold (`#F0AB00`) in the
deck's own palette, so it reads as TechEd family at a glance and as "ours" on a second
look.

**Personality:** a workshop at night with one work lamp on. A deep violet room, one
warm light where the action is. Direct, a bit nerdy, never corporate.

## 2. Audience & the one job

**Audience:** SAP developers, architects, BI/data people and consultants in the
Netherlands and Belgium. Many have attended before. Secondary: partner companies
deciding whether to sponsor, and first-time speakers.

**The one job:** get a returning community member to say "I'm going" within ten
seconds: dates, place, what it is, how to join. Until registration opens the job is
"watch the video (in the hero), save the date, submit a talk"; once it opens, every
page funnels to Register.

**Secondary jobs:** sponsors find the package and a contact; attendees find how to
get to the venue.

## 3. Palette (OKLCH, named roles)

Values are the TechEd deck colours converted to OKLCH, plus the logo orange. The violet
field is the *ground*, not an accent: it is how TechEd looks. Exactly one saturated
accent sits on top of it: orange. No pure black or white anywhere.

| Role | Token | Value | Source | Use |
|---|---|---|---|---|
| Ground top | `--bg-top` | `oklch(0.22 0.10 262)` | TechEd navy `#00144A` | Top of the hero field, footer |
| Ground | `--bg` | `oklch(0.27 0.15 277)` | TechEd mid `#1C0C6E` | Default page canvas, behind body text |
| Ground bottom | `--bg-bottom` | `oklch(0.53 0.27 280)` | TechEd violet `#5D36FF` | Bottom of the hero field only |
| Surface | `--surface` | `oklch(0.37 0.22 275)` | TechEd deep violet `#2C12AD` | Cards (sponsor tiers, venue, event days) |
| Shape | `--shape` | `oklch(0.47 0.245 279)` | TechEd numeral `#4D2FDF` | Giant background numeral, soft background shapes |
| Line | `--line` | `oklch(0.59 0.235 285)` | TechEd violet light `#7858FF` | Hairlines, card borders, long pill of the trio |
| Lilac | `--lilac` | `oklch(0.74 0.15 298)` | TechEd lilac `#B894FF` | Short pill of the trio, quiet chips |
| Ink | `--ink` | `oklch(0.98 0.01 285)` | white, tinted | Primary text |
| Ink soft | `--ink-soft` | `oklch(0.83 0.10 238)` | TechEd sky `#89D1FF` | Dates, labels under stats, captions |
| Accent | `--accent` | `oklch(0.79 0.17 71)` | logo orange `#fca400` (≈ SAP gold) | CTA fill, the full stop, the lead dot of the trio, focus ring |
| On accent | `--on-accent` | `oklch(0.18 0.09 281)` | TechEd navy ink `#0E0537` | Text on orange and lilac |
| Plate | `--plate` | `oklch(0.97 0.02 85)` | cream, from the logo family | Reserved; no longer behind the logo (see §8). Available for a future light surface |

Rules:
- **Hero field:** `linear-gradient(in oklch 180deg, var(--bg-top), var(--bg) 45%, var(--bg-bottom))`,
  exactly the deck's background. It is the only gradient on the site. Bands below the
  hero are flat `--bg` or `--bg-top`.
- Orange covers at most ~5% of any viewport. If two orange things compete in one
  screen, one goes back to ink.
- Body text sits on `--bg`, `--bg-top` or `--surface`, never on `--bg-bottom` (contrast
  drops too far for small text). `--ink-soft` only at 17px and up.
- **Not used** (deviation from the TechEd analysis): teal, pink, SAP blue and pill blue.
  They are TechEd's accents; using them next to orange would give four accents.

## 4. Type

**72**, SAP's public web font (`assets/fonts/72-Regular|Bold|Black.woff2`, carried over
from 2025). It is the web cut of 72 Brand, the deck's typeface, so the lettering is
TechEd's own. Two weights do the work, as in the deck (Book/Medium):

| Role | Face | Weight | Notes |
|---|---|---|---|
| Display | 72 | 700 (Bold) | Thesis and section heads. Leading 1.05, tracking −0.015em. The deck uses Medium; Bold is the nearest cut we have, so we size slightly smaller rather than heavier |
| Body | 72 | 400 (Regular) | 1.0625rem (17px), line-height 1.6, max 66ch |
| Utility | 72 | 400, `tabular-nums` | Dates, times, rooms, nav, labels, in `--ink-soft` |
| Numeral | 72 | 900 (Black) | Only the giant background date numeral |

**Scale** (major third, 1.25, from 17px): 14 · 17 · 21 · 27 · 33 · 42 · 52 · 66 · 82 px.
Thesis `clamp(3rem, 7vw, 5.125rem)` (TechEd's `--tx-display`, capped at 82px).
Section heads 33–42. Giant numeral `clamp(18rem, 45vw, 40rem)`.

Do:
- Set the thesis on two lines, **staggered** the TechEd way: "Build things" left,
  "that matter." indented by roughly a third of the first line's width.
- Short lines, 2–4 words. Sentence case everywhere.
- Dates in numerals when tight (`22–23.01.2027`), in words in sentences.

Don't:
- No all-caps, no italics, no single coloured word in a headline (the orange full
  stop is the only coloured glyph).
- No second typeface, no monospace labels, no Inter.
- No gradient text (deviation: the TechEd tokens clip a gradient into the giant
  numeral; we use solid `--shape` instead).
- No three-line "last line in accent" claims: that is TechEd's pattern with teal and
  would put an orange headline next to orange buttons.

## 5. Layout & spacing

- **Grid:** 12 columns, max width 72rem, 1.5rem gutters. Text on columns 1–7; shapes,
  the numeral and imagery on 8–12, allowed to bleed off the right edge.
- **Alignment:** left-aligned (deviation: the deck centres titles; centred long copy reads
  badly on the web). The staggered thesis provides the TechEd rhythm instead.
- **Spacing (rem):** 0.25 · 0.5 · 0.75 · 1 · 1.5 · 2 · 3 · 4 · 6 · 8. Section padding
  6rem desktop, 4rem mobile.
- **Rhythm:** hero 100vh on the gradient field (logo first, thesis below); then flat bands alternating `--bg` and
  `--bg-top`; the pill trio marks the start of each major section.
- **Radius:** pill `999px` for buttons, chips and the trio; `20px` for cards (TechEd's
  card radius); `0` for bands and images. Three radii, each meaning one thing.
- **Mobile:** single column from 48rem; the staggered indent shrinks to ~1.5em; nav
  becomes a plain disclosure.

## 6. Signature element

**The lamp dot**, three ways, all orange, all from TechEd's vocabulary:
1. **The full stop.** "Build things that matter**.**", the period in `--accent`, as TechEd
   does with teal. The one coloured glyph on the page.
2. **The pill trio, sitNL edition.** TechEd's dot + long pill + short pill, recoloured:
   orange dot (`--accent`), long pill in `--line`, short pill in `--lilac`. Used as the
   section opener and at the bottom right of the hero. Max once per section.
3. **The giant numeral.** "23" in 72 Black, `--shape` on the violet field, behind the
   date block, the deck's "720pt number" device.

Background shapes, structural only: one oversized, softly rounded four-point star (from
`sparkle.svg`) or the half ring (`arc-half-ring`), in `--shape` at low contrast, partly
off-canvas behind the hero. Never animated in a loop, never carrying content.

Not used: the multicolour sparkle cluster and `rings.svg` in their original colours
(teal/pink/blue). If a later step wants them, recolour to `--accent` + `--lilac` + `--line`.

## 7. Motion

**Structural:**
- Hero load: the full-width logo fades in first (opacity, 8px rise, 600ms,
  `cubic-bezier(.2,.8,.2,1)`); then the brand line and thesis lines (360ms, from 280ms),
  second line 90ms later; then the full stop and the trio's dot "switch on" (scale
  0.6 → 1, 240ms). One orchestrated moment.
- Disclosures (mobile nav, agenda filters, FAQ): height/opacity 200ms.
- Video: plain inline embed, no facade; it fades in with the hero (from 600ms).

**Polish:** button fill lightens on hover (accent L 0.79 → 0.84), 120ms; link underline
offset 2px → 4px. Added in the polish pass (step 5):
- **Film grain:** static, faint (7% opacity) noise over the whole page. Not animated.
- **Magnetic hero button:** the primary hero button ("Add to calendar") drifts at most 6px toward the pointer.
  Fine pointers only; off under reduced motion and on touch.
- **Quiet bar:** once the hero has scrolled away, a solid `--bg-top` bar with a hairline
  slides in (200ms) with the small logo and nav (links only on phones). No blur.
- **Nav links:** orange underline (`--accent`) fades in on hover, offset 6px → 4px.
- **Deliberately not added:** custom cursor (the one effect removed), Lenis inertia
  scrolling, GSAP scroll-scrub, hero parallax, hover states on non-interactive rows.

**Motion style: Calm** (chosen in build step 2). On top of the hero moment, each
section's intro and its first content block fade in once as they enter the viewport
(opacity 0 → 1, 12px rise, 500ms, same easing). Nothing else moves on scroll.

**Restraint rule:** no reveal on every element (only section intros and their first
block), no parallax, no inertia scrolling, no pinned sections, no looping animation
(deviation: the TechEd tokens include `.tx-float`; we drop it). Under
`prefers-reduced-motion: reduce` everything is static and the dots render on.

## 8. Components

- **Page spine:** hero (with video) → the two days → what to expect → call for papers →
  agenda (hidden until live) → sponsors (logo wall, links to `sponsors.html`) → getting there → register / save-the-date →
  footer. Bands alternate `--bg` / `--bg-top`.
- **Header:** small logo on a transparent background (left, no plate), nav right:
  `register`, `agenda`, `call for papers`, `sponsoring`, `getting there`. Items that
  aren't live are left out. Now showing: `call for papers`, `sponsoring`, `getting there`. Nav links never wrap internally; on phones the nav drops below the logo.
- **Call for papers:** section intro with the trio, then two columns: what we're
  looking for (topics, session formats as outline chips in `--line`, primary "Submit a
  talk" to Andreas Kunz's planning tool, secondary mailto) and a `--surface` card with
  key dates (submissions open, deadline, speakers notified, event). Remove the nav item
  and the section's primary button when the call closes.
- **Hero:** on the gradient field, full viewport height. First the **sitNL logo at full
  content width** (transparent, directly on the field), then the brand line in `--ink-soft` ("SAP Inside Track
  Netherlands 2027", sentence case, not an eyebrow), staggered thesis, two-sentence lede,
  actions (primary "Add to calendar", secondary "Submit a talk"), the video on the
  right, pill trio bottom right, background shape behind.
- **Date block:** "22 and 23 January 2027, 's-Hertogenbosch" with the giant "23"
  numeral behind; add-to-calendar links beneath.
- **Video:** the 2027 call-to-join video (`https://youtu.be/rvsy4NfkXWg`) embedded
  inline in the hero, right of the thesis (owner's choice): a 16:9
  `youtube-nocookie.com` iframe with a one-line caption in `--ink-soft`. Two columns
  from 64rem (thesis scales to `clamp(3rem, 4.6vw, 4.25rem)` there); stacked below the
  actions on smaller screens. No separate video section.
- **Primary button:** pill, `--accent` fill, `--on-accent` text, 72 Bold. (Deviation:
  TechEd's primary is a white pill with violet text; ours is orange because orange is the
  action colour.) Labels are verbs: "Register", "Become a sponsor", "Open route".
- **Secondary button:** pill, 1.5px `--ink` outline, ink text.
- **Event day pair:** two `--surface` cards, 20px radius, 1px `--line` border, solid (no
  glass, no blur): "22 January, pre-event (codejams)" and "23 January, sitNL".
- **Chips (agenda filters, later):** pill chips; tracks in `--lilac` with `--on-accent`
  text, rooms as outline chips in `--line`. No teal/pink chips.
- **Agenda (later, hidden at launch):** data from Andreas Kunz's planning tool
  (`sitnl.cfapps.eu12.hana.ondemand.com`); time-ordered list, times in tabular utility type.
- **Sponsors:** live since 20 November 2026: home-page section with a logo wall (each logo
  links to its entry on `sponsors.html`), "Become a sponsor" and "Meet our sponsors". Logos in their own colours (owner's choice, Nov 2026), grouped by tier when there is more than one. The site
  stays static and independent of the sponsor registration app: a script periodically
  pulls the sponsors from that BTP app and publishes them as a static `js/sponsors.js`
  (`var sponsorLineupJson = [...]`, fields `name`, `logo`, `website`, `sponsortype`,
  `description`, social handles). The 2027 page reads that same format, so the export
  script keeps working unchanged. Logos in `assets/sponsors/` (`sprite.svg` + `sap.png`, from 2025). Page: `sponsors.html`; logos keep their own colours on a transparent background; sprite logos that use colour variables get the dark-background values of the 2025 site (white replaced by `--ink`).
- **Getting there:** venue card, an embedded OpenStreetMap map (iframe, marker on
  Amerikastraat 10 at 51.71252, 5.33944, 20px radius, `--line` border, "View larger map"
  link), then short lists for car, parking, public transport.
- **Footer:** "Reach out to us" (`info@sitnl.nl`), "Code of conduct",
  copyright. On `--bg-top`.
- **Focus:** 3px `--accent` outline, 3px offset, on every interactive element.

## 9. Copy rules

- English throughout. Sentence case. No exclamation marks in UI copy.
- Say what happens: "Register for sitNL", "Register for the codejam", "Download the
  sponsor deck". Never a bare "Click here" or "Learn more".
- Facts first: date, place, free, in person. Then the why.
- "We" is the organising team, "you" the visitor. No marketing verbs (unlock, leverage,
  seamless, empower).
- SAP product names written correctly. Code of conduct (`codeofconduct.html`, 2027
  theme): 2025 text with light edits (UI5 topics → SAP technology topics, refund clause
  dropped because the event is free, typos and gendered wording fixed).

## 10. Assets plan

TechEd source assets live in the SitNL folder under `artifacts/teched-2026-visuals/`
and in the project under `design/`. They come from SAP's keynote deck: check the SAP
brand guidelines before publishing any of them on sitnl.nl.

| Asset | Look | Source | Status |
|---|---|---|---|
| sitNL 2027 logo | Heraldic lions, five flat colours: brown `#7f5200`, dark orange `#b27400`, mid orange `#e59502`, orange `#fca400`, cream `#ffe4b3`. Shield interior is transparent. Used on a transparent background directly on the navy/violet field (owner's choice, Nov 2026): the brown lions read with less contrast than on cream, which is accepted. Full content width in the hero, 6.5rem in the header | `assets/sitnl-2027-logo.svg` (vector trace, 23 KB); raster `sitnl-2027-logo.png/.webp`; project copy `design/logo/sitnl-2027-logo.svg` | In place |
| Font 72 | Regular, Bold, Black woff2 | `assets/fonts/` (from 2025 site) | In place |
| Pill trio | Orange dot, violet-light long pill, lilac short pill | Recoloured from `design/motifs/pill-trio.svg`, inline SVG or CSS | Later step |
| Background shape | Oversized rounded four-point star or half ring, `--shape` at low contrast | `motifs/sparkle.svg` (currentColor) or `backgrounds/arc-half-ring.webp` | Later step |
| Favicon / app icon | Orange dot on navy square | Draw as SVG (the logo SVG is used for now) | Later step |
| Social share image | 1200×630 gradient field, staggered thesis, logo on transparent background | Generate from the hero | Later step |
| Call-to-join video | YouTube embed in the hero | `rvsy4NfkXWg` via youtube-nocookie | Live |
| Isometric illustrations | Lilac/violet isometric scenes (people, code, data network) | `illustrations/iso-*.webp`, at most one per teaser, never as decoration in the hero | Optional, after brand check |
| Kasimir loops | Short animated loops of SAP's Kasimir character | `kasimir/` (+ `KASIMIR.md` on usage and rights), project file `SAP_Kasimir_Idle_Final.mp4` | Optional, only if the rights in `KASIMIR.md` allow it; never more than one on the page |
| Event photos | Last year's crowd shots, duotoned navy shadows / warm highlights | 2025 sitNL photo album | Optional |
| Speaker photos | Square, same duotone | Andreas Kunz's planning tool | When agenda opens |
| Sponsor logos | Own colours on the navy field (owner's choice, Nov 2026) | 2025 `images/sponsors/sprite.svg` + static `js/sponsors.js` export | Carry over |
| Getting-there icons | Line icons in ink, or one isometric illustration | Hand-drawn SVG sprite or `iso-*` | Later step |

Not used: SAP product screenshots, the Joule diamond, third-party logos, the deck's
videos, the `reference/` slides (inspiration only).

## 11. Conversion essentials

- Above the fold: the full-width logo (owner's choice). On wide screens the dates, city,
  "free" and the action now sit just below the fold; on phones they stay in the first
  screen. Revisit if the hero action needs to be visible immediately.
- Before registration opens the hero actions are "Add to calendar" and "Submit a talk"
  (the video plays in the hero itself); after, "Register" to the SAP Community event
  page and the planning tool.
- Separate registration for the 22 January pre-event codejams (limited seats) and the
  23 January main event.
- "Add to calendar" (.ics plus Google/Outlook links) next to the dates.
- Sponsor path: one nav item, one home-page block, a mailto with subject
  `[sitnl 2027] Sponsoring`.
- Trust signals: "organised by the community", edition count, sponsor logos.
- Performance: no framework for static pages; home page under 300 KB before the video;
  72 preloaded (Regular + Bold), Black loaded only where the numeral appears.
- SEO/OG: title "SAP Inside Track Netherlands 2027", description, share image, `Event`
  structured data with both dates and the venue.

## 12. Anti-patterns

- Teal, pink or SAP blue as accents next to orange.
- Glass cards, backdrop blur, glow effects (deviation from `.tx-card`, which uses blur).
- Gradient text, including the giant numeral.
- Floating/looping animation (`.tx-float`).
- Copying SAP artwork wholesale or impersonating SAP marketing: we use the system,
  not their slides, product UI or Joule branding.
- Cards for everything; boxes inside boxes.
- Tracked-out caps eyebrows, 01/02/03 markers for non-sequences, meta strings joined
  with middle dots.
- Pure `#000`/`#fff`, untinted greys.
- Orange anywhere other than actions, the lamp dot (stop, trio dot, focus) and the logo.
- Nav items for content that isn't live.

## 13. References

- TechEd 2026 analysis in the claude.ai project: `design/teched-2026-visuele-stijl.md`,
  `design/sitnl-teched-tokens.css`, `design/motifs/*.svg`, `design/demo.html`.
- Reference slides: `artifacts/teched-2026-visuals/reference/` (e.g. slide 02,
  "Build things that matter.", the staggered headline with the accent full stop).
- SAP TechEd 2026 key visual, `sap.com/events/teched.html`.
- sitNL 2027 logo, and the 2025 site archived at `archive/sitnl2025/`.
- Bringhurst, *The Elements of Typographic Style*.
- Registration and planning: `https://sitnl.cfapps.eu12.hana.ondemand.com/`.
