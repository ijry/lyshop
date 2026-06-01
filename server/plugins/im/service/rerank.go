package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	immodel "github.com/ijry/lyshop/plugins/im/model"
)

// rrfK is the standard Reciprocal Rank Fusion constant. 60 is the value from
// the original Cormack et al. paper and the de-facto industry default.
const rrfK = 60.0

// fuseRRF combines several ranked candidate lists (each ordered best→worst) into
// one ranking via Reciprocal Rank Fusion: score(d) = Σ 1/(k + rank_i(d)).
// Lists may overlap or be disjoint; an item's contributions across lists add up.
// Returns knowledge IDs ordered by fused score, best first.
func fuseRRF(lists ...[]uint64) []uint64 {
	score := map[uint64]float64{}
	order := []uint64{} // first-seen order, for stable tie-breaking
	seen := map[uint64]bool{}
	for _, list := range lists {
		for rank, id := range list {
			score[id] += 1.0 / (rrfK + float64(rank+1))
			if !seen[id] {
				seen[id] = true
				order = append(order, id)
			}
		}
	}
	sort.SliceStable(order, func(i, j int) bool {
		return score[order[i]] > score[order[j]]
	})
	return order
}

// rerank reorders candidate documents against the query using an external
// cross-encoder reranker, returning the indices into `docs` ordered best→worst.
// On any error it returns nil so the caller keeps the pre-rerank order.
//
// The request/response shape follows the Cohere / Jina / HuggingFace TEI
// "/rerank" convention: {model, query, documents:[...]} →
// {results:[{index, relevance_score}, ...]}.
func (cfg AIConfig) rerank(ctx context.Context, query string, docs []string) []int {
	if cfg.RerankURL == "" || len(docs) == 0 {
		return nil
	}
	body := map[string]any{
		"query":     query,
		"documents": docs,
		"top_n":     len(docs),
	}
	if cfg.RerankModel != "" {
		body["model"] = cfg.RerankModel
	}
	raw, _ := json.Marshal(body)

	to := cfg.Timeout
	if to <= 0 {
		to = 30 * time.Second
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.RerankURL, bytes.NewReader(raw))
	if err != nil {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.RerankAPIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.RerankAPIKey)
	}
	resp, err := (&http.Client{Timeout: to}).Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil
	}
	var out struct {
		Results []struct {
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relevance_score"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil || len(out.Results) == 0 {
		return nil
	}
	// Results may already be sorted, but don't rely on it.
	sort.SliceStable(out.Results, func(i, j int) bool {
		return out.Results[i].RelevanceScore > out.Results[j].RelevanceScore
	})
	idx := make([]int, 0, len(out.Results))
	for _, r := range out.Results {
		if r.Index >= 0 && r.Index < len(docs) {
			idx = append(idx, r.Index)
		}
	}
	return idx
}

// rerankEntries reorders knowledge entries with the cross-encoder and trims to
// topK. Returns the input unchanged (trimmed) when reranking is unavailable.
func (cfg AIConfig) rerankEntries(ctx context.Context, query string, entries []immodel.ImKnowledge, topK int) []immodel.ImKnowledge {
	if len(entries) == 0 {
		return entries
	}
	if cfg.RerankURL != "" {
		docs := make([]string, len(entries))
		for i, e := range entries {
			docs[i] = rerankDoc(e)
		}
		if order := cfg.rerank(ctx, query, docs); len(order) > 0 {
			out := make([]immodel.ImKnowledge, 0, topK)
			for _, i := range order {
				out = append(out, entries[i])
				if len(out) >= topK {
					break
				}
			}
			return out
		}
	}
	if len(entries) > topK {
		return entries[:topK]
	}
	return entries
}

// rerankDoc builds the text representation of an entry passed to the reranker.
func rerankDoc(e immodel.ImKnowledge) string {
	s := strings.TrimSpace(e.Title + "。" + e.Content)
	if t := strings.TrimSpace(e.Tags); t != "" {
		s += " [" + t + "]"
	}
	return s
}

// TestRerankConnection performs a minimal rerank round-trip to validate config.
func TestRerankConnection(ctx context.Context) (string, error) {
	cfg := LoadAIConfig()
	if cfg.RerankURL == "" {
		return "", fmt.Errorf("未配置重排服务地址")
	}
	order := cfg.rerank(ctx, "测试", []string{"无关内容", "测试文档"})
	if len(order) == 0 {
		return "", fmt.Errorf("重排服务无有效返回")
	}
	return "连接正常", nil
}
