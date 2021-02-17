package simplemap

type IMapDataSource interface {
	GetDataItemValue(path string, control IMapControl)
	LoadContent(itemUrl string, control IMapControl)
	GetWidgets(filter string, offset int, count int, toolbox IMapToolbox)
}
