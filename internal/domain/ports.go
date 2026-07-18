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
	ID   MercadoID
	Nome string
}

type Fonte struct {
	ID                             FonteID
	MercadoID                      MercadoID
	URL                            string
	FiltroNomeDocumento            string // optional; empty = no filter
	FallbackDataExpiracaoFilename  bool   // default false (ADR 0013)
}

type Produto struct {
	ID         ProdutoID
	Nome       string
	NomeNorm   string
	Categorias []string
}

type Marca struct {
	ID       MarcaID
	Nome     string
	NomeNorm string
}

type Documento struct {
	ID          DocumentoID
	FonteID     FonteID
	MercadoID   MercadoID
	Filename    string
	Dia         string // YYYY-MM-DD America/Sao_Paulo discovery day
	Estado      EstadoDocumento
	Atualizado  time.Time
}

// Oferta is the persisted price observation (after match-or-create).
type Oferta struct {
	ID            OfertaID
	DocumentoID   DocumentoID
	ProdutoID     ProdutoID
	MarcaID       *MarcaID // optional (ADR 0015)
	MercadoID     MercadoID
	Valor         float64
	Quantidade    float64
	Medida        Medida
	DataExpiracao string
	Promocao      *Promocao
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
	SaveAll(ctx context.Context, ofertas []Oferta) error
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

type ArtefatoStore interface {
	SavePDF(ctx context.Context, doc Documento, pdf []byte) error
	SaveImages(ctx context.Context, doc Documento, images []PageImage) error
	SaveRawExtrator(ctx context.Context, doc Documento, raw []byte) error
	SaveValidated(ctx context.Context, doc Documento, ofertas []Oferta, falhas []FalhaExtracao) error
}

type FonteClient interface {
	DiscoverPDFs(ctx context.Context, fonte Fonte) (filenames []string, err error)
	DownloadPDF(ctx context.Context, fonte Fonte, filename string) (pdf []byte, err error)
}

type Rasterizer interface {
	Rasterize(ctx context.Context, pdf []byte) ([]PageImage, error)
}

// FilenameDateParser parses dataExpiracao from a Documento filename when Fonte enables fallback.
type FilenameDateParser interface {
	Parse(filename string) (dateYYYYMMDD string, ok bool)
}
