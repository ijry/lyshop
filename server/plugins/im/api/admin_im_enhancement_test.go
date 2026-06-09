package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type imAPIResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func TestAdminEnhancementRoutesAndAnalytics(t *testing.T) {
	router, gdb := setupImAdminRouter(t)
	now := time.Date(2026, 6, 9, 10, 0, 0, 0, time.Local)
	rows := []immodel.ImEventLog{
		{Event: immodel.ImEventSessionCreated, SessionID: 1, UserID: 10, Source: immodel.ImEventSourceUser, Success: 1},
		{Event: immodel.ImEventMessageSent, SessionID: 1, UserID: 10, Source: immodel.ImEventSourceUser, Success: 1},
		{Event: immodel.ImEventAIReply, SessionID: 1, Source: immodel.ImEventSourceAI, Success: 1},
		{Event: immodel.ImEventRAGHit, SessionID: 1, Source: immodel.ImEventSourceAI, Success: 1},
		{Event: immodel.ImEventToHuman, SessionID: 1, UserID: 10, Source: immodel.ImEventSourceUser, Success: 1},
	}
	for i := range rows {
		rows[i].CreatedAt = now
		require.NoError(t, gdb.Create(&rows[i]).Error)
	}

	resp := doImAdminReq(t, router, http.MethodGet, "/admin/im/analytics?from=2026-06-01&to=2026-06-30", "", "*")
	require.Equal(t, 0, resp.Code)
	var payload struct {
		Summary map[string]int64 `json:"summary"`
		Trend   []map[string]any `json:"trend"`
	}
	require.NoError(t, json.Unmarshal(resp.Data, &payload))
	require.Equal(t, int64(1), payload.Summary["sessions"])
	require.Equal(t, int64(1), payload.Summary["messages"])
	require.Equal(t, int64(1), payload.Summary["ai_replies"])
	require.Equal(t, int64(1), payload.Summary["rag_hits"])
	require.Equal(t, int64(1), payload.Summary["to_human"])
	require.Len(t, payload.Trend, 1)
	require.Equal(t, "2026-06-09", payload.Trend[0]["date"])

	logResp := doImAdminReq(t, router, http.MethodGet, "/admin/im/logs?event=to_human", "", "*")
	require.Equal(t, 0, logResp.Code)
	require.Contains(t, string(logResp.Data), immodel.ImEventToHuman)
}

func TestAdminEnhancementPermissionDenied(t *testing.T) {
	router, _ := setupImAdminRouter(t)
	resp := doImAdminReq(t, router, http.MethodGet, "/admin/im/analytics", "", "im:reply")
	require.Equal(t, 403, resp.Code)
}

func setupImAdminRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:im_admin_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(&immodel.ImEventLog{}))
	r := gin.New()
	admin := r.Group("/admin", func(c *gin.Context) {
		raw := strings.TrimSpace(c.GetHeader("X-Test-Perms"))
		if raw == "" {
			c.Set("perms", []string{"*"})
		} else {
			c.Set("perms", strings.Split(raw, ","))
		}
		c.Set("user_id", uint64(99))
		c.Next()
	})
	RegisterAdminRoutes(admin)
	return r, gdb
}

func doImAdminReq(t *testing.T, router *gin.Engine, method, path, body, perms string) imAPIResp {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Perms", perms)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var out imAPIResp
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	return out
}
