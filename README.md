# sitnl.nl

Website of SAP Inside Track Netherlands (#sitNL). Static site, no build step.
The 2027 edition: Friday 22 January 2027 (pre-event, codejams) and Saturday 23 January
2027 (sitNL), SAP Nederland, Amerikastraat 10, 5232 BE 's-Hertogenbosch.

Live since 20 November 2026. Design system and all design decisions: [`design.md`](design.md).
Open items: `design/todo.md` in the claude.ai project "SitNL".

## Where things are

| Path | What |
|---|---|
| `index.html` | The whole home page: HTML, CSS (tokens at the top of `<style>`) and a small inline script (reveals, quiet bar, magnetic button). No framework. |
| `codeofconduct.html` | Code of conduct, same theme (its own copy of the tokens). |
| `design.md` | Design system: palette, type, layout, motion, components, copy rules. Source of truth for any change. |
| `assets/` | `sitnl-2027-logo.svg` (vector, use this), `.png`/`.webp` fallbacks, `fonts/72-*.woff2`, `sitnl-2027.ics` (calendar file). |
| `_redirects` | Netlify redirects for old 2025 URLs (`location.html`, `sponsors.html`). |
| `agenda/main.go` | Data script: pulls agenda, speakers, speaker photos and sponsors from the apps (see below). |
| `bruno/` | Bruno API requests for the planning tool (login, proposals, speakers, photos). |
| `artifacts/` | Logos, t-shirt artwork and the TechEd 2026 visual kit (`teched-2026-visuals/`). Not served on purpose; check SAP brand guidelines before using keynote assets. |
| `archive/` | Previous editions (`sitnl2018`, `sitnl2019`, `sitnl2024`, `sitnl2025`), each self-contained. |

## Preview locally

```bash
cd ~/Library/CloudStorage/OneDrive-SAPSE/Work/SitNL
python3 -m http.server 8000
# open http://localhost:8000
```

The hero video is a local poster (`images/video-poster.webp/.jpg`, taken from the
YouTube thumbnail) that swaps in the YouTube player on click, so nothing from YouTube
loads before someone presses play. Replace the poster if the video changes.

Open via a local server, not by double-clicking the file: the YouTube video and the
map are iframes and need http.

## Deploy

Netlify deploys `main` of `github.com/jpenninkhof/sitnl` to sitnl.nl automatically
(repo root is the publish folder, no build command).

```bash
git add -A
git commit -m "..."
git push
```

To check a change before it goes live, push a branch (`git switch -c my-change`,
`git push -u origin my-change`) and use the Netlify deploy preview, if branch deploys
are enabled for the site.

## Data: event settings, agenda, speakers, sponsors

`agenda/main.go` logs in to Andreas Kunz's planning tool
(`sitnl.cfapps.eu12.hana.ondemand.com`) and the sponsor app
(`presales-nl---build-apps-dev-sitnl-srv.cfapps.eu12.hana.ondemand.com`), then writes:

- `js/event.js` (`var eventConfigJson = {...}`): the public event settings from
  `/api/events/sitnl`, such as the call for papers deadline. The planning tool's event
  record also holds mail server settings and e-mail templates, so the script copies only
  the fields listed in `eventPublicFields` in `main.go`. Add a field there to publish it.
- `js/agenda.js` (`var proposalLineupJson = [...]`)
- `js/sponsors.js` (`var sponsorLineupJson = [...]`)
- speaker photos into `images/speakers/`

Run it from the `agenda/` folder, with the planning-tool credentials in environment
variables (never commit them):

```bash
cd agenda
APP_USERNAME=... APP_PASSWORD=... go run .
```

Important for 2027:
- The paths are relative to `agenda/`, so the output lands in `js/` and `images/` in the
  repo root. The 2025 versions of those folders now live in `archive/sitnl2025/`. Do not
  write 2027 data into the archive, and do not let 2025 data leak into the root.
- The 2027 `index.html` does not read `js/agenda.js` or `js/sponsors.js` yet. The agenda
  and sponsors sections are built but hidden (`hidden` attribute). The 2025 rendering code
  (Vue 2, `archive/sitnl2025/js/main.js`) is the reference for the data fields. Wiring the
  data in is an open item.

## Common changes

- **Times, dates, deadlines:** search `index.html` for "to be announced" (event days,
  venue card, call-for-papers key dates). The call for papers deadline comes from the
  planning tool: an element with `data-event="cfpDeadlineAsString"` gets its text from
  `js/event.js` when that file exists. Keep the HTML text in step with the tool, because
  it is what visitors see when `js/event.js` isn't deployed.
- **Open registration:** add `register` (and later `agenda`) to both navs (the header nav
  and the quiet bar `#bar`), point the codejam/sitNL buttons and the `#register` band at
  the real registration pages, change the hero actions per `design.md` §11.
- **Show agenda or sponsors:** remove `hidden` from `#agenda` / `#sponsors`, add the nav
  item in both navs. Then check that visible bands still alternate between `band` (violet)
  and `band band-top` (navy), ending on a violet `.register` band above the navy footer.
  Current order: days (navy), expect (violet), call for papers (navy), sponsors (violet),
  getting there (navy), register (violet). Showing `#agenda` means flipping the bands
  after it.
- **Close the call for papers:** remove the "call for papers" nav item (both navs) and
  the "Submit a talk" buttons (hero and section), or turn the section into a
  "speakers announced" note.
- **New colours, type or motion:** change `design.md` first, then the tokens in
  `index.html` and `codeofconduct.html` (both have the `:root` block).

## Gotchas

- **OneDrive:** keep this folder on "Always keep on this device". Cloud-only files make
  git fail with "Resource deadlock avoided", and OneDrive has rolled edited files back to
  older versions at least once. Check `git diff` before committing.
- **Old 2025 pages** are reachable under `/archive/sitnl2025/`, so don't delete them
  without updating `_redirects`.
- **Copy rules:** English, sentence case, no em dashes, no exclamation marks, no
  invented facts or testimonials (`design.md` §9).
- **The claude.ai preview** of this page blocks YouTube and map iframes and downloads;
  judge those only on a local server or the Netlify preview.
