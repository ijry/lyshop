package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	productmodel "github.com/ijry/lyshop/plugins/product/model"
)

// AIConfig holds the resolved local-LLM settings for the IM plugin.
// Everything is sourced from the plugin's config_kv namespace ("im") so it can
// be edited from the admin 配置中心 without code changes.
type AIConfig struct {
	Enabled       bool
	BaseURL       string // OpenAI-compatible base, e.g. http://localhost:11434/v1 (Ollama) or http://localhost:8000/v1 (vLLM)
	APIKey        string // optional; many local servers ignore it
	ChatModel     string // e.g. qwen2.5:7b
	EmbedModel    string // e.g. bge-m3 / nomic-embed-text；留空则关闭向量召回，退化为关键词召回
	SystemPrompt  string // persona / guardrails
	HumanKeywords []string
	TopK          int
	Temperature   float64
	ProductSearch bool
	Timeout       time.Duration

	// Vector store (Qdrant). When QdrantURL is empty the plugin falls back to
	// in-DB retrieval (and to keyword retrieval when no embed model is set).
	QdrantURL        string  // e.g. http://localhost:6333
	QdrantAPIKey     string  // optional
	QdrantCollection string  // default "im_knowledge"
	ScoreThreshold   float64 // drop hits below this cosine score (0 disables)

	// Retrieval pipeline tuning.
	Hybrid  bool // fuse vector + keyword recall via RRF before rerank
	RecallK int  // candidate pool size per recall channel (defaults to 4×TopK)

	// Reranker (cross-encoder). When RerankURL is set, candidates are reordered
	// by a /rerank service (Cohere / Jina / TEI compatible) and trimmed to TopK.
	RerankURL    string
	RerankAPIKey string
	RerankModel  string

	// Query rewriting. Mode: "" (off) | "rewrite" (simple LLM rewrite) |
	// "hyde" (Hypothetical Document Embeddings) | "multi" (N variants + RRF).
	QueryRewrite  string
	QueryRewriteN int // number of variants for "multi" mode (default 3)

	// Evaluation. When AutoEval is true, AIAnswer scores its own output using
	// the LLM-as-judge pattern and stores the result in ImFeedback.
	AutoEval bool
}

const aiConfigPlugin = "im"

// loadCfg reads a single config_kv value for the im plugin.
func loadCfg(key, def string) string {
	var cfg struct{ Value string }
	err := db.DB.Table("configs").Select("value").
		Where("plugin = ? AND key = ?", aiConfigPlugin, key).
		Scan(&cfg).Error
	if err == nil {
		if v := strings.TrimSpace(cfg.Value); v != "" {
			return v
		}
	}
	return def
}

func loadBool(key string, def bool) bool {
	v := strings.ToLower(loadCfg(key, ""))
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "on" || v == "yes"
}

func loadInt(key string, def int) int {
	if v, err := strconv.Atoi(loadCfg(key, "")); err == nil {
		return v
	}
	return def
}

func loadFloat(key string, def float64) float64 {
	if v, err := strconv.ParseFloat(loadCfg(key, ""), 64); err == nil {
		return v
	}
	return def
}

