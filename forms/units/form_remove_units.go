package units

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"golang.org/x/image/colornames"
)

type FormRemoveUnits struct {
	uicontrols.Dialog
	client      *client.Client
	txtUnitName *uicontrols.TextBox
	lvUnits     *uicontrols.ListView
	units       []*nodeinterface.UnitStateAllResponseItem
}

func NewFormRemoveUnits(parent uiinterfaces.Widget, client *client.Client, units []*nodeinterface.UnitStateAllResponseItem) *FormRemoveUnits {
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

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialiconsoutlined_48dp_1x_outline_clear_black_48dp_png, c.AccentColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	lblConfirmation := pRight.AddTextBlockOnGrid(0, 1, "Do you want to remove the units?")
	lblConfirmation.SetForeColor(colornames.Red)
	lblAllDataWillBeDestroyed := pRight.AddTextBlockOnGrid(0, 2, "All unit's data will be destroyed!")
	lblAllDataWillBeDestroyed.SetForeColor(c.AccentColor())

	c.lvUnits = pRight.AddListViewOnGrid(0, 4)
	c.lvUnits.AddColumn("Name", 200)
	c.lvUnits.AddColumn("Type", 200)

	for _, unit := range c.units {
		item := c.lvUnits.AddItem(unit.UnitName)
		item.SetValue(1, unit.Type)
	}

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	btnOK.SetEnabled(false)
	chkConfirm := pRight.AddCheckBoxOnGrid(0, 3, "Confirm")
	chkConfirm.OnCheckedChanged = func(checkBox *uicontrols.CheckBox, checked bool) {
		if checked {
			btnOK.SetEnabled(true)
		} else {
			btnOK.SetEnabled(false)
		}
	}

	c.TryAccept = func() bool {
		unitIDs := make([]string, 0)
		for _, sInfo := range c.units {
			unitIDs = append(unitIDs, sInfo.UnitId)
		}

		c.client.RemoveUnit(unitIDs, func(err error) {
			if err == nil {
				c.TryAccept = nil
				c.Accept()
			} else {
				uicontrols.ShowErrorMessage(&c, err.Error(), "Error")
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
