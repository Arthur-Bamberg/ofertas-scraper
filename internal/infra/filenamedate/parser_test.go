package filenamedate_test

import (
	"testing"

	"ofertas-scraper/internal/infra/filenamedate"
)

func TestParser_extraiDataDoFilename(t *testing.T) {
	p := filenamedate.Parser{}
	got, ok := p.Parse("ofertas_atacado_20260725.pdf")
	if !ok || got != "2026-07-25" {
		t.Fatalf("got %q ok=%v", got, ok)
	}
	got, ok = p.Parse("flyer-2026-07-18-final.pdf")
	if !ok || got != "2026-07-18" {
		t.Fatalf("got %q ok=%v", got, ok)
	}
	_, ok = p.Parse("sem-data.pdf")
	if ok {
		t.Fatal("não deveria achar data")
	}
}