// LoadAIConfig resolves the current AI configuration.
func LoadAIConfig() AIConfig {
	kw := loadCfg("ai_human_keywords", "人工,转人工,人工客服,真人,客服人工")
	var keywords []string
	for _, k := range strings.Split(kw, ",") {
		if k = strings.TrimSpace(k); k != "" {
			keywords = append(keywords, k)
		}
	}
	return AIConfig{
		Enabled:       loadBool("ai_enabled", false),
		BaseURL:       strings.TrimRight(loadCfg("ai_base_url", "http://localhost:11434/v1"), "/"),
		APIKey:        loadCfg("ai_api_key", ""),
		ChatModel:     loadCfg("ai_chat_model", "qwen2.5:7b"),
		EmbedModel:    loadCfg("ai_embed_model", ""),
		SystemPrompt:  loadCfg("ai_system_prompt", defaultSystemPrompt),
		HumanKeywords: keywords,
		TopK:          loadInt("ai_top_k", 3),
		Temperature:   loadFloat("ai_temperature", 0.3),
		ProductSearch: loadBool("ai_product_search", true),
		Timeout:       time.Duration(loadInt("ai_timeout_sec", 30)) * time.Second,

		QdrantURL:        strings.TrimRight(loadCfg("ai_qdrant_url", ""), "/"),
		QdrantAPIKey:     loadCfg("ai_qdrant_api_key", ""),
		QdrantCollection: loadCfg("ai_qdrant_collection", "im_knowledge"),
		ScoreThreshold:   loadFloat("ai_score_threshold", 0),

		Hybrid:  loadBool("ai_hybrid", false),
		RecallK: loadInt("ai_recall_k", 0),

		RerankURL:    strings.TrimRight(loadCfg("ai_rerank_url", ""), "/"),
		RerankAPIKey: loadCfg("ai_rerank_api_key", ""),
		RerankModel:  loadCfg("ai_rerank_model", ""),

		QueryRewrite:  loadCfg("ai_query_rewrite", ""),
		QueryRewriteN: loadInt("ai_query_rewrite_n", 3),

		AutoEval: loadBool("ai_auto_eval", false),
	}
}

const defaultSystemPrompt = "你是本商城的智能客服助手。请基于提供的【知识库】和【商品信息】用中文简洁、礼貌地回答用户问题。" +
	"若资料中没有相关信息，请如实说明并建议用户输入“人工”转接人工客服，不要编造价格、库存或政策。"

// AIEnabled reports whether AI first-line service is turned on.
func AIEnabled() bool { return LoadAIConfig().Enabled }

