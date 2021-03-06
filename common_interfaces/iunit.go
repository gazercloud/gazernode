package common_interfaces

type IUnit interface {
	Id() string
	SetId(unitId string)
	Type() string
	SetType(unitType string)
	Name() string
	SetName(unitType string)
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
}
