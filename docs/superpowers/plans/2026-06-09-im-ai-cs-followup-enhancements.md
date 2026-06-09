# IM AI-CS Follow-Up Enhancements Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the five selected AI-CS-inspired enhancements to the existing LYShop IM plugin: visitor context, web embed script, richer log center, live typing draft, and optional web search.

**Architecture:** Extend the current `server/plugins/im` domain in place. `ImSession` remains the conversation aggregate, WebSocket keeps the existing frame envelope, `ImEventLog` becomes the operational log center, and AI web search is disabled by default through IM config keys.

**Tech Stack:** Go 1.26, Gin, GORM, Vue 3, Pinia, UniApp Vue 3, uview-plus, Vite, docs-site markdown.

---

## File Structure

- Modify `server/plugins/im/model/im.go`: add visitor context fields to `ImSession`, add log level/category/trace fields to `ImEventLog`, and add new event constants.
- Modify `server/plugins/im/service/session.go`: add visitor context upsert logic, forward `typing_draft`/`typing_stop`, include session context in staff assignment frames, and preserve frame compatibility.
- Modify `server/plugins/im/service/event.go`: accept/query log level/category/trace fields and include log-center-friendly defaults.
- Modify `server/plugins/im/service/ai.go`: add disabled-by-default web search config and append search snippets to AI context when enabled.
- Create `server/plugins/im/service/web_search.go`: provider interface plus Serper-compatible HTTP provider.
- Modify `server/plugins/im/api/front.go`: accept visitor context on current `/api/v1/im/session`.
- Modify `server/plugins/im/api/admin.go`: add log filters for level/category/trace.
- Create `web/public/im-widget.js`: embeddable iframe widget script targeting the existing Web chat page.
- Modify `web/src/stores/chat.ts`, `web/src/components/ChatDialog.vue`, `web/src/views/Chat.vue`: send visitor context and live typing draft frames.
- Modify `admin/src/views/im/SessionList.vue` and `eapp/pages/im/chat.vue`: show visitor context and incoming typing draft.
- Modify `admin/src/views/im/EventLogs.vue`: expose richer log filters/details.
- Modify `docs-site/docs/api/im.md`, `docs/im-api-reference.md`, and `docs/im-feature-matrix.md`: document latest interfaces, frames, widget, logs, and config.

## Implementation Tasks

### Task 1: Visitor Context On Existing Session API

**Files:**
- Modify: `server/plugins/im/model/im.go`
- Modify: `server/plugins/im/service/session.go`
- Modify: `server/plugins/im/api/front.go`
- Modify: `admin/src/views/im/SessionList.vue`

- [ ] Add nullable/string fields to `ImSession`: `VisitorID`, `VisitorIP`, `VisitorLocation`, `VisitorBrowser`, `VisitorOS`, `VisitorLanguage`, `VisitorReferrer`, `VisitorURL`, `VisitorDevice`, `VisitorExtra`.
- [ ] Add `SessionContextInput` and `GetOrCreateSessionWithContext(ctx, userID, input)`; keep `GetOrCreateSession(ctx, userID)` as a wrapper for compatibility.
- [ ] In `front.go`, bind optional query/body context on `GET /api/v1/im/session`, infer IP from request, and call `GetOrCreateSessionWithContext`.
- [ ] In `assign` frames sent to staff, include a `visitor` payload with safe context fields.
- [ ] Render visitor context in Admin session detail without changing existing session list behavior.

### Task 2: Web Embed Script

**Files:**
- Create: `web/public/im-widget.js`
- Modify: `web/src/views/Chat.vue`
- Modify: `web/src/stores/chat.ts`

- [ ] Add `window.LYShopIMWidget.init(options)` with options: `baseUrl`, `token`, `position`, `title`, `launcherText`, `width`, `height`, `context`.
- [ ] Script creates a fixed launcher button and an iframe pointing to `/chat?embed=1&token=...`.
- [ ] Use `postMessage` to send visitor context to the iframe after it loads.
- [ ] In Web chat, detect `embed=1`, remove duplicate launcher/chrome, and merge posted context before calling `/api/v1/im/session`.

### Task 3: Log Center Enhancements

**Files:**
- Modify: `server/plugins/im/model/im.go`
- Modify: `server/plugins/im/service/event.go`
- Modify: `server/plugins/im/api/admin.go`
- Modify: `admin/src/views/im/EventLogs.vue`

- [ ] Extend `ImEventLog` with `Level`, `Category`, `TraceID`, `Message`, and `Meta`.
- [ ] Update `EventInput` and `RecordEvent` to default `level=info`, `category=im`, `trace_id` from input, `message` from input, and `meta` from structured data.
- [ ] Update `EventLogQuery` and `/admin/api/im/logs` to filter by `level`, `category`, `trace_id`, and keyword over `event/message`.
- [ ] Update EventLogs page with filters and detail display for message/meta.

### Task 4: Real-Time Typing Draft

**Files:**
- Modify: `server/plugins/im/service/session.go`
- Modify: `web/src/stores/chat.ts`
- Modify: `web/src/components/ChatDialog.vue`
- Modify: `web/src/views/Chat.vue`
- Modify: `admin/src/views/im/SessionList.vue`
- Modify: `eapp/pages/im/chat.vue`

- [ ] Add WS handling for `typing_draft` and `typing_stop`; forward only to the opposite side of the active session.
- [ ] Payload format: `{sender_type, sender_id, draft, updated_at}`; server truncates draft to 500 runes.
- [ ] Web user chat sends debounced `typing_draft` while editing and `typing_stop` on send/clear.
- [ ] Admin and Eapp staff chat show the latest user draft and hide it on `typing_stop`.
- [ ] Staff draft forwarding is supported by the same frame shape for future UX without requiring another backend change.

### Task 5: Optional Web Search For AI

**Files:**
- Modify: `server/plugins/im/service/ai.go`
- Create: `server/plugins/im/service/web_search.go`
- Modify: `server/plugins/im/service/event.go`

- [ ] Add AI config keys: `ai_web_search_enabled`, `ai_web_search_provider`, `ai_web_search_api_key`, `ai_web_search_endpoint`, `ai_web_search_top_k`.
- [ ] Implement Serper-compatible search provider with timeout inherited from AI config.
- [ ] In `AIAnswer`, when enabled, search with the rewritten retrieval query and append concise snippets as `【联网搜索】`.
- [ ] Record `web_search` event with success/failure and result count; failures do not block AI answer.
- [ ] Keep disabled-by-default behavior when config/API key is missing.

### Task 6: Documentation And Verification

**Files:**
- Modify: `docs-site/docs/api/im.md`
- Modify: `docs/im-api-reference.md`
- Modify: `docs/im-feature-matrix.md`

- [ ] Document current session context fields and `/api/v1/im/session` optional parameters.
- [ ] Document `web/public/im-widget.js` usage and iframe embed constraints.
- [ ] Document log fields, filters, and event meanings.
- [ ] Document `typing_draft`/`typing_stop` frames.
- [ ] Document web search config and deployment impact.
- [ ] Run `cd server; go test ./plugins/im/...`.
- [ ] Run frontend builds where package scripts are available: `admin`, `web`, `app`, and `eapp`.

## Final Review Checklist

- [ ] Existing session/message/upload/reply routes remain compatible.
- [ ] No independent AI-CS-style visitor/conversation subsystem is introduced.
- [ ] Web search is disabled unless explicitly configured.
- [ ] docs-site covers feature description, API changes, and deployment/config impact.
- [ ] Git status only contains intended files plus any pre-existing unrelated untracked files.
