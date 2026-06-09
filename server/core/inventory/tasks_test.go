package inventory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnqueueIntegrationTaskCreatesPendingTask(t *testing.T) {
	task := NewIntegrationTask("external_wms", "order", "ORD-10", "reserve", TaskPayload{
		Items: []ReserveItem{{SkuID: 1, Qty: 2}},
	})
	require.Equal(t, "pending", task.Status)
	require.Equal(t, "external_wms", task.Provider)
	require.Equal(t, "reserve", task.Action)
}
