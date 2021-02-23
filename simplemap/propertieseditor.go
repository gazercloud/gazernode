package simplemap

import (
	"encoding/base64"
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/dialogs"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/coreforms"
	"github.com/gazercloud/gazerui/filedialogs"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiproperties"
	"image/color"
	"io/ioutil"
	"reflect"
	"time"
)

type PropertiesEditor struct {
	uicontrols.Panel

	client *client.Client

	iPropertiesContainer uiproperties.IPropertiesContainer
	propControls         map[string]uiinterfaces.Widget
	propsMap             map[string]*uiproperties.Property

	loading bool
}

func NewPropertiesEditor(parent uiinterfaces.Widget, client *client.Client) *PropertiesEditor {
	var c PropertiesEditor
	c.client = client
	c.InitControl(parent, &c)
	c.SetMinWidth(270)
	c.SetPanelPadding(0)
	return &c
}

func (c *PropertiesEditor) ControlType() string {
	return "PropertiesEditor"
}

func (c *PropertiesEditor) Dispose() {
	if c.iPropertiesContainer != nil {
		c.iPropertiesContainer.SetPropertyChangeNotifier(nil)
	}
	c.propControls = nil
	c.propsMap = nil
	c.iPropertiesContainer = nil
	c.Panel.Dispose()
}

func (c *PropertiesEditor) SetPropertiesContainer(propertiesContainer uiproperties.IPropertiesContainer) {
	logger.Println("SetPropertiesContainer")

	c.propControls = make(map[string]uiinterfaces.Widget)
	c.propsMap = make(map[string]*uiproperties.Property)

	if c.iPropertiesContainer != nil && !reflect.ValueOf(c.iPropertiesContainer).IsNil() {
		c.iPropertiesContainer.SetPropertyChangeNotifier(nil)
	}

	c.iPropertiesContainer = propertiesContainer
	if c.iPropertiesContainer != nil && !reflect.ValueOf(c.iPropertiesContainer).IsNil() {
		c.iPropertiesContainer.SetPropertyChangeNotifier(c.OnPropertyChanged)
	}

	c.RebuildInterface()
}

