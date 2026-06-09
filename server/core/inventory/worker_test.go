package inventory

import (
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ijry/lyshop/core/db"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type workerStubProvider struct {
	processErr error
}

func (s *workerStubProvider) ProcessTask(_ *gorm.DB, _ *InventoryIntegrationTask, _ time.Time) error {
	return s.processErr
}

func TestClaimDueTaskSkipsLockedRecords(t *testing.T) {
	testDB := newInventoryTestDB(t)
	now := time.Now()

	task := InventoryIntegrationTask{
		Provider:      "external_wms",
		Action:        "reserve",
		BizType:       "order",
		BizNo:         "T1001",
		Status:        TaskStatusPending,
		NextRetryAt:   &now,
		LockExpiresAt: ptrTime(now.Add(2 * time.Minute)),
		LockOwner:     "worker-a",
	}
	require.NoError(t, testDB.Create(&task).Error)

	claimed, err := ClaimDueTask(testDB, "worker-b", now)
	require.NoError(t, err)
	require.Nil(t, claimed)
}

func TestProcessTaskMarksSuccess(t *testing.T) {
	testDB := newInventoryTestDB(t)
	provider := &workerStubProvider{processErr: nil}
	now := time.Now()

	task := InventoryIntegrationTask{
		Provider:      "external_wms",
		Action:        "deduct",
		BizType:       "order",
		BizNo:         "T1002",
		Status:        TaskStatusPending,
		MaxAttempts:   3,
		BackoffSeconds: 30,
		NextRetryAt:   &now,
	}
	require.NoError(t, testDB.Create(&task).Error)

	err := ProcessTask(testDB, provider, &task, "worker-a", now)
	require.NoError(t, err)

	var latest InventoryIntegrationTask
	require.NoError(t, testDB.First(&latest, task.ID).Error)
	require.Equal(t, TaskStatusSuccess, latest.Status)
	require.Equal(t, 1, latest.AttemptCount)
	require.NotNil(t, latest.CompletedAt)
}

func TestProcessTaskSchedulesRetryOnTransientFailure(t *testing.T) {
	testDB := newInventoryTestDB(t)
	provider := &workerStubProvider{processErr: ErrInventoryBusy}
	now := time.Now()

	task := InventoryIntegrationTask{
		Provider:       "external_wms",
		Action:         "reserve",
		BizType:        "order",
		BizNo:          "T1003",
		Status:         TaskStatusPending,
		MaxAttempts:    3,
		BackoffSeconds: 30,
		NextRetryAt:    &now,
	}
	require.NoError(t, testDB.Create(&task).Error)

	err := ProcessTask(testDB, provider, &task, "worker-a", now)
	require.NoError(t, err)

	var latest InventoryIntegrationTask
	require.NoError(t, testDB.First(&latest, task.ID).Error)
	require.Equal(t, TaskStatusPending, latest.Status)
	require.Equal(t, 1, latest.AttemptCount)
	require.NotNil(t, latest.NextRetryAt)
	require.True(t, latest.NextRetryAt.After(now))
	require.Empty(t, latest.LockOwner)
	require.Nil(t, latest.LockExpiresAt)
}

func newInventoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:inventory_worker_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)

	old := db.DB
	db.DB = gdb
	t.Cleanup(func() { db.DB = old })

	require.NoError(t, gdb.AutoMigrate(&InventoryIntegrationTask{}))
	return gdb
}

func ptrTime(v time.Time) *time.Time {
	return &v
}
