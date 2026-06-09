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
