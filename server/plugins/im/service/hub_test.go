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
