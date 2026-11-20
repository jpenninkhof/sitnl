# SitNL 2027: visuele elementen uit de TechEd 2026 keynote

Bron: `2026-TechEd-Keynote_Part1_FINAL.pptx` (87 slides, thema "Build things that matter"). Doel: een herkenbare, actuele look voor de SitNL-site (23 januari 2027) die aansluit op de TechEd-stijl, zonder SAP-productschermen of logo's van derden over te nemen.

## 1. De stijl in het kort
- **Donker, diep en lichtgevend**: elke slide staat op een verloop van navy (#00144A) naar violet (#5D36FF), vaak met een lichtere gloed onderin of een grote zachte boog/ring.
- **Grote witte tekst, één accentkleur**: kop in wit, één woord of de afsluitende punt in teal (#67EDD5). Voorbeeld: "Build things that matter**.**"
- **Gigantische cijfers op de achtergrond**: een cijfer van 720pt in indigo (#4D2FDF -> #5D36FF), nauwelijks zichtbaar, met daarvoor het echte statistiek-getal en een label.
- **Terugkerend motief**: de *pill trio* (teal stip, lange blauwe pil, korte roze pil) als scheiding/ondertekening, en de *sparkle cluster* (violet stip, teal 4-puntige ster, roze stip) als opener boven koppen.
- **Chips en tag-clouds**: afgeronde labels in blauw/teal/roze/lila voor technologieën en thema's.
- **Glass cards en bento-grids**: halftransparante violette kaarten met dunne lila rand.
- **Concentrische ringen** in teal/roze/blauw rond een gloeiend middelpunt.
- **Isometrische illustraties** in lila/violet (gebouwen, mensen, schild, data, code, robot).

## 2. Kleuren
| Token | Hex | Gebruik in deck |
|---|---|---|
| `--tx-violet` | #5D36FF | Hoofdkleur, eind van verloop, kaarten |
| `--tx-violet-light` | #7858FF | Gloed, hover, randen |
| `--tx-violet-deep` | #2C12AD | Donkere kaartvlakken |
| `--tx-indigo-numeral` | #4D2FDF | Gigantische achtergrondcijfers |
| `--tx-navy` | #00144A | Start van verloop, donkerste vlak |
| `--tx-navy-ink` | #0E0537 | Tekst op lichte chips |
| `--tx-blue` | #1A98FF | Lange pil, chips, lijnen |
| `--tx-blue-sap` | #0070F2 | Links/knoppen (SAP-blauw) |
| `--tx-sky` | #89D1FF | Labels onder stats, subteksten |
| `--tx-teal` | #67EDD5 | Accentwoord, eyebrow, punt, stip |
| `--tx-pink` | #FF8BF0 | Korte pil, chips, ringen |
| `--tx-lilac` | #B894FF | Chips, ringen, randen |
| `--tx-gold` | #F0AB00 | Spaarzaam (SAP-goud) |

Frequentie in het deck: violet en blauw domineren, daarna sky, gold, teal en roze als accenten. Houd per scherm maximaal één felle accentkleur dominant.

**Achtergrondverloop**: `linear-gradient(180deg, #00144A 0%, #1C0C6E 45%, #5D36FF 100%)`. Variant met gloed onderin: `radial-gradient(120% 80% at 50% 110%, #7858FF, #5D36FF 30%, #1C0C6E 65%, #00144A)`.

**Contrast**: wit op #5D36FF haalt ruim 4.5:1. Teal en roze zijn alleen geschikt als tekstkleur op navy/donkerviolet, niet op #5D36FF bij kleine tekst. Op teal/roze/lila chips: donkere tekst (#0E0537).

## 3. Typografie
Het deck gebruikt **72 Brand** (SAP-merkfont) in twee gewichten:
- **72 Brand Medium**: koppen, grote cijfers, titels.
- **72 Brand Book**: eyebrows, labels, bodytekst.

Op de site: de bestaande `fonts/72-Regular|Bold|Black.woff2` zijn de webversie van 72. Gebruik Regular voor "Book" en Bold voor "Medium" (Bold oogt iets zwaarder; voeg 72 Semibold toe als je de deck-look exact wilt).

| Rol | Deck (op 1920px slide) | Web (CSS token) | Gewicht | Kleur |
|---|---|---|---|---|
| Achtergrondcijfer | 520 tot 720pt | `--tx-numeral` clamp(18rem, 45vw, 40rem) | Medium | #4D2FDF -> #5D36FF, ~55% opacity |
| Statistiek | 180pt ("140") | `--tx-stat` clamp(4.5rem, 12vw, 11rem) | Medium | wit |
| Display-kop | 60pt | `--tx-display` clamp(3rem, 7vw, 6rem) | Medium | wit, accent teal |
| Sectiekop | 40pt | `--tx-h2` clamp(2rem, 4vw, 3rem) | Medium | wit |
| Eyebrow / kicker | 24pt | `--tx-eyebrow` clamp(1rem, 1.6vw, 1.5rem) | Book | teal |
| Body | 20pt | 1.125rem | Book | wit / sky |
| Chip/label | 14pt | 0.875rem | Book | wit of navy-ink |

Overige regels uit het deck:
- Koppen gecentreerd, korte regels (2 tot 4 woorden per regel), regelafstand ~1.05.
- Lichte negatieve letterspatiëring op grote koppen (spc -150 = ca. -0.015em).
- Patroon "eyebrow in teal boven witte kop" (bijv. "Available on" / "macOS & Windows", "Today's" / "technological shift").
- Patroon "drieregelige claim met laatste regel in teal" ("Better outcomes. Faster. **At lower cost.**").
- Geen hoofdletters-only, geen cursief.

## 4. Componenten (CSS in `sitnl-teched-tokens.css`)
- `.tx-surface`: navy-violet achtergrond met witte 72-tekst.
- `.tx-eyebrow`, `.tx-display`, `.tx-accent`: kopopbouw.
- `.tx-stat` met `data-numeral`: groot getal met gigantisch cijfer erachter. Ideaal voor datum (23 januari), aantal sessies, sprekers of deelnemers.
- `.tx-pills`: pill trio als scheiding (`<span class="tx-pills"><i></i></span>`).
- `.tx-chips` / `.tx-chip--teal|pink|lilac|ghost`: tracks en technologieën in de agenda-filter.
- `.tx-card`, `.tx-bento`: glass cards in bento-grid voor agenda/locatie/sponsors/speakers.
- `.tx-btn`, `.tx-btn--outline`: witte pilknop met violette tekst, of outline met teal rand.
- `.tx-float`: subtiele zweefanimatie (alleen zonder reduced-motion).

`demo.html` laat alles samen zien.

## 5. Assets (map `artifacts/teched-2026-visuals/` in de SitNL-map)
**motifs/** (zelf nagebouwd als vector, vrij te schalen en te kleuren)
- `pill-trio.svg`, `sparkle-cluster.svg`, `sparkle.svg` (currentColor), `rings.svg`.

**backgrounds/** (WebP, max 2560px)
- Volledige achtergronden: `gradient-band-1`, `gradient-wide`, `globe-arcs` (wereldbol met lichtbogen), `dot-network` (netwerk van stippen, transparant), `arc-half-ring` (violette halve ring), `text-pattern-build-things` (ring gevuld met "Build things that matter"-tekstpatroon).
- Overlays met transparantie om over een vlak te leggen: `gradient-band-2`, `gradient-horizon-1/2`, `gradient-floor-glow`, `gradient-side`, `gradient-vertical-1/2`, `grid-floor` (perspectiefraster).

**shapes/** (PNG + WebP, transparant)
- `star-glass-1/2`, `star-outline-glass`: glazen 4-puntige sterren.
- `glass-disc-1/2`, `glass-cylinder-1/2`: glazen platforms/sokkels om illustraties op te zetten.

**illustrations/** (PNG + WebP, transparant, isometrisch)
- `iso-analytics`, `iso-finance-docs`, `iso-logistics-ship`, `iso-collaboration`, `iso-people`, `iso-factory`, `iso-shield`, `iso-data-network`, `iso-code`, `iso-robot`, `iso-robot-arm`. Mooi als vervanging van de huidige iconen in de agenda-slots of bij sponsor/locatie-teasers.

**kasimir/**: zeven geoptimaliseerde Kasimir-loops met posters, plus `kasimir.css`, `kasimir.js` en `KASIMIR.md` (gebruik, inbouw, rechten).

**reference/**: 17 sleutelslides als WebP ter inspiratie (opening, "Build things that matter", groot cijfer, chip-cloud, ringen, bento, value-cards).

## 6. Toepassing op de site (suggesties)
- **Hero**: radial-achtergrond, sparkle cluster, eyebrow "SAP Inside Track Netherlands", kop "Build things that matter." (of een eigen SitNL-claim met teal punt), datum in sky, witte "Registreer"-pil, pill trio eronder, ringen rechtsboven.
- **Datum-blok**: `.tx-stat` met 23 en "januari 2027", achtergrondcijfer "27".
- **Agenda-filter**: huidige filters als chips; tracks krijgen elk een vaste chipkleur.
- **Teasers (agenda/locatie/sponsors)**: bento met glass cards en een isometrische illustratie per kaart.
- **Sectiescheiding**: pill trio in plaats van lijnen.
- **Oranje SitNL-logo**: het oranje (#FF8000) botst met violet. Opties: logo in wit op donker gebruiken, of oranje bewust als extra accent alleen in het logo houden.

## 7. Wat bewust niet is meegenomen
- SAP-productscreenshots (Joule Work, Outlook, devices), Joule-diamant, app-iconen en logo's van derden (AWS, Salesforce, browsers etc.): niet relevant voor een community-site en merkrechtelijk lastig.
- Video's uit het deck (6 stuks).
- Het SAP-merkfont 72 Brand zelf; de site gebruikt al de publieke 72-webfont.
- Check bij twijfel de SAP-brand-richtlijnen voordat je keynote-assets publiek op sitnl.nl zet.
