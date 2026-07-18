package application_test

import (
	"testing"

	"ofertas-scraper/internal/application"
	"ofertas-scraper/internal/domain"
)

type stubDateParser struct {
	date string
	ok   bool
}

func (s stubDateParser) Parse(string) (string, bool) { return s.date, s.ok }

func TestValidarExtracao_listaVaziaEhFalhou(t *testing.T) {
	_, falhas, estado := application.ValidarExtracao(nil)
	if len(falhas) != 0 {
		t.Fatalf("falhas: %d", len(falhas))
	}
	if estado != domain.EstadoFalhou {
		t.Fatalf("estado: got %q", estado)
	}
}

func TestValidarExtracao_mistoParcial(t *testing.T) {
	candidatos := []domain.CandidatoOferta{
		{Produto: "Arroz", Valor: 10, Quantidade: 1000, Medida: "g", DataExpiracao: "2026-07-20"},
		{Produto: "Feijão", Valor: 8, Quantidade: 1, Medida: "kg", DataExpiracao: "2026-07-20"},
	}
	validas, falhas, estado := application.ValidarExtracao(candidatos)
	if len(validas) != 1 || len(falhas) != 1 {
		t.Fatalf("validas=%d falhas=%d", len(validas), len(falhas))
	}
	if estado != domain.EstadoParcial {
		t.Fatalf("estado: got %q", estado)
	}
}

func TestAplicarFallbackDataExpiracao(t *testing.T) {
	c := domain.CandidatoOferta{Produto: "X", Valor: 1, Quantidade: 1, Medida: "unidade"}
	fonte := domain.Fonte{FallbackDataExpiracaoFilename: true}
	got := application.AplicarFallbackDataExpiracao(c, fonte, "encarte-2026-07-25.pdf", stubDateParser{date: "2026-07-25", ok: true})
	if got.DataExpiracao != "2026-07-25" {
		t.Fatalf("got %q", got.DataExpiracao)
	}

	comExtrator := c
	comExtrator.DataExpiracao = "2026-07-01"
	got = application.AplicarFallbackDataExpiracao(comExtrator, fonte, "encarte-2026-07-25.pdf", stubDateParser{date: "2026-07-25", ok: true})
	if got.DataExpiracao != "2026-07-01" {
		t.Fatalf("Extrator deve vencer, got %q", got.DataExpiracao)
	}
}
