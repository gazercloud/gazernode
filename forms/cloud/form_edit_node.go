package cloud

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormEditNode struct {
	uicontrols.Dialog
	client  *client.Client
	nodeId  string
	oldName string
	txtName *uicontrols.TextBox
	btnOK   *uicontrols.Button
}

func NewFormEditNode(parent uiinterfaces.Widget, client *client.Client, nodeId string, oldName string) *FormEditNode {
	var c FormEditNode
	c.client = client
	c.nodeId = nodeId
	c.oldName = oldName
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_maps_layers_materialiconsoutlined_48dp_1x_outline_layers_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Node name:")
	c.txtName = pRight.AddTextBoxOnGrid(1, 0)
	c.txtName.SetText(c.oldName)

	pRight.AddVSpacerOnGrid(0, 10)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		if c.txtName.Text() == "" || len(c.txtName.Text()) > 50 {
			uicontrols.ShowErrorMessage(&c, "wrong name", "Error")
			return false
		}

		c.btnOK.SetEnabled(false)
		c.client.CloudUpdateNode(c.nodeId, c.txtName.Text(), func(resp nodeinterface.CloudUpdateNodeResponse, err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				c.btnOK.SetEnabled(true)
				uicontrols.ShowErrorMessage(&c, err.Error(), "error")
			}
		})
		return false
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.OnShow = func() {
		c.txtName.Focus()
		c.txtName.SelectAllText()
	}

	return &c
}

func (c *FormEditNode) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit node")
	c.Resize(400, 200)
}
