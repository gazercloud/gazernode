package simplemap

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gazercloud/gazernode/actions"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiproperties"
	"github.com/nfnt/resize"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
	"time"
)

type MapControlText struct {
	MapControl

	textColor            *uiproperties.Property
	text                 *uiproperties.Property
	textHAlign           *uiproperties.Property
	textVAlign           *uiproperties.Property
	textFontSize         *uiproperties.Property
	textAdaptiveFontSize *uiproperties.Property

	dataSourceFormat *uiproperties.Property

	borderType         *uiproperties.Property
	borderCircleBegin_ *uiproperties.Property
	borderCircleEnd_   *uiproperties.Property

	backgroundImage_ *uiproperties.Property
	backgroundColor  *uiproperties.Property
	borderColor      *uiproperties.Property
	borderWidth      *uiproperties.Property

	action *uiproperties.Property

	img *image.RGBA
}

func NewMapControlText(mapWidget *MapWidget, parent IMapControl) *MapControlText {
	var c MapControlText
	c.initMapControl(&c, mapWidget, parent)
	c.SetX(1)
	c.SetY(1)
	c.SetWidth(50)
	c.SetHeight(50)

	// Text properties
	c.type_.SetOwnValue("text")

	// Background properties
	c.backgroundImage_ = AddPropertyToControl(&c, "background_image", "Image", uiproperties.PropertyTypeString, "Background", "file")
	c.backgroundImage_.DefaultValue = make([]byte, 0)

	c.backgroundImage_.OnChanged = func(property *uiproperties.Property, oldValue interface{}, newValue interface{}) {
		imgData, err := base64.StdEncoding.DecodeString(c.backgroundImage_.String())
		if err == nil {
			img, _, err := image.Decode(bytes.NewBuffer(imgData))
			if err == nil {
				c.img = image.NewRGBA(img.Bounds())
				draw.Draw(c.img, c.img.Bounds(), img, image.Point{}, draw.Src)
			}
		}
	}
	c.backgroundColor = AddPropertyToControl(&c, "background_color", "Color", uiproperties.PropertyTypeColor, "Background", "")
	c.backgroundColor.SetOwnValue(color.RGBA{})
	c.backgroundColor.DefaultValue = c.backgroundColor.ValueOwn()

	c.text = AddPropertyToControl(&c, "text", "Text", uiproperties.PropertyTypeMultiline, "Text", "")
	c.textColor = AddPropertyToControl(&c, "text_color", "Color", uiproperties.PropertyTypeColor, "Text", "")
	c.textColor.SetOwnValue(colornames.Gray)
	c.textFontSize = AddPropertyToControl(&c, "text_font_size", "Font Size", uiproperties.PropertyTypeInt32, "Text", "")
	c.textFontSize.SetOwnValue(20)
	c.textAdaptiveFontSize = AddPropertyToControl(&c, "text_adaptive_font_size", "Auto Font Size", uiproperties.PropertyTypeBool, "Text", "")
	c.textHAlign = AddPropertyToControl(&c, "text_horizontal_align", "HAlign", uiproperties.PropertyTypeString, "Text", "horizontal-align")
	c.textHAlign.SetOwnValue("center")
	c.textVAlign = AddPropertyToControl(&c, "text_vertical_align", "VAlign", uiproperties.PropertyTypeString, "Text", "vertical-align")
	c.textVAlign.SetOwnValue("center")

	// Datasource properties
	c.dataSource_ = AddPropertyToControl(&c, "data_source", "Path", uiproperties.PropertyTypeString, "DataSource", "datasource")
	c.dataSourceFormat = AddPropertyToControl(&c, "data_source_format", "Format", uiproperties.PropertyTypeString, "DataSource", "data_source_format")
	c.dataSourceFormat.SetOwnValue("{v}")
	c.dataSourceFormat.DefaultValue = c.dataSourceFormat.ValueOwn()

	// Border properties
	c.borderType = AddPropertyToControl(&c, "border_type", "Type", uiproperties.PropertyTypeString, "Border", "border_type")
	c.borderType.SetOwnValue("rect")
	c.borderType.DefaultValue = c.borderType.ValueOwn()
	c.borderColor = AddPropertyToControl(&c, "border_color", "Color", uiproperties.PropertyTypeColor, "Border", "")
	c.borderColor.SetOwnValue(colornames.Gray)
	c.borderColor.DefaultValue = c.borderColor.ValueOwn()
	c.borderWidth = AddPropertyToControl(&c, "border_width", "Width", uiproperties.PropertyTypeInt32, "Border", "")
	c.borderWidth.SetOwnValue(1)
	c.borderWidth.DefaultValue = c.borderWidth.ValueOwn()

	// Circle border properties
	c.borderCircleBegin_ = AddPropertyToControl(&c, "border_circle_begin", "Begin", uiproperties.PropertyTypeInt32, "Border Circle", "")
	c.borderCircleBegin_.SetOwnValue(int32(0))
	c.borderCircleBegin_.DefaultValue = c.borderCircleBegin_.ValueOwn()
	c.borderCircleBegin_.SetVisible(false)

	c.borderCircleEnd_ = AddPropertyToControl(&c, "border_circle_end", "End", uiproperties.PropertyTypeInt32, "Border Circle", "")
	c.borderCircleEnd_.SetOwnValue(int32(360))
	c.borderCircleEnd_.DefaultValue = c.borderCircleEnd_.ValueOwn()
	c.borderCircleEnd_.SetVisible(false)

	c.action = AddPropertyToControl(&c, "action", "Action", uiproperties.PropertyTypeString, "Action", "action")
	c.action.SetOwnValue("")
	c.action.DefaultValue = c.action.ValueOwn()

	c.type_.SetVisible(false)

	return &c
}

