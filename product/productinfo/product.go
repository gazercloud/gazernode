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
	img, _ := canvas.Decode(bytes.NewBuffer(resources.R_files_favicon_ico))
	return img
}

func Icon64() image.Image {
	img, _, _ := image.Decode(bytes.NewBuffer(resources.R_files_mainicon64_png))
	return img
}
