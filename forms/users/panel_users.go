package users

import (
	"github.com/gazercloud/gazernode/client"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiresources"
	"time"
)

type PanelUsers struct {
	uicontrols.Panel
	client *client.Client

	lvUsers              *uicontrols.ListView
	firstTimeStateLoaded bool

	btnAdd     *uicontrols.Button
	btnEdit    *uicontrols.Button
	btnRemove  *uicontrols.Button
	btnRefresh *uicontrols.Button

	txtHeaderChartGroup *uicontrols.TextBlock

	btnRemoveSession *uicontrols.Button

	timer       *uievents.FormTimer
	lvSessions  *uicontrols.ListView
	currentUser string
}

func NewPanelUsers(parent uiinterfaces.Widget, client *client.Client) *PanelUsers {
	var c PanelUsers
	c.client = client
	c.InitControl(parent, &c)
	return &c
}

func (c *PanelUsers) OnInit() {
	//pHeader := c.AddPanelOnGrid(0, 0)
	//txtHeader := pHeader.AddTextBlockOnGrid(0, 0, "Users")
	//txtHeader.SetFontSize(24)

	pContent := c.AddPanelOnGrid(0, 1)
	pContent.SetPanelPadding(0)
	splitter := pContent.AddSplitContainerOnGrid(0, 0)
	splitter.SetPosition(360)
	splitter.SetYExpandable(true)

	pUnitsList := splitter.Panel1.AddPanelOnGrid(0, 0)
	pUnitsList.SetPanelPadding(0)

	txtHeader := pUnitsList.AddTextBlockOnGrid(0, 0, "Users")
	txtHeader.SetFontSize(24)

	pButtons := pUnitsList.AddPanelOnGrid(0, 1)
	pButtons.SetPanelPadding(0)

	c.btnAdd = pButtons.AddButtonOnGrid(0, 0, "", func(event *uievents.Event) {
		f := NewFormAddUser(c, c.client)
		f.ShowDialog()
		f.OnAccept = func() {
			c.loadUsers()
		}
	})
	c.btnAdd.SetTooltip("Add new user")

	c.btnEdit = pButtons.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
		f := NewFormEditUser(c, c.client, c.lvUsers.SelectedItem().UserData("userName").(string))
		f.ShowDialog()
		f.OnAccept = func() {
			c.loadUsers()
		}
	})
	c.btnEdit.SetTooltip("Edit user")

	c.btnRemove = pButtons.AddButtonOnGrid(2, 0, "", func(event *uievents.Event) {
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove user?", "Confirmation", func() {
			userName := c.lvUsers.SelectedItem().TempData
			c.client.UserRemove(userName, func(response nodeinterface.UserRemoveResponse, err error) {
				c.loadUsers()
			})
		}, nil)
	})
	c.btnRemove.SetTooltip("Remove user")

	pButtons.AddHSpacerOnGrid(3, 0)

	c.btnRefresh = pButtons.AddButtonOnGrid(4, 0, "", func(event *uievents.Event) {
		c.loadUsers()
	})
	c.btnRefresh.SetTooltip("Refresh")

	c.lvUsers = pUnitsList.AddListViewOnGrid(0, 2)
	c.lvUsers.AddColumn("Name", 200)
	c.lvUsers.OnSelectionChanged = c.loadSelected

	menu := uicontrols.NewPopupMenu(c.lvUsers)
	menu.AddItem("Remove user", func(event *uievents.Event) {
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove user?", "Confirmation", func() {
			userName := c.lvUsers.SelectedItem().TempData
			c.client.UserRemove(userName, func(response nodeinterface.UserRemoveResponse, err error) {
				c.loadUsers()
			})
		}, nil)
	}, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_open_in_browser_materialicons_48dp_1x_baseline_open_in_browser_black_48dp_png, c.ForeColor()), "")
	c.lvUsers.SetContextMenu(menu)

	pHeaderRight := splitter.Panel2.AddPanelOnGrid(0, 0)
	pHeaderRight.SetPanelPadding(0)
	c.txtHeaderChartGroup = pHeaderRight.AddTextBlockOnGrid(0, 0, "")
	c.txtHeaderChartGroup.SetFontSize(24)

	pItems := splitter.Panel2.AddPanelOnGrid(0, 1)
	pItems.SetPanelPadding(0)

	pButtonsRight := pItems.AddPanelOnGrid(0, 0)
	pButtonsRight.SetPanelPadding(0)

	// LINK CONTROL
	pLink := pButtonsRight.AddPanelOnGrid(0, 0)
	pLink.SetPanelPadding(0)

	pLink.AddHSpacerOnGrid(0, 1)
	c.btnRemoveSession = pLink.AddButtonOnGrid(3, 1, "", func(event *uievents.Event) {
		uicontrols.ShowQuestionMessageOKCancel(c, "Remove selected sessions?", "Confirmation", func() {
			items := c.SelectedUsers()
			for _, item := range items {
				c.client.SessionRemove(item, nil)
			}
		}, nil)
	})
	c.btnRemoveSession.SetTooltip("Remove session")

	c.lvSessions = pItems.AddListViewOnGrid(0, 1)
	c.lvSessions.AddColumn("User", 150)
	c.lvSessions.AddColumn("Date/Time", 200)
	c.lvSessions.AddColumn("SessionToken", 500)

	c.timer = c.Window().NewTimer(1000, c.timerUpdate)
	c.timer.StartTimer()

	c.loadUsers()
	c.UpdateStyle()
	c.loadSelected()
}

