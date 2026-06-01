package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
)

// ---- LLM-as-judge scoring ---------------------------------------------------
//
// AutoScore estimates two RAGAS-inspired dimensions for one AI turn:
//
//   - Faithfulness  (0-5): is the answer grounded in the retrieved context?
//   - Relevance     (0-5): does it actually address the user's question?
//
// The scores are stored asynchronously in ImFeedback (source="auto"); failures
// are silently ignored so they never affect answer latency.
//
// Call site: end of AIAnswer, guarded by cfg.AutoEval.

// AutoScore embeds and stores an auto-eval result for one AI answer.
// It is always called in a goroutine; all errors are silently discarded.
func AutoScore(context_ string, query, reply string, sessionID uint64) {
	cfg := LoadAIConfig()
	if cfg.BaseURL == "" || cfg.ChatModel == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	f, r, err := scoreOneTurn(ctx, cfg, context_, query, reply)
	if err != nil {
		return
	}
	fb := &immodel.ImFeedback{
		SessionID:    sessionID,
		Source:       immodel.FeedbackSourceAuto,
		Faithfulness: f,
		Relevance:    r,
		Query:        query,
		Answer:       reply,
	}
	db.DB.Create(fb)
}

// scoreOneTurn calls the LLM to judge faithfulness and relevance.
func scoreOneTurn(ctx context.Context, cfg AIConfig, context_, query, reply string) (faithfulness, relevance float64, err error) {
	prompt := fmt.Sprintf(`你是一个 RAG 系统评估专家。请对下面的 AI 客服回答进行评分，输出 JSON。

【用户问题】
%s

【检索到的参考资料】
%s

【AI 回答】
%s

请输出如下 JSON，faithfulness 和 relevance 各取 0-5 的整数分（5 最高）：
{"faithfulness": <int>, "relevance": <int>, "reason": "<一句话说明>"}
只输出 JSON，不要多余内容。`, query, context_, reply)

	msgs := []chatMessage{
		{Role: "system", Content: "你是一个严格、公正的 RAG 评估者，只输出要求格式的 JSON。"},
		{Role: "user", Content: prompt},
	}
	out, e := cfg.ChatComplete(ctx, msgs)
	if e != nil {
		return 0, 0, e
	}
	// Extract the first {...} block in case the model adds preamble.
	raw := extractJSON(out)
	var result struct {
		Faithfulness int    `json:"faithfulness"`
		Relevance    int    `json:"relevance"`
		Reason       string `json:"reason"`
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return 0, 0, fmt.Errorf("parse score: %w", err)
	}
	return float64(result.Faithfulness), float64(result.Relevance), nil
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end <= start {
		return s
	}
	return s[start : end+1]
}

// ---- Feedback CRUD ----------------------------------------------------------

// SaveFeedback persists user-submitted feedback for an AI answer.
func SaveFeedback(ctx context.Context, fb *immodel.ImFeedback) error {
	return db.DB.WithContext(ctx).Create(fb).Error
}

// ListFeedback returns paginated feedback records, optionally filtered by session.
func ListFeedback(ctx context.Context, sessionID uint64, page, size int) ([]immodel.ImFeedback, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&immodel.ImFeedback{})
	if sessionID > 0 {
		tx = tx.Where("session_id = ?", sessionID)
	}
	var total int64
	tx.Count(&total)
	var list []immodel.ImFeedback
	err := tx.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

// FeedbackStats returns aggregate evaluation metrics.
func FeedbackStats(ctx context.Context) (map[string]any, error) {
	var rows []struct {
		Source       string
		Count        int64
		AvgFaith     *float64
		AvgRelevance *float64
		AvgRating    *float64
	}
	err := db.DB.WithContext(ctx).Raw(`
		SELECT source,
		       COUNT(*) AS count,
		       AVG(NULLIF(faithfulness, 0)) AS avg_faith,
		       AVG(NULLIF(relevance, 0))    AS avg_relevance,
		       AVG(NULLIF(rating, 0))       AS avg_rating
		FROM im_feedbacks
		GROUP BY source
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(rows))
	for _, r := range rows {
		out[r.Source] = map[string]any{
			"count":         r.Count,
			"avg_faith":     r.AvgFaith,
			"avg_relevance": r.AvgRelevance,
			"avg_rating":    r.AvgRating,
		}
	}
	return out, nil
}