func (c *PropertiesEditor) RebuildInterface() {
	logger.Println("RebuildInterface")
	//c.BeginUpdate()
	c.loading = true
	c.RemoveAllWidgets()
	if c.iPropertiesContainer == nil || reflect.ValueOf(c.iPropertiesContainer).IsNil() {
		c.loading = false
		return
	}

	c.SetCellPadding(3)

	lastGroup := ""
	groupIndex := 0
	var groupPanel *uicontrols.Panel
	indexInGroup := 0
	if true {
		for _, property := range c.iPropertiesContainer.GetProperties() {
			if !property.Visible() {
				continue
			}

			if lastGroup != property.GroupName {
				lblGroupName := c.AddTextBlockOnGrid(0, groupIndex*2, property.GroupName)
				lblGroupName.SetName("Group name " + property.Name)
				//lblGroupName.SetFontSize(10)
				lblGroupName.TextHAlign = canvas.HAlignLeft
				//lblGroupName.SetUnderline(true)
				lblGroupName.SetBorderBottom(1, c.AccentColor())
				lblGroupName.SetForeColor(c.AccentColor())
				lastGroup = property.GroupName
				groupPanel = c.AddPanelOnGrid(0, groupIndex*2+1)
				groupPanel.SetPanelPadding(0)
				groupPanel.SetCellPadding(2)
				groupIndex++
			}

			if groupPanel == nil {
				break
			}

			c.propsMap[property.Name] = property

			//panelPropEditor := c.AddPanelOnGrid(0, index)
			//panelPropEditor.SetPanelPadding(0)
			lblName := groupPanel.AddTextBlockOnGrid(0, indexInGroup, "  "+property.DisplayName+":")
			lblName.SetName("Prop " + property.Name)
			lblName.TextHAlign = canvas.HAlignLeft

			if property.Type == uiproperties.PropertyTypeBool {
				numEditor := groupPanel.AddCheckBoxOnGrid(1, indexInGroup, "")
				numEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				numEditor.OnCheckedChanged = c.CheckBoxChanged
				numEditor.SetUserData("propName", property.Name)
				numEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = numEditor
			}

			if property.Type == uiproperties.PropertyTypeInt {
				numEditor := groupPanel.AddSpinBoxOnGrid(1, indexInGroup)
				numEditor.SetPrecision(0)
				numEditor.SetIncrement(1)
				numEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				numEditor.OnValueChanged = c.SpinBoxChanged
				numEditor.SetUserData("propName", property.Name)
				numEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = numEditor
			}
			if property.Type == uiproperties.PropertyTypeInt32 {
				numEditor := groupPanel.AddSpinBoxOnGrid(1, indexInGroup)
				numEditor.SetPrecision(0)
				numEditor.SetIncrement(1)
				numEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				numEditor.OnValueChanged = c.SpinBoxChanged
				numEditor.SetUserData("propName", property.Name)
				numEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = numEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "" {
				txtEditor := groupPanel.AddTextBoxOnGrid(1, indexInGroup)
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnTextChanged = c.TextBoxChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "datasource" {
				txtEditor := groupPanel.AddTextBoxExtOnGrid(1, indexInGroup, "", func(textBoxExt *uicontrols.TextBoxExt) {
					dialogs.LookupDataItem(c, c.client, func(key string) {
						textBoxExt.SetText(key)
					})
				})
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnTextChanged = c.TextBoxExtChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "data_source_format" {
				txtEditor := groupPanel.AddTextBoxExtOnGrid(1, indexInGroup, "", func(textBoxExt *uicontrols.TextBoxExt) {
					EditDataSourceFormat(c, textBoxExt.Text(), func(key string) {
						textBoxExt.SetText(key)
					})
				})
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnTextChanged = c.TextBoxExtChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			continue
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "action" {
				txtEditor := groupPanel.AddButtonOnGrid(1, indexInGroup, "Edit ...", nil)
				txtEditor.SetOnPress(func(ev *uievents.Event) {
					dialog := NewActionEditor(c, txtEditor.TempData)
					dialog.OnAccept = func() {
						txtEditor.TempData = dialog.ActionText()
						if c.loading {
							return
						}
						propName := txtEditor.UserData("propName").(string)
						if _, ok := c.propControls[propName]; ok {
							c.iPropertiesContainer.SetPropertyValue(propName, dialog.ActionText())
						}
					}
					dialog.ShowDialog()
				})
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "file" {
				txtEditor := groupPanel.AddButtonOnGrid(1, indexInGroup, "Browse ...", nil)
				txtEditor.SetOnPress(func(ev *uievents.Event) {
					filedialogs.ShowOpenFile(c, func(filePath string) {
						data, err := ioutil.ReadFile(filePath)
						if err == nil {
							txtEditor.SetUserData("data", data)
							if c.loading {
								return
							}
							propName := txtEditor.UserData("propName").(string)
							if _, ok := c.propControls[propName]; ok {
								stringData := base64.StdEncoding.EncodeToString(data)
								c.iPropertiesContainer.SetPropertyValue(propName, stringData)
							}
						}
					})
				})
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "horizontal-align" {
				txtEditor := groupPanel.AddComboBoxOnGrid(1, indexInGroup)
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnCurrentIndexChanged = c.ComboBoxChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				txtEditor.AddItem("Left", "left")
				txtEditor.AddItem("Center", "center")
				txtEditor.AddItem("Right", "right")
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "vertical-align" {
				txtEditor := groupPanel.AddComboBoxOnGrid(1, indexInGroup)
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnCurrentIndexChanged = c.ComboBoxChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				txtEditor.AddItem("Top", "top")
				txtEditor.AddItem("Center", "center")
				txtEditor.AddItem("Bottom", "bottom")
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "border_type" {
				txtEditor := groupPanel.AddComboBoxOnGrid(1, indexInGroup)
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnCurrentIndexChanged = c.ComboBoxChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				txtEditor.AddItem("Rect", "rect")
				txtEditor.AddItem("Circle", "circle")
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeString && property.SubType == "edges" {
				txtEditor := groupPanel.AddComboBoxOnGrid(1, indexInGroup)
				txtEditor.OnCurrentIndexChanged = c.ComboBoxChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				txtEditor.AddItem("Square", "square")
				txtEditor.AddItem("Round", "round")
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeMultiline {
				txtEditor := groupPanel.AddButtonOnGrid(1, indexInGroup, "Edit ...", nil)
				txtEditor.SetOnPress(func(ev *uievents.Event) {
					dialog := coreforms.NewMultilineEditor(c, txtEditor.TempData)
					dialog.OnAccept = func() {
						c.MultilineChanged(txtEditor, dialog.ResultText())
					}
					dialog.OnTextChanged = func(txtMultiline *coreforms.MultilineEditor, oldText string, newText string) {
						c.MultilineChanged(txtEditor, dialog.Text())
					}
					dialog.ShowDialog()
				})
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeColor {
				txtEditor := uicontrols.NewColorPicker(c)
				txtEditor.SetPos(0, 0)
				txtEditor.SetSize(0, 20)
				txtEditor.OwnWindow = c.OwnWindow
				txtEditor.OnColorChanged = c.ColorPickerChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				groupPanel.AddWidgetOnGrid(txtEditor, 1, indexInGroup)
				c.propControls[property.Name] = txtEditor
			}

			if property.DefaultValue != nil {
				btnSetDefault := groupPanel.AddButtonOnGrid(2, indexInGroup, "", func(event *uievents.Event) {
					if c.loading {
						return
					}
					evButton, ok := event.Sender.(*uicontrols.Button)
					if ok {
						propName := evButton.UserData("propName").(string)
						if propToDefault, ok := c.propsMap[propName]; ok {
							if propToDefault.DefaultValue != nil {
								c.iPropertiesContainer.SetPropertyValue(propName, propToDefault.DefaultValue)
							}
						}
					}
				})
				btnSetDefault.SetUserData("propName", property.Name)
				btnSetDefault.SetUserData("propType", property.Type)
				//btnSetDefault.SetImage(uiresources.ResImageAdjusted("icons/material/navigation/drawable-hdpi/ic_close_black_48dp.png", c.InactiveColor()))
				btnSetDefault.SetImageSize(16, 16)
				btnSetDefault.SetBorders(0, color.RGBA{})
				btnSetDefault.SetMinWidth(24)
				btnSetDefault.SetMinHeight(24)
				btnSetDefault.SetMaxWidth(24)
				btnSetDefault.SetMaxHeight(24)
				btnSetDefault.SetTooltip("Reset to default value")
			}
			indexInGroup++
		}
	}
	c.AddVSpacerOnGrid(0, groupIndex*2+2)
	//c.EndUpdate()

	c.LoadPropertiesValues()
	c.loading = false
}

func (c *PropertiesEditor) LoadPropertiesValues() {
	logger.Println("LoadPropertiesValues")
	t1 := time.Now()

	c.BeginUpdate()
	c.loading = true
	for propName, widget := range c.propControls {
		value := c.iPropertiesContainer.PropertyValue(propName)

		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "" {
			txtBox := widget.(*uicontrols.TextBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetText(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "datasource" {
			txtBox := widget.(*uicontrols.TextBoxExt)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetText(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "action" {
			txtBox := widget.(*uicontrols.Button)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.TempData = fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "data_source_format" {
			txtBox := widget.(*uicontrols.TextBoxExt)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetText(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "horizontal-align" {
			txtBox := widget.(*uicontrols.ComboBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetCurrentItemKey(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "vertical-align" {
			txtBox := widget.(*uicontrols.ComboBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetCurrentItemKey(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "border_type" {
			txtBox := widget.(*uicontrols.ComboBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetCurrentItemKey(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeString && c.propsMap[propName].SubType == "edges" {
			txtBox := widget.(*uicontrols.ComboBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetCurrentItemKey(fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeMultiline {
			txtBox := widget.(*uicontrols.Button)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.TempData = fmt.Sprint(c.iPropertiesContainer.PropertyValue(propName))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeInt {
			txtBox := widget.(*uicontrols.SpinBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetValue(float64(c.iPropertiesContainer.PropertyValue(propName).(int)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeInt32 {
			txtBox := widget.(*uicontrols.SpinBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetValue(float64(value.(int32)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeBool {
			txtBox := widget.(*uicontrols.CheckBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetChecked(c.iPropertiesContainer.PropertyValue(propName).(bool))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeColor {
			txtBox := widget.(*uicontrols.ColorPicker)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetColor(c.iPropertiesContainer.PropertyValue(propName).(color.Color))
			}
		}
	}

	for _, widget := range c.propControls {
		widget.EndUpdate()
	}
	c.loading = false
	c.EndUpdate()

	t2 := time.Now()
	fmt.Println("1:", t2.Sub(t1))
}

func (c *PropertiesEditor) TextBoxChanged(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
	if c.loading {
		return
	}
	propName := txtBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.iPropertiesContainer.SetPropertyValue(propName, newValue)
	}
}

func (c *PropertiesEditor) ComboBoxChanged(ev *uicontrols.ComboBoxEvent) {
	if c.loading {
		return
	}
	txtBox := ev.Sender.(*uicontrols.ComboBox)
	propName := txtBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.iPropertiesContainer.SetPropertyValue(propName, txtBox.CurrentItemKey())
	}
}

func (c *PropertiesEditor) TextBoxExtChanged(txtBox *uicontrols.TextBoxExt, oldValue string, newValue string) {
	if c.loading {
		return
	}
	propName := txtBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.iPropertiesContainer.SetPropertyValue(propName, newValue)
	}
}

func (c *PropertiesEditor) MultilineChanged(txtBox *uicontrols.Button, newValue string) {
	if c.loading {
		return
	}
	propName := txtBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.iPropertiesContainer.SetPropertyValue(propName, newValue)
	}
}

func (c *PropertiesEditor) SpinBoxChanged(spinBox *uicontrols.SpinBox, value float64) {
	if c.loading {
		return
	}
	propName := spinBox.UserData("propName").(string)
	propType := spinBox.UserData("propType").(uiproperties.PropertyType)
	if _, ok := c.propControls[propName]; ok {
		if propType == uiproperties.PropertyTypeInt {
			c.iPropertiesContainer.SetPropertyValue(propName, int(value))
		}
		if propType == uiproperties.PropertyTypeInt32 {
			c.iPropertiesContainer.SetPropertyValue(propName, int32(value))
		}
	}
}

func (c *PropertiesEditor) CheckBoxChanged(checkBox *uicontrols.CheckBox, checked bool) {
	if c.loading {
		return
	}
	propName := checkBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.iPropertiesContainer.SetPropertyValue(propName, checked)
	}
}

func (c *PropertiesEditor) ColorPickerChanged(colorPicker *uicontrols.ColorPicker, color color.Color) {
	if c.loading {
		return
	}
	propName := colorPicker.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.iPropertiesContainer.SetPropertyValue(propName, color)
	}
}

func (c *PropertiesEditor) OnPropertyChanged(prop *uiproperties.Property) {
	logger.Println("OnPropertyChanged")
	c.LoadPropertiesValues()
}
