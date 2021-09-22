package units

import (
	"fmt"
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"github.com/gazercloud/gazerui/uistyles"
	"golang.org/x/image/colornames"
)

type FormAddUnit struct {
	uicontrols.Dialog
	client          *client.Client
	txtUnitFilter   *uicontrols.TextBox
	pUnits          *uicontrols.Panel
	pCategories     *uicontrols.Panel
	txtStat         *uicontrols.TextBlock
	allowAccept     bool
	countOfFound    int
	offset          int
	btnLast         *uicontrols.Button
	btnCategories   []*uicontrols.Button
	currentCategory string

	UnitId string
}

func NewFormAddUnit(parent uiinterfaces.Widget, client *client.Client) *FormAddUnit {
	var c FormAddUnit
	c.client = client
	c.InitControl(parent, &c)
	c.SetName("FormAddUnit")

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(150)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pStat := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 2)

	c.txtStat = pStat.AddTextBlockOnGrid(0, 0, "---")
	c.txtStat.TextHAlign = canvas.HAlignCenter

	/*img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImageAdjusted("icons/material/image/drawable-hdpi/ic_blur_on_black_48dp.png", c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)*/
	c.pCategories = pLeft.AddPanelOnGrid(0, 1)
	c.btnCategories = make([]*uicontrols.Button, 0)

	pRight.AddTextBlockOnGrid(0, 0, "Unit type filter:")
	c.txtUnitFilter = pRight.AddTextBoxOnGrid(0, 1)
	c.txtUnitFilter.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		c.offset = 0
		c.updateUnits()
	}

	c.pUnits = pRight.AddPanelOnGrid(0, 2)
	c.pUnits.SetBorders(1, colornames.Gray)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	c.TryAccept = func() bool {
		if !c.allowAccept {
			if c.countOfFound == 1 {
				if c.btnLast != nil {
					c.btnLast.Press()
				}
			}
		}
		return c.allowAccept
	}

	c.OnShow = func() {
		c.txtUnitFilter.Focus()
	}

	c.updateUnits()
	c.updateCategories()

	return &c
}

func (c *FormAddUnit) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Add unit")
	c.Resize(600, 680)
}

func (c *FormAddUnit) Dispose() {
	c.client = nil
	c.txtUnitFilter = nil
	c.pUnits = nil
	c.btnLast = nil

	c.Dialog.Dispose()
}

func (c *FormAddUnit) UpdateStyle() {
	c.Dialog.UpdateStyle()
	c.updateUnits()
}

func (c *FormAddUnit) addCategory(index int, cat nodeinterface.UnitTypeCategoriesResponseItem, code string) {
	btn := c.pCategories.AddButtonOnGrid(0, index, "   "+cat.DisplayName+" ", func(event *uievents.Event) {
		t, ok := event.Sender.(*uicontrols.Button).UserData("key").(string)
		if ok {
			c.currentCategory = t
		}
		c.offset = 0
		c.updateCategoriesButtons()
		c.updateUnits()
	})
	btn.SetUserData("key", code)
	btn.SetUserData("category", cat)
	btn.SetTextImageVerticalOrientation(false)

	c.btnCategories = append(c.btnCategories, btn)
}

func (c *FormAddUnit) updateCategories() {
	c.currentCategory = ""
	c.pCategories.RemoveAllWidgets()
	c.btnCategories = make([]*uicontrols.Button, 0)
	c.pCategories.AddTextBlockOnGrid(0, 0, "loading ...")
	c.pCategories.AddVSpacerOnGrid(0, 1)
	c.client.UnitCategories(func(infos nodeinterface.UnitTypeCategoriesResponse, err error) {
		c.pCategories.RemoveAllWidgets()
		c.btnCategories = make([]*uicontrols.Button, 0)
		if err != nil {
			c.pCategories.AddTextBlockOnGrid(0, 0, err.Error())
		} else {
			maxI := 0
			for i, cat := range infos.Items {
				c.addCategory(i+1, cat, cat.Name)
				maxI = i
			}
			c.pCategories.AddVSpacerOnGrid(0, maxI+3)
		}
		c.updateCategoriesButtons()
	})
}

