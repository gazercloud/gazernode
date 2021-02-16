package application

type IApplication interface {
	AppName() string
	AppVersion() string
}