// IsHumanRequest reports whether the user's text asks for a human agent.
func IsHumanRequest(cfg AIConfig, text string) bool {
	t := strings.ToLower(strings.TrimSpace(text))
	if t == "" {
		return false
	}
	for _, k := range cfg.HumanKeywords {
		if strings.Contains(t, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

// ---- OpenAI-compatible HTTP calls ----------------------------------------

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (cfg AIConfig) httpClient() *http.Client {
	to := cfg.Timeout
	if to <= 0 {
		to = 30 * time.Second
	}
	return &http.Client{Timeout: to}
}

func (cfg AIConfig) post(ctx context.Context, path string, body any, out any) error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("未配置大模型服务地址")
	}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BaseURL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}
	resp, err := cfg.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("大模型服务返回 %d: %s", resp.StatusCode, strings.TrimSpace(buf.String()))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// ChatComplete calls the chat completions endpoint and returns the reply text.
func (cfg AIConfig) ChatComplete(ctx context.Context, msgs []chatMessage) (string, error) {
	body := map[string]any{
		"model":       cfg.ChatModel,
		"messages":    msgs,
		"temperature": cfg.Temperature,
		"stream":      false,
	}
	var out struct {
		Choices []struct {
			Message chatMessage `json:"message"`
		} `json:"choices"`
	}
	if err := cfg.post(ctx, "/chat/completions", body, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("大模型未返回内容")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

// Embed calls the embeddings endpoint and returns the vector for one input.
func (cfg AIConfig) Embed(ctx context.Context, text string) ([]float64, error) {
	if cfg.EmbedModel == "" {
		return nil, fmt.Errorf("未配置向量模型")
	}
	body := map[string]any{"model": cfg.EmbedModel, "input": text}
	var out struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := cfg.post(ctx, "/embeddings", body, &out); err != nil {
		return nil, err
	}
	if len(out.Data) == 0 || len(out.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("向量服务未返回结果")
	}
	return out.Data[0].Embedding, nil
}

// ---- RAG retrieval --------------------------------------------------------

func cosine(a, b []float64) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

// retrieveKnowledge returns the most relevant knowledge entries for a query
// using a recall → (fuse) → rerank pipeline:
//
//  1. Recall a candidate pool (size RecallK, default 4×TopK):
//     - vector recall: Qdrant ANN, else in-memory cosine over embedded rows;
//     - keyword recall (always when Hybrid, or as the sole channel otherwise).
//  2. When Hybrid and both channels produced results, fuse them with RRF.
//  3. When a reranker is configured, cross-encode the pool against the query
//     and reorder; finally trim to TopK.
//
// Every stage degrades gracefully: no embed model → keyword only; reranker
// down → pre-rerank order; empty pool → nil.
func retrieveKnowledge(ctx context.Context, cfg AIConfig, query string) []immodel.ImKnowledge {
	// Multi-query rewrites fan out retrieval across N variants; this mode takes
	// full ownership of the pipeline and returns here.
	if strings.ToLower(strings.TrimSpace(cfg.QueryRewrite)) == "multi" {
		if hits := retrieveKnowledgeMulti(ctx, cfg, query); hits != nil {
			return hits
		}
		// fall through to single-query path if multi produced nothing
	}

	topK := cfg.TopK
	if topK <= 0 {
		topK = 3
	}
	recallK := cfg.RecallK
	if recallK <= 0 {
		recallK = topK * 4
	}
	if recallK < topK {
		recallK = topK
	}

	vecEntries := recallVector(ctx, cfg, query, recallK)

	// Without hybrid and without a reranker, the vector order is already final.
	if !cfg.Hybrid && cfg.RerankURL == "" {
		if vecEntries != nil {
			return trimEntries(vecEntries, topK)
		}
		// fall through to keyword-only path below
	}

	// Build the candidate pool.
	var pool []immodel.ImKnowledge
	if cfg.Hybrid {
		kwEntries := recallKeyword(ctx, query, recallK)
		switch {
		case len(vecEntries) > 0 && len(kwEntries) > 0:
			pool = fuseEntries(vecEntries, kwEntries)
		case len(vecEntries) > 0:
			pool = vecEntries
		default:
			pool = kwEntries
		}
	} else if vecEntries != nil {
		pool = vecEntries
	} else {
		pool = recallKeyword(ctx, query, recallK)
	}

	if len(pool) == 0 {
		return nil
	}
	return cfg.rerankEntries(ctx, query, pool, topK)
}

// recallVector returns up to `limit` entries by vector similarity, or nil when
// no vector channel is available (so callers can fall back to keyword recall).
func recallVector(ctx context.Context, cfg AIConfig, query string, limit int) []immodel.ImKnowledge {
	// 1. Qdrant ANN.
	if store := VectorStoreFor(cfg); store != nil && cfg.EmbedModel != "" {
		if hits := qdrantRetrieve(ctx, cfg, store, query, limit); hits != nil {
			return hits
		}
	}
	// 2. In-memory cosine over embedded rows.
	if cfg.EmbedModel != "" {
		qvec, err := cfg.Embed(ctx, query)
		if err != nil {
			return nil
		}
		var all []immodel.ImKnowledge
		db.DB.WithContext(ctx).Where("status = 1").Order("sort asc, id asc").Find(&all)
		type scored struct {
			k     immodel.ImKnowledge
			score float64
		}
		var ranked []scored
		for _, k := range all {
			if len(k.Embedding) == 0 {
				continue
			}
			var vec []float64
			if json.Unmarshal(k.Embedding, &vec) != nil {
				continue
			}
			ranked = append(ranked, scored{k, cosine(qvec, vec)})
		}
		if len(ranked) == 0 {
			return nil
		}
		sort.Slice(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
		out := make([]immodel.ImKnowledge, 0, limit)
		for i := 0; i < len(ranked) && i < limit; i++ {
			if ranked[i].score <= 0 || ranked[i].score < cfg.ScoreThreshold {
				break
			}
			out = append(out, ranked[i].k)
		}
		return out
	}
	return nil
}

// recallKeyword returns up to `limit` entries by keyword (token-overlap) score.
func recallKeyword(ctx context.Context, query string, limit int) []immodel.ImKnowledge {
	var all []immodel.ImKnowledge
	db.DB.WithContext(ctx).Where("status = 1").Order("sort asc, id asc").Find(&all)
	if len(all) == 0 {
		return nil
	}
	return keywordRankKnowledge(all, query, limit)
}

// fuseEntries fuses two ranked entry lists with RRF and returns deduped entries
// in fused order.
func fuseEntries(a, b []immodel.ImKnowledge) []immodel.ImKnowledge {
	byID := make(map[uint64]immodel.ImKnowledge, len(a)+len(b))
	idsOf := func(list []immodel.ImKnowledge) []uint64 {
		out := make([]uint64, len(list))
		for i, e := range list {
			out[i] = e.ID
			byID[e.ID] = e
		}
		return out
	}
	fused := fuseRRF(idsOf(a), idsOf(b))
	out := make([]immodel.ImKnowledge, 0, len(fused))
	for _, id := range fused {
		out = append(out, byID[id])
	}
	return out
}

// trimEntries returns at most topK entries.
func trimEntries(entries []immodel.ImKnowledge, topK int) []immodel.ImKnowledge {
	if len(entries) > topK {
		return entries[:topK]
	}
	return entries
}

// qdrantRetrieve embeds the query, searches Qdrant, and loads the matching
// rows from the DB preserving Qdrant's relevance order. Returns nil on any
// error (so the caller can fall back); returns a possibly-empty slice when the
// search succeeds.
func qdrantRetrieve(ctx context.Context, cfg AIConfig, store VectorStore, query string, topK int) []immodel.ImKnowledge {
	qvec, err := cfg.Embed(ctx, query)
	if err != nil {
		return nil
	}
	hits, err := store.Search(ctx, qvec, topK)
	if err != nil {
		return nil
	}
	ids := make([]uint64, 0, len(hits))
	for _, h := range hits {
		if cfg.ScoreThreshold > 0 && h.Score < cfg.ScoreThreshold {
			continue
		}
		ids = append(ids, h.ID)
	}
	if len(ids) == 0 {
		return []immodel.ImKnowledge{}
	}
	var rows []immodel.ImKnowledge
	db.DB.WithContext(ctx).Where("id IN ? AND status = 1", ids).Find(&rows)
	// Reorder rows to match Qdrant's ranking.
	byID := make(map[uint64]immodel.ImKnowledge, len(rows))
	for _, r := range rows {
		byID[r.ID] = r
	}
	out := make([]immodel.ImKnowledge, 0, len(ids))
	for _, id := range ids {
		if r, ok := byID[id]; ok {
			out = append(out, r)
		}
	}
	return out
}

func keywordRankKnowledge(all []immodel.ImKnowledge, query string, topK int) []immodel.ImKnowledge {
	tokens := tokenize(query)
	if len(tokens) == 0 {
		return nil
	}
	type scored struct {
		k     immodel.ImKnowledge
		score int
	}
	var ranked []scored
	for _, k := range all {
		hay := strings.ToLower(k.Title + " " + k.Content + " " + k.Tags)
		s := 0
		for _, t := range tokens {
			if strings.Contains(hay, t) {
				s++
			}
		}
		if s > 0 {
			ranked = append(ranked, scored{k, s})
		}
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	out := make([]immodel.ImKnowledge, 0, topK)
	for i := 0; i < len(ranked) && i < topK; i++ {
		out = append(out, ranked[i].k)
	}
	return out
}

// tokenize splits a query into lowercase tokens. For CJK text (no spaces) it
// also emits 2-gram character windows so keyword overlap still works.
func tokenize(q string) []string {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	add := func(s string) {
		s = strings.TrimSpace(s)
		if len([]rune(s)) >= 2 && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	for _, w := range strings.FieldsFunc(q, func(r rune) bool {
		return r == ' ' || r == ',' || r == '，' || r == '。' || r == '?' || r == '？' || r == '!' || r == '！' || r == '、'
	}) {
		add(w)
		runes := []rune(w)
		for i := 0; i+1 < len(runes); i++ {
			add(string(runes[i : i+2]))
		}
	}
	return out
}

// retrieveProducts returns on-shelf products relevant to the query for the
// "商品信息分析" capability.
func retrieveProducts(ctx context.Context, cfg AIConfig, query string) []productmodel.Product {
	if !cfg.ProductSearch {
		return nil
	}
	tokens := tokenize(query)
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) > 8 { // cap to keep the OR query bounded
		tokens = tokens[:8]
	}
	tx := db.DB.WithContext(ctx).Model(&productmodel.Product{}).Where("status = 1")
	cond := db.DB
	for i, t := range tokens {
		like := "%" + t + "%"
		c := db.DB.Where("title LIKE ? OR subtitle LIKE ?", like, like)
		if i == 0 {
			cond = c
		} else {
			cond = cond.Or("title LIKE ? OR subtitle LIKE ?", like, like)
		}
	}
	var list []productmodel.Product
	tx.Where(cond).Order("sales desc").Limit(5).Find(&list)
	return list
}

// ---- Answer assembly ------------------------------------------------------

// AIAnswer builds a RAG prompt from the knowledge base and product catalog,
// includes recent conversation history, and returns the model's reply.
//
// The caller MUST have already persisted the user's message (so it appears as
// the final turn in recentHistory); userText is used for KB/product retrieval.
func AIAnswer(ctx context.Context, session *immodel.ImSession, userText string) (string, error) {
	cfg := LoadAIConfig()
	if !cfg.Enabled {
		return "", fmt.Errorf("AI 客服未启用")
	}

	// Rewrite the retrieval query (not the generation query) if configured.
	retrievalQuery := applyQueryRewrite(ctx, cfg, userText)

	var ctxParts []string
	if kb := retrieveKnowledge(ctx, cfg, retrievalQuery); len(kb) > 0 {
		var b strings.Builder
		b.WriteString("【知识库】\n")
		for i, k := range kb {
			fmt.Fprintf(&b, "%d. %s：%s\n", i+1, k.Title, strings.TrimSpace(k.Content))
		}
		ctxParts = append(ctxParts, b.String())
	}
	if ps := retrieveProducts(ctx, cfg, retrievalQuery); len(ps) > 0 {
		var b strings.Builder
		b.WriteString("【商品信息】\n")
		for i, p := range ps {
			stock := "有货"
			if p.Stock <= 0 {
				stock = "暂时缺货"
			}
			fmt.Fprintf(&b, "%d. %s｜价格 ¥%.2f｜%s｜已售%d", i+1, p.Title, p.Price, stock, p.Sales)
			if strings.TrimSpace(p.Subtitle) != "" {
				fmt.Fprintf(&b, "｜%s", strings.TrimSpace(p.Subtitle))
			}
			b.WriteString("\n")
		}
		ctxParts = append(ctxParts, b.String())
	}

	msgs := []chatMessage{{Role: "system", Content: cfg.SystemPrompt}}
	if len(ctxParts) > 0 {
		msgs = append(msgs, chatMessage{
			Role:    "system",
			Content: "以下是可参考的资料，请优先依据它们回答：\n\n" + strings.Join(ctxParts, "\n"),
		})
	}
	// Recent conversation history for context (last few turns). The latest
	// user message is already persisted, so it is the final turn here.
	hist := recentHistory(ctx, session.ID, 8)
	if len(hist) == 0 { // safety net if persistence raced
		hist = []chatMessage{{Role: "user", Content: userText}}
	}
	msgs = append(msgs, hist...)

	reply, err := cfg.ChatComplete(ctx, msgs)
	if err != nil {
		return "", err
	}

	// Async LLM-as-judge auto-eval: score faithfulness + relevance and store.
	if cfg.AutoEval {
		context := strings.Join(ctxParts, "\n")
		go AutoScore(context, userText, reply, session.ID)
	}

	return reply, nil
}

// recentHistory returns up to `limit` most recent messages of a session as
// chat messages (user / assistant roles), oldest first.
func recentHistory(ctx context.Context, sessionID uint64, limit int) []chatMessage {
	var list []immodel.ImMessage
	db.DB.WithContext(ctx).
		Where("session_id = ? AND type = ?", sessionID, immodel.MsgTypeText).
		Order("id desc").Limit(limit).Find(&list)
	out := make([]chatMessage, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		m := list[i]
		role := "user"
		if m.SenderType == immodel.SenderAI || m.SenderType == immodel.SenderStaff {
			role = "assistant"
		} else if m.SenderType == immodel.SenderSystem {
			continue
		}
		out = append(out, chatMessage{Role: role, Content: m.Content})
	}
	return out
}

// ReindexKnowledge re-embeds all enabled knowledge entries. Safe to call when
// no embed model is configured (it just marks entries un-indexed).
// ReindexKnowledge re-embeds all knowledge entries and rebuilds the vector
// store. When Qdrant is configured the collection is reset and every enabled
// entry is upserted; the DB embedding column is kept in sync as a local cache /
// fallback. Returns the number of entries successfully (re)indexed.
func ReindexKnowledge(ctx context.Context) (int, error) {
	cfg := LoadAIConfig()
	var all []immodel.ImKnowledge
	db.DB.WithContext(ctx).Find(&all)
	if cfg.EmbedModel == "" {
		db.DB.WithContext(ctx).Model(&immodel.ImKnowledge{}).Where("1 = 1").
			Updates(map[string]any{"embedding": nil, "indexed": 0})
		return 0, fmt.Errorf("未配置向量模型，已退化为关键词召回")
	}

	store := VectorStoreFor(cfg)
	dimReady := false // collection ensured once we know the vector dimension

	done := 0
	for _, k := range all {
		vec, err := cfg.Embed(ctx, embedInput(k))
		if err != nil {
			continue
		}
		raw, _ := json.Marshal(vec)
		db.DB.WithContext(ctx).Model(&immodel.ImKnowledge{}).Where("id = ?", k.ID).
			Updates(map[string]any{"embedding": raw, "indexed": 1})

		if store != nil {
			if !dimReady {
				if err := store.Reset(ctx, len(vec)); err != nil {
					return done, fmt.Errorf("重建向量库失败：%w", err)
				}
				dimReady = true
			}
			// Only enabled entries are searchable; still upsert disabled ones
			// with status payload so toggling status later just needs an update.
			if err := store.Upsert(ctx, k.ID, vec, k.Status); err != nil {
				return done, fmt.Errorf("写入向量库失败：%w", err)
			}
		}
		done++
	}
	return done, nil
}

// embedInput builds the text fed to the embedding model for one entry.
func embedInput(k immodel.ImKnowledge) string {
	return k.Title + "\n" + k.Content + "\n" + k.Tags
}

// EmbedKnowledgeEntry embeds a single entry (best-effort, called on create/
// update) and, when Qdrant is configured, upserts it into the vector store.
func EmbedKnowledgeEntry(ctx context.Context, id uint64) {
	cfg := LoadAIConfig()
	if cfg.EmbedModel == "" {
		return
	}
	var k immodel.ImKnowledge
	if err := db.DB.WithContext(ctx).First(&k, id).Error; err != nil {
		return
	}
	vec, err := cfg.Embed(ctx, embedInput(k))
	if err != nil {
		return
	}
	raw, _ := json.Marshal(vec)
	db.DB.WithContext(ctx).Model(&immodel.ImKnowledge{}).Where("id = ?", id).
		Updates(map[string]any{"embedding": raw, "indexed": 1})

	if store := VectorStoreFor(cfg); store != nil {
		// Ensure collection exists (idempotent) sized to this vector, then upsert.
		if err := store.EnsureCollection(ctx, len(vec)); err == nil {
			store.Upsert(ctx, id, vec, k.Status)
		}
	}
}

// RemoveKnowledgeVector deletes one entry's point from the vector store
// (best-effort; no-op when Qdrant is not configured).
func RemoveKnowledgeVector(ctx context.Context, id uint64) {
	cfg := LoadAIConfig()
	if store := VectorStoreFor(cfg); store != nil {
		store.Delete(ctx, id)
	}
}

// TestAIConnection performs a minimal chat round-trip to validate config.
func TestAIConnection(ctx context.Context) (string, error) {
	cfg := LoadAIConfig()
	if cfg.BaseURL == "" || cfg.ChatModel == "" {
		return "", fmt.Errorf("请先配置服务地址与对话模型")
	}
	return cfg.ChatComplete(ctx, []chatMessage{
		{Role: "system", Content: "你是一个连通性测试助手。"},
		{Role: "user", Content: "请回复：连接正常"},
	})
}
