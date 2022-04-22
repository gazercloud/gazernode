package productinfo

import (
	"bytes"
	"github.com/gazercloud/gazernode/resources"
	//"github.com/gazercloud/gazerui/canvas"
	"image"
)

func Name() string {
	return "GazerNode"
}

func Version() string {
	return "2.4.5"
}

func BuildTime() string {
	return BUILDTIME
}

func Icon() image.Image {
	//img, _ := canvas.Decode(bytes.NewBuffer(resources.R_files_favicon_ico))
	return nil
}

func Icon64() image.Image {
	img, _, _ := image.Decode(bytes.NewBuffer(resources.R_files_mainicon64_png))
	return img
}
