package common_interfaces

import (
	"time"
)

type IDataStorage interface {
	SetItemByName(name string, value string, UOM string, dt time.Time, external bool) error
	SetAllItemsByUnitName(name string, value string, UOM string, dt time.Time, external bool) error
	TouchItem(name string) (*Item, error)
	GetItem(name string) (Item, error)
	GetUnitValues(unitId string) []ItemGetUnitItems
	RenameItems(oldPrefix string, newPrefix string)
	RemoveItemsOfUnit(unitId string) error
	SetProperty(itemName string, propName string, propValue string)
	SetPropertyIfDoesntExist(itemName string, propName string, propValue string)

	//Exec(function string, request []byte, host string) ([]byte, error)

	StatGazerNode() StatGazerNode
	StatGazerCloud() StatGazerCloud

	AddToWatch(unitId string, itemName string)
	RemoveFromWatch(unitId string, itemName string)
}
