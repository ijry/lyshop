# IM Customer Service Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enhance the existing LYShop IM plugin with event logs, analytics, Redis-backed WebSocket fanout, and image/file chat messages.

**Architecture:** Extend the current `server/plugins/im` plugin in place. Keep `ImSession`, `ImMessage`, and `ImStaff` as the domain core; add event logging and analytics around the current flow, use Redis only as a cross-instance delivery bus, and reuse the existing storage driver for attachments.

**Tech Stack:** Go 1.26, Gin, GORM, Redis/go-redis, Vue 3, UniApp Vue 3, uview-plus, Vite, docs-site markdown.

---

## File Structure

- Modify `server/plugins/im/model/im.go`: add `MsgTypeFile`, event constants, and `ImEventLog`.
- Modify `server/plugins/im/plugin.go`: migrate `ImEventLog`.
- Create `server/plugins/im/service/event.go`: event recording, event query, analytics aggregation.
- Create `server/plugins/im/service/event_test.go`: service tests for event logging and aggregation.
- Create `server/plugins/im/service/upload.go`: attachment validation and storage-driver upload helper.
- Create `server/plugins/im/service/upload_test.go`: validation and permission tests.
- Modify `server/plugins/im/service/session.go`: record events during session lifecycle and include message `extra` in WS payloads.
- Modify `server/plugins/im/service/ai.go`: record RAG-hit signal from AI answer pipeline.
- Modify `server/plugins/im/service/hub.go`: add optional Redis bus initialization, publish, subscribe, and no-loop delivery.
- Create `server/plugins/im/service/hub_test.go`: Redis envelope no-loop/unit tests with `miniredis`.
- Modify `server/plugins/im/api/front.go`: add user upload route.
- Modify `server/plugins/im/api/admin.go`: add upload, analytics, and logs routes.
- Create `server/plugins/im/api/admin_im_enhancement_test.go`: route and permission coverage.
- Modify `server/plugins/im/plugin.json`: add Admin menu entries for analytics/logs.
- Modify `admin/src/router/index.ts`: add routes for analytics/logs pages.
- Create `admin/src/views/im/Analytics.vue`: dashboard for IM summary/trend.
- Create `admin/src/views/im/EventLogs.vue`: filterable event log page.
- Modify `admin/src/views/im/SessionList.vue`: render attachments and upload files.
- Modify `admin/src/locales/zh-CN.ts` and `admin/src/locales/en.ts`: labels for new views and attachment UI.
- Modify `app/pages/im/chat.vue`, `app/locales/zh-CN.ts`, `app/locales/en.ts`: user attachment send/render.
- Modify `eapp/pages/im/chat.vue`, `eapp/api/im.ts`, `eapp/components/biz/ChatBubble.vue`: merchant attachment send/render.
- Modify `web/src/stores/chat.ts`, `web/src/components/ChatDialog.vue`, `web/src/views/Chat.vue`, `web/src/api/request.ts`: PC/Web attachment send/render.
- Modify `docs-site/docs/api/im.md`: latest IM API, events, analytics, upload, Redis deployment notes.
- Modify `docs/im-api-reference.md` and `docs/im-feature-matrix.md`: internal documentation parity.

## Implementation Tasks

### Task 1: Event Model And Logging Service

**Files:**
- Modify: `server/plugins/im/model/im.go`
- Modify: `server/plugins/im/plugin.go`
- Create: `server/plugins/im/service/event.go`
- Create: `server/plugins/im/service/event_test.go`

- [ ] **Step 1: Add the failing event service test**

Create `server/plugins/im/service/event_test.go`:

```go
package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRecordEventAndListLogs(t *testing.T) {
	testDB := setupImEventTestDB(t)
	ctx := context.Background()

	require.NoError(t, RecordEvent(ctx, EventInput{
		Event:     immodel.ImEventSessionCreated,
		SessionID: 11,
		UserID:    22,
		Source:    immodel.ImEventSourceUser,
		Success:   true,
		Extra:     map[string]any{"mode": "ai"},
	}))

	logs, total, err := ListEventLogs(ctx, EventLogQuery{SessionID: 11, Page: 1, Size: 20})
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	require.Len(t, logs, 1)
	require.Equal(t, immodel.ImEventSessionCreated, logs[0].Event)
	require.Equal(t, uint64(22), logs[0].UserID)
	require.Equal(t, int8(1), logs[0].Success)
	require.Contains(t, logs[0].Extra, `"mode":"ai"`)

	var count int64
	require.NoError(t, testDB.Model(&immodel.ImEventLog{}).Count(&count).Error)
	require.Equal(t, int64(1), count)
}

func setupImEventTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:im_event_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })
	require.NoError(t, gdb.AutoMigrate(&immodel.ImEventLog{}))
	return gdb
}
```

- [ ] **Step 2: Run the failing test**

Run:

```powershell
cd server
go test ./plugins/im/service -run TestRecordEventAndListLogs -count=1
```

Expected: fail because `ImEventLog`, `EventInput`, `RecordEvent`, and `ListEventLogs` do not exist.

- [ ] **Step 3: Add model constants and migration**

In `server/plugins/im/model/im.go`, add:

