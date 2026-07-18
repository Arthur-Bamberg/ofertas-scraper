package domain

import (
	"strings"
	"time"
)

type Medida string

const (
	MedidaG       Medida = "g"
	MedidaML      Medida = "ml"
	MedidaUnidade Medida = "unidade"
)

const (
	CodigoMedidaInvalida         = "medida_invalida"
	CodigoProdutoInvalido        = "produto_invalido"
	CodigoValorInvalido          = "valor_invalido"
	CodigoQuantidadeInvalida     = "quantidade_invalida"
	CodigoDataExpiracaoInvalida  = "data_expiracao_invalida"
	CodigoPromocaoInvalida       = "promocao_invalida"
)

// CandidatoOferta is the Extrator candidate before match-or-create.
type CandidatoOferta struct {
	Produto       string    `json:"produto"`
	Marca         string    `json:"marca,omitempty"`
	Categorias    []string  `json:"categorias,omitempty"`
	Valor         float64   `json:"valor"`
	Quantidade    float64   `json:"quantidade"`
	Medida        string    `json:"medida"`
	DataExpiracao string    `json:"dataExpiracao"`
	Promocao      *Promocao `json:"promocao,omitempty"`
}

// OfertaValidada is a candidate that passed domain validation (labels, not ids).
type OfertaValidada struct {
	Produto       string
	Marca         string
	Categorias    []string
	Valor         float64
	Quantidade    float64
	Medida        Medida
	DataExpiracao string
	Promocao      *Promocao
}

type Promocao struct {
	Leve               *float64 `json:"leve,omitempty"`
	Pague              *float64 `json:"pague,omitempty"`
	QuantidadePromocao *float64 `json:"quantidadePromocao,omitempty"`
	PromocaoCartao     *bool    `json:"promocaoCartao,omitempty"`
	ValorPromocional   float64  `json:"valorPromocional"`
}

type FalhaExtracao struct {
	Codigo    string          `json:"codigo"`
	Detalhe   string          `json:"detalhe,omitempty"`
	Candidato CandidatoOferta `json:"candidato"`
}

func ValidarCandidato(c CandidatoOferta) (OfertaValidada, *FalhaExtracao) {
	produto := strings.TrimSpace(c.Produto)
	if produto == "" {
		return falha(c, CodigoProdutoInvalido, "produto é obrigatório")
	}

	if c.Valor <= 0 {
		return falha(c, CodigoValorInvalido, "valor deve ser > 0")
	}

	if c.Quantidade < 0 {
		return falha(c, CodigoQuantidadeInvalida, "quantidade deve ser >= 0")
	}

	medida, ok := parseMedida(c.Medida)
	if !ok {
		return falha(c, CodigoMedidaInvalida, "medida deve ser g, ml ou unidade")
	}

	if _, err := time.Parse("2006-01-02", c.DataExpiracao); err != nil {
		return falha(c, CodigoDataExpiracaoInvalida, "dataExpiracao deve ser YYYY-MM-DD")
	}

	if falhaPromo := validarPromocao(c); falhaPromo != nil {
		return OfertaValidada{}, falhaPromo
	}

	return OfertaValidada{
		Produto:       produto,
		Marca:         strings.TrimSpace(c.Marca),
		Categorias:    normalizarCategorias(c.Categorias),
		Valor:         c.Valor,
		Quantidade:    c.Quantidade,
		Medida:        medida,
		DataExpiracao: c.DataExpiracao,
		Promocao:      c.Promocao,
	}, nil
}

func falha(c CandidatoOferta, codigo, detalhe string) (OfertaValidada, *FalhaExtracao) {
	return OfertaValidada{}, &FalhaExtracao{
		Codigo:    codigo,
		Detalhe:   detalhe,
		Candidato: c,
	}
}

func parseMedida(s string) (Medida, bool) {
	switch Medida(s) {
	case MedidaG, MedidaML, MedidaUnidade:
		return Medida(s), true
	default:
		return "", false
	}
}

func validarPromocao(c CandidatoOferta) *FalhaExtracao {
	p := c.Promocao
	if p == nil {
		return nil
	}
	if p.ValorPromocional <= 0 {
		return &FalhaExtracao{Codigo: CodigoPromocaoInvalida, Detalhe: "valorPromocional deve ser > 0", Candidato: c}
	}

	levePague := p.Leve != nil || p.Pague != nil
	qtd := p.QuantidadePromocao != nil
	cartao := p.PromocaoCartao != nil

	switch {
	case levePague && !qtd && !cartao:
		if p.Leve == nil || p.Pague == nil || *p.Leve <= 0 || *p.Pague <= 0 {
			return &FalhaExtracao{Codigo: CodigoPromocaoInvalida, Detalhe: "leve/pague inválidos", Candidato: c}
		}
	case qtd && !levePague && !cartao:
		if *p.QuantidadePromocao <= 0 {
			return &FalhaExtracao{Codigo: CodigoPromocaoInvalida, Detalhe: "quantidadePromocao inválida", Candidato: c}
		}
	case cartao && !levePague && !qtd:
		if !*p.PromocaoCartao {
			return &FalhaExtracao{Codigo: CodigoPromocaoInvalida, Detalhe: "promocaoCartao deve ser true", Candidato: c}
		}
	default:
		return &FalhaExtracao{Codigo: CodigoPromocaoInvalida, Detalhe: "promocao deve ser exatamente um dos três formatos", Candidato: c}
	}
	return nil
}

func normalizarCategorias(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}
	for _, raw := range in {
		n := NormalizarRotulo(raw)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}
