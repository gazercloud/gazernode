package actions

import (
	"bytes"
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/dialogs"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"image"
	"image/png"
)

type OpenMap struct {
	uicontrols.Panel
	client   *client.Client
	txtResId *uicontrols.TextBoxExt
	txtType  *uicontrols.TextBlock
	txtName  *uicontrols.TextBlock
	imgThumb *uicontrols.ImageBox
}

func NewOpenMap(parent uiinterfaces.Widget, client *client.Client) *OpenMap {
	var c OpenMap
	c.client = client
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	c.AddTextBlockOnGrid(0, 0, "ResourceID:")
	c.txtResId = c.AddTextBoxExtOnGrid(1, 0, "", func(textBoxExt *uicontrols.TextBoxExt) {
		dialogs.LookupResource(&c, c.client, func(key string) {
			textBoxExt.SetText(key)
			c.updateResourceProperties()
		})
	})
	c.AddTextBlockOnGrid(0, 1, "Name:")
	c.txtName = c.AddTextBlockOnGrid(1, 1, "")
	c.AddTextBlockOnGrid(0, 2, "Type:")
	c.txtType = c.AddTextBlockOnGrid(1, 2, "")
	c.imgThumb = c.AddImageBoxOnGrid(1, 3, nil)

	c.imgThumb.SetScaling(uicontrols.ImageBoxScaleNoScaleImageInLeftTop)
	c.imgThumb.SetMinHeight(64)
	c.imgThumb.SetMinWidth(64)
	c.imgThumb.SetYExpandable(true)

	return &c
}

func (c *OpenMap) LoadAction(value string) {
	var a OpenMapAction
	err := json.Unmarshal([]byte(value), &a)
	if err != nil {
		return
	}
	c.txtResId.SetText(a.ResId)
	c.updateResourceProperties()
}

func (c *OpenMap) SaveAction() string {
	var a OpenMapAction
	a.ResId = c.txtResId.Text()
	bs, _ := json.MarshalIndent(a, "", " ")
	return string(bs)
}

func (c *OpenMap) updateResourceProperties() {
	c.client.ResGet(c.txtResId.Text(), func(item *common_interfaces.ResourcesItem, err error) {
		if err != nil {
			return
		}
		if item == nil {
			return
		}
		c.txtName.SetText(item.Info.Name)
		typeName := "unknown"
		if item.Info.Type == "simple_map" {
			typeName = "Map"
		}
		if item.Info.Type == "chart_group" {
			typeName = "Chart group"
		}
		c.txtType.SetText(typeName)

		var thumbnail image.Image
		if item.Info.Thumbnail != nil {
			thumbnail, _ = png.Decode(bytes.NewBuffer(item.Info.Thumbnail))
		} else {
			thumbnail = image.NewAlpha(image.Rect(0, 0, 32, 32))
		}
		c.imgThumb.SetImage(thumbnail)
	})
}
