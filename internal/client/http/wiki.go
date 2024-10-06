package http

import (
	"context"
	"fmt"
	"wiki-export/internal/config/clients"
)

type WikiHttpClient interface {
	PagesExportMarkdown(ctx context.Context, pageId int) ([]byte, error)
}

type wikiClient struct {
	http *Client
	cfg  *clients.Wiki
}

func NewWikiHttpClient(cfg *clients.Wiki) WikiHttpClient {
	client := NewClient(cfg.Host, cfg.Port, cfg.Timeout)

	client.GetClient().SetAuthToken(fmt.Sprintf("Token %s:%s", cfg.TokenId, cfg.TokenSecret))

	return &wikiClient{
		http: client,
		cfg:  cfg,
	}
}

func (c *wikiClient) PagesExportMarkdown(ctx context.Context, pageId int) ([]byte, error) {
	path := fmt.Sprintf(c.cfg.Paths.PagesExportMarkdown, pageId)
	resp, err := c.http.GetClient().R().SetContext(ctx).Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to export page with markdown format: %s", resp.Status())
	}

	return resp.Body(), nil
}
