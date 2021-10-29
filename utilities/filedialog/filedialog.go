package filedialog

import (
	"fmt"
	"github.com/gazercloud/gazernode/utilities"
	"github.com/gazercloud/gazernode/utilities/paths"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"golang.org/x/sys/windows"
	"runtime"
	"strings"
	"time"
)

type FileDialog struct {
	uicontrols.Dialog

	currentPath      []string
	isSaveFileDialog bool
	defaultFileName  string
	FileName         string

	lvLeft  *uicontrols.ListView
	lvRight *uicontrols.ListView

	lblCurrentPath *uicontrols.TextBlock
	txtFileName    *uicontrols.TextBox

	btnOK *uicontrols.Button
}

func NewFileDialog(parent uiinterfaces.Widget, isSaveFileDialog bool, defaultFileName string) *FileDialog {
	var c FileDialog
	c.isSaveFileDialog = isSaveFileDialog
	c.defaultFileName = defaultFileName
	c.InitControl(parent, &c)
	c.Resize(700, 550)
	if c.isSaveFileDialog {
		c.SetTitle("Save File")
	} else {
		c.SetTitle("Open File")
	}
	c.currentPath = make([]string, 0)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.lblCurrentPath = pContent.AddTextBlockOnGrid(0, 0, "current path: "+c.buildCurrentPath())

	pLVs := pContent.AddPanelOnGrid(0, 1)
	pLVs.SetPanelPadding(0)

	c.lvLeft = pLVs.AddListViewOnGrid(0, 0)
	c.lvLeft.AddColumn("-", 150)
	c.lvLeft.SetMaxWidth(160)

	c.lvLeft.OnItemClicked = func(item *uicontrols.ListViewItem) {
		path, ok := item.UserData("info").([]string)
		if ok {
			c.currentPath = path
			c.lblCurrentPath.SetText(c.buildCurrentPath())
			c.loadCurrentDirectory()
		}
	}

	c.lvRight = pLVs.AddListViewOnGrid(1, 0)
	c.lvRight.AddColumn("Name", 200)
	c.lvRight.AddColumn("Date modified", 140)
	c.lvRight.AddColumn("Size", 100)
	c.lvRight.OnItemDblClicked = func(item *uicontrols.ListViewItem) {
		c.ItemAction()
	}

	c.lvRight.OnItemClicked = func(item *uicontrols.ListViewItem) {
		c.ItemClick()
	}

	pFileName := pContent.AddPanelOnGrid(0, 2)
	pFileName.SetPanelPadding(0)

	pFileName.AddTextBlockOnGrid(0, 0, "File name:")
	c.txtFileName = pFileName.AddTextBoxOnGrid(1, 0)
	c.txtFileName.SetText(c.defaultFileName)

	pButtons := pContent.AddPanelOnGrid(0, 3)
	pButtons.SetPanelPadding(0)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		filename := c.buildFullFileName()
		c.FileName = filename
		return true
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.loadLeftList()
	c.loadCurrentDirectory()

	return &c
}

func (c *FileDialog) buildFullFileName() string {
	result := c.buildCurrentPath()
	if runtime.GOOS == "windows" {
		result += "\\"
	} else {
		result += "/"
	}
	result += c.txtFileName.Text()
	return result
}

func (c *FileDialog) buildCurrentPath() string {
	path := ""

	if runtime.GOOS == "windows" {
		for i, p := range c.currentPath {
			if i > 0 {
				path += "\\"
			}
			path += p
		}
	} else {
		for _, p := range c.currentPath {
			path += "/"
			path += p
		}
	}

	return path
}

func (c *FileDialog) updateLeftListSelected() {
	if len(c.currentPath) == 0 {
		c.lvLeft.UnselectAllItems()
		return
	}

	index := -1
	indexCountOfMatches := 0

	for i := 0; i < c.lvLeft.ItemsCount(); i++ {
		pathParts := c.lvLeft.Item(i).UserData("info").([]string)
		countOfMatches := 0
		for j := 0; j < len(pathParts); j++ {
			if j >= len(c.currentPath) || c.currentPath[j] != pathParts[j] {
				countOfMatches = 0
				break
			}
			countOfMatches++
		}

		if countOfMatches > indexCountOfMatches {
			indexCountOfMatches = countOfMatches
			index = i
		}
	}

	if index > -1 {
		c.lvLeft.SelectItem(index)
	}
}

