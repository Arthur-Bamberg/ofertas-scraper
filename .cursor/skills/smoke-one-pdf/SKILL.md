---
name: smoke-one-pdf
description: >-
  Roda o job diário limitado a um único Documento/PDF (smoke) e devolve
  caminhos/file:// do PDF e do output bruto do Extrator para validação
  manual. Use when the user asks to smoke-test one PDF, validar extrator
  com um encarte, RUN_MAX_DOCUMENTOS, RUN_FONTE_ID, or revisar artefatos
  de uma tentativa.
disable-model-invocation: true
---

# Smoke: um PDF para validação manual

Fluxo do agente: executar `ofertas-scraper run` com limites de smoke, achar a tentativa mais recente de Artefato e **sempre** entregar ao usuário os links locais do PDF e do JSON bruto do Extrator.

## Pré-requisitos

- Redis/SRH no ar (`docker compose up` se local)
- Seed feito (`go run ./cmd/ofertas-scraper seed`)
- Extrator real: `GEMINI_API_KEY` set e `EXTRATOR_STUB=0` (stub não serve para validar qualidade)
- `pdftoppm` no PATH
- Vars tipicamente em `.env` (não commitar segredos)

## Rodar (um Documento)

```bash
# na raiz do repo; carregar .env se o shell não carregar sozinho
export RUN_FONTE_ID="${RUN_FONTE_ID:-fonte-fort}"   # perguntar se o user não disser
export RUN_MAX_DOCUMENTOS=1
export EXTRATOR_STUB=0
# TZ / UPSTASH_* / GEMINI_* / ARTEFATO_ROOT conforme .env

go run ./cmd/ofertas-scraper run
```

**Comportamento:** processa só a Fonte `RUN_FONTE_ID` e para após o **primeiro** PDF mantido pelo `filtroNomeDocumento` da Fonte (ordem da descoberta). Não há `RUN_FILENAME` — para mirar um nome específico, ajuste temporariamente o filtro da Fonte no Redis/seed (regex no filename) ou peça ao user qual Fonte usar.

**Reprocesso no mesmo dia:** Documentos `concluido` / `parcial` são skipados. Se o log mostrar `skip … (concluido|parcial)` e o user quiser nova tentativa, apague o Documento (e Ofertas/Falhas) daquele Fonte+filename+dia no Redis, ou use um Documento em `falhou`. Depois rode de novo.

## Artefatos (onde olhar)

Layout (ADR 0021):

```text
{ARTEFATO_ROOT}/{fonteId}/{dia}/{filename}/{tentativa}/
  original.pdf          ← PDF baixado
  extrator-raw.json     ← saída bruta do Extrator
  validated.json        ← Ofertas + Falhas (útil, mas o foco do user é o raw)
  images/page-NNN.jpg
```

`ARTEFATO_ROOT` default: `./.data/artefatos`. `dia` = discovery day `America/Sao_Paulo`. `tentativa` = UTC `YYYYMMDDTHHMMSS`.

Achar a tentativa mais recente da Fonte (após o run):

```bash
ROOT="${ARTEFATO_ROOT:-./.data/artefatos}"
FONTE="${RUN_FONTE_ID:-fonte-fort}"
# diretório de tentativa mais novo (mtime)
LATEST=$(find "$ROOT/$FONTE" -type f -name 'extrator-raw.json' -printf '%T@ %p\n' 2>/dev/null | sort -nr | head -1 | cut -d' ' -f2-)
DIR=$(dirname "$LATEST")
```

Se `find -printf` não existir (macOS), use `find … -name extrator-raw.json` + ordenar por mtime de outra forma, ou liste `ls -lt` sob `$ROOT/$FONTE/*/`.

## Resposta obrigatória ao usuário

Depois do run (sucesso ou falha parcial com Artefatos gravados), responda **neste formato**, com caminhos absolutos e links `file://`:

```markdown
## Smoke — 1 PDF

- **Fonte:** `{fonteId}`
- **Documento:** `{filename}`
- **Dia:** `{dia}`
- **Tentativa:** `{tentativa}`
- **Estado (log):** `{concluido|parcial|falhou|skip…}`

### PDF
`{abs}/original.pdf`
[file://{abs}/original.pdf](file://{abs}/original.pdf)

### Extrator (raw)
`{abs}/extrator-raw.json`
[file://{abs}/extrator-raw.json](file://{abs}/extrator-raw.json)

### Validado (opcional)
`{abs}/validated.json`
```

Regras:

1. Não resuma o JSON do Extrator a menos que o user peça — o objetivo é ele validar manualmente.
2. Se `extrator-raw.json` ou `original.pdf` não existirem, diga qual etapa falhou (download / raster / extrator) com base no log e nos arquivos presentes.
3. Não commitar `.env`, Artefatos nem chaves.
4. Preferir `go run ./cmd/ofertas-scraper` na raiz do repo; não inventar flags CLI que não existam.

## Checklist

```
- [ ] Redis/SRH ok
- [ ] Seed ok
- [ ] EXTRATOR_STUB=0 + GEMINI_API_KEY
- [ ] RUN_FONTE_ID + RUN_MAX_DOCUMENTOS=1
- [ ] run executado
- [ ] LATEST tentativa localizada
- [ ] Resposta com file:// do PDF e do extrator-raw.json
```
