package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTPDownloader interface {
	Download(ctx context.Context, u string) (io.ReadCloser, int64, error)
}

func NewHTTPDownloader(httpClient *http.Client) HTTPDownloader {
	return &httpDownloader{
		client: httpClient,
	}
}

type httpDownloader struct {
	client *http.Client
}

func (downloader *httpDownloader) Download(ctx context.Context, u string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create a http request: %w", err)
	}
	resp, err := downloader.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send http request: %w", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return resp.Body, 0, errInvalidHTTPStatusCode
	}
	return resp.Body, resp.ContentLength, nil
}
