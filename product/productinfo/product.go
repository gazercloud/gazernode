package productinfo

import (
	"bytes"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazerui/canvas"
	"image"
)

func Name() string {
	return "gazer_node"
}

func Version() string {
	return "1.2.5 BETA"
}

func BuildTime() string {
	return BUILDTIME
}

func Icon() image.Image {
	iconBin, _ := resources.Asset("files/favicon.ico")
	img, _ := canvas.Decode(bytes.NewBuffer(iconBin))
	return img
}

func Icon64() image.Image {
	iconBin, _ := resources.Asset("files/mainicon64.png")
	img, _, _ := image.Decode(bytes.NewBuffer(iconBin))
	return img
}
