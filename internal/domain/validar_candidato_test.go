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
		DataInicio: "2026-07-18", DataExpiracao: "2026-07-20",
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
		DataInicio: "2026-07-18", DataExpiracao: "2026-07-20",
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
		DataInicio: "2026-07-18", DataExpiracao: "2026-07-20",
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
		DataInicio: "2026-07-18", DataExpiracao: "2026-07-20",
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
		DataInicio: "2020-01-01", DataExpiracao: "2020-01-01",
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
		DataInicio: "2020-01-01", DataExpiracao: "20/01/2020",
	}

	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoDataExpiracaoInvalida {
		t.Fatalf("esperava data_expiracao_invalida, got %#v", falha)
	}
}

func TestValidarCandidato_rejeitaInicioAposExpiracao(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataInicio:    "2026-07-25",
		DataExpiracao: "2026-07-20",
	}

	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoVigenciaInvalida {
		t.Fatalf("esperava vigencia_invalida, got %#v", falha)
	}
}

func TestValidarCandidato_aceitaVigenciaUmDia(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataInicio:    "2026-07-20",
		DataExpiracao: "2026-07-20",
	}
	if _, falha := domain.ValidarCandidato(c); falha != nil {
		t.Fatalf("mesmo dia deve ser válido: %#v", falha)
	}
}

func TestValidarCandidato_rejeitaDataInicioInvalida(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataInicio:    "20/07/2026",
		DataExpiracao: "2026-07-20",
	}
	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoDataInicioInvalida {
		t.Fatalf("esperava data_inicio_invalida, got %#v", falha)
	}
}

func TestValidarCandidato_rejeitaDataExpiracaoVazia(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:    "Banana",
		Valor:      2.5,
		Quantidade: 1,
		Medida:     "unidade",
		DataInicio: "2026-07-18",
	}
	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoDataExpiracaoInvalida {
		t.Fatalf("esperava data_expiracao_invalida, got %#v", falha)
	}
}

func TestValidarCandidato_aceitaPromocaoClube(t *testing.T) {
	clube := true
	c := domain.CandidatoOferta{
		Produto:       "Leite",
		Valor:         5,
		Quantidade:    1000,
		Medida:        "ml",
		DataInicio:    "2026-07-18",
		DataExpiracao: "2026-07-20",
		Promocao: &domain.Promocao{
			PromocaoClube:    &clube,
			ValorPromocional: 4,
		},
	}

	oferta, falha := domain.ValidarCandidato(c)
	if falha != nil {
		t.Fatalf("promocao clube válida: %#v", falha)
	}
	if oferta.Promocao == nil || oferta.Promocao.PromocaoClube == nil || !*oferta.Promocao.PromocaoClube {
		t.Fatalf("promocao: %#v", oferta.Promocao)
	}
}

func TestValidarCandidato_rejeitaPromocaoCartaoEClube(t *testing.T) {
	cartao, clube := true, true
	c := domain.CandidatoOferta{
		Produto:       "Leite",
		Valor:         5,
		Quantidade:    1000,
		Medida:        "ml",
		DataInicio:    "2026-07-18",
		DataExpiracao: "2026-07-20",
		Promocao: &domain.Promocao{
			PromocaoCartao:   &cartao,
			PromocaoClube:    &clube,
			ValorPromocional: 4,
		},
	}
	_, falha := domain.ValidarCandidato(c)
	if falha == nil || falha.Codigo != domain.CodigoPromocaoInvalida {
		t.Fatalf("esperava promocao_invalida, got %#v", falha)
	}
}

func TestValidarCandidato_marcaOpcional(t *testing.T) {
	c := domain.CandidatoOferta{
		Produto:       "Banana",
		Valor:         2.5,
		Quantidade:    1,
		Medida:        "unidade",
		DataInicio: "2026-07-18", DataExpiracao: "2026-07-20",
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
		DataInicio: "2026-07-18", DataExpiracao: "2026-07-20",
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