```go
const (
	MsgTypeFile = "file"

	ImEventSessionCreated = "session_created"
	ImEventMessageSent    = "message_sent"
	ImEventAIReply        = "ai_reply"
	ImEventAIFailed       = "ai_failed"
	ImEventRAGHit         = "rag_hit"
	ImEventToHuman        = "to_human"
	ImEventStaffAccept    = "staff_accept"
	ImEventSessionClose   = "session_close"
	ImEventSessionTransfer = "session_transfer"
	ImEventFileUploaded   = "file_uploaded"

	ImEventSourceUser   = "user"
	ImEventSourceStaff  = "staff"
	ImEventSourceAI     = "ai"
	ImEventSourceSystem = "system"
)

type ImEventLog struct {
	model.Base
	Event     string `gorm:"size:64;not null;index" json:"event"`
	SessionID uint64 `gorm:"not null;default:0;index" json:"session_id"`
	UserID    uint64 `gorm:"not null;default:0;index" json:"user_id"`
	StaffID   uint64 `gorm:"not null;default:0;index" json:"staff_id"`
	MessageID uint64 `gorm:"not null;default:0;index" json:"message_id"`
	Source    string `gorm:"size:16;not null;default:'system';index" json:"source"`
	Success   int8   `gorm:"not null;default:1;index" json:"success"`
	LatencyMS int64  `gorm:"not null;default:0" json:"latency_ms"`
	Extra     string `gorm:"type:json" json:"extra"`
}
```

In `server/plugins/im/plugin.go`, add `&immodel.ImEventLog{}` to `AutoMigrate`.

- [ ] **Step 4: Implement event service**

Create `server/plugins/im/service/event.go`:

```go
package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ijry/lyshop/core/db"
	immodel "github.com/ijry/lyshop/plugins/im/model"
)

type EventInput struct {
	Event     string
	SessionID uint64
	UserID    uint64
	StaffID   uint64
	MessageID uint64
	Source    string
	Success   bool
	LatencyMS int64
	Extra     map[string]any
}

type EventLogQuery struct {
	Event     string
	SessionID uint64
	UserID    uint64
	StaffID   uint64
	Source    string
	Success   *int8
	Page      int
	Size      int
}

func RecordEvent(ctx context.Context, input EventInput) error {
	if strings.TrimSpace(input.Event) == "" {
		return nil
	}
	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = immodel.ImEventSourceSystem
	}
	success := int8(0)
	if input.Success {
		success = 1
	}
	extra := ""
	if len(input.Extra) > 0 {
		raw, _ := json.Marshal(input.Extra)
		extra = string(raw)
	}
	row := &immodel.ImEventLog{
		Event:     input.Event,
		SessionID: input.SessionID,
		UserID:    input.UserID,
		StaffID:   input.StaffID,
		MessageID: input.MessageID,
		Source:    source,
		Success:   success,
		LatencyMS: input.LatencyMS,
		Extra:     extra,
	}
	return db.DB.WithContext(ctx).Create(row).Error
}

func ListEventLogs(ctx context.Context, q EventLogQuery) ([]immodel.ImEventLog, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 || q.Size > 100 {
		q.Size = 20
	}
	tx := db.DB.WithContext(ctx).Model(&immodel.ImEventLog{})
	if q.Event != "" {
		tx = tx.Where("event = ?", q.Event)
	}
	if q.SessionID > 0 {
		tx = tx.Where("session_id = ?", q.SessionID)
	}
	if q.UserID > 0 {
		tx = tx.Where("user_id = ?", q.UserID)
	}
	if q.StaffID > 0 {
		tx = tx.Where("staff_id = ?", q.StaffID)
	}
	if q.Source != "" {
		tx = tx.Where("source = ?", q.Source)
	}
	if q.Success != nil {
		tx = tx.Where("success = ?", *q.Success)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []immodel.ImEventLog
	err := tx.Order("id desc").Offset((q.Page - 1) * q.Size).Limit(q.Size).Find(&list).Error
	return list, total, err
}

type AnalyticsQuery struct {
	From    time.Time
	To      time.Time
	StaffID uint64
}
```

- [ ] **Step 5: Run test and commit**

Run:

```powershell
cd server
go test ./plugins/im/service -run TestRecordEventAndListLogs -count=1
gofmt -w plugins/im/model/im.go plugins/im/plugin.go plugins/im/service/event.go plugins/im/service/event_test.go
go test ./plugins/im/service -run TestRecordEventAndListLogs -count=1
git add plugins/im/model/im.go plugins/im/plugin.go plugins/im/service/event.go plugins/im/service/event_test.go
git commit -m "新增客服事件日志模型" -m "为 IM 插件新增事件日志模型、迁移和基础记录查询服务，为客服报表、审计和后续埋点提供统一事件来源。"
```

Expected: test passes.

### Task 2: Analytics And Admin Log APIs

**Files:**
- Modify: `server/plugins/im/service/event.go`
- Create: `server/plugins/im/api/admin_im_enhancement_test.go`
- Modify: `server/plugins/im/api/admin.go`

- [ ] **Step 1: Add failing admin API tests**

Create `server/plugins/im/api/admin_im_enhancement_test.go` with tests for route registration, analytics aggregation, and log filtering:

```go
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
	now := time.Now()
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
```

- [ ] **Step 2: Run failing tests**

Run:

```powershell
cd server
go test ./plugins/im/api -run TestAdminEnhancement -count=1
```

Expected: fail because routes and analytics function do not exist.

- [ ] **Step 3: Implement analytics aggregation**

