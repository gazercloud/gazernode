package tree_items_parser

import "strings"

type TreeItem struct {
	FullName    string
	ShortName   string
	Children    []*TreeItem
	childrenMap map[string]*TreeItem
}

func newTreeItem(fullName string, shortName string) *TreeItem {
	var c TreeItem
	c.FullName = fullName
	c.ShortName = shortName
	c.Children = make([]*TreeItem, 0)
	c.childrenMap = make(map[string]*TreeItem, 0)
	return &c
}

func ParseItems(items []string) *TreeItem {
	rootItem := newTreeItem("", "Root")

	for _, fullName := range items {
		parts := strings.Split(fullName, "/")
		currentTreeItem := rootItem
		for index, p := range parts {
			if foundItem, ok := currentTreeItem.childrenMap[p]; ok {
				currentTreeItem = foundItem
			} else {
				finalFullName := ""
				if index == len(parts)-1 {
					finalFullName = fullName
				}
				nti := newTreeItem(finalFullName, p)
				currentTreeItem.childrenMap[p] = nti
				currentTreeItem.Children = append(currentTreeItem.Children, nti)

				currentTreeItem = nti
			}
		}
	}

	return rootItem
}
