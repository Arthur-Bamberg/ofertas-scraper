# Backlog

Itens futuros — fora do escopo do MVP atual. Não são compromissos de sprint; são lembretes para revisitarmos.

## Próxima implementação (ordem acordada)

1. [x] Adaptador Gemini do Extrator (prompt cache + structured output) — ADR 0023; stub via `EXTRATOR_STUB=1`.
2. [x] Vigência `dataInicio` + origens + fallbacks sempre ativos; remove flag na Fonte (ADR 0028). Promoção clube (ADR 0029).
3. [ ] Índice Redis `ofertas:produto:{produtoId}` → documentoIds (ADR 0026).
4. [ ] Higiene: remover estado persistido `descoberto` do código/docs restantes (glossário já alinhado).

## Operação e alertas

- [ ] Ferramenta/canal de **notificação** quando o Extrator ou outra infra crítica (Redis, Fonte) ficar indisponível ou quando Documentos falharem de forma crônica — hoje só stdout + estado no Redis (ADR 0022).

## Extração e catálogo

- [ ] **Job de merge de catálogo** (Produto/Marca): fora do job diário; propõe ou aplica fusões de duplicatas por grafia/sentido diferente, usando Gemini como juiz assistido — o match-or-create do pipeline permanece exact normalizado (ADR 0011). Substituí a ideia de fuzzy inline no hot path.
- [ ] Safety-net de Medida no domain (kg/L → g/ml) se o Extrator falhar nisso com frequência (ADR 0004).
- [ ] Estratégia de Fonte além de HTML/JSON genérico, se aparecer Fonte real que não caiba no adaptador atual (ADR 0020).
- [ ] Não usar `?debugbar_time=` (PHP DebugBar) — é artefato de sessão de debug, não API estável.
