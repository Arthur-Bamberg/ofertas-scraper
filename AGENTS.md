# AGENTS.md — ofertas-scraper

Instructions for coding agents working on this repository. Domain language lives in [`CONTEXT.md`](./CONTEXT.md); architectural decisions in [`docs/adr/`](./docs/adr/). Prefer those terms (`Oferta`, `Fonte`, `Documento`, `Extrator`, …) over synonyms.

## What this system does

Daily job (08:00 America/Sao_Paulo) that:

1. Loads **Fontes** from Redis (Upstash)
2. GETs each Fonte, discovers `.pdf` names, applies optional per-Fonte filter
3. Creates **Documentos** (identity: Fonte + filename + discovery day)
4. Downloads PDF → rasterizes pages → downscales images → sends to **Extrator**
5. Validates results in **domain** → persists **Ofertas** and **Falhas de Extração**
6. Always stores **Artefatos** for debug

Processing is **sequential** in the MVP (Fonte by Fonte, Documento by Documento).

## Stack

| Concern | Choice |
|---------|--------|
| Language | Go |
| Persistence | Upstash Redis (REST) |
| Local Redis | Redis + [Serverless Redis HTTP (SRH)](https://upstash.com/docs/redis/sdks/ts/developing) via Docker Compose |
| Extrator impl | Gemini (prompt cache + structured output) |
| Schedule | External cron/systemd timer; binary is a one-shot CLI |
| Timezone | `America/Sao_Paulo` |

## Repository layout (Clean Architecture)

```
cmd/ofertas-scraper/          # main / DI wiring
internal/
  domain/                     # entities + ports (interfaces)
  application/                # use cases / job flows
  presentation/               # CLI entry only (calls application)
  infra/                      # Gemini, Upstash, HTTP Fonte, PDF→image, local Artefatos
schemas/oferta.json           # Extrator candidate / Oferta contract
schemas/extracao.json         # Extrator root response (list of candidates)
prompts/extrator.txt          # stable system prompt (bump only when explicitly versioning)
docker-compose.yml            # Redis + SRH
.githooks/                    # versioned Git hooks (pre-commit → go test)
scripts/install-git-hooks.sh  # once per clone: core.hooksPath=.githooks
CONTEXT.md
docs/adr/
```

### Dependency rule

- `presentation` → `application` → `domain`
- `infra` implements interfaces defined in `domain`
- `domain` never imports `infra`, `application`, or `presentation`
- `application` depends on `domain` interfaces only

### Ports (interfaces in `domain`)

Required ports (names may vary; responsibilities must not):

- **MercadoRepository** — list/save Mercados
- **FonteRepository** — list/save Fontes (each Fonte references a Mercado)
- **ProdutoRepository** — list/save Produtos (catalog identity + categorias)
- **MarcaRepository** — list/save Marcas
- **DocumentoRepository** — track Documento lifecycle
- **OfertaRepository** — persist Ofertas linked to Documento (and Produto + Marca + Mercado)
- **FalhaExtracaoRepository** — persist Falhas de Extração linked to Documento
- **Extrator** — images in → candidate Ofertas (raw) out
- **ArtefatoStore** — save/load Artefatos (local now; bucket later behind same interface)
- **FonteClient** — HTTP GET + PDF name discovery
- **Rasterizer** — PDF bytes → downscaled page images
- **FilenameDateParser** — optional parse of `dataExpiracao` from Documento filename (when Fonte enables fallback)

## Documento lifecycle

States: `descoberto` → `processando` → `concluido` | `parcial` | `falhou`

| State | Meaning |
|-------|---------|
| `descoberto` | Listed from Fonte, not downloaded yet |
| `processando` | Download / raster / Extrator in progress |
| `concluido` | ≥1 Oferta saved, zero Falhas de Extração |
| `parcial` | ≥1 Oferta saved and ≥1 Falha de Extração |
| `falhou` | Hard failure **or** zero Ofertas persistidas (Extrator com `ofertas` vazio, ou só Falhas de Extração) |

Same-day re-run: skip `concluido` and `parcial`; retry `falhou`, `descoberto`, and orphan `processando`.

## Extrator (Gemini adapter)

- System prompt: `prompts/extrator.txt` — must spell out glossary definitions for Produto (sem marca), Marca, and Categoria (taxonômia, não tipo vendável) so extraction stays assertive (ADR 0012)
- Output schema: `schemas/extracao.json` → items conform to `schemas/oferta.json`
- Keep prompt/schema **unversioned** until an explicit version bump is requested (ADR 0014)
- Cache the **stable prompt** (and schema binding) via Gemini context cache; **do not** cache Documento images
- Adapter lives in `infra`; use case only sees `Extrator`
- Domain validates every candidate Oferta after extraction (do not trust the model alone)

### Oferta rules agents must respect

- Extrator candidates include `produto`, optional `marca`, `categorias[]`, plus `valor` / `quantidade` / `medida` / `dataExpiracao` / optional `promocao` (see ADR 0010); persisted Oferta stores `produtoId`, `mercadoId`, and optional `marcaId` after match-or-create (ADR 0015)
- `medida` is only `g` | `ml` | `unidade`
- Extrator must normalize **kg → 1000 g** and **L → 1000 ml** (adjust `quantidade`) before output; domain does **not** convert — any other `medida` is a Falha de Extração (see ADR 0004; hybrid domain safety-net deferred)
- `dataExpiracao` = end of validity on the flyer (not discovery day); prefer Extrator value; if missing and Fonte has filename-fallback enabled, parse from Documento filename; Extrator always wins when present (ADR 0013); past dates are valid (price history — ADR 0016)
- `promocao` is optional and one of three shapes (leve/pague, quantidade+valor, cartão)
- Fonte field `fallbackDataExpiracaoFilename` defaults to `false` (opt-in per Fonte)
- Domain match-or-create for Produto/Marca uses normalized exact label match only (ADR 0011); no fuzzy matching in the MVP

## Artefatos

For every processed Documento, always persist:

1. Original PDF
2. Images sent to Extrator
3. Raw Extrator response
4. Validated result (Ofertas + Falhas de Extração)

Store via `ArtefatoStore` only — never write files ad hoc from use cases.

## Local environment

- `docker compose up` starts Redis + SRH (Upstash-compatible REST)
- App uses `UPSTASH_REDIS_REST_URL` + `UPSTASH_REDIS_REST_TOKEN` (point at SRH locally; real Upstash in cloud)
- Artefatos default to a local directory (e.g. `./.data/artefatos`)
- Cron example: `0 8 * * *` with `TZ=America/Sao_Paulo` calling the CLI once

### Git hooks (pre-commit tests)

Go has no Husky. This repo versions hooks under `.githooks/` and points Git at them once per clone (ADR 0024):

```bash
./scripts/install-git-hooks.sh
```

After that, every `git commit` runs `go test ./...` and aborts on failure. Do not use `--no-verify` unless explicitly required.

### Suggested env vars

```
UPSTASH_REDIS_REST_URL=
UPSTASH_REDIS_REST_TOKEN=
GEMINI_API_KEY=
EXTRATOR_STUB=1
SEED_PATH=./seed/fontes.json
ARTEFATO_ROOT=./.data/artefatos
RASTER_MAX_EDGE_PX=1280
RASTER_JPEG_QUALITY=80
TZ=America/Sao_Paulo
```

CLI: `ofertas-scraper seed` (upsert Mercados/Fontes from `SEED_PATH`) then `ofertas-scraper run`. Rasterizer needs `pdftoppm` (poppler-utils) on PATH.

## Agent do's and don'ts

**Do**

- Use glossary terms from `CONTEXT.md`
- Add/change persistence and Extrator only behind `domain` interfaces
- Keep job orchestration in `application`
- Keep prompt + schema in sync; only introduce versioned filenames when explicitly asked to bump the Extrator contract (ADR 0014)
- Prefer small, sequential changes with tests around domain validation
- Install Git hooks after cloning (`./scripts/install-git-hooks.sh`) so pre-commit runs `go test ./...`

**Don't**

- Call Gemini or Redis from `domain` or `application` directly
- Put ports in `presentation` (CLI only)
- Send raw PDF bytes to the Extrator (images only)
- Hardcode Fonte URLs (they live in Redis)
- Parallelize Fontes/Documentos in the MVP without an explicit decision
- Commit secrets (`.env`, API keys)
- Skip pre-commit with `--no-verify` unless the user explicitly asks

## Pipeline reference

```
cron 08:00 America/Sao_Paulo
  → CLI (presentation)
    → RunDailyJob (application)
      → FonteRepository.List()
      → for each Fonte (sequential):
          FonteClient.DiscoverPDFs(+ filter)
          → for each new Documento:
              download PDF
              ArtefatoStore.Save(pdf)
              Rasterizer → images (max edge 1280px)
              ArtefatoStore.Save(images)
              Extrator.Extract(images)
              ArtefatoStore.Save(raw)
              domain.Validate → Ofertas | Falhas
              persist + ArtefatoStore.Save(validated)
              update Documento state
```
