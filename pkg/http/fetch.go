package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func Fetch(ctx context.Context, url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err // colocar error customizado
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err // colocar error customizado
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode)) // colocar erro customizado
	}

	resp.Header.Set("User-Agent", "")

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // colocar error customizado
	}

	return bodyBytes, nil
}
