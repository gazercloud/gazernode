package common_interfaces

import (
	"time"
)

type IDataStorage interface {
	SetItem(name string, value string, UOM string, dt time.Time, flags string) error
	TouchItem(name string) error
	GetItem(name string) (Item, error)
	GetUnitValues(unitId string) []ItemGetUnitItems
	RenameItems(oldPrefix string, newPrefix string)
}