func (c *PanelUsers) Dispose() {
	c.client = nil

	c.lvUsers = nil

	c.btnAdd = nil
	c.btnEdit = nil
	c.btnRemove = nil

	c.btnRemoveSession = nil

	c.lvSessions = nil
	c.Panel.Dispose()
}

func (c *PanelUsers) loadSelected() {
	selectedItem := c.lvUsers.SelectedItem()
	if selectedItem != nil {
		name := c.lvUsers.SelectedItem().UserData("userName").(string)
		c.txtHeaderChartGroup.SetText("Sessions of user: " + name)
		c.SetCurrentUser(name)
	} else {
		c.txtHeaderChartGroup.SetText("no user selected")
		c.SetCurrentUser("")
	}
}

func (c *PanelUsers) FullRefresh() {
	c.loadUsers()
}

func (c *PanelUsers) loadUsers() {
	c.client.UserList(func(response nodeinterface.UserListResponse, err error) {
		if err != nil {
			return
		}
		if c.lvUsers == nil {
			return
		}

		c.firstTimeStateLoaded = true
		c.lvUsers.RemoveItems()
		for _, s := range response.Items {
			lvItem := c.lvUsers.AddItem(s)
			lvItem.TempData = s
			lvItem.SetUserData("userName", s)
		}
	})
}

func (c *PanelUsers) SelectedUsers() []string {
	items := make([]string, 0)
	for _, item := range c.lvSessions.SelectedItems() {
		name := item.TempData
		items = append(items, name)
	}
	return items
}

func (c *PanelUsers) SetCurrentUser(name string) {
	c.lvSessions.RemoveItems()
	if len(name) > 0 {
		c.currentUser = name
	} else {
		c.currentUser = name
	}
}

func (c *PanelUsers) UpdateStyle() {
	c.Panel.UpdateStyle()

	activeColor := c.AccentColor()
	inactiveColor := c.InactiveColor()

	c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, activeColor))
	c.btnEdit.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, activeColor))
	c.btnRemove.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, activeColor))
	c.btnRemoveSession.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, activeColor))

	c.btnAdd.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, inactiveColor))
	c.btnEdit.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_create_materialicons_48dp_1x_baseline_create_black_48dp_png, inactiveColor))
	c.btnRemove.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, inactiveColor))
	c.btnRemoveSession.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_clear_materialicons_48dp_1x_baseline_clear_black_48dp_png, inactiveColor))

	c.btnRefresh.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, activeColor))
	c.btnRefresh.SetImageDisabled(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, inactiveColor))
}

func (c *PanelUsers) timerUpdate() {
	if c.Disposed() {
		return
	}

	if !c.IsVisibleRec() {
		return
	}

	if !c.firstTimeStateLoaded {
		c.loadUsers()
	}

	if len(c.lvUsers.SelectedItems()) > 0 {
		if len(c.lvUsers.SelectedItems()) == 1 {
			c.btnEdit.SetEnabled(true)
			c.btnRemove.SetEnabled(true)
		} else {
			c.btnEdit.SetEnabled(false)
			c.btnRemove.SetEnabled(false)
		}
	} else {
		c.btnEdit.SetEnabled(false)
		c.btnRemove.SetEnabled(false)
	}

	itemsSelected := c.lvSessions.SelectedItems()

	if len(itemsSelected) > 0 {
		c.btnRemoveSession.SetEnabled(true)

	} else {
		c.btnRemoveSession.SetEnabled(false)
	}

	if len(c.currentUser) > 0 {
		c.client.SessionList(c.currentUser, func(response nodeinterface.SessionListResponse, err error) {

			if len(response.Items) != c.lvSessions.ItemsCount() {
				c.lvSessions.RemoveItems()
				for i := 0; i < len(response.Items); i++ {
					c.lvSessions.AddItem("---")
				}
			}
			for index, di := range response.Items {

				c.lvSessions.Item(index).TempData = di.SessionToken
				c.lvSessions.SetItemValue(index, 0, di.UserName)
				c.lvSessions.SetItemValue(index, 1, time.Unix(0, di.SessionOpenTime*1000).Format("2006-01-02 15-04-05"))
				c.lvSessions.SetItemValue(index, 2, di.SessionToken)
			}
		})
	}
}
