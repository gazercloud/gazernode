package forms

import (
	"encoding/json"
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"time"
)

type FormUnitEdit struct {
	uicontrols.Dialog
	unitId        string
	unitType      string
	unitTypeName  string
	client        *client.Client
	panelName     *uicontrols.Panel
	txtHelp       *uicontrols.TextBlock
	txtName       *uicontrols.TextBox
	chkAutoName   *uicontrols.CheckBox
	configObj     interface{}
	configMetaObj []*units_common.UnitConfigItem
	help          string
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
					uicontrols.ShowInformationMessage(c.ContentPanel(), err.Error(), "error")
				}
			})
		} else {
			configBytes, _ := json.MarshalIndent(c.configObj, "", "")
			c.client.AddUnit(c.unitType, c.txtName.Text(), string(configBytes), func(unitId string, err error) {
				if err == nil {
					c.TryAccept = nil
					c.Accept()
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
			dialog := NewFormHelp(&c, c.unitTypeName, c.help)
			dialog.Resize(900, 500)
			dialog.ShowDialog()
		}
		c.txtHelp.SetForeColor(c.AccentColor())
		c.txtHelp.SetUnderline(true)
		c.txtHelp.SetMouseCursor(ui.MouseCursorPointer)
	}

	if c.unitId != "" {
		c.client.GetUnitConfig(unitId, func(name string, config string, configMeta string, unitType string, err error) {
			c.configMetaObj = units_common.LoadUnitConfigItems(configMeta)
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

				c.chkAutoName = c.panelName.AddCheckBoxOnGrid(1, 1, "AutoName")
				c.chkAutoName.SetChecked(true)

				pan := NewPanelUnitConfigItems(&c, c.configMetaObj, c.configObj, c.client)
				pan.OnConfigChanged = c.onConfigChanged
				pRight.AddWidgetOnGrid(pan, 0, 1)

				makeHelpButton(c.panelName)

				c.client.UnitTypes("", "", 0, 10000000, func(types nodeinterface.UnitTypeListResponse, err error) {
					for _, ut := range types.Types {
						if ut.Type == unitType {
							c.help = ut.Help
							c.unitTypeName = ut.DisplayName
						}
					}
				})
			}
		})
	} else {
		c.client.GetUnitConfigByType(c.unitType, func(name string, configMeta string, err error) {
			c.configMetaObj = units_common.LoadUnitConfigItems(configMeta)
			config := `{}`
			err = json.Unmarshal([]byte(config), &c.configObj)
			if err == nil {
				c.panelName = pRight.AddPanelOnGrid(0, 0)
				c.panelName.AddTextBlockOnGrid(0, 0, "Name:")
				c.txtName = c.panelName.AddTextBoxOnGrid(1, 0)
				c.txtName.SetText(name)
				c.txtName.Focus()

				c.chkAutoName = c.panelName.AddCheckBoxOnGrid(1, 1, "AutoName")
				c.chkAutoName.SetChecked(true)

				pan := NewPanelUnitConfigItems(&c, c.configMetaObj, c.configObj, c.client)
				pan.OnConfigChanged = c.onConfigChanged
				pRight.AddWidgetOnGrid(pan, 0, 1)
				makeHelpButton(c.panelName)
			}

			c.client.UnitTypes("", "", 0, 10000000, func(types nodeinterface.UnitTypeListResponse, err error) {
				for _, ut := range types.Types {
					if ut.Type == unitType {
						c.help = ut.Help
						c.unitTypeName = ut.DisplayName
					}
				}
			})
		})
	}

	c.OnShow = func() {
	}

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

type ConfigForName struct {
	Address string
}

func (c *FormUnitEdit) onConfigChanged() {
	c.SetTitle("Edit unit " + time.Now().String())
	if c.chkAutoName.IsChecked() {
		propName := ""
		if c.configMetaObj != nil {
			for _, item := range c.configMetaObj {
				if item.ItemIsDisplayName {
					propName = item.Name
				}
			}
		}

		if propName != "" {
			if conf, ok := c.configObj.(map[string]interface{}); ok {
				if propValue, ok := conf[propName]; ok {
					c.txtName.SetText(c.unitTypeName + " " + fmt.Sprint(propValue))
				}
			}
		} else {
			c.chkAutoName.SetChecked(false)
		}
	}
}
