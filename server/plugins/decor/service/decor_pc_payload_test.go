package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultComponentsForPagePC(t *testing.T) {
	got := defaultComponentsForPage("pc")
	require.JSONEq(t, `{"pageStyle":{"background":{"mode":"solid","solidColor":"#f8fafc","overlay":{"enabled":false,"color":"#000000","opacity":0.2}},"content":{"maxWidth":1280,"gutterX":24,"sectionGap":24},"surface":{"radius":12,"shadow":"none"}},"components":[]}`, string(got))
}

func TestDefaultComponentsForPageNonPC(t *testing.T) {
	got := defaultComponentsForPage("index")
	require.JSONEq(t, `[]`, string(got))
}