func (c *MapControlText) GetProperties() []*uiproperties.Property {
	props := c.PropertiesContainer.GetProperties()
	return props
}

func (c *MapControlText) OnMouseDown(x, y int) {
	fmt.Println("OnMouseDown Text ", c.Name(), x, y)
	fmt.Println()

	if c.action.String() != "" {
		c.ExecuteAction()
	}
}

func (c *MapControlText) ExecuteAction() {
	var a actions.Action
	err := json.Unmarshal([]byte(c.action.String()), &a)
	if err == nil {
		_ = c.ExecAction(&a)
	}
}

func (c *MapControlText) HasAction() bool {
	if c.action.String() != "" {
		return true
	}
	return false
}

func (c *MapControlText) drawBackImage(ctx ui.DrawContext) {
	if c.width_.Int32() > 0 && c.height_.Int32() > 0 {
		if c.img != nil {
			wX, wY := c.MapWidget.RectClientAreaOnWindow()

			clipXo := int32(ctx.State().TranslateX - wX)
			clipYo := int32(ctx.State().TranslateY - wY)
			clipX := int32(ctx.State().TranslateX - wX)
			clipY := int32(ctx.State().TranslateY - wY)
			clipW := c.scaleValue(c.Width())
			clipH := c.scaleValue(c.Height())

			if clipX < 0 {
				clipW += clipX
				clipX = 0
			}
			if clipY < 0 {
				clipH += clipY
				clipY = 0
			}
			if clipX+clipW > int32(c.MapWidget.Width()) {
				clipW = clipW - (clipX + clipW - int32(c.MapWidget.Width()))
			}
			if clipY+clipH > int32(c.MapWidget.Height()) {
				clipH = clipH - (clipY + clipH - int32(c.MapWidget.Height()))
			}
			if clipW < 0 {
				clipW = 0
			}
			if clipH < 0 {
				clipH = 0
			}

			percentsX := 0.0
			if clipXo < 0 {
				percentsX = -(float64(clipXo) / float64(c.scaleValue(c.Width())))
			}
			percentsY := 0.0
			if clipYo < 0 {
				percentsY = -(float64(clipYo) / float64(c.scaleValue(c.Height())))
			}
			percentsW := float64(clipW) / float64(c.scaleValue(c.Width()))
			percentsH := float64(clipH) / float64(c.scaleValue(c.Height()))

			srcPosX := int(float64(c.img.Bounds().Max.X) * percentsX)
			srcPosY := int(float64(c.img.Bounds().Max.Y) * percentsY)
			srcPosW := int(float64(c.img.Bounds().Max.X) * percentsW)
			srcPosH := int(float64(c.img.Bounds().Max.Y) * percentsH)

			if clipW > 0 && clipH > 0 {
				qq := c.img.SubImage(image.Rect(srcPosX, srcPosY, srcPosX+srcPosW, srcPosY+srcPosH))
				img := resize.Resize(uint(clipW), uint(clipH), qq, resize.Bilinear)
				ctx.DrawImage(int(percentsX*float64(c.scaleValue(c.Width()))), int(percentsY*float64(c.scaleValue(c.Height()))), 0, 0, img)
			}

		}
	}

}

