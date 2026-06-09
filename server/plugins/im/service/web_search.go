package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	immodel "github.com/ijry/lyshop/plugins/im/model"
)

type webSearchResult struct {
	Title   string
	Link    string
	Snippet string
}

func retrieveWebSearch(ctx context.Context, cfg AIConfig, session *immodel.ImSession, query string) string {
	if !cfg.WebSearchEnabled {
		return ""
	}
	start := time.Now()
	results, err := searchWeb(ctx, cfg, query)
	if err != nil {
		recordEventBestEffort(ctx, EventInput{
			Event:     immodel.ImEventWebSearch,
			Level:     "warn",
			SessionID: session.ID,
			UserID:    session.UserID,
			Source:    immodel.ImEventSourceAI,
			Success:   false,
			LatencyMS: time.Since(start).Milliseconds(),
			Message:   "联网搜索失败",
			Meta:      map[string]any{"query": query, "error": err.Error()},
		})
		return ""
	}
	recordEventBestEffort(ctx, EventInput{
		Event:     immodel.ImEventWebSearch,
		SessionID: session.ID,
		UserID:    session.UserID,
		Source:    immodel.ImEventSourceAI,
		Success:   true,
		LatencyMS: time.Since(start).Milliseconds(),
		Message:   "联网搜索完成",
		Meta:      map[string]any{"query": query, "count": len(results), "provider": cfg.WebSearchProvider},
	})
	if len(results) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("【联网搜索】\n")
	for i, item := range results {
		fmt.Fprintf(&b, "%d. %s", i+1, item.Title)
		if item.Link != "" {
			fmt.Fprintf(&b, "｜%s", item.Link)
		}
		if item.Snippet != "" {
			fmt.Fprintf(&b, "｜%s", item.Snippet)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func searchWeb(ctx context.Context, cfg AIConfig, query string) ([]webSearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}
	switch strings.ToLower(strings.TrimSpace(cfg.WebSearchProvider)) {
	case "", "serper":
		return searchSerper(ctx, cfg, query)
	default:
		return nil, fmt.Errorf("不支持的联网搜索提供方：%s", cfg.WebSearchProvider)
	}
}

func searchSerper(ctx context.Context, cfg AIConfig, query string) ([]webSearchResult, error) {
	if strings.TrimSpace(cfg.WebSearchAPIKey) == "" {
		return nil, fmt.Errorf("未配置 ai_web_search_api_key")
	}
	endpoint := cfg.WebSearchEndpoint
	if endpoint == "" {
		endpoint = "https://google.serper.dev/search"
	}
	limit := cfg.WebSearchTopK
	if limit <= 0 {
		limit = 3
	}
	if limit > 8 {
		limit = 8
	}
	body := map[string]any{"q": query, "num": limit}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", cfg.WebSearchAPIKey)
	client := cfg.httpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("搜索服务返回 %d", resp.StatusCode)
	}
	var out struct {
		Organic []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"organic"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	results := make([]webSearchResult, 0, len(out.Organic))
	for _, item := range out.Organic {
		if strings.TrimSpace(item.Title) == "" && strings.TrimSpace(item.Snippet) == "" {
			continue
		}
		results = append(results, webSearchResult{
			Title:   trimField(item.Title, 160),
			Link:    trimField(item.Link, 256),
			Snippet: trimField(item.Snippet, 300),
		})
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}
