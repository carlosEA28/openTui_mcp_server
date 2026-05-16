package main

import (
	"context"
	"fmt"
	"time"

	"github.com/carlosEA28/openTui_mcp_server/pkg/http"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fetch, err := http.Fetch(ctx, "https://opentui.com/docs/getting-started/")
	if err != nil {
		return
	}

	fmt.Print(string(fetch))
}
