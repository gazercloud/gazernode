package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormRemoveUnits struct {
	uicontrols.Dialog
	client      *client.Client
	txtUnitName *uicontrols.TextBox
	txtError    *uicontrols.TextBlock
	lvUnits     *uicontrols.ListView
	units       []*units_common.UnitInfo
}

func NewFormRemoveUnits(parent uiinterfaces.Widget, client *client.Client, units []*units_common.UnitInfo) *FormRemoveUnits {
	var c FormRemoveUnits
	c.client = client
	c.units = units
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImageAdjusted("icons/material/image/drawable-hdpi/ic_blur_on_black_48dp.png", c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	//pRight.AddTextBlockOnGrid(0, 0, "Select cloud channel:")

	c.lvUnits = pRight.AddListViewOnGrid(0, 1)
	c.lvUnits.AddColumn("Name", 200)
	c.lvUnits.AddColumn("Type", 200)

	for _, unit := range c.units {
		item := c.lvUnits.AddItem(unit.Name)
		item.SetValue(1, unit.Type)
	}

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		unitIDs := make([]string, 0)
		for _, sInfo := range c.units {
			unitIDs = append(unitIDs, sInfo.Id)
		}

		c.client.RemoveUnit(unitIDs, func(err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				c.txtError.SetText(err.Error())
			}
		})
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	return &c
}

func (c *FormRemoveUnits) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Remove units")
	c.Resize(600, 400)
}