func (c *FileDialog) ItemClick() {
	if len(c.lvRight.SelectedItems()) != 1 {
		return
	}
	item := c.lvRight.SelectedItems()[0]
	entry, ok := item.UserData("info").(utilities.FileInfo)
	if !ok {
		return
	}

	if !entry.Dir {
		c.txtFileName.SetText(entry.Name)
	}
}

func (c *FileDialog) ItemAction() {
	if len(c.lvRight.SelectedItems()) != 1 {
		return
	}

	item := c.lvRight.SelectedItems()[0]
	entry, ok := item.UserData("info").(utilities.FileInfo)
	if !ok {
		if len(c.currentPath) > 0 {
			c.currentPath = c.currentPath[:len(c.currentPath)-1]
		}

		c.lblCurrentPath.SetText(c.buildCurrentPath())
		c.loadCurrentDirectory()
		return
	}

	if entry.Dir {
		// INTO DIRECTORY

		c.currentPath = append(c.currentPath, entry.Name)

		c.lblCurrentPath.SetText(c.buildCurrentPath())
		c.loadCurrentDirectory()
		return
	}

	c.txtFileName.SetText(entry.Name)
}

func (c *FileDialog) fileSizeToString(size int64) string {
	if size < 1024 {
		return fmt.Sprint(size)
	}
	if size >= 1024 && size < 1024*1024 {
		return fmt.Sprint(size/1024) + "k"
	}
	if size >= 1024*1024 && size < 1024*1024*1024 {
		return fmt.Sprint(size/(1024*1024)) + "M"
	}
	if size >= 1024*1024*1024 {
		return fmt.Sprint(size/(1024*1024*1024)) + "G"
	}
	return ""
}

func (c *FileDialog) loadCurrentDirectory() {
	dirs, err := c.getDir(c.buildCurrentPath())
	c.lvRight.RemoveItems()

	if len(c.currentPath) > 0 {
		lvItem := c.lvRight.AddItem("[..]")
		lvItem.SetUserData("info", "..")
	}

	if err == nil {
		for _, entry := range dirs {
			if entry.Dir {
				lvItem := c.lvRight.AddItem("[" + entry.Name + "]")
				lvItem.SetUserData("info", entry)
				lvItem.SetValue(1, entry.Date.Format("2006-01-02 15:04"))
				lvItem.SetValue(2, "DIR")
			} else {
				lvItem := c.lvRight.AddItem(entry.Name)
				lvItem.SetUserData("info", entry)
				lvItem.SetValue(1, entry.Date.Format("2006-01-02 15:04"))
				lvItem.SetValue(2, c.fileSizeToString(entry.Size))
			}
		}
	}

	c.updateLeftListSelected()
}

func (c *FileDialog) addItemToLeftList(name string, path string) {
	pathParts := strings.FieldsFunc(strings.ReplaceAll(path, "\\", "/"), func(r rune) bool {
		return r == '/'
	})

	lvItem := c.lvLeft.AddItem(name)
	lvItem.SetUserData("info", pathParts)
}

func (c *FileDialog) loadLeftList() {
	c.lvLeft.RemoveItems()

	c.addItemToLeftList("Home", paths.HomeFolder())
	c.addItemToLeftList("Documents", paths.DocumentsFolder())
	c.addItemToLeftList("Pictures", paths.PicturesFolder())
	c.addItemToLeftList("Download", paths.DownloadsFolder())

	if runtime.GOOS == "windows" {
		drives := c.drives()
		for _, d := range drives {
			c.addItemToLeftList(d+":", d+":")
		}
	}
}

func (c *FileDialog) bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return
}

func (c *FileDialog) drives() []string {
	drives := make([]string, 0)
	drivesBits, err := windows.GetLogicalDrives()
	if err == nil {
		drives = c.bitsToDrives(drivesBits)
	}
	return drives
}

func (c *FileDialog) getDir(path string) (result []utilities.FileInfo, err error) {
	if path == "" && runtime.GOOS == "windows" {
		result = make([]utilities.FileInfo, 0)
		for _, d := range c.drives() {
			var fi utilities.FileInfo
			fi.Dir = true
			fi.Name = d + ":"
			fi.Size = 0
			fi.Ext = ""
			fi.Path = d + ":"
			fi.Attr = ""
			fi.Date = time.Now()
			result = append(result, fi)
		}
	} else {
		if runtime.GOOS == "windows" {
			path += "\\"
		}

		result, err = utilities.GetDir(path)
	}
	return
}
