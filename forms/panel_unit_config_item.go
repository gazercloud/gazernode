package forms

import (
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/dialogs"
	"github.com/gazercloud/gazernode/protocols/lookup"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"strconv"
)

type PanelUnitConfigItem struct {
	uicontrols.Panel
	client *client.Client

	item   *units_common.UnitConfigItem
	config interface{}

	lblName     *uicontrols.TextBlock
	txtValue    *uicontrols.TextBox
	numValue    *uicontrols.SpinBox
	chkValue    *uicontrols.CheckBox
	lvItems     *uicontrols.ListView
	innerPanel  *uicontrols.Panel
	innerWidget *PanelUnitConfigItems

	OnChanged func()
}

func NewPanelUnitConfigItem(parent uiinterfaces.Widget, item *units_common.UnitConfigItem, config interface{}, client *client.Client) *PanelUnitConfigItem {
	var c PanelUnitConfigItem
	c.item = item
	c.client = client

	if _, ok := config.(map[string]interface{}); !ok {
		config = make(map[string]interface{})
	}

	c.config = config
	c.InitControl(parent, &c)

	if c.item.Type == "string" {
		c.lblName = c.AddTextBlockOnGrid(0, 0, item.DisplayName+":")
		var value interface{}
		var ok bool
		c.txtValue = c.AddTextBoxOnGrid(1, 0)
		c.txtValue.SetName("unitConfigItem" + item.DisplayName)
		if _, ok = c.config.(map[string]interface{})[c.item.Name]; !ok {
			c.config.(map[string]interface{})[c.item.Name] = item.DefaultValue
		}

		if value, ok = c.config.(map[string]interface{})[c.item.Name]; ok {
			valueString, ok := value.(string)
			if !ok {
				c.config.(map[string]interface{})[c.item.Name] = "default string (1)"
				valueString = c.config.(map[string]interface{})[c.item.Name].(string)
			}
			c.txtValue.SetText(valueString)
			c.txtValue.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
				c.config.(map[string]interface{})[c.item.Name] = txtBox.Text()
				c.NotifyChanges()
			}

			if item.Format != "" {
				lookUpButton := c.AddButtonOnGrid(2, 0, "Select ...", func(event *uievents.Event) {
					if item.Format == "data-item" {
						dialogs.LookupDataItem(&c, c.client, func(key string) {
							c.txtValue.SetText(key)
						})
					} else {
						c.client.Lookup(item.Format, func(result lookup.Result, err error) {
							LookupDialog(&c, c.client, item.Format, func(key string) {
								c.txtValue.SetText(key)
							})
						})
					}
				})
				lookUpButton.SetUserData("entity", item.Format)
			}

		}
	}

	if c.item.Type == "bool" {
		var value interface{}
		var ok bool
		c.chkValue = c.AddCheckBoxOnGrid(0, 1, c.item.DisplayName)
		if _, ok = c.config.(map[string]interface{})[c.item.Name]; !ok {
			c.config.(map[string]interface{})[c.item.Name] = item.DefaultValue == "true"
		}

		if value, ok = c.config.(map[string]interface{})[c.item.Name]; ok {
			valueBool, ok := value.(bool)
			if !ok {
				c.config.(map[string]interface{})[c.item.Name] = item.DefaultValue == "true"
				valueBool = c.config.(map[string]interface{})[c.item.Name].(bool)
			}
			c.chkValue.SetChecked(valueBool)
			c.chkValue.OnCheckedChanged = func(checkBox *uicontrols.CheckBox, checked bool) {
				c.config.(map[string]interface{})[c.item.Name] = checkBox.IsChecked()
				c.NotifyChanges()
			}
		}
	}

	if c.item.Type == "num" {
		c.lblName = c.AddTextBlockOnGrid(0, 0, item.DisplayName+":")
		//c.SetLabelWidth(150)

		var value interface{}
		var ok bool
		c.numValue = c.AddSpinBoxOnGrid(1, 0)
		c.numValue.SetMaxWidth(100)
		//c.AddHSpacerOnGrid(2, 0)

		{
			minValue, err := strconv.ParseFloat(item.MinValue, 64)
			if err == nil {
				c.numValue.SetMinValue(minValue)
			} else {
				c.numValue.SetMinValue(0)
			}
		}

		{
			maxValue, err := strconv.ParseFloat(item.MaxValue, 64)
			if err == nil {
				c.numValue.SetMaxValue(maxValue)
			} else {
				c.numValue.SetMaxValue(100)
			}
		}

		{
			precision, _ := strconv.ParseInt(item.Format, 10, 64)
			c.numValue.SetPrecision(int(precision))
		}

		if _, ok = c.config.(map[string]interface{})[c.item.Name]; !ok {
			floatValue, _ := strconv.ParseFloat(item.DefaultValue, 64)
			c.config.(map[string]interface{})[c.item.Name] = floatValue
		}

		if value, ok = c.config.(map[string]interface{})[c.item.Name]; ok {
			valueFloat, ok := value.(float64)
			if !ok {
				floatValue, _ := strconv.ParseFloat(item.DefaultValue, 64)
				c.config.(map[string]interface{})[c.item.Name] = floatValue
				valueFloat = c.config.(map[string]interface{})[c.item.Name].(float64)
			}
			c.numValue.SetValue(valueFloat)
			c.numValue.OnValueChanged = func(spinBox *uicontrols.SpinBox, value float64) {
				c.config.(map[string]interface{})[c.item.Name] = spinBox.Value()
				c.NotifyChanges()
			}
		}
	}

	if c.item.Type == "table" {
		c.lblName = c.AddTextBlockOnGrid(0, 0, item.DisplayName+":")
		panelTable := c.AddPanelOnGrid(0, 1)
		panelButtons := panelTable.AddPanelOnGrid(0, 0)
		panelButtons.AddButtonOnGrid(0, 0, "Add", func(event *uievents.Event) {
			v := c.config.(map[string]interface{})[c.item.Name]
			arr := v.([]interface{})

			obj := make(map[string]interface{})
			for _, vv := range c.item.Children {
				obj[vv.Name] = vv.DefaultValue
			}
			arr = append(arr, obj)
			c.config.(map[string]interface{})[c.item.Name] = arr

			c.reloadTable()
		})
		panelButtons.AddButtonOnGrid(1, 0, "Remove", func(event *uievents.Event) {
			selectedIndex := c.lvItems.SelectedItemIndex()
			if selectedIndex >= 0 {
				v := c.config.(map[string]interface{})[c.item.Name]
				arr1 := v.([]interface{})
				arr1 = append(arr1[:selectedIndex], arr1[selectedIndex+1:]...)
				c.config.(map[string]interface{})[c.item.Name] = arr1
				c.reloadTable()
				c.NotifyChanges()
			}
		})
		c.lvItems = panelTable.AddListViewOnGrid(0, 1)
		c.innerPanel = panelTable.AddPanelOnGrid(1, 1)
		c.reloadTable()

		c.lvItems.OnSelectionChanged = func() {
			c.innerPanel.RemoveAllWidgets()
			var value interface{}
			var ok bool
			var valueArray []interface{}
			valueArray = make([]interface{}, 0)
			if value, ok = c.config.(map[string]interface{})[c.item.Name]; ok {
				valueArray, ok = value.([]interface{})
				if !ok {
					c.config.(map[string]interface{})[c.item.Name] = make([]interface{}, 0)
					valueArray = c.config.(map[string]interface{})[c.item.Name].([]interface{})
				}
			}
			selectedItemIndex := c.lvItems.SelectedItemIndex()
			if selectedItemIndex >= 0 {
				c.innerWidget = NewPanelUnitConfigItems(&c, c.item.Children, valueArray[selectedItemIndex], c.client)
				c.innerPanel.AddWidgetOnGrid(c.innerWidget, 0, 0)
				c.innerWidget.OnChanged = func() {
					c.loadTable()
				}
			}
		}
	}

	return &c
}

