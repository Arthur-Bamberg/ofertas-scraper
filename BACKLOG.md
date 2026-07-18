# Backlog

Itens futuros — fora do escopo do MVP atual. Não são compromissos de sprint; são lembretes para revisitarmos.

## Operação e alertas

- [ ] Ferramenta/canal de **notificação** quando o Extrator ou outra infra crítica (Redis, Fonte) ficar indisponível ou quando Documentos falharem de forma crônica — hoje só stdout + estado no Redis (ADR 0022).

## Extração e catálogo

- [ ] Adaptador Gemini do Extrator (prompt cache + structured output) — hoje só stub (ADR 0023).
- [ ] Match fuzzy de Produto/Marca se duplicatas por grafia diferente ficarem caras (hoje: match exato normalizado, ADR 0011).
- [ ] Safety-net de Medida no domain (kg/L → g/ml) se o Extrator falhar nisso com frequência (ADR 0004).
- [ ] Estratégia de Fonte além de HTML/JSON genérico, se aparecer Fonte real que não caiba no adaptador atual (ADR 0020).
- [ ] Não usar `?debugbar_time=` (PHP DebugBar) — é artefato de sessão de debug, não API estável.