Append to `server/plugins/im/service/event.go`:

```go
type AnalyticsResult struct {
	Summary map[string]int64      `json:"summary"`
	Trend   []map[string]any      `json:"trend"`
}

func emptyAnalyticsCounters() map[string]int64 {
	return map[string]int64{
		"sessions": 0, "messages": 0, "ai_replies": 0, "ai_failed": 0,
		"rag_hits": 0, "to_human": 0, "accepts": 0, "closes": 0,
		"transfers": 0, "files": 0,
	}
}

func incrementAnalyticsCounter(counters map[string]int64, event string, count int64) {
	switch event {
	case immodel.ImEventSessionCreated:
		counters["sessions"] += count
	case immodel.ImEventMessageSent:
		counters["messages"] += count
	case immodel.ImEventAIReply:
		counters["ai_replies"] += count
	case immodel.ImEventAIFailed:
		counters["ai_failed"] += count
	case immodel.ImEventRAGHit:
		counters["rag_hits"] += count
	case immodel.ImEventToHuman:
		counters["to_human"] += count
	case immodel.ImEventStaffAccept:
		counters["accepts"] += count
	case immodel.ImEventSessionClose:
		counters["closes"] += count
	case immodel.ImEventSessionTransfer:
		counters["transfers"] += count
	case immodel.ImEventFileUploaded:
		counters["files"] += count
	}
}

func Analytics(ctx context.Context, q AnalyticsQuery) (*AnalyticsResult, error) {
	tx := db.DB.WithContext(ctx).Model(&immodel.ImEventLog{})
	if !q.From.IsZero() {
		tx = tx.Where("created_at >= ?", q.From)
	}
	if !q.To.IsZero() {
		tx = tx.Where("created_at < ?", q.To)
	}
	if q.StaffID > 0 {
		tx = tx.Where("staff_id = ?", q.StaffID)
	}
	var rows []struct {
		Event string
		Count int64
	}
	if err := tx.Select("event, COUNT(*) AS count").Group("event").Scan(&rows).Error; err != nil {
		return nil, err
	}
	summary := emptyAnalyticsCounters()
	for _, r := range rows {
		incrementAnalyticsCounter(summary, r.Event, r.Count)
	}

	trendTx := db.DB.WithContext(ctx).Model(&immodel.ImEventLog{})
	if !q.From.IsZero() {
		trendTx = trendTx.Where("created_at >= ?", q.From)
	}
	if !q.To.IsZero() {
		trendTx = trendTx.Where("created_at < ?", q.To)
	}
	if q.StaffID > 0 {
		trendTx = trendTx.Where("staff_id = ?", q.StaffID)
	}
	var events []immodel.ImEventLog
	if err := trendTx.Order("created_at asc").Find(&events).Error; err != nil {
		return nil, err
	}
	byDay := map[string]map[string]int64{}
	var days []string
	for _, event := range events {
		day := event.CreatedAt.Format("2006-01-02")
		if _, ok := byDay[day]; !ok {
			byDay[day] = emptyAnalyticsCounters()
			days = append(days, day)
		}
		incrementAnalyticsCounter(byDay[day], event.Event, 1)
	}
	trend := make([]map[string]any, 0, len(days))
	for _, day := range days {
		row := map[string]any{"date": day}
		for key, value := range byDay[day] {
			row[key] = value
		}
		trend = append(trend, row)
	}
	return &AnalyticsResult{Summary: summary, Trend: trend}, nil
}
```

- [ ] **Step 4: Add handlers and routes**

In `server/plugins/im/api/admin.go`, register:

```go
g.GET("/im/analytics", middleware.RequirePermission("im:view"), adminAnalytics)
g.GET("/im/logs", middleware.RequirePermission("im:view"), adminListEventLogs)
```

Add handlers:

```go
func adminAnalytics(c *gin.Context) {
	from, _ := time.ParseInLocation("2006-01-02", c.Query("from"), time.Local)
	to, _ := time.ParseInLocation("2006-01-02", c.Query("to"), time.Local)
	if !to.IsZero() {
		to = to.Add(24 * time.Hour)
	}
	staffID, _ := strconv.ParseUint(c.Query("staff_id"), 10, 64)
	result, err := imsvc.Analytics(c.Request.Context(), imsvc.AnalyticsQuery{From: from, To: to, StaffID: staffID})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, result)
}

func adminListEventLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	sessionID, _ := strconv.ParseUint(c.Query("session_id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	staffID, _ := strconv.ParseUint(c.Query("staff_id"), 10, 64)
	var success *int8
	if raw := c.Query("success"); raw != "" {
		v, _ := strconv.ParseInt(raw, 10, 8)
		vv := int8(v)
		success = &vv
	}
	list, total, err := imsvc.ListEventLogs(c.Request.Context(), imsvc.EventLogQuery{
		Event: c.Query("event"), SessionID: sessionID, UserID: userID, StaffID: staffID,
		Source: c.Query("source"), Success: success, Page: page, Size: size,
	})
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, response.PageData{List: list, Total: total, Page: page, Size: size})
}
```

Add `time` import to `admin.go`.

- [ ] **Step 5: Run tests and commit**

Run:

