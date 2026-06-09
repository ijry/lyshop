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
