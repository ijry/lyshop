package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCartFieldEncodeDecode(t *testing.T) {
	field := formatCartField(101, 9001)
	require.Equal(t, "101:9001", field)

	skuID, activityProductID, err := parseCartField(field)
	require.NoError(t, err)
	require.Equal(t, uint64(101), skuID)
	require.Equal(t, uint64(9001), activityProductID)
}

func TestParseCartField_BackwardCompatible(t *testing.T) {
	skuID, activityProductID, err := parseCartField("101")
	require.NoError(t, err)
	require.Equal(t, uint64(101), skuID)
	require.Equal(t, uint64(0), activityProductID)
}

