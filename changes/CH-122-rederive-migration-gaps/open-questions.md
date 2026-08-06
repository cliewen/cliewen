---
id: CH-122-questions
type: open-questions
status: open
links: [CH-122, P-012, M-057, CAP-003, ADR-038]
title: Blokerende spørgsmål for CH-122
---

# CH-122 — åbne spørgsmål

## Q1 — Hvordan skal dette repository oprette identitetsledgeren, når et urelateret fund om en ældre carrier blokerer `clue migrate --apply`?

`go run ./cmd/clue migrate` planlægger `MIG-008 .clue/id-ledger.yaml: seed the identity ledger with 149 live id(s) from the current corpus scan`. Den påkrævede kommando `go run ./cmd/clue migrate --apply` nægter at anvende den sikre ændring, fordi repositoryet mangler `.github/workflows/clue.yml`. Den rapporterer: `MIG-003 ... thin CI caller is missing; run clue init before migrating`.

Repositoryet har i stedet sine nuværende CI-filer under `.github/workflows/`, og M-057 angiver ikke, om vi skal tilføje den ældre caller, ændre migrationens atomaritet så urelaterede fund ikke forhindrer MIG-008, eller udsætte repositoryets oprettelse af ledgeren. At tilføje calleren eller ændre migrationens fejlsemantik kan ændre CI- eller migrationskontrakten, så ingen af delene antages her.

**Beslutning nødvendig:** Godkend én af disse veje, og angiv om den hører til M-057 eller kræver en separat afgrænset kerneændring.