func (c *PanelUnitConfigItem) Dispose() {
	c.item = nil
	c.config = nil

	c.lblName = nil
	c.txtValue = nil
	c.numValue = nil
	c.chkValue = nil
	c.lvItems = nil
	c.innerPanel = nil
	c.innerWidget = nil

	c.Panel.Dispose()
}

func (c *PanelUnitConfigItem) SetLabelWidth(width int) {
	if c.lblName != nil {
		c.lblName.SetMinWidth(width)
	}
}

func (c *PanelUnitConfigItem) LabelWidth() int {
	if c.lblName == nil {
		return 100
	}
	return c.lblName.MinWidth()
}

func (c *PanelUnitConfigItem) Save() interface{} {
	if c.item.Type == "string" {
		return c.txtValue.Text()
	}
	if c.item.Type == "bool" {
		return c.chkValue.IsChecked()
	}
	return ""
}

func (c *PanelUnitConfigItem) reloadTable() {
	colIndexByName := make(map[string]int)
	colNameByIndex := make(map[int]string)

	c.lvItems.RemoveItems()
	c.lvItems.RemoveColumns()

	for colIndex, col := range c.item.Children {
		c.lvItems.AddColumn(col.DisplayName, 100)
		colIndexByName[col.Name] = colIndex
		colNameByIndex[colIndex] = col.Name
	}

	var value interface{}
	var ok bool
	if _, ok = c.config.(map[string]interface{})[c.item.Name]; !ok {
		c.config.(map[string]interface{})[c.item.Name] = make([]interface{}, 0)
	}

	if value, ok = c.config.(map[string]interface{})[c.item.Name]; ok {
		valueArray, ok := value.([]interface{})
		if !ok {
			c.config.(map[string]interface{})[c.item.Name] = make([]interface{}, 0)
			valueArray = c.config.(map[string]interface{})[c.item.Name].([]interface{})
		}

		rowCount := len(valueArray)
		for i := 0; i < rowCount; i++ {
			item := c.lvItems.AddItem("--")
			arrayItem := valueArray[i]
			arrayItemAsMap, ok := arrayItem.(map[string]interface{})
			if !ok {
				valueArray[i] = make(map[string]interface{})
				arrayItemAsMap, _ = arrayItem.(map[string]interface{})
			}

			for colIndex, value := range colNameByIndex {
				var cellValue string
				cellValueInterface, ok := arrayItemAsMap[value]
				if ok {
					stringV, ok := cellValueInterface.(string)
					if ok {
						cellValue = stringV
					} else {
						stringer, ok := cellValueInterface.(fmt.Stringer)
						if ok {
							cellValue = stringer.String()
						} else {
							cellValue = fmt.Sprint(cellValueInterface)
						}
					}
				}
				item.SetValue(colIndex, cellValue)
			}
		}
	}
}

