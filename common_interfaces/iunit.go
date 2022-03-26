package common_interfaces

type IUnit interface {
	Init()
	Id() string
	SetId(unitId string)
	Type() string
	SetType(unitType string)
	DisplayName() string
	SetDisplayName(unitDisplayName string)
	SetIUnit(iUnit IUnit)
	MainItem() string
	Start(iDataStorage IDataStorage) error
	Stop()
	IsStarted() bool
	SetConfig(config string)
	GetConfig() string
	GetConfigMeta() string

	InternalUnitStart() error
	InternalUnitStop()
	ItemChanged(itemName string, value ItemValue)

	InternalInitItems()
	InternalDeInitItems()

	PropSet(props []ItemProperty)
	PropGet() []ItemProperty
	Prop(name string) string
	PropSetIfNotExists(name string, value string)
}
