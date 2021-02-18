package forms

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type PanelUsers struct {
	uicontrols.Panel

	lvUsers *uicontrols.ListView
	client  *client.Client
}

func NewPanelUsers(parent uiinterfaces.Widget, client *client.Client) *PanelUsers {
	var c PanelUsers
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *PanelUsers) OnInit() {
	txtHeader := c.AddTextBlockOnGrid(0, 0, "Users")
	txtHeader.SetFontSize(24)

	pButtons := c.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)
	/*btnAdd := pButtons.AddButtonOnGrid(0, 0, "Add", func(event *uievents.Event) {
	})
	btnAdd.SetMinWidth(70)
	btnAdd.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_add_black_48dp.png", c.ForeColor()))
	btnEdit := pButtons.AddButtonOnGrid(1, 0, "Edit", func(event *uievents.Event) {
	})
	btnEdit.SetMinWidth(70)
	btnEdit.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_create_black_48dp.png", c.ForeColor()))
	btnRemove := pButtons.AddButtonOnGrid(2, 0, "Remove", func(event *uievents.Event) {
	})
	btnRemove.SetMinWidth(70)
	btnRemove.SetImage(uiresources.ResImageAdjusted("icons/material/content/drawable-hdpi/ic_clear_black_48dp.png", c.ForeColor()))
	*/
	pButtons.AddHSpacerOnGrid(3, 0)

	c.lvUsers = c.AddListViewOnGrid(0, 2)
	c.lvUsers.AddColumn("Name", 200)
	c.lvUsers.AddColumn("Id", 200)
}
