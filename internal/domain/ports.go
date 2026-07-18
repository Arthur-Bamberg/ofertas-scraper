package domain

import (
	"context"
	"time"
)

type MercadoID string
type FonteID string
type ProdutoID string
type MarcaID string
type DocumentoID string
type OfertaID string

type Mercado struct {
	ID   MercadoID `json:"id"`
	Nome string    `json:"nome"`
}

type Fonte struct {
	ID                            FonteID   `json:"id"`
	MercadoID                     MercadoID `json:"mercadoId"`
	URL                           string    `json:"url"`
	FiltroNomeDocumento           string    `json:"filtroNomeDocumento"`
	FallbackDataExpiracaoFilename bool      `json:"fallbackDataExpiracaoFilename"`
}

type Produto struct {
	ID         ProdutoID `json:"id"`
	Nome       string    `json:"nome"`
	NomeNorm   string    `json:"nomeNorm"`
	Categorias []string  `json:"categorias"`
}

type Marca struct {
	ID       MarcaID `json:"id"`
	Nome     string  `json:"nome"`
	NomeNorm string  `json:"nomeNorm"`
}

type Documento struct {
	ID         DocumentoID     `json:"id"`
	FonteID    FonteID         `json:"fonteId"`
	MercadoID  MercadoID       `json:"mercadoId"`
	Filename   string          `json:"filename"`
	Dia        string          `json:"dia"`
	Estado     EstadoDocumento `json:"estado"`
	UltimoErro string          `json:"ultimoErro,omitempty"`
	Atualizado time.Time       `json:"atualizado"`
}

// Oferta is the persisted price observation (after match-or-create).
type Oferta struct {
	ID            OfertaID  `json:"id"`
	DocumentoID   DocumentoID `json:"documentoId"`
	ProdutoID     ProdutoID `json:"produtoId"`
	MarcaID       *MarcaID  `json:"marcaId,omitempty"`
	MercadoID     MercadoID `json:"mercadoId"`
	Valor         float64   `json:"valor"`
	Quantidade    float64   `json:"quantidade"`
	Medida        Medida    `json:"medida"`
	DataExpiracao string    `json:"dataExpiracao"`
	Promocao      *Promocao `json:"promocao,omitempty"`
}

type MercadoRepository interface {
	List(ctx context.Context) ([]Mercado, error)
	Save(ctx context.Context, m Mercado) error
	Get(ctx context.Context, id MercadoID) (Mercado, error)
}

type FonteRepository interface {
	List(ctx context.Context) ([]Fonte, error)
	Save(ctx context.Context, f Fonte) error
}

type ProdutoRepository interface {
	GetByNomeNorm(ctx context.Context, nomeNorm string) (Produto, bool, error)
	Save(ctx context.Context, p Produto) error
}

type MarcaRepository interface {
	GetByNomeNorm(ctx context.Context, nomeNorm string) (Marca, bool, error)
	Save(ctx context.Context, m Marca) error
}

type DocumentoRepository interface {
	GetByIdentity(ctx context.Context, fonteID FonteID, filename, dia string) (Documento, bool, error)
	Save(ctx context.Context, d Documento) error
}

type OfertaRepository interface {
	SaveAll(ctx context.Context, documentoID DocumentoID, ofertas []Oferta) error
	ListByDocumento(ctx context.Context, documentoID DocumentoID) ([]Oferta, error)
}

type FalhaExtracaoRepository interface {
	SaveAll(ctx context.Context, documentoID DocumentoID, falhas []FalhaExtracao) error
}

type PageImage struct {
	Page int
	JPEG []byte
}

type Extrator interface {
	Extract(ctx context.Context, images []PageImage) (candidatos []CandidatoOferta, raw []byte, err error)
}

// PDFDescoberto is a PDF link found on a Fonte page (ADR 0020).
type PDFDescoberto struct {
	Filename string // last path segment — Documento identity
	URL      string // absolute download URL
}

type ArtefatoStore interface {
	SavePDF(ctx context.Context, doc Documento, tentativa string, pdf []byte) error
	SaveImages(ctx context.Context, doc Documento, tentativa string, images []PageImage) error
	SaveRawExtrator(ctx context.Context, doc Documento, tentativa string, raw []byte) error
	SaveValidated(ctx context.Context, doc Documento, tentativa string, ofertas []Oferta, falhas []FalhaExtracao) error
}

type FonteClient interface {
	DiscoverPDFs(ctx context.Context, fonte Fonte) ([]PDFDescoberto, error)
	DownloadPDF(ctx context.Context, url string) (pdf []byte, err error)
}

type Rasterizer interface {
	Rasterize(ctx context.Context, pdf []byte) ([]PageImage, error)
}

// FilenameDateParser parses dataExpiracao from a Documento filename when Fonte enables fallback.
type FilenameDateParser interface {
	Parse(filename string) (dateYYYYMMDD string, ok bool)
}
