package main

import (
	"context"
	"log"

	blevestore "github.com/carlosEA28/openTui_mcp_server/pkg/lib/bleve"
	"github.com/carlosEA28/openTui_mcp_server/pkg/tools"
	"github.com/carlosEA28/openTui_mcp_server/pkg/types"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const indexPath = "/home/carloseduardo/Desktop/pastaProjetosGit/go/openTUI_mcp/data/index"

func main() {
	store, err := blevestore.Open(indexPath)
	if err != nil {
		log.Fatalf("failed to open index: %v", err)
	}
	defer store.Close()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "opentui-mcp",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_docs",
		Description: "Searches the OpenTUI documentation. Returns the top-k most relevant chunks of text from the docs. Call once with a specific query — do not call repeatedly with similar queries. If results are returned, use them directly to answer.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args types.SearchDocsArgs) (*mcp.CallToolResult, any, error) {
		result, _, err := tools.SearchDocs(ctx, req, args, store)
		return result, nil, err
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
