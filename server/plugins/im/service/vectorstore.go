package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// VectorStore abstracts an external vector database. Only the operations the
// IM knowledge base needs are exposed. Implementations must be safe for
// concurrent use.
type VectorStore interface {
	// EnsureCollection creates the collection if missing (idempotent), sized to
	// the given vector dimension and cosine distance.
	EnsureCollection(ctx context.Context, dim int) error
	// Upsert stores/overwrites one point keyed by knowledge ID, with the vector
	// and a small payload used for filtering (status) and post-load (none else).
	Upsert(ctx context.Context, id uint64, vector []float64, status int8) error
	// Delete removes a point by knowledge ID (no error if absent).
	Delete(ctx context.Context, id uint64) error
	// Search returns up to topK knowledge IDs most similar to the query vector,
	// filtered to enabled (status=1) points, each with its similarity score.
	Search(ctx context.Context, vector []float64, topK int) ([]ScoredID, error)
	// Reset drops and recreates the collection (used by full reindex).
	Reset(ctx context.Context, dim int) error
}

// ScoredID pairs a knowledge ID with its similarity score (cosine, higher=closer).
type ScoredID struct {
	ID    uint64
	Score float64
}

// qdrantStore talks to Qdrant over its REST API. No SDK dependency is needed —
// the handful of endpoints we use are simple JSON calls, matching the style of
// the OpenAI-compatible client in ai.go.
type qdrantStore struct {
	baseURL    string
	apiKey     string
	collection string
	timeout    time.Duration
}

// newQdrantStore builds a store from config. Returns nil when no base URL is
// configured, signalling callers to fall back to in-DB / keyword retrieval.
func newQdrantStore(cfg AIConfig) *qdrantStore {
	if strings.TrimSpace(cfg.QdrantURL) == "" {
		return nil
	}
	col := cfg.QdrantCollection
	if col == "" {
		col = "im_knowledge"
	}
	to := cfg.Timeout
	if to <= 0 {
		to = 30 * time.Second
	}
	return &qdrantStore{
		baseURL:    strings.TrimRight(cfg.QdrantURL, "/"),
		apiKey:     cfg.QdrantAPIKey,
		collection: col,
		timeout:    to,
	}
}

// VectorStoreFor returns a configured VectorStore, or nil if Qdrant is not set.
func VectorStoreFor(cfg AIConfig) VectorStore {
	if s := newQdrantStore(cfg); s != nil {
		return s
	}
	return nil
}

func (q *qdrantStore) do(ctx context.Context, method, path string, body any, out any) error {
	var reader *bytes.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, q.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("api-key", q.apiKey)
	}
	client := &http.Client{Timeout: q.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("qdrant %s %s 返回 %d: %s", method, path, resp.StatusCode, strings.TrimSpace(buf.String()))
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// collectionExists reports whether the collection is already present.
func (q *qdrantStore) collectionExists(ctx context.Context) bool {
	err := q.do(ctx, http.MethodGet, "/collections/"+q.collection, nil, nil)
	return err == nil
}

func (q *qdrantStore) EnsureCollection(ctx context.Context, dim int) error {
	if dim <= 0 {
		return fmt.Errorf("无效的向量维度")
	}
	if q.collectionExists(ctx) {
		return nil
	}
	return q.createCollection(ctx, dim)
}

func (q *qdrantStore) createCollection(ctx context.Context, dim int) error {
	body := map[string]any{
		"vectors": map[string]any{"size": dim, "distance": "Cosine"},
	}
	return q.do(ctx, http.MethodPut, "/collections/"+q.collection, body, nil)
}

func (q *qdrantStore) Reset(ctx context.Context, dim int) error {
	// Delete is idempotent in Qdrant (404 tolerated by treating as ok below).
	_ = q.do(ctx, http.MethodDelete, "/collections/"+q.collection, nil, nil)
	return q.createCollection(ctx, dim)
}

func (q *qdrantStore) Upsert(ctx context.Context, id uint64, vector []float64, status int8) error {
	body := map[string]any{
		"points": []map[string]any{{
			"id":      id,
			"vector":  vector,
			"payload": map[string]any{"status": status},
		}},
	}
	return q.do(ctx, http.MethodPut, "/collections/"+q.collection+"/points?wait=true", body, nil)
}

func (q *qdrantStore) Delete(ctx context.Context, id uint64) error {
	body := map[string]any{"points": []uint64{id}}
	return q.do(ctx, http.MethodPost, "/collections/"+q.collection+"/points/delete?wait=true", body, nil)
}

func (q *qdrantStore) Search(ctx context.Context, vector []float64, topK int) ([]ScoredID, error) {
	if topK <= 0 {
		topK = 3
	}
	body := map[string]any{
		"vector": vector,
		"limit":  topK,
		"filter": map[string]any{
			"must": []map[string]any{{
				"key":   "status",
				"match": map[string]any{"value": 1},
			}},
		},
	}
	var out struct {
		Result []struct {
			ID    uint64  `json:"id"`
			Score float64 `json:"score"`
		} `json:"result"`
	}
	if err := q.do(ctx, http.MethodPost, "/collections/"+q.collection+"/points/search", body, &out); err != nil {
		return nil, err
	}
	res := make([]ScoredID, 0, len(out.Result))
	for _, r := range out.Result {
		res = append(res, ScoredID{ID: r.ID, Score: r.Score})
	}
	return res, nil
}
