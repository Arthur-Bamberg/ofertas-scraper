# Ofertas Scraper

Sistema que coleta PDFs de Fontes de ofertas, rasteriza páginas em imagens, extrai Ofertas via Extrator e persiste Produtos, Mercados e histórico de preços com Artefatos para debug.

## Language

**Oferta**:
Observação de preço extraída de um Documento: valor, quantidade, medida, data de expiração (fim da validade no encarte; pode estar no passado para histórico) e promoção opcional; refere-se a um Produto em um Mercado, com Marca opcional. O histórico de preços de um Produto é o conjunto de Ofertas ao longo do tempo (via Documentos); o “período atual” é filtro do consumidor da base, não um estado embutido no Produto.
_Avoid_: Deal, item, listing, produto, histórico de produto (como entidade separada)

**Produto**:
Identidade de catálogo do que está à venda, sem marca (ex.: “Arroz integral”, “Arroz branco parboilizado”); carrega categorias taxonômicas para navegação e agregação. Distinta por tipo vendável; N Marcas aparecem via Ofertas, não como lista fixa no Produto.
_Avoid_: Oferta, item, SKU, variante (como substituto de Produto), marca

**Marca**:
Identidade comercial do fabricante ou rótulo (ex.: Camil, Tio João), independente do Produto e do Mercado. Opcional na Oferta quando o encarte não traz marca. Quando presente, o Extrator a identifica nas imagens e o domain casa ou cria a Marca.
_Avoid_: fabricante como texto solto na Oferta, brand

**Categoria**:
Rótulo taxonômico de um Produto para filtrar e agrupar (ex.: mercearia, grãos, arroz). Não distingue tipos vendáveis — isso são Produtos diferentes. Novas categorias vistas na extração unem-se às já existentes do Produto.
_Avoid_: tipo, variante, tag solta na Oferta

**Fonte**:
URL configurada e persistida cuja resposta revela URLs de PDFs (HTML/JS/JSON); o nome do arquivo identifica o Documento e o download usa a URL do link. Pode apontar para página HTML ou endpoint JSON de ofertas. Pertence a um Mercado; pode incluir filtro opcional por regex sobre o nome do Documento. Pode obter `dataExpiracao` a partir do nome do arquivo quando o Extrator omite a data (`fallbackDataExpiracaoFilename` na Fonte, não no Mercado): a flag começa desligada e liga-se automaticamente ao aparecer um candidato sem data; o fallback aplica-se já na mesma tentativa. Se o Extrator envia a data, ela prevalece.
_Avoid_: Site, link, URL, origem

**Mercado**:
Identidade comercial (rede ou bandeira) à qual uma Fonte pertence; sujeito da comparação de Ofertas e do histórico de preços entre estabelecimentos. Não é extraído das imagens — vem da configuração da Fonte.
_Avoid_: loja, supermercado, site, Fonte

**Documento**:
PDF identificado em uma Fonte pelo nome do arquivo e pelo dia da descoberta; rastreado ao longo do processamento (download, rasterização, extração). Estados persistidos: processando, concluído, parcial, falhou. Entra em processando no início do tratamento (não há estado persistido “listado mas ainda não baixado”). Sem pelo menos uma Oferta persistida (lista vazia do Extrator ou só Falhas de Extração), o Documento termina em falhou. Em falha dura de processamento (incluindo Extrator indisponível após retentativas), permanece consultável com o motivo do último erro.
_Avoid_: PDF, arquivo, anexo, descoberto (como estado persistido)

**Extrator**:
Capacidade de obter candidatos a Oferta a partir das imagens de um Documento (rótulos de produto, marca e categorias, mais preço, data de expiração e promoção). Espera-se `dataExpiracao` em todo candidato; ausência é tratada via fallback de filename na Fonte.
_Avoid_: Gemini, IA, conversor, parser, LLM

**Medida**:
Unidade de quantidade de uma Oferta, sempre normalizada para `g`, `ml` ou `unidade` (kg → 1000 g; L → 1000 ml).
_Avoid_: unidade de medida, kg, litro, L

**Promoção**:
Condição comercial opcional de uma Oferta: leve/pague, quantidade com valor promocional, ou preço exclusivo de cartão.
_Avoid_: desconto, oferta especial, deal

**Falha de Extração**:
Registro de uma tentativa de Oferta que não passou na validação, vinculada ao Documento de origem: código estável do motivo, detalhe livre opcional e cópia do candidato rejeitado.
_Avoid_: export com erro, erro de IA, rejeição

**Artefato**:
Material obtido ou gerado em uma tentativa de processamento de um Documento e retido para debug — tipicamente PDF original, imagens enviadas ao Extrator, resposta bruta do Extrator e resultado validado (Ofertas e Falhas de Extração). Em falha dura, persiste-se só o que a tentativa chegou a produzir (best-effort); não se fabricam placeholders para etapas que não rodaram. Cada reprocessamento acrescenta uma nova tentativa; tentativas anteriores permanecem. Ofertas e Falhas de Extração persistidas no estado atual do Documento são substituídas na nova tentativa — o histórico de tentativas vive nos Artefatos. Hoje em armazenamento local; depois em bucket.
_Avoid_: arquivo, blob, export, attachment, log

## Qualidade local

Antes de cada commit, o hook versionado em `.githooks/pre-commit` executa `go test ./...` (equivalente a Husky em projetos Node). Instalar uma vez por clone com `./scripts/install-git-hooks.sh`. Detalhes em [`AGENTS.md`](./AGENTS.md) e ADR 0024.
