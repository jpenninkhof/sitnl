# Kasimir op de SitNL-site

Kasimir, de TechEd-kat, is een knipoog naar de community: ontstaan als grap van de productieploeg tijdens de eerste virtuele TechEd (2020), daarna een "Where's Waldo"-spel voor kijkers, en inmiddels de officiële mascotte. De huidige 3D-versie (gestileerd, met vacht, rig en animaties) is in 2024 gemaakt door Koen Vissers (kunemann.com) samen met Goebel + Mattes, in opdracht van SAP.

## Rechten eerst
Koen vermeldt dat de animaties gemaakt zijn voor interne SAP-events en merchandise. Voor gebruik op de publieke site sitnl.nl is dus toestemming nodig van SAP (TechEd/community-team) en een akkoord van Koen. Host de bestanden zelf (map `kasimir/`), niet via een directe link naar files.koen.work, en zet een credit in de footer, bijvoorbeeld: "Kasimir 3D: Koen Vissers (kunemann.com) met Goebel + Mattes, met dank aan SAP TechEd."

## Clips (web-geoptimaliseerd, zonder geluid, H.264, met poster-JPG)
| Bestand | Duur | Formaat | Achtergrond | Wat je ziet | Voorstel op de site |
|---|---|---|---|---|---|
| `kasimir-idle.mp4` | 6,7 s | 720x720 | navy | Zit, knippert, kijkt rond | Hero naast de kop, of de "peek" rechtsonder |
| `kasimir-early-bird.mp4` | 11 s | 1280x720 | zwart | In vogelpak, klapt laptop open | Early-bird registratie-blok |
| `kasimir-money.mp4` | 4 s | 720x720 | navy | Zwaait met geld | Sponsorpagina |
| `kasimir-nails.mp4` | 7 s | 720x720 | navy | Vijlt nagels, verveeld | "Coming soon" / agenda nog niet bekend / CfP gesloten |
| `kasimir-beard.mp4` | 10,5 s | 720x720 | navy | Krijgt blauwe baard, close-up | Speakers- of 404-pagina |
| `kasimir-robot.mp4` | 27 s | 1280x720 | navy | Rijdt mee op een robotstofzuiger | "Agent Kasimir" / AI-track, of "That's a wrap" |
| `kasimir-lofi.mp4` | 16,7 s | 820x568 | zwart | Koptelefoon op, typt op laptop | Agenda / hands-on sessies, "lo-fi beats to code to" |

Totaal ca. 2,8 MB. Originelen (met geluid bij robot en lofi) staan niet in de site-map.

## Inbouwen
```html
<link rel="stylesheet" href="kasimir/kasimir.css">

<!-- navy clip: randen lossen op in de achtergrond -->
<div class="kasimir kasimir--blend">
  <video data-kasimir autoplay muted loop playsinline poster="kasimir/kasimir-money.jpg" src="kasimir/kasimir-money.mp4"></video>
</div>

<!-- zwarte clip: als ingelijst scherm -->
<div class="kasimir kasimir--wide kasimir--screen">
  <video data-kasimir autoplay muted loop playsinline poster="kasimir/kasimir-early-bird.jpg" src="kasimir/kasimir-early-bird.mp4"></video>
</div>

<!-- optioneel: Kasimir gluurt na 20 s rechtsonder binnen ("spot Kasimir") -->
<button class="kasimir-peek" data-src="kasimir/kasimir-idle.mp4" data-delay="20000"></button>

<script src="kasimir/kasimir.js"></script>
```

- `kasimir.js` laadt een video pas als hij in beeld komt en pauzeert hem daarbuiten. Bij "reduce motion" blijft alleen de poster staan.
- Navy clips zijn gerenderd op #001949 (boven) naar #002869 (onder). Zet ze in een sectie met `background: #001949` of `--tx-navy` voor een naadloze overgang.
- De peek telt hoe vaak iemand Kasimir heeft gespot (localStorage). Leuk als knipoog naar het Where's Waldo-spel uit de virtuele TechEd-tijd.

## Ideeën
- Verstop elk jaar een kleine Kasimir op een andere pagina; wie hem vindt, krijgt iets bij de registratiebalie.
- Gebruik `kasimir-nails` zolang de agenda nog niet compleet is en wissel naar `kasimir-lofi` zodra hij live staat.
- Vraag Koen of er een versie met transparante achtergrond (WebM met alfa) bestaat; dan kan Kasimir los over de hero lopen.