func (c *PanelUnitConfigItem) loadTable() {
	colIndexByName := make(map[string]int)
	colNameByIndex := make(map[int]string)

	for colIndex, col := range c.item.Children {
		colIndexByName[col.Name] = colIndex
		colNameByIndex[colIndex] = col.Name
	}

	var value interface{}
	var ok bool
	if _, ok = c.config.(map[string]interface{})[c.item.Name]; !ok {
		c.config.(map[string]interface{})[c.item.Name] = make([]interface{}, 0)
	}

	if value, ok = c.config.(map[string]interface{})[c.item.Name]; ok {
		valueArray, ok := value.([]interface{})
		if !ok {
			c.config.(map[string]interface{})[c.item.Name] = make([]interface{}, 0)
			valueArray = c.config.(map[string]interface{})[c.item.Name].([]interface{})
		}

		rowCount := len(valueArray)
		for i := 0; i < rowCount; i++ {
			if i >= c.lvItems.ItemsCount() {
				continue
			}
			item := c.lvItems.Item(i)
			arrayItem := valueArray[i]
			arrayItemAsMap, ok := arrayItem.(map[string]interface{})
			if !ok {
				valueArray[i] = make(map[string]interface{})
				arrayItemAsMap, _ = arrayItem.(map[string]interface{})
			}

			for colIndex, value := range colNameByIndex {
				var cellValue string
				cellValueInterface, ok := arrayItemAsMap[value]
				if ok {
					stringV, ok := cellValueInterface.(string)
					if ok {
						cellValue = stringV
					} else {
						stringer, ok := cellValueInterface.(fmt.Stringer)
						if ok {
							cellValue = stringer.String()
						} else {
							cellValue = fmt.Sprint(cellValueInterface)
						}
					}
				}
				item.SetValue(colIndex, cellValue)
			}
		}
	}
}

func (c *PanelUnitConfigItem) NotifyChanges() {
	if c.OnChanged != nil {
		c.OnChanged()
	}
}
