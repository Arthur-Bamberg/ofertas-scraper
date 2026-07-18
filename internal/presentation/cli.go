package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"ofertas-scraper/internal/application"
	"ofertas-scraper/internal/domain"
	"ofertas-scraper/internal/infra/artefato"
	"ofertas-scraper/internal/infra/extrator"
	"ofertas-scraper/internal/infra/filenamedate"
	"ofertas-scraper/internal/infra/fontehttp"
	"ofertas-scraper/internal/infra/raster"
	"ofertas-scraper/internal/infra/upstash"
)

// Env holds process configuration.
type Env struct {
	UpstashURL   string
	UpstashToken string
	ArtefatoRoot string
	RasterMaxPx  int
	RasterJPEGQ  int
	SeedPath     string
	UseStubExtrator bool
}

func LoadEnv() Env {
	maxPx, _ := strconv.Atoi(os.Getenv("RASTER_MAX_EDGE_PX"))
	jpegQ, _ := strconv.Atoi(os.Getenv("RASTER_JPEG_QUALITY"))
	stub := os.Getenv("EXTRATOR_STUB") == "1" || os.Getenv("GEMINI_API_KEY") == ""
	return Env{
		UpstashURL:      os.Getenv("UPSTASH_REDIS_REST_URL"),
		UpstashToken:    os.Getenv("UPSTASH_REDIS_REST_TOKEN"),
		ArtefatoRoot:    envOr("ARTEFATO_ROOT", "./.data/artefatos"),
		RasterMaxPx:     maxPx,
		RasterJPEGQ:     jpegQ,
		SeedPath:        envOr("SEED_PATH", "./seed/fontes.json"),
		UseStubExtrator: stub,
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

type seedFile struct {
	Mercados []domain.Mercado `json:"mercados"`
	Fontes   []domain.Fonte   `json:"fontes"`
}

// RunSeed upserts Mercados/Fontes from seed JSON into Redis.
func RunSeed(ctx context.Context, env Env) error {
	if env.UpstashURL == "" || env.UpstashToken == "" {
		return fmt.Errorf("UPSTASH_REDIS_REST_URL and UPSTASH_REDIS_REST_TOKEN required")
	}
	raw, err := os.ReadFile(env.SeedPath)
	if err != nil {
		return err
	}
	var seed seedFile
	if err := json.Unmarshal(raw, &seed); err != nil {
		return err
	}
	client := upstash.NewClient(env.UpstashURL, env.UpstashToken, nil)
	mercados := upstash.NewMercadoRepo(client)
	fontes := upstash.NewFonteRepo(client)
	for _, m := range seed.Mercados {
		if m.ID == "" {
			m.ID = domain.MercadoID(application.NewID())
		}
		if err := mercados.Save(ctx, m); err != nil {
			return err
		}
		log.Printf("seed mercado %s (%s)", m.ID, m.Nome)
	}
	for _, f := range seed.Fontes {
		if f.ID == "" {
			f.ID = domain.FonteID(application.NewID())
		}
		if err := fontes.Save(ctx, f); err != nil {
			return err
		}
		log.Printf("seed fonte %s → %s", f.ID, f.URL)
	}
	return nil
}

// RunDaily wires adapters and runs the daily job.
func RunDaily(ctx context.Context, env Env) error {
	if env.UpstashURL == "" || env.UpstashToken == "" {
		return fmt.Errorf("UPSTASH_REDIS_REST_URL and UPSTASH_REDIS_REST_TOKEN required")
	}
	client := upstash.NewClient(env.UpstashURL, env.UpstashToken, nil)

	var ext domain.Extrator
	if env.UseStubExtrator {
		log.Printf("using Extrator stub (Gemini adapter not wired yet)")
	} else {
		log.Printf("GEMINI_API_KEY set but Gemini Extrator not implemented yet; using stub")
	}
	ext = extrator.Stub{Candidatos: nil}

	deps := application.RunDailyJobDeps{
		Fontes:     upstash.NewFonteRepo(client),
		Documentos: upstash.NewDocumentoRepo(client),
		Produtos:   upstash.NewProdutoRepo(client),
		Marcas:     upstash.NewMarcaRepo(client),
		Ofertas:    upstash.NewOfertaRepo(client),
		Falhas:     upstash.NewFalhaRepo(client),
		FonteHTTP:  fontehttp.New(http.DefaultClient),
		Raster:     raster.NewPdftoppm(env.RasterMaxPx, env.RasterJPEGQ),
		Extrator:   ext,
		Artefatos:  artefato.NewLocalStore(env.ArtefatoRoot),
		Dates:      filenamedate.Parser{},
		Log:        log.Default(),
		ExtratorRetryBudget: time.Hour,
	}
	return application.RunDailyJob(ctx, deps)
}