func (c *FormAddUnit) updateCategoriesButtons() {
	for _, btn := range c.btnCategories {
		btn.SetBackColor(colornames.Black)
		cat, ok := btn.UserData("category").(nodeinterface.UnitTypeCategoriesResponseItem)
		if ok {
			t, ok := btn.UserData("key").(string)
			if ok {
				btn.SetImageSize(32, 32)
				if t == c.currentCategory {
					btn.SetForeColor(uistyles.DefaultBackColor)
					btn.SetBackColor(c.ForeColor())
					if cat.Image != nil {
						btn.SetImage(uiresources.ImageFromBinAdjusted(cat.Image, uistyles.DefaultBackColor))
					}
				} else {
					btn.SetForeColor(nil)
					btn.SetBackColor(nil)
					if cat.Image != nil {
						btn.SetImage(uiresources.ImageFromBinAdjusted(cat.Image, btn.AccentColor()))
					}
				}

				btn.SetMouseCursor(ui.MouseCursorPointer)

			}
		}
	}
}

func (c *FormAddUnit) updateUnits() {
	countOnPage := 7

	c.client.UnitTypes(c.currentCategory, c.txtUnitFilter.Text(), c.offset, countOnPage, func(unitTypes nodeinterface.UnitTypeListResponse, err error) {
		c.pUnits.RemoveAllWidgets()

		c.txtStat.SetText(fmt.Sprintf("shown %d units out of %d", len(unitTypes.Types), unitTypes.InFilterCount))

		if err == nil {
			x := 0
			y := 0
			c.countOfFound = 0
			for _, unitType := range unitTypes.Types {
				c.countOfFound++
				btn := c.pUnits.AddButtonOnGrid(x, y, "  "+unitType.DisplayName, func(event *uievents.Event) {
					obj, ok := event.Sender.(*uicontrols.Button).UserData("object").(nodeinterface.UnitTypeListResponseItem)
					if ok {
						f := NewFormUnitEdit(c, c.client, "", obj.Type)
						f.ShowDialog()
						f.OnAccept = func() {
							logger.Println("OnAccept NewFormUnitEdit")
							c.UnitId = f.CreatedUnitId
							c.allowAccept = true
							c.TryAccept = nil
							c.Accept()
						}
					}
				})
				btn.SetImageSize(48, 48)
				btn.SetImage(uiresources.ImageFromBinAdjusted(unitType.Image, c.AccentColor()))
				btn.SetMinHeight(56)
				btn.SetMaxHeight(56)
				btn.SetMinWidth(390)
				btn.SetUserData("unitType", unitType.Type)
				btn.SetUserData("object", unitType)
				btn.SetTextImageVerticalOrientation(false)
				btn.SetMouseCursor(ui.MouseCursorPointer)
				c.btnLast = btn
				x++
				if x > 0 {
					x = 0
					y++
				}

			}
			c.pUnits.AddVSpacerOnGrid(0, countOnPage)
			pNavButtons := c.pUnits.AddPanelOnGrid(0, countOnPage+1)
			btnNavLeft := pNavButtons.AddButtonOnGrid(0, 0, "<", func(event *uievents.Event) {
				c.offset -= countOnPage
				if c.offset < 0 {
					c.offset = 0
				}
				c.updateUnits()
			})
			if c.offset == 0 {
				btnNavLeft.SetEnabled(false)
			}
			maxOffset := unitTypes.InFilterCount / countOnPage
			maxOffset *= countOnPage

			btnNavRight := pNavButtons.AddButtonOnGrid(1, 0, ">", func(event *uievents.Event) {
				c.offset += countOnPage
				if c.offset > maxOffset {
					c.offset = maxOffset
				}
				c.updateUnits()
			})

			if c.offset >= maxOffset {
				btnNavRight.SetEnabled(false)
			}

		}
	})

}
