package domain_test

import (
	"testing"

	"ofertas-scraper/internal/domain"
)

func TestValidarCandidato_aceitaCandidatoValido(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Arroz integral",
		Valor:         12.9,
		Quantidade:    1000,
		Medida:        "g",
		DataExpiracao: "2026-07-20",
		Categorias:    []string{"mercearia", "arroz"},
	}

	oferta, falha := domain.ValidarCandidato(c)
	if falha != nil {
		t.Fatalf("esperava sem Falha, obteve %#v", falha)
	}
	if oferta.Produto != "Arroz integral" {
		t.Fatalf("produto: got %q", oferta.Produto)
	}
	if oferta.Medida != domain.MedidaG {
		t.Fatalf("medida: got %q", oferta.Medida)
	}
}

func TestValidarCandidato_rejeitaMedidaNaoNormalizada(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Arroz integral",
		Valor:         12.9,
		Quantidade:    1,
		Medida:        "kg",
		DataExpiracao: "2026-07-20",
	}

	_, falha := domain.ValidarCandidato(c)
	if falha == nil {
		t.Fatal("esperava Falha de Extração para medida kg")
	}
	if falha.Codigo != domain.CodigoMedidaInvalida {
		t.Fatalf("codigo: got %q want %q", falha.Codigo, domain.CodigoMedidaInvalida)
	}
	if falha.Candidato.Medida != "kg" {
		t.Fatalf("candidato deve preservar medida rejeitada, got %q", falha.Candidato.Medida)
	}
}

func TestValidarCandidato_rejeitaProdutoVazio(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "  ",
		Valor:         1,
		Quantidade:    1,
		Medida:        "unidade",
		DataExpiracao: "2026-07-20",
	}

	_, falha := domain.ValidarCandidato(c)
	if falha == nil {
		t.Fatal("esperava Falha para produto vazio")
	}
	if falha.Codigo != domain.CodigoProdutoInvalido {
		t.Fatalf("codigo: got %q", falha.Codigo)
	}
}

func TestValidarCandidato_rejeitaValorNaoPositivo(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         0,
		Quantidade:    1,
		Medida:        "unidade",
		DataExpiracao: "2026-07-20",
	}

	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoValorInvalido {
		t.Fatalf("esperava valor_invalido, got %#v", falha)
	}
}

func TestValidarCandidato_aceitaDataExpiracaoPassada(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataExpiracao: "2020-01-01",
	}

	oferta, falha := domain.ValidarCandidato(c)
	if falha != nil {
		t.Fatalf("data passada deve ser aceita para histórico: %#v", falha)
	}
	if oferta.DataExpiracao != "2020-01-01" {
		t.Fatalf("data: got %q", oferta.DataExpiracao)
	}
}

func TestValidarCandidato_rejeitaDataExpiracaoInvalida(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataExpiracao: "20/01/2020",
	}

	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoDataExpiracaoInvalida {
		t.Fatalf("esperava data_expiracao_invalida, got %#v", falha)
	}
}

func TestValidarCandidato_marcaOpcional(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataExpiracao: "2026-07-20",
	}

	oferta, falha := domain.ValidarCandidato(c)
	if falha != nil {
		t.Fatalf("marca omitida deve ser válida: %#v", falha)
	}
	if oferta.Marca != "" {
		t.Fatalf("marca: got %q", oferta.Marca)
	}
}

func TestValidarCandidato_aceitaPromocaoLevePague(t *testing.T) {
	leve, pague := 3.0, 2.0
	c := domain.CandidatoOferta{
		Produto:       "Refrigerante",
		Marca:         "Cola",
		Valor:         8,
		Quantidade:    2000,
		Medida:        "ml",
		DataExpiracao: "2026-07-20",
		Promocao: &domain.Promocao{
			Leve:             &leve,
			Pague:            &pague,
			ValorPromocional: 12,
		},
	}

	oferta, falha := domain.ValidarCandidato(c)
	if falha != nil {
		t.Fatalf("promocao válida: %#v", falha)
	}
	if oferta.Promocao == nil || oferta.Promocao.ValorPromocional != 12 {
		t.Fatalf("promocao: %#v", oferta.Promocao)
	}
}