func (c *MapControlText) drawControl(ctx ui.DrawContext) {
	if c.width_.Int32() > 0 && c.height_.Int32() > 0 {
		_, _, _, backA := c.backgroundColor.Color().RGBA()
		if backA > 0 {
			ctx.SetColor(c.backgroundColor.Color())
			ctx.FillRect(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())))
		}

		c.drawBackImage(ctx)

		ctx.Save()
		ctx.ClipIn(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())))
		ha := canvas.HAlignCenter
		va := canvas.VAlignCenter

		if c.textHAlign.String() == "left" {
			ha = canvas.HAlignLeft
		}
		if c.textHAlign.String() == "right" {
			ha = canvas.HAlignRight
		}
		if c.textVAlign.String() == "top" {
			va = canvas.VAlignTop
		}
		if c.textVAlign.String() == "bottom" {
			va = canvas.VAlignBottom
		}

		ctx.SetTextAlign(ha, va)
		ctx.SetColor(c.textColor.Color())
		ctx.SetFontFamily("Roboto")
		fontSize := c.textFontSize.Int32()
		ctx.SetFontSize(float64(fontSize))
		w, h := ctx.MeasureText(c.text.String())

		if c.textAdaptiveFontSize.Bool() {
			kW := float64(w) / float64(c.Width())
			kH := float64(h) / float64(c.Height())
			targetFontSizeW := 0.0
			targetFontSizeH := 0.0
			if kW != 0 {
				targetFontSizeW = float64(fontSize) * (1 / kW)
				targetFontSizeW -= targetFontSizeW / 5
				if targetFontSizeW > 288 {
					targetFontSizeW = 288
				}
			}
			if kH != 0 {
				targetFontSizeH = float64(fontSize) * (1 / kH)
				targetFontSizeH -= targetFontSizeH / 5
				if targetFontSizeH > 288 {
					targetFontSizeH = 288
				}
			}

			fontSize = int32(math.Min(targetFontSizeW, targetFontSizeH))
		}

		ctx.SetFontSize(float64(c.scaleValue(fontSize)))
		ctx.DrawText(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())), c.text.String())
		ctx.Load()

		if c.borderType.String() == "rect" {
			_, _, _, borderA := c.borderColor.Color().RGBA()
			if borderA > 0 && c.borderWidth.Int32() > 0 {
				ctx.SetColor(c.borderColor.Color())
				ctx.SetStrokeWidth(int(c.scaleValue(c.borderWidth.Int32())))
				ctx.DrawRect(0, 0, int(c.scaleValue(c.Width())), int(c.scaleValue(c.Height())))
			}
		}
		if c.borderType.String() == "circle" {
			cc := ctx.GG()
			cc.Push()
			cc.SetLineCapSquare()
			cc.Translate(float64(ctx.State().TranslateX), float64(ctx.State().TranslateY))
			cc.SetColor(c.borderColor.Color())
			cc.SetLineWidth(c.scaleValueFloat(c.borderWidth.Int32()))
			//cc.DrawEllipse(float64(c.scaleValue(c.Width() / 2)), float64(c.scaleValue(c.Height() / 2)), float64(c.scaleValue(c.Width() / 2)), float64(c.scaleValue(c.Height() / 2)))
			cc.DrawEllipticalArc(float64(c.scaleValueFloat(c.Width()/2)), float64(c.scaleValueFloat(c.Height()/2)), float64(c.scaleValueFloat(c.Width()/2)), float64(c.scaleValueFloat(c.Height()/2)), gg.Radians(float64(c.borderCircleBegin_.Int32())), gg.Radians(float64(c.borderCircleEnd_.Int32())))
			cc.Stroke()
			cc.Pop()
		}
	}

}

func (c *MapControlText) SetBorderTypeRect() {
	c.borderType.SetOwnValue("rect")
}

func (c *MapControlText) SetBorderTypeCircle() {
	c.borderType.SetOwnValue("circle")
}

func (c *MapControlText) SetBorderWidth(borderWidth int) {
	c.borderWidth.SetOwnValue(borderWidth)
}

func (c *MapControlText) TypeName() string {
	return "text"
}

func (c *MapControlText) adaptiveSize() bool {
	return true
}

func (c *MapControlText) Tick() {
	if c.dataSource() != "" {
		if c.mapDataSource != nil {
			c.mapDataSource.GetDataItemValue(c.FullDataSource(), c.imapControl)
		}

		txt := c.dataSourceFormat.String()

		txt = strings.ReplaceAll(txt, "{v}", c.value.Value)
		txt = strings.ReplaceAll(txt, "{uom}", c.value.UOM)
		txt = strings.ReplaceAll(txt, "{d}", time.Unix(0, c.value.DT*1000).Format("02.01.2006"))
		txt = strings.ReplaceAll(txt, "{t}", time.Unix(0, c.value.DT*1000).Format("15:04:05"))

		c.text.SetOwnValue(txt)
	}
}

func (c *MapControlText) Action() *actions.Action {
	if c.action.String() == "" {
		return nil
	}

	var a actions.Action
	err := json.Unmarshal([]byte(c.action.String()), &a)
	if err == nil {
		return &a
	}
	return nil
}
