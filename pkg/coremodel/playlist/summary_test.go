package playlist

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlaylistSummary(t *testing.T) {
	builder := GetObjectSummaryBuilder()

	// Do not parse invalid input
	_, _, err := builder(context.Background(), "abc", []byte("{invalid json"))
	require.Error(t, err)

	playlist := Model{
		Interval: "30s",
		Name:     "test",
		Items: &[]PlaylistItem{
			{Type: PlaylistItemTypeDashboardByUid, Value: "D1"},
			{Type: PlaylistItemTypeDashboardByTag, Value: "tagA"},
			{Type: PlaylistItemTypeDashboardByUid, Value: "D3"},
		},
	}
	out, err := json.Marshal(playlist)
	require.NoError(t, err)
	require.NotNil(t, out)

	// Do not parse invalid input
	summary, body, err := builder(context.Background(), "abc", out)
	require.NoError(t, err)
	require.Equal(t, "test", summary.Name)
	require.Equal(t, 2, len(summary.References))
	require.Equal(t, map[string]string{"tagA": ""}, summary.Labels)
	require.True(t, json.Valid(body))
}
