package public_channels

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormEditCloudChannel struct {
	uicontrols.Dialog
	client      *client.Client
	channelId   string
	txtUnitName *uicontrols.TextBox
	txtError    *uicontrols.TextBlock
}

func NewFormEditCloudChannel(parent uiinterfaces.Widget, client *client.Client, channelId string, name string) *FormEditCloudChannel {
	var c FormEditCloudChannel
	c.client = client
	c.channelId = channelId
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_file_cloud_upload_materialiconsoutlined_48dp_1x_outline_cloud_upload_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Channel name:")
	c.txtUnitName = pRight.AddTextBoxOnGrid(1, 0)
	c.txtUnitName.SetText(name)
	c.txtError = pRight.AddTextBlockOnGrid(1, 3, "")
	c.txtError.SetForeColor(c.AccentColor())

	pRight.AddVSpacerOnGrid(0, 10)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		if c.txtUnitName.Text() == "" || len(c.txtUnitName.Text()) > 50 {
			uicontrols.ShowErrorMessage(&c, "wrong name", "Error")
			return false
		}

		c.client.EditCloudChannel(c.channelId, c.txtUnitName.Text(), func(err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				uicontrols.ShowErrorMessage(&c, err.Error(), "error")
			}
		})
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	c.OnShow = func() {
		c.txtUnitName.Focus()
	}

	return &c
}

func (c *FormEditCloudChannel) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit cloud channel")
	c.Resize(400, 200)
}
