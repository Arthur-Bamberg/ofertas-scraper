package upstash_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"ofertas-scraper/internal/domain"
	"ofertas-scraper/internal/infra/upstash"
)

func TestMercadoRepo_SaveGetList(t *testing.T) {
	store := map[string]string{}
	sets := map[string]map[string]struct{}{}
	var mu sync.Mutex

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cmd []any
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			t.Fatal(err)
		}
		mu.Lock()
		defer mu.Unlock()
		op := cmd[0].(string)
		switch op {
		case "SET":
			store[cmd[1].(string)] = cmd[2].(string)
			writeResult(w, "OK")
		case "GET":
			v, ok := store[cmd[1].(string)]
			if !ok {
				writeRaw(w, "null")
				return
			}
			writeResult(w, v)
		case "SADD":
			key := cmd[1].(string)
			if sets[key] == nil {
				sets[key] = map[string]struct{}{}
			}
			sets[key][cmd[2].(string)] = struct{}{}
			writeResult(w, 1)
		case "SMEMBERS":
			key := cmd[1].(string)
			var members []string
			for m := range sets[key] {
				members = append(members, m)
			}
			b, _ := json.Marshal(members)
			writeRaw(w, string(b))
		default:
			t.Fatalf("unexpected op %s", op)
		}
	}))
	defer srv.Close()

	client := upstash.NewClient(srv.URL, "token", srv.Client())
	repo := upstash.NewMercadoRepo(client)
	ctx := context.Background()

	m := domain.Mercado{ID: "m1", Nome: "Atacadão"}
	if err := repo.Save(ctx, m); err != nil {
		t.Fatal(err)
	}
	got, err := repo.Get(ctx, "m1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Nome != "Atacadão" {
		t.Fatalf("got %+v", got)
	}
	list, err := repo.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list=%v err=%v", list, err)
	}
}

func TestOfertaRepo_SaveAllReplacesIncludingEmpty(t *testing.T) {
	store := map[string]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cmd []any
		_ = json.NewDecoder(r.Body).Decode(&cmd)
		switch cmd[0].(string) {
		case "SET":
			store[cmd[1].(string)] = cmd[2].(string)
			writeResult(w, "OK")
		case "GET":
			v, ok := store[cmd[1].(string)]
			if !ok {
				writeRaw(w, "null")
				return
			}
			writeResult(w, v)
		}
	}))
	defer srv.Close()

	repo := upstash.NewOfertaRepo(upstash.NewClient(srv.URL, "t", srv.Client()))
	ctx := context.Background()
	docID := domain.DocumentoID("d1")

	err := repo.SaveAll(ctx, docID, []domain.Oferta{{ID: "o1", DocumentoID: docID, Valor: 1.5}})
	if err != nil {
		t.Fatal(err)
	}
	list, _ := repo.ListByDocumento(ctx, docID)
	if len(list) != 1 {
		t.Fatalf("want 1 got %d", len(list))
	}
	if err := repo.SaveAll(ctx, docID, nil); err != nil {
		t.Fatal(err)
	}
	list, _ = repo.ListByDocumento(ctx, docID)
	if len(list) != 0 {
		t.Fatalf("want empty after replace, got %d", len(list))
	}
}

func writeResult(w http.ResponseWriter, v any) {
	b, _ := json.Marshal(map[string]any{"result": v})
	w.Write(b)
}

func writeRaw(w http.ResponseWriter, resultJSON string) {
	w.Write([]byte(`{"result":` + resultJSON + `}`))
}