```powershell
cd server
gofmt -w plugins/im/service/event.go plugins/im/api/admin.go plugins/im/api/admin_im_enhancement_test.go
go test ./plugins/im/api -run TestAdminEnhancement -count=1
go test ./plugins/im/service -run TestRecordEventAndListLogs -count=1
git add plugins/im/service/event.go plugins/im/api/admin.go plugins/im/api/admin_im_enhancement_test.go
git commit -m "新增客服报表与事件日志接口" -m "基于 IM 事件日志提供后台报表聚合和日志查询接口，补充权限校验和路由测试。"
```

Expected: tests pass.

### Task 3: Attachment Upload Service And APIs

**Files:**
- Create: `server/plugins/im/service/upload.go`
- Create: `server/plugins/im/service/upload_test.go`
- Modify: `server/plugins/im/api/front.go`
- Modify: `server/plugins/im/api/admin.go`

- [ ] **Step 1: Add failing upload validation test**

Create `server/plugins/im/service/upload_test.go`:

```go
package service

import (
	"testing"

	immodel "github.com/ijry/lyshop/plugins/im/model"
	"github.com/stretchr/testify/require"
)

func TestClassifyAttachment(t *testing.T) {
	img, err := ClassifyAttachment("photo.png", "image/png", 1024)
	require.NoError(t, err)
	require.Equal(t, immodel.MsgTypeImage, img.MessageType)

	doc, err := ClassifyAttachment("policy.pdf", "application/pdf", 1024)
	require.NoError(t, err)
	require.Equal(t, immodel.MsgTypeFile, doc.MessageType)

	_, err = ClassifyAttachment("shell.exe", "application/octet-stream", 1024)
	require.ErrorContains(t, err, "不支持")

	_, err = ClassifyAttachment("huge.png", "image/png", 11<<20)
	require.ErrorContains(t, err, "文件过大")
}
```

- [ ] **Step 2: Run failing upload test**

Run:

```powershell
cd server
go test ./plugins/im/service -run TestClassifyAttachment -count=1
```

Expected: fail because upload helper does not exist.

- [ ] **Step 3: Implement upload helper**

Create `server/plugins/im/service/upload.go`:

```go
package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	driverStorage "github.com/ijry/lyshop/core/driver/storage"
	immodel "github.com/ijry/lyshop/plugins/im/model"
)

const maxIMUploadSize int64 = 10 << 20

type AttachmentInfo struct {
	URL         string `json:"url"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Mime        string `json:"mime"`
	MessageType string `json:"message_type"`
}

var imageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
var fileExts = map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".txt": true, ".csv": true, ".md": true, ".zip": true}

