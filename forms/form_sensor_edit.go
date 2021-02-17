package forms

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazerui/coreforms"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type FormUnitEdit struct {
	uicontrols.Dialog
	unitId    string
	unitType  string
	client    *client.Client
	panelName *uicontrols.Panel
	txtHelp   *uicontrols.TextBlock
	txtName   *uicontrols.TextBox
	configObj interface{}
	help      string
}

func NewFormUnitEdit(parent uiinterfaces.Widget, client *client.Client, unitId string, unitType string) *FormUnitEdit {
	var c FormUnitEdit
	c.client = client
	c.unitId = unitId
	c.unitType = unitType
	c.InitControl(parent, &c)
	c.SetName("FormUnitEdit")
	c.Resize(800, 750)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.TryAccept = func() bool {
		if c.txtName == nil {
			return false
		}

		if c.txtName.Text() == "" {
			return false
		}

		if c.unitId != "" {
			b, _ := json.MarshalIndent(c.configObj, "", "")
			c.client.SetUnitConfig(c.unitId, c.txtName.Text(), string(b), func(err error) {
				if err == nil {
					c.TryAccept = nil
					c.Accept()
				} else {
					//c.txtError.SetText(err.Error())
				}
			})
		} else {
			c.client.AddUnit(c.unitType, c.txtName.Text(), func(unitId string, err error) {
				if err == nil {
					b, _ := json.MarshalIndent(c.configObj, "", "")
					c.client.SetUnitConfig(unitId, c.txtName.Text(), string(b), func(err error) {
						if err == nil {
							c.TryAccept = nil
							c.Accept()
						} else {
							//c.txtError.SetText(err.Error())
						}
					})
				} else {
					uicontrols.ShowInformationMessage(c.ContentPanel(), err.Error(), "error")
				}
			})

		}
		return false
	}

	c.SetAcceptButton(btnOK)
	c.SetRejectButton(btnCancel)

	makeHelpButton := func(panel *uicontrols.Panel) {
		c.txtHelp = panel.AddTextBlockOnGrid(2, 0, "Help")
		c.txtHelp.OnClick = func(ev *uievents.Event) {
			dialog := coreforms.NewMultilineEditor(&c, c.help)
			dialog.Resize(900, 500)
			dialog.ShowDialog()
		}
		c.txtHelp.SetForeColor(c.AccentColor())
		c.txtHelp.SetUnderline(true)
		c.txtHelp.SetMouseCursor(ui.MouseCursorPointer)
	}

	if c.unitId != "" {
		c.client.GetUnitConfig(unitId, func(name string, config string, configMeta string, unitType string, err error) {
			configMetaObj := units_common.LoadUnitConfigItems(configMeta)
			if len(config) == 0 {
				config = `{}`
			}
			err = json.Unmarshal([]byte(config), &c.configObj)
			if err == nil {
				c.panelName = pRight.AddPanelOnGrid(0, 0)
				c.panelName.AddTextBlockOnGrid(0, 0, "Name:")
				c.txtName = c.panelName.AddTextBoxOnGrid(1, 0)
				c.txtName.SetText(name)
				c.txtName.Focus()

				pan := NewPanelUnitConfigItems(&c, configMetaObj, c.configObj, c.client)
				pRight.AddWidgetOnGrid(pan, 0, 1)

				makeHelpButton(c.panelName)

				c.client.UnitTypes("", "", 0, 10000000, func(types common_interfaces.UnitTypes, err error) {
					for _, ut := range types.Types {
						if ut.Type == unitType {
							c.help = ut.Help
						}
					}
				})
			}
		})
	} else {
		c.client.GetUnitConfigByType(c.unitType, func(name string, configMeta string, err error) {
			configMetaObj := units_common.LoadUnitConfigItems(configMeta)
			config := `{}`
			err = json.Unmarshal([]byte(config), &c.configObj)
			if err == nil {
				c.panelName = pRight.AddPanelOnGrid(0, 0)
				c.panelName.AddTextBlockOnGrid(0, 0, "Name:")
				c.txtName = c.panelName.AddTextBoxOnGrid(1, 0)
				c.txtName.SetText(name)
				c.txtName.Focus()

				pan := NewPanelUnitConfigItems(&c, configMetaObj, c.configObj, c.client)
				pRight.AddWidgetOnGrid(pan, 0, 1)
				makeHelpButton(c.panelName)
			}
		})
	}

	c.OnShow = func() {
	}

	c.client.UnitTypes("", "", 0, 10000000, func(types common_interfaces.UnitTypes, err error) {
		for _, ut := range types.Types {
			if ut.Type == c.unitType {
				c.help = ut.Help
			}
		}
	})

	return &c
}

func (c *FormUnitEdit) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit unit")
	c.Resize(400, 200)
}

func (c *FormUnitEdit) Dispose() {
	c.client = nil
	c.txtName = nil
	c.configObj = nil
	c.Dialog.Dispose()
}
