package application

import "ofertas-scraper/internal/domain"

// ValidarExtracao validates Extrator candidates and derives Documento terminal state.
func ValidarExtracao(candidatos []domain.CandidatoOferta) (validas []domain.OfertaValidada, falhas []domain.FalhaExtracao, estado domain.EstadoDocumento) {
	for _, c := range candidatos {
		oferta, falha := domain.ValidarCandidato(c)
		if falha != nil {
			falhas = append(falhas, *falha)
			continue
		}
		validas = append(validas, oferta)
	}
	estado = domain.EstadoAposValidacao(len(validas), len(falhas))
	return validas, falhas, estado
}

// AplicarFallbackDataExpiracao fills empty dataExpiracao from filename when Fonte allows it (ADR 0013).
func AplicarFallbackDataExpiracao(c domain.CandidatoOferta, fonte domain.Fonte, filename string, parser domain.FilenameDateParser) domain.CandidatoOferta {
	if c.DataExpiracao != "" || !fonte.FallbackDataExpiracaoFilename || parser == nil {
		return c
	}
	if date, ok := parser.Parse(filename); ok {
		c.DataExpiracao = date
	}
	return c
}