func ClassifyAttachment(filename, mime string, size int64) (*AttachmentInfo, error) {
	if size > maxIMUploadSize {
		return nil, fmt.Errorf("文件过大，最大支持 10MB")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	switch {
	case imageExts[ext] && strings.HasPrefix(strings.ToLower(mime), "image/"):
		return &AttachmentInfo{Name: filepath.Base(filename), Size: size, Mime: mime, MessageType: immodel.MsgTypeImage}, nil
	case fileExts[ext]:
		return &AttachmentInfo{Name: filepath.Base(filename), Size: size, Mime: mime, MessageType: immodel.MsgTypeFile}, nil
	default:
		return nil, fmt.Errorf("不支持的文件类型：%s", ext)
	}
}

func UploadAttachment(ctx context.Context, fh *multipart.FileHeader) (*AttachmentInfo, error) {
	info, err := ClassifyAttachment(fh.Filename, fh.Header.Get("Content-Type"), fh.Size)
	if err != nil {
		return nil, err
	}
	driver, err := driverStorage.Get()
	if err != nil {
		return nil, err
	}
	res, err := driver.Upload(ctx, fh)
	if err != nil {
		return nil, err
	}
	info.URL = res.URL
	info.Path = res.Path
	if info.Mime == "" {
		info.Mime = res.Mime
	}
	if info.Size == 0 {
		info.Size = res.Size
	}
	return info, nil
}
```

- [ ] **Step 4: Add upload routes**

In `front.go`, register `auth.POST("/im/upload", uploadAttachment)` and implement:

```go
func uploadAttachment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID, _ := strconv.ParseUint(c.PostForm("session_id"), 10, 64)
	if sessionID == 0 {
		response.Fail(c, 400, "session_id required")
		return
	}
	session, err := imsvc.GetSession(c.Request.Context(), sessionID)
	if err != nil || session.UserID != userID.(uint64) {
		response.Fail(c, 403, "无权上传到该会话")
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	info, err := imsvc.UploadAttachment(c.Request.Context(), fh)
	if err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	response.OK(c, info)
}
```

Add `GetSession(ctx, sessionID)` helper in `session.go`:

```go
func GetSession(ctx context.Context, sessionID uint64) (*immodel.ImSession, error) {
	var session immodel.ImSession
	err := db.DB.WithContext(ctx).First(&session, sessionID).Error
	return &session, err
}
```

In `admin.go`, register `g.POST("/im/upload", middleware.RequirePermission("im:reply"), adminUploadAttachment)` and implement the same upload flow without user ownership check. Validate that the session exists.

- [ ] **Step 5: Run tests and commit**

Run:

```powershell
cd server
gofmt -w plugins/im/service/upload.go plugins/im/service/upload_test.go plugins/im/service/session.go plugins/im/api/front.go plugins/im/api/admin.go
go test ./plugins/im/service -run TestClassifyAttachment -count=1
go test ./plugins/im/api -run TestAdminEnhancement -count=1
git add plugins/im/service/upload.go plugins/im/service/upload_test.go plugins/im/service/session.go plugins/im/api/front.go plugins/im/api/admin.go
git commit -m "新增客服附件上传接口" -m "复用现有存储驱动为 IM 客服提供用户端和后台附件上传能力，增加类型大小校验和会话权限校验。"
```

Expected: tests pass.

### Task 4: Event Instrumentation And Attachment Message Payloads

**Files:**
- Modify: `server/plugins/im/service/session.go`
- Modify: `server/plugins/im/service/ai.go`
- Modify: `server/plugins/im/api/admin.go`
- Modify: `server/plugins/im/api/front.go`

- [ ] **Step 1: Add failing service behavior test**

Extend `server/plugins/im/service/event_test.go`:

```go
func TestSaveMessageRecordsMessageEvent(t *testing.T) {
	testDB := setupImEventTestDB(t)
	require.NoError(t, testDB.AutoMigrate(&immodel.ImSession{}, &immodel.ImMessage{}))
	ctx := context.Background()
	session := immodel.ImSession{UserID: 7, Mode: immodel.SessionModeHuman, Status: immodel.SessionStatusOngoing}
	require.NoError(t, testDB.Create(&session).Error)

	msg := &immodel.ImMessage{SessionID: session.ID, SenderType: immodel.SenderUser, SenderID: 7, Type: immodel.MsgTypeText, Content: "hello"}
	require.NoError(t, SaveMessage(ctx, msg))

	var count int64
	require.NoError(t, testDB.Model(&immodel.ImEventLog{}).Where("event = ?", immodel.ImEventMessageSent).Count(&count).Error)
	require.Equal(t, int64(1), count)
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./plugins/im/service -run TestSaveMessageRecordsMessageEvent -count=1
```

Expected: fail because `SaveMessage` does not record event.

- [ ] **Step 3: Instrument lifecycle events**

Update `session.go`:

- After new session creation in `GetOrCreateSession`, call `RecordEvent` with `session_created`.
- In `SaveMessage`, after DB create, call `RecordEvent` with `message_sent`; set source from sender type.
- In `SwitchToHuman`, record `to_human`.
- In `AcceptSession`, record `staff_accept`.
- In `TransferSession`, record `session_transfer`.
- In `CloseSession`, record `session_close`.
- In `answerWithAI`, record `ai_reply` on successful reply and `ai_failed` when `AIAnswer` returns error/blank.
- When forwarding WS `msg`, include `extra` from payload into `ImMessage.Extra`, and include `extra` in outgoing frames.

Use best-effort helper:

```go
func recordEventBestEffort(ctx context.Context, input EventInput) {
	if err := RecordEvent(ctx, input); err != nil {
		// keep IM flow available even when event logging fails
		return
	}
}
```

- [ ] **Step 4: Add attachment event on upload handlers**

In upload handlers, after successful upload:

```go
_ = imsvc.RecordEvent(c.Request.Context(), imsvc.EventInput{
	Event:     immodel.ImEventFileUploaded,
	SessionID: sessionID,
	UserID:    session.UserID,
	StaffID:   staffID,
	Source:    immodel.ImEventSourceStaff,
	Success:   true,
	Extra: map[string]any{
		"name": info.Name,
		"size": info.Size,
		"type": info.MessageType,
	},
})
```

For user upload, use `Source: immodel.ImEventSourceUser` and no `StaffID`.

- [ ] **Step 5: Run tests and commit**

Run:

```powershell
cd server
gofmt -w plugins/im/service/session.go plugins/im/service/ai.go plugins/im/api/admin.go plugins/im/api/front.go plugins/im/service/event_test.go
go test ./plugins/im/service -run "TestRecordEventAndListLogs|TestSaveMessageRecordsMessageEvent|TestClassifyAttachment" -count=1
go test ./plugins/im/api -run TestAdminEnhancement -count=1
git add plugins/im/service/session.go plugins/im/service/ai.go plugins/im/api/admin.go plugins/im/api/front.go plugins/im/service/event_test.go
git commit -m "补充客服事件埋点" -m "在 IM 会话、消息、AI 回复、转人工、接入、关闭、转接和附件上传流程中记录事件，支撑后台报表和排障日志。"
```

Expected: tests pass.

### Task 5: Redis WebSocket Fanout

**Files:**
- Modify: `server/plugins/im/service/hub.go`
- Create: `server/plugins/im/service/hub_test.go`

- [ ] **Step 1: Add failing Redis envelope test**

Create `server/plugins/im/service/hub_test.go`:

```go
package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHubEnvelopeIgnoresOwnNode(t *testing.T) {
	h := NewHub()
	h.nodeID = "node-a"
	raw, err := json.Marshal(hubEnvelope{NodeID: "node-a", TargetID: "user_1", Data: json.RawMessage(`{"type":"ping"}`)})
	require.NoError(t, err)
	require.False(t, h.shouldDeliverRemote(raw))
}

func TestHubEnvelopeAcceptsRemoteNode(t *testing.T) {
	h := NewHub()
	h.nodeID = "node-a"
	raw, err := json.Marshal(hubEnvelope{NodeID: "node-b", TargetID: "user_1", Data: json.RawMessage(`{"type":"ping"}`)})
	require.NoError(t, err)
	require.True(t, h.shouldDeliverRemote(raw))
}
```

- [ ] **Step 2: Run failing test**

Run:

```powershell
cd server
go test ./plugins/im/service -run TestHubEnvelope -count=1
```

Expected: fail because `hubEnvelope`, `nodeID`, and `shouldDeliverRemote` do not exist.

- [ ] **Step 3: Extend Hub with Redis bus**

Modify `hub.go`:

- Add imports: `context`, `crypto/rand`, `encoding/hex`, `encoding/json`, `log/slog`, `os`, `time`, and `github.com/ijry/lyshop/core/cache`.
- Extend `Hub` with `nodeID string`, `remote chan *delivery`, `redisEnabled bool`.
- Add `hubEnvelope` type.
- Change `NewHub` to set `nodeID`.
- Add `InitRedisBus(ctx context.Context)` that subscribes `lyshop:im:ws`.
- Add `publishRemote`.
- Change `Send` to local deliver plus publish.
- Add `sendLocal` for remote deliveries to avoid republish.

Core snippets:

```go
const imWSChannel = "lyshop:im:ws"

type hubEnvelope struct {
	NodeID    string          `json:"node_id"`
	TargetID  string          `json:"target_id"`
	Data      json.RawMessage `json:"data"`
	CreatedAt int64          `json:"created_at"`
}

func newNodeID() string {
	var buf [4]byte
	_, _ = rand.Read(buf[:])
	host, _ := os.Hostname()
	return host + "-" + hex.EncodeToString(buf[:])
}
```

`shouldDeliverRemote` should unmarshal and compare `NodeID`.

In `imPlugin.Install()`, after starting Hub:

```go
go imsvc.GlobalHub.Run()
imsvc.GlobalHub.InitRedisBus(context.Background())
```

Add `context` import to `plugin.go`.

- [ ] **Step 4: Run tests and commit**

Run:

```powershell
cd server
gofmt -w plugins/im/service/hub.go plugins/im/service/hub_test.go plugins/im/plugin.go
go test ./plugins/im/service -run TestHubEnvelope -count=1
go test ./plugins/im/service -run "TestRecordEventAndListLogs|TestSaveMessageRecordsMessageEvent|TestClassifyAttachment|TestHubEnvelope" -count=1
git add plugins/im/service/hub.go plugins/im/service/hub_test.go plugins/im/plugin.go
git commit -m "支持客服 WebSocket 跨实例广播" -m "在现有 IM Hub 上增加 Redis Pub/Sub 广播适配，保持本机投递优先并通过节点标识避免消息回环。"
```

Expected: tests pass.

### Task 6: Admin UI Analytics, Logs, And Attachments

**Files:**
- Modify: `server/plugins/im/plugin.json`
- Modify: `admin/src/router/index.ts`
- Create: `admin/src/views/im/Analytics.vue`
- Create: `admin/src/views/im/EventLogs.vue`
- Modify: `admin/src/views/im/SessionList.vue`
- Modify: `admin/src/locales/zh-CN.ts`
- Modify: `admin/src/locales/en.ts`

- [ ] **Step 1: Add menu and router entries**

In `server/plugins/im/plugin.json`, add children:

```json
{ "title": "客服报表", "title_key": "menu.im.analytics", "path": "/im/analytics", "permission": "im:view" },
{ "title": "事件日志", "title_key": "menu.im.logs", "path": "/im/logs", "permission": "im:view" }
```

In `admin/src/router/index.ts`, add:

```ts
{ path: 'im/analytics', name: 'nav.imAnalytics', component: () => import('@/views/im/Analytics.vue'), meta: { titleKey: 'nav.imAnalytics' } },
{ path: 'im/logs', name: 'nav.imLogs', component: () => import('@/views/im/EventLogs.vue'), meta: { titleKey: 'nav.imLogs' } },
```

- [ ] **Step 2: Create Analytics page**

Create `admin/src/views/im/Analytics.vue`:

```vue
<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-slate-900">{{ $t('imAnalytics.title') }}</h1>
        <p class="text-sm text-slate-400 mt-1">{{ $t('imAnalytics.subtitle') }}</p>
      </div>
      <button class="px-4 py-2 rounded-lg bg-red-600 text-white text-sm" @click="load">{{ $t('common.refresh') }}</button>
    </div>
    <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
      <div v-for="item in cards" :key="item.key" class="bg-white rounded-xl border border-slate-100 p-4">
        <div class="text-xs text-slate-400">{{ item.label }}</div>
        <div class="text-2xl font-semibold text-slate-900 mt-2">{{ item.value }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import request from '@/api/request'

const summary = ref<Record<string, number>>({})
const labels: Record<string, string> = {
  sessions: '会话', messages: '消息', ai_replies: 'AI回复', ai_failed: 'AI失败',
  rag_hits: 'RAG命中', to_human: '转人工', accepts: '接入', closes: '关闭',
  transfers: '转接', files: '文件',
}

const cards = computed(() => Object.keys(labels).map((key) => ({
  key,
  label: labels[key],
  value: summary.value[key] || 0,
})))

async function load() {
  const data: any = await request.get('/im/analytics')
  summary.value = data?.summary || {}
}

onMounted(load)
</script>
```

- [ ] **Step 3: Create EventLogs page**

Create `admin/src/views/im/EventLogs.vue`:

```vue
<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-semibold text-slate-900">{{ $t('imLogs.title') }}</h1>
      <button class="px-4 py-2 rounded-lg bg-red-600 text-white text-sm" @click="load">{{ $t('common.search') }}</button>
    </div>
    <div class="bg-white rounded-xl border border-slate-100 p-4 flex gap-3">
      <input v-model="filters.event" class="border rounded-lg px-3 py-2 text-sm" placeholder="event" />
      <input v-model="filters.session_id" class="border rounded-lg px-3 py-2 text-sm" placeholder="session_id" />
    </div>
    <div class="bg-white rounded-xl border border-slate-100 overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-slate-50 text-slate-500">
          <tr><th class="p-3 text-left">ID</th><th class="p-3 text-left">Event</th><th class="p-3 text-left">Session</th><th class="p-3 text-left">Source</th><th class="p-3 text-left">Time</th></tr>
        </thead>
        <tbody>
          <tr v-for="row in list" :key="row.id" class="border-t">
            <td class="p-3">{{ row.id }}</td>
            <td class="p-3">{{ row.event }}</td>
            <td class="p-3">{{ row.session_id }}</td>
            <td class="p-3">{{ row.source }}</td>
            <td class="p-3">{{ row.created_at }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import request from '@/api/request'

const filters = ref({ event: '', session_id: '' })
const list = ref<any[]>([])

async function load() {
  const params: any = {}
  if (filters.value.event) params.event = filters.value.event
  if (filters.value.session_id) params.session_id = filters.value.session_id
  const data: any = await request.get('/im/logs', { params })
  list.value = data?.list || []
}

onMounted(load)
</script>
```

- [ ] **Step 4: Enhance SessionList attachment send/render**

In `admin/src/views/im/SessionList.vue`:

- Add hidden file input and upload button beside send.
- Add `uploading` state.
- Add `sendAttachment(file: File)` that posts `FormData` to `/im/upload`, then sends `/im/sessions/:id/reply` with `type`, `content`, and `extra`.
- Render `m.type === 'image'` as `<img>`.
- Render `m.type === 'file'` as an `<a>`.

Use:

```ts
async function uploadAttachment(file: File) {
  if (!activeSession.value) return
  const form = new FormData()
  form.append('session_id', String(activeSession.value.id))
  form.append('file', file)
  const info: any = await request.post('/im/upload', form, { headers: { 'Content-Type': 'multipart/form-data' } })
  await request.post(`/im/sessions/${activeSession.value.id}/reply`, {
    type: info.message_type,
    content: info.name,
    extra: JSON.stringify({
      file_url: info.url,
      file_path: info.path,
      file_name: info.name,
      file_size: info.size,
      mime: info.mime,
    }),
  })
}
```

Update `adminReply` in `server/plugins/im/api/admin.go` to bind `Extra string` and assign `msg.Extra = req.Extra`, so attachment metadata can be persisted by the same reply endpoint.

- [ ] **Step 5: Add locales and run Admin build**

Add keys in `admin/src/locales/zh-CN.ts` and `en.ts`:

```ts
'menu.im.analytics': '客服报表',
'menu.im.logs': '事件日志',
'nav.imAnalytics': '客服报表',
'nav.imLogs': '事件日志',
'imAnalytics.title': '客服报表',
'imAnalytics.subtitle': '统计客服会话、AI、RAG 与附件消息',
'imLogs.title': '事件日志',
```

Run:

```powershell
cd admin
npm run build
git add ../server/plugins/im/plugin.json src/router/index.ts src/views/im/Analytics.vue src/views/im/EventLogs.vue src/views/im/SessionList.vue src/locales/zh-CN.ts src/locales/en.ts
git commit -m "完善后台客服报表与附件界面" -m "在后台客服中心增加报表和事件日志页面，并增强客服会话页的图片文件消息展示与上传发送能力。"
```

Expected: Admin build succeeds.

### Task 7: App, Eapp, And Web Attachment UI

**Files:**
- Modify: `app/pages/im/chat.vue`
- Modify: `app/locales/zh-CN.ts`
- Modify: `app/locales/en.ts`
- Modify: `eapp/pages/im/chat.vue`
- Modify: `eapp/components/biz/ChatBubble.vue`
- Modify: `web/src/api/request.ts`
- Modify: `web/src/stores/chat.ts`
- Modify: `web/src/components/ChatDialog.vue`
- Modify: `web/src/views/Chat.vue`

- [ ] **Step 1: Update App user chat**

In `app/pages/im/chat.vue`:

- Import `upload` from `@/utils/request`.
- Add upload button using uview-plus in the fixed input bar.
- Use `uni.chooseImage` for App first-phase uploads. File upload on App native/miniprogram can be added later through platform-specific file pickers; Web and Admin cover ordinary files in this phase.
- Upload to `/api/v1/im/upload` with `session_id`.
- Send WS frame:

```ts
payload: {
  msg_type: info.message_type,
  content: info.name,
  extra: {
    file_url: info.url,
    file_path: info.path,
    file_name: info.name,
    file_size: info.size,
    mime: info.mime,
  },
}
```

Render image/file in the message bubble area. Use `uni.previewImage({ urls: [url] })` for image preview.

- [ ] **Step 2: Update Eapp merchant chat**

In `eapp/pages/im/chat.vue`:

- Add `upload` helper if missing in `eapp/utils/request.ts` using the same pattern as `app/utils/request.ts`.
- Add upload button with `up-button`.
- Upload to `/im/upload` because eapp requests normalize to `/admin/api`.
- Send reply via existing `sendImMessage` or direct API with `type` and `extra`.

In `eapp/components/biz/ChatBubble.vue`, render:

```vue
<image v-if="message.type === 'image'" :src="fileUrl(message)" mode="widthFix" class="chat-image" @click="preview(fileUrl(message))" />
<view v-else-if="message.type === 'file'" class="file-card" @click="openFile(fileUrl(message))">{{ fileName(message) }}</view>
```

- [ ] **Step 3: Update Web store and dialog**

In `web/src/api/request.ts`, update upload signature to support extra form fields:

```ts
export async function upload<T = any>(url: string, file: File, fields?: Record<string, string>): Promise<T> {
  if (MOCK_ENABLED) return mockRequest<T>('POST', url, { name: file.name, size: file.size, ...(fields || {}) })
  const form = new FormData()
  Object.entries(fields || {}).forEach(([key, value]) => form.append(key, value))
  form.append('file', file)
  return http.post(url, form, { headers: { 'Content-Type': 'multipart/form-data' } }) as Promise<T>
}
```

In `web/src/stores/chat.ts`, extend message type:

```ts
messages: [] as Array<{ id: number; sender_type: number; content: string; type?: string; extra?: any }>,
```

Add action `sendAttachment(file: File)`:

```ts
async sendAttachment(file: File) {
  if (!this.sessionID) return
  const { upload } = await import('@/api/request')
  const info: any = await upload('/api/v1/im/upload', file, { session_id: String(this.sessionID) })
  const extra = { file_url: info.url, file_path: info.path, file_name: info.name, file_size: info.size, mime: info.mime }
  this.messages.push({ id: Date.now(), sender_type: 1, content: info.name, type: info.message_type, extra })
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'msg', session_id: this.sessionID, payload: { msg_type: info.message_type, content: info.name, extra } }))
  }
}
```

Render attachments in `ChatDialog.vue` and `Chat.vue`.

- [ ] **Step 4: Run front-end builds**

Run:

```powershell
cd web
npm run build
cd ..\app
npm run build:h5
cd ..\eapp
npm run build:h5
```

Expected: all builds succeed. If app/eapp build depends on platform-specific tooling unavailable in the current shell, run `npm run build:h5` and document the exact failure.

- [ ] **Step 5: Commit**

Run:

```powershell
git add app/pages/im/chat.vue app/locales/zh-CN.ts app/locales/en.ts eapp/pages/im/chat.vue eapp/components/biz/ChatBubble.vue eapp/utils/request.ts web/src/api/request.ts web/src/stores/chat.ts web/src/components/ChatDialog.vue web/src/views/Chat.vue
git commit -m "支持多端客服附件消息" -m "为用户端、商家端和 Web 客服会话增加图片文件上传、发送、预览和文件卡片展示，复用 IM 上传接口。"
```

### Task 8: Documentation And Final Verification

**Files:**
- Modify: `docs-site/docs/api/im.md`
- Modify: `docs/im-api-reference.md`
- Modify: `docs/im-feature-matrix.md`

- [ ] **Step 1: Update docs-site IM API**

Update `docs-site/docs/api/im.md` to describe latest architecture:

- User upload: `POST /api/v1/im/upload`.
- Admin upload: `POST /admin/api/im/upload`.
- Admin analytics: `GET /admin/api/im/analytics`.
- Admin logs: `GET /admin/api/im/logs`.
- Message types: `text`、`image`、`file`、`product_card`、`order_card`、`system`.
- Attachment `extra` JSON format.
- Redis deployment: all backend replicas must share the same external Redis for cross-instance WS; embedded Redis is single-instance only.
- Storage impact: uses active storage driver; Docker local storage uses existing `./data/uploads:/app/uploads`.

- [ ] **Step 2: Update internal docs**

Update `docs/im-api-reference.md` and `docs/im-feature-matrix.md` with:

- First-phase feature matrix rows for analytics, event logs, Redis fanout, file/image messages.
- API examples matching implemented routes.
- Permission notes: `im:view` and `im:reply`.

- [ ] **Step 3: Run final verification**

Run:

```powershell
cd server
go test ./plugins/im/... ./core/cache/... ./core/driver/storage/...
cd ..\admin
npm run build
cd ..\web
npm run build
cd ..\eapp
npm test
```

Expected:

- Go tests pass.
- Admin and Web builds pass.
- Eapp tests pass.

If a build/test fails for an environment prerequisite, record the command and exact error in the final handoff.

- [ ] **Step 4: Commit docs**

Run:

```powershell
git add docs-site/docs/api/im.md docs/im-api-reference.md docs/im-feature-matrix.md
git commit -m "更新客服增强功能文档" -m "同步 IM 客服最新接口、附件消息格式、事件日志、报表统计口径和 Redis 多实例部署影响。"
```

## Final Review Checklist

- [ ] `git log -3 --pretty=%B` contains only Chinese commit messages and no `Co-Authored-By`.
- [ ] `git status --short` is clean.
- [ ] `docs-site/docs/api/im.md` includes function description, API changes, and deployment/configuration impact.
- [ ] No anonymous visitor widget code is included in this phase.
- [ ] Existing `/api/v1/im/session`, `/api/v1/im/messages`, `/ws/im`, and admin reply/session routes remain compatible.
