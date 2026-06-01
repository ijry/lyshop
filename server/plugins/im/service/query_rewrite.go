package service

import (
	"context"
	"fmt"
	"strings"

	immodel "github.com/ijry/lyshop/plugins/im/model"
)

// applyQueryRewrite transforms the raw user query into a retrieval-optimised
// form based on the configured mode. The original text is always used for the
// final LLM generation turn; this only affects what is sent to the vector
// store / keyword scorer.
//
// Modes:
//   - "" / "none"  — return query unchanged (no overhead).
//   - "rewrite"    — ask the LLM to expand colloquial or abbreviated queries
//     into a clean, complete Chinese sentence.
//   - "hyde"       — ask the LLM to draft a hypothetical answer and use that
//     as the embedding text (Gao et al. 2022).
//   - "multi"      — generate QueryRewriteN variants, retrieve for each, fuse
//     with RRF. Returns the original query but stores the
//     multi-retrieve results inside the pipeline via
//     retrieveKnowledgeMulti.
//
// Failures silently fall back to the original query so the pipeline never
// breaks because of a rewrite error.
func applyQueryRewrite(ctx context.Context, cfg AIConfig, query string) string {
	switch strings.ToLower(strings.TrimSpace(cfg.QueryRewrite)) {
	case "rewrite":
		if q := rewriteQuery(ctx, cfg, query); q != "" {
			return q
		}
	case "hyde":
		if doc := hydeQuery(ctx, cfg, query); doc != "" {
			return doc
		}
	}
	// "multi" is handled by the retrieval layer directly (it needs to fan out
	// several retrieval calls and fuse them); the query itself is unchanged here.
	return query
}

// rewriteQuery asks the LLM to clean up the query into a canonical form.
func rewriteQuery(ctx context.Context, cfg AIConfig, query string) string {
	prompt := "请将下面这句客户咨询改写为一句完整、清晰的中文问题，保留原意，不要添加额外信息，只输出改写后的问题。\n\n原始问题：" + query
	msgs := []chatMessage{
		{Role: "system", Content: "你是一个查询改写助手，只输出改写后的问题，不做其他解释。"},
		{Role: "user", Content: prompt},
	}
	out, err := cfg.ChatComplete(ctx, msgs)
	if err != nil || strings.TrimSpace(out) == "" {
		return ""
	}
	return strings.TrimSpace(out)
}

// hydeQuery asks the LLM to draft a hypothetical answer, whose embedding
// typically lies closer to real answer documents than the question itself.
func hydeQuery(ctx context.Context, cfg AIConfig, query string) string {
	prompt := "请为下面这个客服问题写一段简短的假设性回答（约100字以内），用于帮助检索相关文档，不需要准确，只需要语义接近即可。\n\n问题：" + query
	msgs := []chatMessage{
		{Role: "system", Content: "你是一个文档检索助手，只输出假设性回答文本，不做其他解释。"},
		{Role: "user", Content: prompt},
	}
	out, err := cfg.ChatComplete(ctx, msgs)
	if err != nil || strings.TrimSpace(out) == "" {
		return ""
	}
	return strings.TrimSpace(out)
}

// multiQueryVariants generates QueryRewriteN syntactically distinct variants
// of the query using the LLM. Each variant should emphasise a different aspect
// of the intent. Returns at most n non-empty variants.
func multiQueryVariants(ctx context.Context, cfg AIConfig, query string) []string {
	n := cfg.QueryRewriteN
	if n <= 0 {
		n = 3
	}
	prompt := fmt.Sprintf("请为下面这个客服问题生成 %d 种不同的改写版本（每行一个，不加编号，不做解释）。\n\n原始问题：%s", n, query)
	msgs := []chatMessage{
		{Role: "system", Content: "你是一个查询扩展助手，每行输出一个改写问题，不做其他输出。"},
		{Role: "user", Content: prompt},
	}
	out, err := cfg.ChatComplete(ctx, msgs)
	if err != nil || strings.TrimSpace(out) == "" {
		return nil
	}
	var variants []string
	for _, line := range strings.Split(out, "\n") {
		if v := strings.TrimSpace(line); v != "" && v != query {
			variants = append(variants, v)
		}
		if len(variants) >= n {
			break
		}
	}
	return variants
}

// retrieveKnowledgeMulti retrieves candidates for each query variant and fuses
// them with RRF. Only used when QueryRewrite == "multi".
func retrieveKnowledgeMulti(ctx context.Context, cfg AIConfig, query string) []immodel.ImKnowledge {
	variants := multiQueryVariants(ctx, cfg, query)
	if len(variants) == 0 {
		return nil
	}
	// Always include the original query as one of the sources.
	queries := append([]string{query}, variants...)

	topK := cfg.TopK
	if topK <= 0 {
		topK = 3
	}
	recallK := cfg.RecallK
	if recallK <= 0 {
		recallK = topK * 4
	}

	byID := map[uint64]immodel.ImKnowledge{}
	lists := make([][]uint64, 0, len(queries))
	for _, q := range queries {
		hits := recallVector(ctx, cfg, q, recallK)
		if hits == nil {
			hits = recallKeyword(ctx, q, recallK)
		}
		ids := make([]uint64, 0, len(hits))
		for _, h := range hits {
			byID[h.ID] = h
			ids = append(ids, h.ID)
		}
		if len(ids) > 0 {
			lists = append(lists, ids)
		}
	}
	if len(lists) == 0 {
		return nil
	}
	fused := fuseRRF(lists...)
	out := make([]immodel.ImKnowledge, 0, topK)
	for _, id := range fused {
		if k, ok := byID[id]; ok {
			out = append(out, k)
			if len(out) >= recallK {
				break
			}
		}
	}
	return cfg.rerankEntries(ctx, query, out, topK)
}
