package simplemap

import (
	"bytes"
	"encoding/base64"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiproperties"
	"github.com/nfnt/resize"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/colornames"
	"image"
)

type MapControlSvg struct {
	MapControl

	code_      *uiproperties.Property
	color_     *uiproperties.Property
	lineWidth_ *uiproperties.Property
	img        image.Image
}

func NewMapControlSvg(mapWidget *MapWidget, parent IMapControl) *MapControlSvg {
	var c MapControlSvg
	c.initMapControl(&c, mapWidget, parent)
	c.SetX(1)
	c.SetY(1)
	c.SetWidth(50)
	c.SetHeight(50)
	c.type_.SetOwnValue("svg")

	c.code_ = AddPropertyToControl(&c, "code", "Code", uiproperties.PropertyTypeString, "Svg", "file")
	c.code_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		imgData, err := base64.StdEncoding.DecodeString(c.code_.String())
		if err == nil {
			reader := bytes.NewBuffer(imgData)

			icon, _ := oksvg.ReadIconStream(reader)
			w := int(icon.ViewBox.W)
			h := int(icon.ViewBox.H)
			icon.SetTarget(0, 0, float64(w), float64(h))
			rgba := image.NewRGBA(image.Rect(0, 0, w, h))
			icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
			c.img = rgba
			//ctx.DrawImage(0, 0, w, h, img)

			/*img, err := png.Decode(bytes.NewBuffer(imgData))
			if err == nil {
				img := resize.Resize(uint(c.scaleValue(c.Width())), uint(c.scaleValue(c.Height())), img, resize.Bilinear)
				ctx.DrawImage(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Width())), img)
			}*/
		}
	}

	c.color_ = AddPropertyToControl(&c, "color", "Color", uiproperties.PropertyTypeColor, "Rect", "")
	c.color_.SetOwnValue(colornames.Gray)

	c.lineWidth_ = AddPropertyToControl(&c, "lineWidth", "LineWidth", uiproperties.PropertyTypeInt32, "Rect", "")
	c.lineWidth_.SetOwnValue(int32(1))

	return &c
}

func (c *MapControlSvg) drawControl(ctx ui.DrawContext) {
	//leftX := c.scaleValue(0)
	//topY := c.scaleValue(0)
	if c.width_.Int32() > 0 && c.height_.Int32() > 0 {
		//cc := ctx.GG()

		/*svgContent := `
		<?xml version="1.0" encoding="UTF-8" standalone="no"?>
		<svg version = "1.1"
		     baseProfile="full"
		     xmlns = "http://www.w3.org/2000/svg"
		     xmlns:xlink = "http://www.w3.org/1999/xlink"
		     xmlns:ev = "http://www.w3.org/2001/xml-events"
		     height = "400px"  width = "400px">
		     <rect x="0" y="0" width="400" height="400"
		          fill="none" stroke="black" stroke-width="5px" stroke-opacity="0.5"/>
		     <g fill-opacity="0.6" stroke="black" stroke-width="0.5px">
		        <circle cx="200px" cy="200px" r="104px" fill="red"   transform="translate(  0,-52)" />
		        <circle cx="200px" cy="200px" r="104px" fill="blue"  transform="translate( 60, 52)" />
		        <circle cx="200px" cy="200px" r="104px" fill="green" transform="translate(-60, 52)" />
		     </g>
		</svg>`*/

		img := resize.Resize(uint(c.scaleValue(c.Width())), uint(c.scaleValue(c.Height())), c.img, resize.Bicubic)
		ctx.DrawImage(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())), img)

		/*imgData, err := base64.StdEncoding.DecodeString(c.code_.String())
		if err == nil {
			reader := bytes.NewBuffer(imgData)

			icon, _ := oksvg.ReadIconStream(reader)
			w := int(icon.ViewBox.W)
			h := int(icon.ViewBox.H)
			icon.SetTarget(0, 0, float64(w), float64(h))
			rgba := image.NewRGBA(image.Rect(0, 0, w, h))
			icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
			img := resize.Resize(uint(c.scaleValue(c.Width())), uint(c.scaleValue(c.Height())), rgba, resize.Bicubic)
			ctx.DrawImage(0, 0, w, h, img)
		}*/

		/*ctx.SetColor(c.color_.Color())
		ctx.SetStrokeWidth(int(c.scaleValue(c.lineWidth_.Int32())))
		ctx.DrawRect(int(leftX), int(topY), int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())))*/
	}
}

func (c *MapControlSvg) TypeName() string {
	return "svg"
}

func (c *MapControlSvg) adaptiveSize() bool {
	return true
}
