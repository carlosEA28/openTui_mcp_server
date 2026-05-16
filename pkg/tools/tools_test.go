package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	blevestore "github.com/carlosEA28/openTui_mcp_server/pkg/lib/bleve"
	"github.com/carlosEA28/openTui_mcp_server/pkg/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newTestStore(t *testing.T) *blevestore.Store {
	t.Helper()
	indexPath := filepath.Join(t.TempDir(), "index")
	store, err := blevestore.Open(indexPath)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})
	return store
}

func indexDocs(t *testing.T, store *blevestore.Store, count int, text string) {
	t.Helper()
	for i := 0; i < count; i++ {
		chunk := types.Chunk{
			Title: fmt.Sprintf("Doc %d", i),
			Text:  text,
		}
		if err := store.Add(fmt.Sprintf("id-%d", i), fmt.Sprintf("https://example.com/%d", i), chunk); err != nil {
			t.Fatalf("add doc %d: %v", i, err)
		}
	}
}

func newToolRequest() *mcp.CallToolRequest {
	return &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{Name: "search_docs"},
	}
}

func TestSearchDocsReturnsResultsForValidQuery(t *testing.T) {
	store := newTestStore(t)
	indexDocs(t, store, 1, "install guide")

	result, output, err := SearchDocs(context.Background(), newToolRequest(), types.SearchDocsArgs{Query: "install"}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatalf("expected non-nil result")
	}
	if output.Query != "install" {
		t.Fatalf("expected query to be preserved, got %q", output.Query)
	}
	if len(output.Results) == 0 {
		t.Fatalf("expected at least one result")
	}
}

func TestSearchDocsTrimsWhitespaceFromQuery(t *testing.T) {
	store := newTestStore(t)
	indexDocs(t, store, 1, "getting started")

	_, output, err := SearchDocs(context.Background(), newToolRequest(), types.SearchDocsArgs{Query: "  getting started  "}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Query != "getting started" {
		t.Fatalf("expected trimmed query, got %q", output.Query)
	}
	if len(output.Results) == 0 {
		t.Fatalf("expected results for trimmed query")
	}
}

func TestSearchDocsDefaultsLimitWhenNonPositive(t *testing.T) {
	store := newTestStore(t)
	indexDocs(t, store, 10, "install")

	_, output, err := SearchDocs(context.Background(), newToolRequest(), types.SearchDocsArgs{Query: "install", Limit: 0}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(output.Results))
	}
}

func TestSearchDocsCapsLimitAtTwenty(t *testing.T) {
	store := newTestStore(t)
	indexDocs(t, store, 25, "openTUI server")

	_, output, err := SearchDocs(context.Background(), newToolRequest(), types.SearchDocsArgs{Query: "openTUI", Limit: 50}, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(output.Results) != 20 {
		t.Fatalf("expected 20 results, got %d", len(output.Results))
	}
}

func TestSearchDocsRejectsEmptyQuery(t *testing.T) {
	store := newTestStore(t)

	_, _, err := SearchDocs(context.Background(), newToolRequest(), types.SearchDocsArgs{Query: "  "}, store)
	if err == nil {
		t.Fatalf("expected error for empty query")
	}
	if !strings.Contains(err.Error(), "query is required") {
		t.Fatalf("expected query validation error, got %v", err)
	}
}

func TestSearchDocsRejectsMissingParams(t *testing.T) {
	store := newTestStore(t)

	_, _, err := SearchDocs(context.Background(), &mcp.CallToolRequest{}, types.SearchDocsArgs{Query: "install"}, store)
	if err == nil {
		t.Fatalf("expected error for missing params")
	}
	if !strings.Contains(err.Error(), "missing tool params") {
		t.Fatalf("expected missing params error, got %v", err)
	}
}
